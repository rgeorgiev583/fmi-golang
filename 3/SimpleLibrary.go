package main

import (
    "fmt"
    "json"
    "xml"
)

type SimpleLibrary struct {
    books map[string]*Book
    registeredCopyCount map[string]int
    availableCopyCount map[string]int
    librarians chan struct{}
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

type NotAvaiableBookError struct {
    BookError
}

type AllCopiesAvailableBookError struct {
    BookError
}

func (l *Library) MarshalJSON() ([]byte, error) {
    return json.Marshal(l)
}

func (l *Library) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, l)
}

func (l *Library) MarshalXML(e *Encoder, start StartElement) error {
    return e.EncodeElement(l, start)
}

func (l *Library) UnmarshalXML(d *Decoder, start StartElement) error {
    return d.DecodeElement(l, &start)

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
        err = TooManyCopiesBookError{ISBN: book.ISBN}
    } else {
        sl.books[book.ISBN] = book
        registeredCopyCount = ++sl.registeredCopyCount[book.ISBN]
        ++sl.availableCopyCount[book.ISBN]
    }

    return
}

func (sl *SimpleLibrary) AddBookJSON(data []byte) (int, error) {
    var book *Book
    json.Unmarshal(data, book)
    return el.addBook(book)
}

func (sl *SimpleLibrary) AddBookXML(data []byte) (int, error) {
    var book *Book
    xml.Unmarshal(data, book)
    return el.addBook(book)
}

func (sl *SimpleLibrary) Hello(requests chan<- LibraryRequest, responses <-chan LibraryResponse) {
    for request := range requests {
        <-sl.librarians
        isbn := request.GetISBN()
        response := new(SimpleLibraryResponse)

        switch request.GetType() {
        case 1:
            if book, isBookRegistered := sl.books[isbn]; isBookRegistered && sl.availableCopyCount[isbn] > 0 {
                response.book = book
                --sl.availableCopyCount[isbn]
            } else if !isBookRegistered {
                response.err = &NotFoundBookError{BookError{isbn}}
            } else {
                response.err = &NotAvailableBookError{BookError{isbn}}
            }

        case 2:
            if _, isBookRegistered := sl.books[isbn]; isBookRegistered && sl.availableCopyCount[isbn] < sl.registeredCopyCount[isbn] {
                ++sl.availableCopyCount[isbn]
            } else if !isBookRegistered {
                response.err = &NotFoundBookError{BookError{isbn}}
            } else {
                response.err = &AllCopiesAvailableBookError{BookError{isbn}}
            }

        case 3:
            response.registeredCopyCount = sl.registeredCopyCount[isbn]
            response.availableCopyCount = sl.availableCopyCount[isbn]
        }

        responses <- response
        sl.librarians <- struct{}{}
    }
}

func NewLibrary(librarians int) Library {
    return &SimpleLibrary{
        books: make(map[string]*Book),
        registeredCopyCount: make(map[string]int),
        availableCopyCount: make(map[string]int),
        librarians: make(chan struct{}, librarians)
    }
}
