// main
package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4 * runtime.NumCPU())

	if len(os.Args) < 2 {
		printHelpInfo()
		os.Exit(0)
	}

	switch os.Args[1] {
	default:
		printHelpInfo()
		os.Exit(0)
	case "split":
		splitMain(os.Args[2:])
	case "combine":
		// TODO 添加文件合并功能
	case "md5":
		md5Main(os.Args[2:])
	}
}

func printHelpInfo() {
	// TODO 添加使用说明
	fmt.Println(``)
}
