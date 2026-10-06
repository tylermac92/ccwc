package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type counts struct {
	lines int64
	bytes int64
}

func count(r io.Reader) (counts, error) {
	var c counts
	br := bufio.NewReader(r)
	buf := make([]byte, 32*1024)
	for {
		n, err := br.Read(buf)
		c.bytes += int64(n)
		c.lines += int64(bytes.Count(buf[:n], []byte{'\n'}))
		if errors.Is(err, io.EOF) {
			return c, nil
		}
		if err != nil {
			return c, err
		}
	}
}

func main() {
	countBytes := flag.Bool("c", false, "print the byte counts")
	countLines := flag.Bool("l", false, "print the newline counts")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: ccwc [-c] [-l] file")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 || (!*countBytes && !*countLines) {
		flag.Usage()
		os.Exit(2)
	}

	name := flag.Arg(0)
	f, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ccwc: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	c, err := count(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ccwc: %s: %v\n", name, err)
		os.Exit(1)
	}

	// Same column order as wc: lines, then bytes.
	if *countLines {
		fmt.Printf("%8d", c.lines)
	}
	if *countBytes {
		fmt.Printf("%8d", c.bytes)
	}
	fmt.Printf(" %s\n", name)
}
