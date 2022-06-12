package main

import (
	"fmt"
	"strings"
)

func Print(result []WordResult, scheme string) {
	for _, wordResult := range result {
		scheme = strings.Replace(scheme, "$word", "%[1]s", -1)
		scheme = strings.Replace(scheme, "$count", "%[2]d", -1)
		scheme = strings.Replace(scheme, "\\n", "\n", -1)
		scheme = strings.Replace(scheme, "\\r", "\r", -1)
		scheme = strings.Replace(scheme, "\\t", "\t", -1)
		scheme = strings.Replace(scheme, "\\a", "\a", -1)
		fmt.Printf(scheme, wordResult.Word, wordResult.Count)
	}
}
