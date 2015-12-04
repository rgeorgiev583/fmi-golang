package main

type SimpleLibraryResponse struct {
    book Book
    registeredCopyCount int
    availableCopyCount int
    err error
}

func (r *SimpleLibraryResponse) GetBook(fmt.Stringer, error) {
    return r.book, r.err
}

func (r *SimpleLibraryResponse) GetAvailability(int, int) {
    return availableCopyCount, registeredCopyCount
}
