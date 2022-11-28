package main

import (
	"testing"
)

var config = Config{}

func TestBootstrap(t *testing.T) {
	err := bootstrap(config)
	if err != nil {
		t.Fatal(err.Error())
	}
}
