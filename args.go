package main

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

type Args struct {
	Start    time.Duration
	End      time.Duration
	Case     Case
	Scheme   string
	FileName string
}

func ParseArgs() Args {
	a := Args{
		Start:    0,
		End:      1000 * time.Hour,
		Case:     LowerCase,
		Scheme:   `$word\t$count\n`,
		FileName: "",
	}

	const timeFmtExplanation = "a sequence of decimal numbers with unit suffix, like \"100s\", \"2.3h\" or \"4h35m10s5ms\""

	startHelp := fmt.Sprintf(
		"the time starting from which (inclusive) subtitles should be parsed;\nformat: %s;\ndefault value: 0s.",
		timeFmtExplanation,
	)
	flag.Func("start", startHelp, func(val string) error {
		if start, err := time.ParseDuration(val); err != nil {
			return errors.New("start time is invalid")
		} else {
			a.Start = start
			return nil
		}
	})

	endHelp := fmt.Sprintf(
		"the time before which (exclusive) subtitles should be parsed;\nformat: %s;\ndefault value: 1000h.",
		timeFmtExplanation,
	)
	flag.Func("end", endHelp, func(val string) error {
		if end, err := time.ParseDuration(val); err != nil {
			return errors.New("end time is invalid")
		} else {
			a.End = end
			return nil
		}
	})

	caseHelp := "how should the letter case change;\nformat: one of: lower, upper, original, names-with-capital;\ndefault value: lower."
	flag.Func("case", caseHelp, func(val string) error {
		found := false
		for _, c := range AllCases {
			if c == Case(val) {
				found = true
			}
		}
		if !found {
			return errors.New("unsupported value")
		}
		a.Case = Case(val)
		return nil
	})

	schemeHelp := "a scheme of how to output the result for each word;\n" +
		"format: a string that can contain any ordinary characters, as well as special ones: \\n, \\r, \\t etc.\n" +
		"\tthis string can also contain optional arguments:\n" +
		"\t- $word argument is replaced by the current processed word from the subtitles file,\n" +
		"\t- $count argument is replaced by the total number of occurrences of this word in the subtitles file;\n" +
		"default value: $word\\t$count\\n."
	flag.Func("scheme", schemeHelp, func(val string) error {
		a.Scheme = val
		return nil
	})

	fileHelp := "the path to a subtitles file;\nformat: any valid path to an existing file;\nif it is not specified, stdin will be used."
	flag.Func("file", fileHelp, func(val string) error {
		a.FileName = val
		return nil
	})

	flag.Parse()

	return a
}

type Case string

var AllCases = []Case{LowerCase, UpperCase, OriginalCase, NamesWithCapital}

const (
	LowerCase        Case = "lower"
	UpperCase        Case = "upper"
	OriginalCase     Case = "original"
	NamesWithCapital Case = "names-with-capital"
)

