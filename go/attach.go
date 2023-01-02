package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func attach(config Config, cwd string) {
	// TODO this should be passed in from main.go
	var localRoot = filepath.Join(cwd, "root.x86_64")

	// Get effective UID in case they are using sudo
	uid := os.Geteuid()
	if uid != 0 {
		panic(fmt.Errorf("attach should be called as root from the host"))
	}

	// TODO copy .ssh folder?

	var archChroot = filepath.Join(localRoot, "bin", "arch-chroot")
	var cmd = exec.Command(archChroot, localRoot, "su", "--login", config.UserName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	check(
		err,
		"running arch-chroot",
	)
	fmt.Println("Cleaning up from chroot...")
}

func createHashFile(hash string, filePath string) {
	var err = os.WriteFile(
		filePath,
		[]byte(hash),
		// rw-r--r--
		0644,
	)
	check(
		err,
		fmt.Sprintf("writing sha256sum to %s", filePath),
	)
}
