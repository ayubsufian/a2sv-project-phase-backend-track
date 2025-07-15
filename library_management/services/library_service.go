package services

import (
	"fmt"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID, memberID int) error
	ReturnBook(bookID, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
	AddMember(member models.Member) error
}

type Library struct {
	books   map[int]models.Book
	members map[int]*models.Member
}

func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]*models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	if book.Status == "Borrowed" {
		return fmt.Errorf("book with ID %d is already borrowed", book.ID)
	}
	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	book.Status = "Borrowed"
	l.books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	return nil
}

func (l *Library) ReturnBook(bookID, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	idx := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("member with ID %d did not borrow book with ID %d", memberID, bookID)
	}

	member.BorrowedBooks = append(member.BorrowedBooks[:idx], member.BorrowedBooks[idx+1:]...)

	book.Status = "Available"
	l.books[bookID] = book
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	availableBooks := []models.Book{}
	for _, book := range l.books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, ok := l.members[memberID]
	if !ok {
		return nil, fmt.Errorf("member with ID %d not found", memberID)
	}
	return member.BorrowedBooks, nil
}

func (l *Library) AddMember(member models.Member) error {
	if _, exists := l.members[member.ID]; exists {
		return fmt.Errorf("member with ID %d already exists", member.ID)
	}
	l.members[member.ID] = &member
	return nil
}
