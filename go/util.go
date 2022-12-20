package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
)

func check(maybeError error, message string) {
	if maybeError != nil {
		err := fmt.Errorf("Error %s\n\n%s", message, maybeError.Error())
		panic(err)
	}
}

// Substitute regex pattern matches for replacement in the given file.
//
// A sed-replacement
func processFile(filePath string, pattern string, replacement string) {
	regex, err := regexp.Compile(pattern)
	check(
		err,
		fmt.Sprintf("trying to compile regex %s", pattern),
	)
	// For simplicity, read-only first
	file, err := os.Open(filePath)
	check(
		err,
		fmt.Sprintf("trying to open file %s", filePath),
	)
	scanner := bufio.NewScanner(file)
	scanner.Split(fullSplitter)
	var lines = []string{}
	for {
		ok := scanner.Scan()
		if !ok {
			if err := scanner.Err(); err != nil {
				panic(err)
			}
			// done scanning
			break
		}
		var line = scanner.Text()
		line = regex.ReplaceAllString(line, replacement)
		lines = append(lines, line)
	}
	file.Close()
	// Permissions should be ignored since this file already exists
	// https://github.com/golang/go/issues/33605
	file, err = os.OpenFile(filePath, os.O_WRONLY, 0666)
	check(
		err,
		fmt.Sprintf("trying to open file %s for truncation", filePath),
	)

	for _, line := range lines {
		idx, err := file.WriteString(line)
		check(
			err,
			fmt.Sprintf("writing line %d: \"%s\"", idx+1, line),
		)
	}
	file.Close()
}

// Like the default, but does not remove newlines.
func fullSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i + 1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
