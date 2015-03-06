// combine
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type sliceFile struct {
	abs  string
	id   int
	size int64
}

func (self *sliceFile) output(f *os.File) error {
	r, err := os.Open(self.abs)
	if err != nil {
		return err
	}
	defer r.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	if n != self.size {
		return fmt.Errorf("size is not %d but %d", self.size, n)
	}

	return nil
}

func (self *sliceFile) String() string {
	return self.abs
}

type sliceFiles []*sliceFile

func (self sliceFiles) Len() int {
	return len(self)
}

func (self sliceFiles) Less(i, j int) bool {
	return self[i].id < self[j].id
}

func (self sliceFiles) Swap(i, j int) {
	temp := self[i]
	self[i] = self[j]
	self[j] = temp
}

func combineMain(args []string) {
	cmd := flag.NewFlagSet("combine", flag.ExitOnError)
	flagFile := cmd.String("f", "", "需要合并的文件")
	err := cmd.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cmd.PrintDefaults()
		os.Exit(1)
	}

	abs, err := filepath.Abs(*flagFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	parentPath := filepath.Dir(abs)
	bp := filepath.Base(*flagFile)

	regx := regexp.MustCompile(`([^<>\?/\\\*]+)\.split\.(\d+)`)
	groups := regx.FindStringSubmatch(bp)
	if len(groups) == 0 {
		fmt.Fprintf(os.Stderr, "%s 不是一个切分文件！", abs)
		os.Exit(1)
	}

	baseName := groups[1]
	dir, err := os.Open(parentPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	regx = regexp.MustCompile(baseName + `\.split\.(\d+)`)
	var fileFragments []*sliceFile
	infos, err := dir.Readdir(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var id int
	var total int64
	var size int64
	var block int64
	for _, info := range infos {
		groups = regx.FindStringSubmatch(info.Name())
		if len(groups) == 0 {
			continue
		}

		id, err = strconv.Atoi(groups[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		size = info.Size()

		if block < size {
			block = size
		}

		fileFragments = append(fileFragments, &sliceFile{
			filepath.Join(parentPath, info.Name()),
			id,
			size,
		})

		total += size
	}

	sort.Sort(sliceFiles(fileFragments))
	for i := range fileFragments {
		if fileFragments[i].id != i {
			fmt.Fprintln(os.Stderr, "合并文件缺失",
				fileFragments[i].abs, " 不是第", i, "个")
			os.Exit(1)
		}

		if fileFragments[i].size != block && i != len(fileFragments)-1 {
			fmt.Fprintln(os.Stderr, "分割文件大小不符合要求 expext:",
				block, "but:", fileFragments[i].size)
			os.Exit(1)
		}
	}
	fmt.Println("碎片文件校验成功……\n总大小：",
		float64(total)/1024.0/1024.0, "M")

	output, err := os.Create(baseName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer output.Close()

	for _, f := range fileFragments {
		fmt.Println("合并", f.abs, "……")
		if err = f.output(output); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Println("合并成功！")
}
