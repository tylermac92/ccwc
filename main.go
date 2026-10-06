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
	words int64
	chars int64
	bytes int64
}

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	}
	return false
}

func count(r io.Reader) (counts, error) {
	var c counts
	inWord := false
	br := bufio.NewReader(r)
	buf := make([]byte, 32*1024)
	for {
		n, err := br.Read(buf)
		c.bytes += int64(n)
		c.lines += int64(bytes.Count(buf[:n], []byte{'\n'}))
		for _, b := range buf[:n] {
			// Count every byte that is not a UTF-8 continuation byte, so a
			// multibyte character split across reads is still counted once.
			if b&0xC0 != 0x80 {
				c.chars++
			}
			if isSpace(b) {
				inWord = false
			} else if !inWord {
				inWord = true
				c.words++
			}
		}
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
	countWords := flag.Bool("w", false, "print the word counts")
	countChars := flag.Bool("m", false, "print the character counts")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: ccwc [-c] [-l] [-m] [-w] file")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 || (!*countBytes && !*countLines && !*countWords && !*countChars) {
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

	// Same column order as wc: lines, words, characters, then bytes.
	if *countLines {
		fmt.Printf("%8d", c.lines)
	}
	if *countWords {
		fmt.Printf("%8d", c.words)
	}
	if *countChars {
		fmt.Printf("%8d", c.chars)
	}
	if *countBytes {
		fmt.Printf("%8d", c.bytes)
	}
	fmt.Printf(" %s\n", name)
}
