package main

import "strings"

func IsPalindrome(input string) bool {
	cleaned := valid.FindAllString(input, -1)
	words := strings.ToLower(strings.Join(cleaned, ""))
	runes := []rune(words)
	i, j := 0, len(runes)-1

	for i < j {
		if runes[i] != runes[j] {
			return false
		}
		i++
		j--
	}
	return true
}
