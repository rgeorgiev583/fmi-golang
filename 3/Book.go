package main

import (
    "fmt"
    "json"
    "xml"
)

type Person struct {
    FirstName string `json:"first_name", xml:"first_name"`
    LastName string `json:"last_name", xml:"last_name"`
}

type Book struct {
    XMLName xml.Name `xml:"book"`
    ISBN string `json:"isbn", xml:"isbn,attr"`
    Title string `json:"title", xml:"title"`
    Author Person `json:"author", xml:"author"`
    Ratings []struct {
        Rating uint8 `xml:"rating"`
    } `json:"ratings", xml:"ratings"`
    availableCopyCount int
    registeredCopyCount int
}

func (p *Person) Person() string {
    return fmt.Sprintf("%v %v", p.FirstName, p.LastName)

func (b *Book) String() string {
    return fmt.Sprintf("[%v] %v от %v", b.ISBN, b.Title, b.Author)

func (b *Book) MarshalJSON() ([]byte, error) {
    return json.Marshal(b)
}

func (b *Book) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, b)
}

func (b *Book) MarshalXML(e *Encoder, start StartElement) error {
    return e.EncodeElement(b, start)
}

func (b *Book) UnmarshalXML(d *Decoder, start StartElement) error {
    return d.DecodeElement(b, &start)
