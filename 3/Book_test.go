package main

import (
    "fmt"
    "encoding/xml"
)

func ExampleUnmarshal() {
    bookXml := `
<book isbn="0954540018">
  <title>Who said the race is Over?</title>
  <author>
    <first_name>Anno</first_name>
    <last_name>Birkin</last_name>
  </author>
  <genre>poetry</genre>
  <pages>80</pages>
  <ratings>
    <rating>5</rating>
    <rating>4</rating>
    <rating>4</rating>
    <rating>5</rating>
    <rating>3</rating>
  </ratings>
</book>
`
    book := &Book{}
    err := xml.Unmarshal([]byte(bookXml), book)

    if err != nil {
        fmt.Printf("error: %#v\n", err)
    }

    fmt.Printf("%#v\n", book)
    fmt.Println()
}
