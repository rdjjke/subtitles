package main

import "strings"

type Word string

func SplitToWords(subtitles chan Subtitle, letterCase Case) chan Word {
	words := make(chan Word)
	go func() {
		defer close(words)
		for subtitle := range subtitles {
			ws := splitSentence(subtitle.Text)
			for _, w := range ws {
				if w == "" {
					break
				}
				switch letterCase {
				case CaseOriginal:
					// Do not change
				case CaseLower:
					w = strings.ToLower(w)
				case CaseUpper:
					w = strings.ToUpper(w)
				}
				words <- Word(w)
			}
		}
	}()
	return words
}

var delimeters = []rune(" \r\n\t.,;:!?-\"")

func splitSentence(sentence string) []string {
	var words []string
	var currWord []rune
	for _, r := range sentence {
		isDelimeter := false
		for _, d := range delimeters {
			if r == d {
				isDelimeter = true
				break
			}
		}
		if isDelimeter {
			words = append(words, string(currWord))
			currWord = currWord[:0]
		} else {
			currWord = append(currWord, r)
		}
	}
	return words
}
