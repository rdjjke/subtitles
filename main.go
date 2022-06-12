package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	args := ParseArgs()

	filePath := args.FilePath
	var file io.ReadCloser
	if filePath == "" {
		file = os.Stdin
		filePath = "stdin"
	} else {
		var err error
		file, err = os.Open(filePath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "can't open: %s", err.Error())
			os.Exit(1)
		}
		defer func() { _ = file.Close() }()
	}

	subtitles, errors := ReadSubtitles(file, args.Start, args.End)

	words := SplitToWords(subtitles, args.Case)

	result := MakeResult(words)

	Print(result, args.Scheme)

	err, ok := <-errors
	if ok {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
