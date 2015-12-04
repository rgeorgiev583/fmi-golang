package main

type SimpleLibraryRequest struct {
    requestType int
    bookISBN string
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
