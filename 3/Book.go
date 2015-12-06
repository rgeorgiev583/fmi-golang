package main

import (
	"encoding/xml"
	"fmt"
)

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

func (b *Book) String() string {
	return fmt.Sprintf("[%v] %v от %v %v", b.ISBN, b.Title, b.Author.FirstName, b.Author.LastName)
}
