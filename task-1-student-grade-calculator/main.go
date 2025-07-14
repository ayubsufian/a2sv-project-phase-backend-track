package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Subject struct {
	name  string
	grade float64
}

type Student struct {
	name     string
	subjects []Subject
}

func getInput(scanner *bufio.Scanner, prompt string) string {
	for {
		fmt.Print(prompt)
		scanner.Scan()
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			fmt.Println("Input can't be empty. Try again.")
			continue
		}
		return text
	}
}

func getNonNegativeInt(scanner *bufio.Scanner, prompt string) int {
	for {
		input := getInput(scanner, prompt)
		num, err := strconv.Atoi(input)
		if err != nil || num < 0 {
			fmt.Println("Invalid input. Please enter a non‑negative integer.")
			continue
		}
		return num
	}
}

func getGrade(scanner *bufio.Scanner, prompt string) float64 {
	for {
		input := getInput(scanner, prompt)
		g, err := strconv.ParseFloat(input, 64)
		if err != nil || g < 0 || g > 100 {
			fmt.Println("Invalid grade. Please enter a number between 0 and 100.")
			continue
		}
		return g
	}
}

func (s *Student) AddSubject(subject Subject) {
	s.subjects = append(s.subjects, subject)
}

func (s *Student) CalculateAverageGrade() float64 {
	if len(s.subjects) == 0 {
		return 0.0
	}

	total := 0.0
	for _, subject := range s.subjects {
		total += subject.grade
	}
	return total / float64(len(s.subjects))
}

func (s *Student) DisplayInformation() {
	fmt.Println("\n--- Student Report ---")
	fmt.Printf("Name: %s\n", s.name)
	fmt.Println("Subjects and Grades:")
	for _, subject := range s.subjects {
		fmt.Printf("  - %-20s : %.2f\n", subject.name, subject.grade)
	}
	fmt.Printf("\nAverage Grade: %.2f\n", s.CalculateAverageGrade())
}

func main() {
	var s1 Student
	scanner := bufio.NewScanner(os.Stdin)

	s1.name = getInput(scanner, "Please enter your name: ")
	n := getNonNegativeInt(scanner, "Please enter the number of subjects you have taken: ")
	if n == 0 {
		fmt.Println("\nYou haven’t entered any subjects. No grades to compute.")
		return
	}

	var sub Subject

	for i := 0; i < n; i++ {
		sub.name = getInput(scanner, "Please enter the name of the subject: ")
		sub.grade = getGrade(scanner, "Please enter the grade you obtained in the subject: ")
		s1.AddSubject(sub)
	}
	s1.DisplayInformation()
}
