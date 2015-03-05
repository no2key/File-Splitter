// md5
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func md5Main(args []string) {
	cmd := flag.NewFlagSet("md5", flag.ExitOnError)
	flagFile := cmd.String("f", "", "需要计算校验和的文件")
	err := cmd.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cmd.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*flagFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	m, err := copyAndMd5(ioutil.Discard, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(m)
}
