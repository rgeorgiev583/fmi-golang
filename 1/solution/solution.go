package main

import (
	"bufio"
	"bytes"
	"strings"
)

func ExtractColumn(logContents string, column uint8) string {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(logContents))

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		cols := strings.Split(line, " ")

		switch column {
		case 0:
			buffer.WriteString(cols[0])
			buffer.WriteString(" ")
			buffer.WriteString(cols[1])
		case 1:
			buffer.WriteString(cols[2])
		case 2:
			buffer.WriteString(cols[3])
            colCount := len(cols)

			for j := 4; j < colCount; j++ {
				buffer.WriteString(" ")
				buffer.WriteString(cols[j])
			}
		}

		buffer.WriteString("\n")
	}

	return buffer.String()
}
