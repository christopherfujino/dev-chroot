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
	return // TODO this is getting skipped if the hard-coded output dir already exists
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
			panic(fmt.Errorf("Foo\n"))
		},
		"https://cloud.storage/stuff.tar.gz",
		"stuff.tar.gz",
	)
	t.Error("Expected a panic")
}

func TestGetRoot(t *testing.T) {
	var root = getRoot("root.x86_64/var/lib/dbus")
	if root != "root.x86_64" {
		t.Errorf("Expected \"%s\" to be \"root.x86_64\"", root)
	}

	root = getRoot("/root.x86_64/var/lib/dbus")
	if root != "/" {
		t.Errorf("Expected \"%s\" to be \"/\"", root)
	}
}
