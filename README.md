# ccwc

A small clone of the Unix `wc` command, written in Go. It counts the lines,
words, characters and bytes in a file or in standard input.

## Requirements

- Go 1.24 or later

## Build

```sh
go build -o ccwc .
```

This produces a `ccwc` binary in the current directory. You can also run it
without building:

```sh
go run . -l test.txt
```

## Usage

```
ccwc [-c] [-l] [-m] [-w] [file]
```

| Option | Description                                  |
| ------ | -------------------------------------------- |
| `-c`   | Print the number of bytes                    |
| `-l`   | Print the number of lines (newline characters) |
| `-w`   | Print the number of words                    |
| `-m`   | Print the number of characters               |

- With **no options**, `ccwc` behaves like `-l -w -c` and prints lines, words
  and bytes.
- With **no file**, `ccwc` reads from standard input and prints no filename.
- Options may be combined by passing each separately (`-l -w`). Output is
  always in the same order as `wc`: lines, words, characters, bytes,
  regardless of the order the flags were given in.
- Options must come before the filename.

Each count is right-aligned in a column 8 characters wide, followed by the
filename if one was given.

## Examples

The repository includes `test.txt`, a sample file with 7145 lines, 58164
words, 339292 characters and 342190 bytes.

```
$ ccwc -c test.txt
  342190 test.txt

$ ccwc -l test.txt
    7145 test.txt

$ ccwc -w test.txt
   58164 test.txt

$ ccwc -m test.txt
  339292 test.txt

$ ccwc test.txt
    7145   58164  342190 test.txt

$ cat test.txt | ccwc -l
    7145
```

## How counting works

- **Bytes** are the total number of bytes read.
- **Lines** are the number of `\n` characters, so a final line without a
  trailing newline is not counted (the same as `wc -l`).
- **Words** are runs of non-whitespace characters separated by ASCII
  whitespace (space, tab, newline, vertical tab, form feed, carriage return).
- **Characters** are counted as UTF-8 characters: every byte that is not a
  UTF-8 continuation byte counts as one character. This means a multibyte
  character split across two reads is still counted once.

The input is streamed in 32 KiB chunks, so large files and pipes are handled
without loading them into memory.

## Differences from `wc`

- Only a single file (or standard input) is accepted; multiple files and a
  `total` line are not supported.
- `-m` always counts UTF-8 characters. It does not consult the locale, so it
  never falls back to matching `-c` the way `wc -m` does in a non-UTF-8
  locale.
- Column widths are fixed at 8, whereas GNU `wc` sizes columns to fit the
  input.
- Flags cannot be bundled (use `-l -w`, not `-lw`), and the long options such
  as `--lines` are not supported.

## Exit status

| Status | Meaning                                           |
| ------ | ------------------------------------------------- |
| `0`    | Success                                           |
| `1`    | The file could not be opened or read              |
| `2`    | Invalid usage (for example, more than one file)   |
