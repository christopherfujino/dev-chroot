package main

import (
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
				t.Errorf("Expected an error, got an %v", r)
			}
			if err.Error() != "Foo" {
				t.Errorf("Expected \"Foo\", got %v", err.Error())
			}
			// success
			check(
				os.Remove("stuff.tar.gz"),
				"Failed cleaning up test",
			)
		}
	}()

	downloadTarball(
		func(url string) (*http.Response, error) {
			panic(fmt.Errorf("Foo"))
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
