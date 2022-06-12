package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const srtTimeRegexpTmpl = `^(?P<H>[0-9]{2}):(?P<M>[0-5][0-9]):(?P<S>[0-5][0-9])[.,](?P<MS>[0-9]{3})`

var srtTimeRegexp = regexp.MustCompile(srtTimeRegexpTmpl)
var srtTimeRegexpIdxByNames map[string]int

func init() {
	subexpNames := srtTimeRegexp.SubexpNames()
	srtTimeRegexpIdxByNames = make(map[string]int, len(subexpNames))
	for i, subexpName := range subexpNames {
		srtTimeRegexpIdxByNames[subexpName] = i
	}
}

type Subtitle struct {
	Sequence int
	Time     time.Duration
	Text     string
}

const ChanSizeSubtitles = 100

func ReadSubtitles(file io.ReadCloser, start, end time.Duration) (chan Subtitle, chan error) {
	subtitles := make(chan Subtitle, ChanSizeSubtitles)
	errs := make(chan error)

	go func() {
		defer close(errs)
		defer close(subtitles)

		bufReader := bufio.NewReaderSize(file, os.Getpagesize())
		firstRune, _, err := bufReader.ReadRune()
		if err != nil {
			errs <- fmt.Errorf("can't read subtitles file: %w", err)
			return
		}

		const utf8BOM = '\uFEFF'
		if firstRune != utf8BOM {
			err = bufReader.UnreadRune()
			if err != nil {
				panic("can't unread BOM bytes")
			}
		}

		var line string
		var lineNum int

		for {
			nextLine := func() bool {
				var lineBytes []byte
				for {
					bytes, isPrefix, err := bufReader.ReadLine()
					if err == io.EOF {
						return false
					}
					if err != nil {
						errs <- fmt.Errorf("error while reading subtitles file: %w", err)
						return false
					}
					lineBytes = append(lineBytes, bytes...)
					if !isPrefix {
						break
					}
				}
				line = string(lineBytes)
				lineNum++
				return true
			}

			subtitle := Subtitle{}

			if !nextLine() {
				return
			}
			subtitle.Sequence, err = strconv.Atoi(line)
			if err != nil {
				errs <- fmt.Errorf("error on line %d: expected to be a number: %w", lineNum, err)
				return
			}

			if !nextLine() {
				return
			}
			matches := srtTimeRegexp.FindStringSubmatch(line)
			if len(matches) == 0 {
				errs <- fmt.Errorf("error on line %d: expected to contain a display time", lineNum)
				return
			}
			hours, errH := strconv.Atoi(matches[srtTimeRegexpIdxByNames["H"]])
			minutes, errM := strconv.Atoi(matches[srtTimeRegexpIdxByNames["M"]])
			seconds, errS := strconv.Atoi(matches[srtTimeRegexpIdxByNames["S"]])
			milliseconds, errMS := strconv.Atoi(matches[srtTimeRegexpIdxByNames["MS"]])
			if errH != nil || errM != nil || errS != nil || errMS != nil {
				panic("wrong regexp")
			}
			subtitle.Time = time.Hour*time.Duration(hours) +
				time.Minute*time.Duration(minutes) +
				time.Second*time.Duration(seconds) +
				time.Millisecond*time.Duration(milliseconds)

			var textLines []string
			for nextLine() && strings.TrimSpace(line) != "" {
				textLines = append(textLines, line)
			}
			if len(textLines) == 0 {
				return
			}
			subtitle.Text = strings.Join(textLines, "\n")

			if subtitle.Time >= start && subtitle.Time < end {
				subtitles <- subtitle
			}
		}
	}()

	return subtitles, errs
}
