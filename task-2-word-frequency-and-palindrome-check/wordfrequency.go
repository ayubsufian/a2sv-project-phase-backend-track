package main

import (
	"regexp"
	"strings"
)

var valid = regexp.MustCompile(`[A-Za-z0-9]+`)

func CountWords(input string) map[string]int {
	words := valid.FindAllString(input, -1)
	counts := make(map[string]int)

	for _, word := range words {
		word = strings.ToLower(word)
		counts[word]++
	}
	return counts
}
