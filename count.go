package main

import "sort"

type WordResult struct {
	Word  string
	Count int
}

func MakeResult(words chan Word) []WordResult {
	wordCounts := map[Word]int{}
	for word := range words {
		wordCounts[word]++
	}
	var result []WordResult
	for word, count := range wordCounts {
		result = append(result, WordResult{
			Word:  string(word),
			Count: count,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}
