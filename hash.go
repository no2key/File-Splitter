// hash
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
)

func hashMain(args []string) {
	cmd := flag.NewFlagSet("hash", flag.ExitOnError)
	flagFile := cmd.String("f", "", "需要计算校验和的文件")
	flagHash := cmd.String("hash", "md5", "校验算法，默认md5")
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

	var h hash.Hash
	switch *flagHash {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	default:
		fmt.Fprintf(os.Stderr, "Neither md5 nor sha1 hash...")
		os.Exit(1)
	}
	m, err := copyAndHash(ioutil.Discard, file, h)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(m)
}
