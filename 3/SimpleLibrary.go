package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

const (
	bufferSize = 100
)

const (
	_ = iota
	TakeBook
	ReturnBook
	GetAvailability
)

type SimpleLibrary struct {
	books                                   map[string]*Book
	registeredCopyCount, availableCopyCount map[string]int
	librarians, sync                        chan struct{}
}

type BookError struct {
	ISBN string
}

type TooManyCopiesBookError struct {
	BookError
}

type NotFoundBookError struct {
	BookError
}

type NotAvailableBookError struct {
	BookError
}

type AllCopiesAvailableBookError struct {
	BookError
}

func (e *TooManyCopiesBookError) Error() string {
	return fmt.Sprintf("Има 4 копия на книга %v", e.ISBN)
}

func (e *NotFoundBookError) Error() string {
	return fmt.Sprintf("Непозната книга %v", e.ISBN)
}

func (e *NotAvailableBookError) Error() string {
	return fmt.Sprintf("Няма наличност на книга %v", e.ISBN)
}

func (e *AllCopiesAvailableBookError) Error() string {
	return fmt.Sprintf("Всички копия са налични %v", e.ISBN)
}

func (sl *SimpleLibrary) addBook(book *Book) (registeredCopyCount int, err error) {
	if sl.registeredCopyCount[book.ISBN] >= 4 {
		err = &TooManyCopiesBookError{BookError{book.ISBN}}
	} else {
		sl.Books[book.ISBN] = book
		sl.registeredCopyCount[book.ISBN]++
		sl.availableCopyCount[book.ISBN]++
	}

	registeredCopyCount = sl.registeredCopyCount[book.ISBN]
	return
}

func (sl *SimpleLibrary) AddBookJSON(data []byte) (int, error) {
	book := &Book{}
	json.Unmarshal(data, book)
	return sl.addBook(book)
}

func (sl *SimpleLibrary) AddBookXML(data []byte) (int, error) {
	book := &Book{}
	xml.Unmarshal(data, book)
	return sl.addBook(book)
}

func (sl *SimpleLibrary) Hello() (chan<- LibraryRequest, <-chan LibraryResponse) {
	requests := make(chan LibraryRequest, bufferSize)
	responses := make(chan LibraryResponse, bufferSize)

	<-sl.librarians

	go func() {
		for request := range requests {
			isbn := request.GetISBN()

			<-sl.sync
			book, isBookRegistered := sl.books[isbn]
			sl.sync <- struct{}{}

			response := &SimpleLibraryResponse{}

			if !isBookRegistered {
				response.err = &NotFoundBookError{BookError{isbn}}
				responses <- response
				return
			}

			<-sl.sync

			switch request.GetType() {
			case TakeBook:
				if sl.availableCopyCount[isbn] > 0 {
					sl.availableCopyCount[isbn]--
					response.book = book
				} else {
					response.err = &NotAvailableBookError{BookError{isbn}}
				}

			case ReturnBook:
				if sl.availableCopyCount[isbn] < sl.registeredCopyCount[isbn] {
					sl.availableCopyCount[isbn]++
					response.book = book
				} else {
					response.err = &AllCopiesAvailableBookError{BookError{isbn}}
				}

			case GetAvailability:
				response.book = book
			}

			response.registeredCopyCount = sl.registeredCopyCount[isbn]
			response.availableCopyCount = sl.availableCopyCount[isbn]
			sl.sync <- struct{}{}

			responses <- response
		}

		sl.librarians <- struct{}{}
	}()

	return requests, responses
}

func NewLibrary(librarians int) Library {
	sl := &SimpleLibrary{
		Books:               make(map[string]*Book),
		registeredCopyCount: make(map[string]int),
		availableCopyCount:  make(map[string]int),
		librarians:          make(chan struct{}, librarians),
		sync:                make(chan struct{}, 1),
	}

	for i := 0; i < librarians; i++ {
		sl.librarians <- struct{}{}
	}

	return sl
}
