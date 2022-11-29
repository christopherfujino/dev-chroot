package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
)

var config = Config{}

func TestBootstrapDownload(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				t.Errorf("Expected an os.ErrNotExist, got a %v", r)
			}
			if !errors.Is(err, os.ErrNotExist) {
				t.Errorf("Expected an os.ErrNotExist, got a %v", err)
			}
			// success
		}
	}()
	downloadTarball(
		func(url string) (*http.Response, error) {
			return nil, fmt.Errorf("Foo\n")
		},
		"https://cloud.storage/stuff.tar.gz",
		"stuff.tar.gz",
	)
	t.Error("Expected a panic")
}
