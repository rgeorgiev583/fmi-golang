package main

import (
    "fmt"
    "encoding/json"
    "encoding/xml"
)

const (
    _ = iota
    TakeBook
    ReturnBook
    GetAvailability
)

type SimpleLibrary struct {
    Books map[string]*Book
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

type NotAvailableBookError struct {
    BookError
}

type AllCopiesAvailableBookError struct {
    BookError
}

func (sl *SimpleLibrary) MarshalJSON() ([]byte, error) {
    return json.Marshal(sl)
}

func (sl *SimpleLibrary) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, sl)
}

func (sl *SimpleLibrary) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return e.EncodeElement(sl, start)
}

func (sl *SimpleLibrary) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    return d.DecodeElement(sl, &start)
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
        err = TooManyCopiesBookError{ISBN: book.ISBN}
    } else {
        sl.Books[book.ISBN] = book
        sl.registeredCopyCount[book.ISBN]++
        sl.availableCopyCount[book.ISBN]++
        registeredCopyCount = sl.registeredCopyCount[book.ISBN]
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

func (sl *SimpleLibrary) Hello() (requests chan<- LibraryRequest, responses <-chan LibraryResponse) {
    requests = make(chan<- LibraryRequest)
    responses = make(<-chan LibraryResponse)

    go func() {
        for request := range requests {
            <-sl.librarians
            isbn := request.GetISBN()
            response := new(SimpleLibraryResponse)

            switch request.GetType() {
            case TakeBook:
                if book, isBookRegistered := sl.Books[isbn]; isBookRegistered && sl.availableCopyCount[isbn] > 0 {
                    response.book = book
                    sl.availableCopyCount[isbn]--
                } else if !isBookRegistered {
                    response.err = &NotFoundBookError{BookError{isbn}}
                } else {
                    response.err = &NotAvailableBookError{BookError{isbn}}
                }

            case ReturnBook:
                if _, isBookRegistered := sl.Books[isbn]; isBookRegistered && sl.availableCopyCount[isbn] < sl.registeredCopyCount[isbn] {
                    sl.availableCopyCount[isbn]++
                } else if !isBookRegistered {
                    response.err = &NotFoundBookError{BookError{isbn}}
                } else {
                    response.err = &AllCopiesAvailableBookError{BookError{isbn}}
                }

            case GetAvailability:
                response.registeredCopyCount = sl.registeredCopyCount[isbn]
                response.availableCopyCount = sl.availableCopyCount[isbn]
            }

            responses <- response
            sl.librarians <- struct{}{}
        }
    }()

    return
}

func NewLibrary(librarians int) Library {
    return &SimpleLibrary{
        Books: make(map[string]*Book),
        registeredCopyCount: make(map[string]int),
        availableCopyCount: make(map[string]int),
        librarians: make(chan struct{}, librarians),
    }
}
