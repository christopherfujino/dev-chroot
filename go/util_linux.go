//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Exec(binary string, argv []string) {
	absoluteBinary, err := exec.LookPath("uname")
	check(
		err,
		"looking up binary \"uname\"",
	)
	if err = syscall.Exec(
		absoluteBinary,
		// First element in argv should be name of binary
		append([]string{binary}, argv...),
		os.Environ(),
	); err != nil {
		panic(fmt.Errorf("Received error %v\n", err))
	}
	panic("unreachable")
}
