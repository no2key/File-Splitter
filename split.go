// split
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

type fileHeader struct {
	total  int64
	offset int64
	length int64
	id     int

	file *os.File
}

func newFileHeader(file *os.File, id int, maxId int, blockSize int64, totalSize int64) *fileHeader {
	var length int64

	if id == maxId {
		length = totalSize % totalSize
	} else {
		length = blockSize
	}

	return &fileHeader{
		totalSize,
		int64(id) * blockSize,
		length,
		id,
		file,
	}
}

func (self *fileHeader) output() {
	output_file, err := os.Create(fmt.Sprintf("%s.split.%d", self.file.Name(), self.id))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer output_file.Close()

	bufW := bufio.NewWriter(output_file)
	_, err = io.Copy(bufW, io.NewSectionReader(self.file, self.offset, self.length))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		err := bufW.Flush()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("生成", output_file.Name())
	}
}

func splitMain(args []string) {
	cmd := flag.NewFlagSet("split", flag.ExitOnError)

	var flag_length = cmd.Int64("len", 650, "[可选]单位文件的长度，默认 650M")
	var flag_file = cmd.String("f", "", "需要切割的文件")

	err := cmd.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(*flag_file) == 0 {
		cmd.PrintDefaults()
		os.Exit(1)
	}

	if *flag_length < 0 {
		cmd.PrintDefaults()
		os.Exit(1)
	}

	*flag_length <<= 20

	file, err := os.Open(*flag_file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "\n")
		os.Exit(1)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "\n")
		os.Exit(1)
	}

	// 文件的总长度
	total_length := info.Size()
	if total_length <= *flag_length {
		fmt.Println(os.Stderr, "文件长度过短，不需要切割\n")
		os.Exit(1)
	}

	// 遍历的次数
	maxId := int(math.Ceil(float64(total_length) / float64(*flag_length)))
	wg := sync.WaitGroup{}
	wg.Add(maxId)

	for i := 0; i < maxId; i++ {
		go func(id int) {
			defer wg.Done()
			newFileHeader(file, id, maxId, *flag_length, total_length).output()
		}(i)
	}

	wg.Wait()
	fmt.Println("输出成功！")
}
