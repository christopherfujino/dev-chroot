package main

import (
	"fmt"
)

func check(maybeError error, message string) {
	if maybeError != nil {
		err := fmt.Errorf("Error: %s\n\n%s", message, maybeError.Error())
		panic(err)
	}
}
