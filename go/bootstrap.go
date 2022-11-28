package main

import (
	"fmt"
	"os"
)

// Should be run as root on the host
func bootstrap(config Config) error {
	// Get effective UID in case they are using sudo
	uid := os.Geteuid()
	if uid == 0 {
		return fmt.Errorf("bootstrap should be called as a regular user from the host")
	}
	fmt.Println("Bootstrapping chroot locally")
	return nil
}
