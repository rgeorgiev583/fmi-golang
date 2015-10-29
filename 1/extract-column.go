package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func ExtractColumn(logContents string, column uint8) string {
	var buffer bytes.Buffer
	lines := strings.Split(logContents, "\n")

	for i := range lines {
		cols := strings.Split(lines[i], " ")
		buffer.WriteString(cols[column])

		if column == 0 {
			buffer.WriteString(" ")
			buffer.WriteString(cols[1])
		}

		buffer.WriteString("\n")
	}

	return buffer.String()
}

func main() {
	var (
		column uint8
		buffer bytes.Buffer
	)

	_, err := fmt.Scan(&column)

	if err != nil {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
	}

	fmt.Printf(ExtractColumn(buffer.String(), column))
}
