package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestFullSplitter(t *testing.T) {
	const input = `
Hello, world!

Goodbye, world!
`
	var reader = strings.NewReader(input)
	var scanner = bufio.NewScanner(reader)
	var buffer = []byte{}
	scanner.Split(fullSplitter)
	for {
		ok := scanner.Scan()
		if !ok {
			if err := scanner.Err(); err != nil {
				panic(err)
			}
			// done scanning
			break
		}
		buffer = append(buffer, scanner.Bytes()...)
	}
	var output = string(buffer)
	if output != input {
		t.Errorf("\"%s\"", output)
	}
}
