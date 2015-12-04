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
