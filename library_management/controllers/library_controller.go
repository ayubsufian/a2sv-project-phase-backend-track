package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

func promptInt(scanner *bufio.Scanner, prompt string) (int, bool) {
	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			return 0, false
		}
		text := strings.TrimSpace(scanner.Text())
		v, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}
		if v <= 0 {
			fmt.Println("Number must be positive.")
			continue
		}
		return v, true
	}
}

func promptString(scanner *bufio.Scanner, prompt string) (string, bool) {
	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			return "", false
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			fmt.Println("Input cannot be empty.")
			continue
		}
		return text, true
	}
}

func StartCLI(lib services.LibraryManager) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(`
Library Menu:
1) Add book
2) Remove book
3) Add member
4) Borrow book
5) Return book
6) List available books
7) List books borrowed by member
0) Exit`)
		fmt.Println()
		fmt.Print("Enter choice: ")
		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			addBook(lib, scanner)
		case "2":
			removeBook(lib, scanner)
		case "3":
			addMember(lib, scanner)
		case "4":
			borrowBook(lib, scanner)
		case "5":
			returnBook(lib, scanner)
		case "6":
			listAvailable(lib)
		case "7":
			listBorrowed(lib, scanner)
		case "0":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func addBook(lib services.LibraryManager, scanner *bufio.Scanner) {
	bookID, ok := promptInt(scanner, "Enter book ID: ")
	if !ok {
		return
	}

	title, ok := promptString(scanner, "Enter title: ")
	if !ok {
		return
	}

	author, ok := promptString(scanner, "Enter author: ")
	if !ok {
		return
	}

	book := models.Book{ID: bookID, Title: title, Author: author, Status: "Available"}
	lib.AddBook(book)
	fmt.Println("Book added.")
}

func removeBook(lib services.LibraryManager, scanner *bufio.Scanner) {
	id, ok := promptInt(scanner, "Enter book ID to remove: ")
	if !ok {
		return
	}
	lib.RemoveBook(id)
	fmt.Println("Book removed (if it existed).")
}

func addMember(lib services.LibraryManager, scanner *bufio.Scanner) {
	id, ok := promptInt(scanner, "Enter member ID: ")
	if !ok {
		return
	}

	name, ok := promptString(scanner, "Enter name: ")
	if !ok {
		return
	}

	member := models.Member{ID: id, Name: name}
	if err := lib.AddMember(member); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Member added.")
	}
}

func borrowBook(lib services.LibraryManager, scanner *bufio.Scanner) {
	bookID, ok := promptInt(scanner, "Enter book ID to borrow: ")
	if !ok {
		return
	}
	memberID, ok := promptInt(scanner, "Enter your member ID: ")
	if !ok {
		return
	}

	if err := lib.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed.")
	}
}

func returnBook(lib services.LibraryManager, scanner *bufio.Scanner) {
	bookID, ok := promptInt(scanner, "Enter book ID to return: ")
	if !ok {
		return
	}
	memberID, ok := promptInt(scanner, "Enter your member ID: ")
	if !ok {
		return
	}

	if err := lib.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned.")
	}
}

func listAvailable(lib services.LibraryManager) {
	books := lib.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	fmt.Println("Available books:")
	for _, b := range books {
		fmt.Printf("#%d: %s by %s\n", b.ID, b.Title, b.Author)
	}
}

func listBorrowed(lib services.LibraryManager, scanner *bufio.Scanner) {
	memberID, ok := promptInt(scanner, "Enter member ID to list borrowed books: ")
	if !ok {
		return
	}

	books, err := lib.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(books) == 0 {
		fmt.Println("No books borrowed by this member.")
		return
	}
	fmt.Printf("Books borrowed by member with ID %d:\n", memberID)
	for _, b := range books {
		fmt.Printf("#%d: %s by %s\n", b.ID, b.Title, b.Author)
	}
}
