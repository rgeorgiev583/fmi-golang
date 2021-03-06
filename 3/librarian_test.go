package main

import (
	"fmt"
	"sync"
	"time"
)

type exampleEmpty struct{}
type exampleSemaphore chan exampleEmpty

func ExampleNewLibrary() {
	_ = NewLibrary(2)

	// Output:
}

func ExampleHello() {
	n := 10
	l := NewLibrary(n)
	c := make(exampleSemaphore)

	for i := 1; i <= 2; i++ {
		go func() {
			_, _ = l.Hello()
			c <- exampleEmpty{}
		}()

		select {
		case <-c:
		case <-time.After(time.Second * 1):
			fmt.Printf("Call to Hello() timeouted")
		}
	}

	// Output:
}

func ExampleHello_second() {
	l := NewLibrary(2)
	c := make(exampleSemaphore)

	for i := 1; i <= 3; i++ {
		go func() {
			_, _ = l.Hello()
			c <- exampleEmpty{}
		}()

		select {
		case <-c:
		case <-time.After(time.Second * 1):
			fmt.Printf("Call to Hello() timeouted")
		}
	}

	// Output:
	// Call to Hello() timeouted
}

var hamlet = map[string]string{
	"isbn": "9781420922530",

	"json": `{
					"isbn": "9781420922530",
					"title": "Hamlet",
					"author": {
						"first_name": "William",
						"last_name": "Shakespeare"
					},
					"genre": "poetry",
					"pages": 104,
					"ratings": [4]
				}`,

	"xml": `
<book isbn="9781420922530">
  <title>Hamlet</title>
  <author>
    <first_name>William</first_name>
    <last_name>Shakespeare</last_name>
  </author>
  <genre>poetry</genre>
  <pages>104</pages>
  <ratings>
    <rating>4</rating>
  </ratings>
</book>
`,
}

var macbeth = map[string]string{
	"isbn": "9781853260353",

	"json": `{
					"isbn": "9781853260353",
					"title": "Macbeth",
					"author": {
						"first_name": "William",
						"last_name": "Shakespeare"
					},
					"genre": "poetry",
					"pages": 128,
					"ratings": [4, 5]
				}`,
}

func ExampleAddBookJSON() {
	l := NewLibrary(1)

	n, err := l.AddBookJSON([]byte(hamlet["json"]))
	fmt.Printf("%d %v", n, err)

	// Output:
	// 1 <nil>
}

func ExampleAddBookJSON_second() {
	l := NewLibrary(1)
	var (
		n   int
		err error
	)

	for i := 1; i <= 5; i++ {
		n, err = l.AddBookJSON([]byte(hamlet["json"]))
		fmt.Printf("%d %v\n", n, err)
	}

	// Output:
	// 1 <nil>
	// 2 <nil>
	// 3 <nil>
	// 4 <nil>
	// 4 Има 4 копия на книга 9781420922530
}

func ExampleAddBookJSON_third() {
	l := NewLibrary(1)

	for i := 1; i <= 5; i++ {
		go func() {
			_, _ = l.AddBookJSON([]byte(hamlet["json"]))
		}()
	}

	// Output:
}

func ExampleAddBookXML() {
	l := NewLibrary(1)

	n, err := l.AddBookXML([]byte(hamlet["xml"]))
	fmt.Printf("%d %v", n, err)

	// Output:
	// 1 <nil>
}

func ExampleAddBookXML_second() {
	l := NewLibrary(1)
	var (
		n   int
		err error
	)

	for i := 1; i <= 5; i++ {
		n, err = l.AddBookXML([]byte(hamlet["xml"]))
		fmt.Printf("%d %v\n", n, err)
	}

	// Output:
	// 1 <nil>
	// 2 <nil>
	// 3 <nil>
	// 4 <nil>
	// 4 Има 4 копия на книга 9781420922530
}

func ExampleAddBookXML_third() {
	l := NewLibrary(1)

	for i := 1; i <= 5; i++ {
		go func() {
			_, _ = l.AddBookXML([]byte(hamlet["xml"]))
		}()
	}

	// Output:
}

type ExampleLibraryRequest struct {
	Type int
	ISBN string
}

func (lr *ExampleLibraryRequest) GetType() int {
	return lr.Type
}

func (lr *ExampleLibraryRequest) GetISBN() string {
	return lr.ISBN
}

func responseWithTimeoutRestriction(
	responses <-chan LibraryResponse, called string) {
	select {
	case resp := <-responses:
		book, err := resp.GetBook()
		available, registered := resp.GetAvailability()
		fmt.Printf("%s %v %d %d\n", book, err, available, registered)
	case <-time.After(time.Second * 1):
		fmt.Printf("Response from %s timeouted", called)
	}
}

func responseWithTimeoutRestrictionOmittingBook(
	responses <-chan LibraryResponse, called string) {
	select {
	case resp := <-responses:
		_, err := resp.GetBook()
		available, registered := resp.GetAvailability()
		fmt.Printf("%v %d %d\n", err, available, registered)
	case <-time.After(time.Second * 1):
		fmt.Printf("Response from %s timeouted", called)
	}
}

func ExampleBorrowAvailableBook() {
	l := NewLibrary(1)
	l.AddBookJSON([]byte(hamlet["json"]))
	requests, responses := l.Hello()

	requests <- &ExampleLibraryRequest{
		Type: 1,
		ISBN: hamlet["isbn"],
	}
	responseWithTimeoutRestriction(responses, "Library")
	// Output:
	// [9781420922530] Hamlet от William Shakespeare <nil> 0 1
}

func ExampleBorrowTwoAvailableBooks() {
	l := NewLibrary(1)
	l.AddBookJSON([]byte(hamlet["json"]))
	l.AddBookJSON([]byte(macbeth["json"]))
	requests, responses := l.Hello()
	req := &ExampleLibraryRequest{
		Type: 1,
		ISBN: hamlet["isbn"],
	}
	requests <- req
	responseWithTimeoutRestriction(responses, "Library")
	req.ISBN = macbeth["isbn"]
	requests <- req
	responseWithTimeoutRestriction(responses, "Library")

	// Output:
	// [9781420922530] Hamlet от William Shakespeare <nil> 0 1
	// [9781853260353] Macbeth от William Shakespeare <nil> 0 1
}

func ExampleConcurrentBorrowBook() {
	l := NewLibrary(2)
	l.AddBookJSON([]byte(hamlet["json"]))
	requests1, responses1 := l.Hello()
	requests2, responses2 := l.Hello()
	req := &ExampleLibraryRequest{
		Type: 1,
		ISBN: hamlet["isbn"],
	}
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		requests1 <- req
		responseWithTimeoutRestriction(responses1, "Library")
		wg.Done()
	}()

	go func() {
		requests2 <- req
		responseWithTimeoutRestriction(responses2, "Library")
		wg.Done()
	}()
	wg.Wait()

	// Output:
	// [9781420922530] Hamlet от William Shakespeare <nil> 0 1
	// <nil> Няма наличност на книга 9781420922530 0 1
}

func ExampleBorrowUnavailableBook() {
	l := NewLibrary(1)
	l.AddBookJSON([]byte(hamlet["json"]))
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 1,
		ISBN: hamlet["isbn"],
	}

	requests <- req
	<-responses
	requests <- req

	responseWithTimeoutRestriction(responses, "Library")

	// Output:
	// <nil> Няма наличност на книга 9781420922530 0 1
}

func ExampleBorrowUnexistingBook() {
	l := NewLibrary(1)
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 1,
		ISBN: hamlet["isbn"],
	}
	requests <- req

	responseWithTimeoutRestriction(responses, "Library")

	// Output:
	// <nil> Непозната книга 9781420922530 0 0
}

func ExampleReturnUnavailableBook() {
	l := NewLibrary(1)
	l.AddBookJSON([]byte(hamlet["json"]))
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 1,
		ISBN: hamlet["isbn"],
	}
	requests <- req
	<-responses
	req.Type = 2
	requests <- req

	responseWithTimeoutRestrictionOmittingBook(responses, "Library")

	// Output:
	// <nil> 1 1
}

func ExampleReturnAvailableBook() {
	l := NewLibrary(1)
	l.AddBookJSON([]byte(hamlet["json"]))
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 2,
		ISBN: hamlet["isbn"],
	}
	requests <- req

	responseWithTimeoutRestrictionOmittingBook(responses, "Library")

	// Output:
	// Всички копия са налични 9781420922530 1 1
}

func ExampleReturnUnexistingBook() {
	l := NewLibrary(1)
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 2,
		ISBN: hamlet["isbn"],
	}
	requests <- req

	responseWithTimeoutRestrictionOmittingBook(responses, "Library")

	// Output:
	// Непозната книга 9781420922530 0 0
}

func ExampleAvailabilityOfExistingBook() {
	l := NewLibrary(1)
	l.AddBookJSON([]byte(hamlet["json"]))
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 3,
		ISBN: hamlet["isbn"],
	}
	requests <- req

	responseWithTimeoutRestriction(responses, "Library")

	// Output:
	// [9781420922530] Hamlet от William Shakespeare <nil> 1 1
}

func ExampleAvailabilityOfUnexistingBook() {
	l := NewLibrary(1)
	requests, responses := l.Hello()

	req := &ExampleLibraryRequest{
		Type: 3,
		ISBN: hamlet["isbn"],
	}
	requests <- req

	responseWithTimeoutRestrictionOmittingBook(responses, "Library")

	// Output:
	// Непозната книга 9781420922530 0 0
}

func ExampleScenarioWithTwoLibrarians() {
	n := 2
	l := NewLibrary(n)

	l.AddBookJSON([]byte(hamlet["json"]))
	for i := 1; i <= 4; i++ {
		l.AddBookJSON([]byte(macbeth["json"]))
	}

	requests := make([]chan<- LibraryRequest, n)
	responses := make([]<-chan LibraryResponse, n)

	for i := 0; i < n; i++ {
		requests[i], responses[i] = l.Hello()
	}

	req := &ExampleLibraryRequest{
		Type: 3,
		ISBN: hamlet["isbn"],
	}

	requests[0] <- req
	responseWithTimeoutRestrictionOmittingBook(responses[0], "Library")

	req.Type = 1
	requests[1] <- req
	responseWithTimeoutRestriction(responses[1], "Library")

	requests[0] <- req
	responseWithTimeoutRestriction(responses[0], "Library")

	req.Type = 2
	requests[1] <- req
	responseWithTimeoutRestriction(responses[1], "Library")

	req.Type = 1
	requests[0] <- req
	responseWithTimeoutRestriction(responses[0], "Library")

	req.ISBN = macbeth["isbn"]
	requests[0] <- req
	responseWithTimeoutRestriction(responses[0], "Library")

	requests[1] <- req
	responseWithTimeoutRestriction(responses[1], "Library")

	c := make(exampleSemaphore)

	go func() {
		_, _ = l.Hello()
		c <- exampleEmpty{}
	}()
	select {
	case <-c:
	case <-time.After(time.Second * 1):
		fmt.Printf("Call to Hello() timeouted\n")
	}

	close(requests[0])

	go func() {
		_, _ = l.Hello()
		c <- exampleEmpty{}
	}()

	select {
	case <-c:
	case <-time.After(time.Second * 1):
		fmt.Printf("Call to Hello() timeouted\n")
	}

	// Output:
	// <nil> 1 1
	// [9781420922530] Hamlet от William Shakespeare <nil> 0 1
	// <nil> Няма наличност на книга 9781420922530 0 1
	// [9781420922530] Hamlet от William Shakespeare <nil> 1 1
	// [9781420922530] Hamlet от William Shakespeare <nil> 0 1
	// [9781853260353] Macbeth от William Shakespeare <nil> 3 4
	// [9781853260353] Macbeth от William Shakespeare <nil> 2 4
	// Call to Hello() timeouted
}
