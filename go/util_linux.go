//go:build linux

package main

import (
	"path/filepath"
	"fmt"
	"os"
	"os/exec"
	"syscall" // Consider switching to https://pkg.go.dev/golang.org/x/sys/unix
)

func Exec(binary string, argv []string) {
	absoluteBinary, err := exec.LookPath(binary)
	check(
		err,
		"looking up binary \"uname\"",
	)
	fmt.Printf("%s %v", absoluteBinary, append([]string{filepath.Base(binary)}, argv...))
	if err = syscall.Exec(
		absoluteBinary,
		// First element in argv should be name of binary
		append([]string{filepath.Base(binary)}, argv...),
		os.Environ(),
	); err != nil {
		panic(fmt.Errorf("Received error %v\n", err))
	}
	panic("unreachable")
}
