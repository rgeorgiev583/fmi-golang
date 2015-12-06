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

type Library interface {

	// Добавя книга от json
	// Oтговаря с общия брой копия в библиотеката (не само наличните).
	// Aко са повече от 4 - връща грешка
	AddBookJSON(data []byte) (int, error)

	// Добавя книга от xml
	// Oтговаря с общия брой копия в библиотеката (не само наличните).
	// Ако са повече от 4 - връщаме грешка
	AddBookXML(data []byte) (int, error)

	// Ангажира свободен "библиотекар" да ни обработва заявките.
	// Библиотекарите са фиксиран брой - подават се като параметър на NewLibrary
	// Блокира ако всички библиотекари са заети.
	// Връщат се два канала:
	// първият е само за писане -  по него ще изпращаме заявките
	// вторият е само за четене - по него ще получаваме отговорите.
	// Ако затворим канала със заявките - освобождаваме библиотекаря.
	Hello() (chan<- LibraryRequest, <-chan LibraryResponse)
}

type LibraryRequest interface {
	// Тип на заявката:
	// 1 - Borrow book
	// 2 - Return book
	// 3 - Get availability information about book
	GetType() int

	// Връща isbn на книгата, за която се отнася Request-a
	GetISBN() string
}

type LibraryResponse interface {
	// Ако книгата съществува/налична е - обект имплементиращ Stringer (повече информация по-долу)
	// Aко книгата не съществува първият резултат е nil.
	// Връща се и подобаващa грешка (виж по-долу) - ако такава е възникнала.
	// Когато се е резултат на заявка от тип 2 (Return book) - не е нужно да я закачаме към отговора.
	GetBook() (fmt.Stringer, error)

	// available - Колко наличности от книгата имаме останали след изпълнението на заявката.
	// Тоест, ако сме имали 3 копия от Х и това е отговор на Take заявка - тук ще има 2.
	// registered - Колко копия от тази книга има регистрирани в библиотеката (макс 4).
	GetAvailability() (available int, registered int)
}

type Book struct {
	XMLName xml.Name `xml:"book"`
	ISBN    string   `json:"isbn" xml:"isbn,attr"`
	Title   string   `json:"title" xml:"title"`
	Author  struct {
		FirstName string `json:"first_name" xml:"first_name"`
		LastName  string `json:"last_name" xml:"last_name"`
	} `json:"author" xml:"author"`
	Ratings []uint8 `json:"ratings" xml:"ratings>rating"`
}

type SimpleLibrary struct {
	Books               map[string]*Book
	registeredCopyCount map[string]int
	availableCopyCount  map[string]int
	librarians          chan struct{}
}

type SimpleLibraryRequest struct {
	requestType int
	bookISBN    string
}

type SimpleLibraryResponse struct {
	book                *Book
	registeredCopyCount int
	availableCopyCount  int
	err                 error
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

func (b *Book) String() string {
	return fmt.Sprintf("[%v] %v от %v %v", b.ISBN, b.Title, b.Author.FirstName, b.Author.LastName)
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

func (r *SimpleLibraryRequest) GetType() int {
	return r.requestType
}

func (r *SimpleLibraryRequest) GetISBN() string {
	return r.bookISBN
}

func (r *SimpleLibraryRequest) SetType(t int) {
	r.requestType = t
}

func (r *SimpleLibraryRequest) SetISBN(isbn string) {
	r.bookISBN = isbn
}

func (r *SimpleLibraryResponse) GetBook() (fmt.Stringer, error) {
	return r.book, r.err
}

func (r *SimpleLibraryResponse) GetAvailability() (int, int) {
	return r.availableCopyCount, r.registeredCopyCount
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
            book, isBookRegistered := sl.Books[isbn]
            response := &SimpleLibraryResponse{}

            if !isBookRegistered {
                response.err = &NotFoundBookError{BookError{isbn}}
                responses <- response
                return
            }

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
	}

	for i := 0; i < librarians; i++ {
		sl.librarians <- struct{}{}
	}

	return sl
}
