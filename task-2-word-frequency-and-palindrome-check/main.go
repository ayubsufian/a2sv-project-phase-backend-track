package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	input := getInput("Enter a string for word frequency count: ")
	fmt.Println("Word Frequencies: ", CountWords(input))

	input = getInput("Enter a string to check palindrome: ")
	if IsPalindrome(input) {
		fmt.Println("It is a palindrome")
	} else {
		fmt.Println("It is not a palindrome")
	}
}
