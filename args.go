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
	FilePath string
}

func ParseArgs() Args {
	a := Args{
		Start:    0,
		End:      1000 * time.Hour,
		Case:     CaseLower,
		Scheme:   `$word\t$count\n`,
		FilePath: "",
	}

	const timeFmtExpl = "a sequence of decimal numbers with unit suffix, like \"100s\", \"2.3h\" or \"4h35m10s5ms\""

	startHelp := fmt.Sprintf(
		"the time starting from which (inclusive) subtitles should be parsed;\nformat: %s;\ndefault value: 0s.",
		timeFmtExpl,
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
		timeFmtExpl,
	)
	flag.Func("end", endHelp, func(val string) error {
		if end, err := time.ParseDuration(val); err != nil {
			return errors.New("end time is invalid")
		} else {
			a.End = end
			return nil
		}
	})

	caseHelp := fmt.Sprintf(
		"how should the letter case change;\nformat: one of %v;\ndefault value: %v.",
		AllCases,
		CaseLower,
	)
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
		fmt.Sprintf("default value: %s.", a.Scheme)
	flag.Func("scheme", schemeHelp, func(val string) error {
		a.Scheme = val
		return nil
	})

	fileHelp := "the path to a subtitles file;\n" +
		"format: any valid path to an existing file;\n" +
		"if it is not specified, stdin will be used."
	flag.Func("file", fileHelp, func(val string) error {
		a.FilePath = val
		return nil
	})

	flag.Parse()

	return a
}

type Case string

var AllCases = []Case{
	CaseLower,
	CaseUpper,
	CaseOriginal,
	/*CaseNamesWithCapital,*/
}

const (
	CaseLower    Case = "lower"
	CaseUpper    Case = "upper"
	CaseOriginal Case = "original"
	//CaseNamesWithCapital Case = "names-with-capital"
)
