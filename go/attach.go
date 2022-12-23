package main

import (
	"fmt"
	"os"
	"path/filepath"
)

//   # chroot
//   sudo "$LOCAL_DIR/bin/arch-chroot" "$LOCAL_DIR/" "/$RUNNER" 'initialize-root'
// }

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
	Exec(archChroot, []string{localRoot, "su", "--login", config.UserName})
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
