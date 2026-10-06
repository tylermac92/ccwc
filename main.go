package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	countBytes := flag.Bool("c", false, "print the byte counts")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: ccwc [-c] file")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	if !*countBytes {
		flag.Usage()
		os.Exit(2)
	}

	name := flag.Arg(0)
	info, err := os.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ccwc: %v\n", err)
		os.Exit(1)
	}
	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "ccwc: %s: is a directory\n", name)
		os.Exit(1)
	}

	fmt.Printf("%8d %s\n", info.Size(), name)
}
