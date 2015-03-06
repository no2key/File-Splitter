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
		combineMain(os.Args[2:])
	case "md5", "sha1":
		var args []string
		args = append(args, "-hash", os.Args[1])
		args = append(args, os.Args[2:]...)
		hashMain(args)
	}
}

func printHelpInfo() {
	// TODO 添加使用说明
	fmt.Println(``)
}
