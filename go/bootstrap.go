package main

import (
	"fmt"
	"os"
)

// run as user on host
//function bootstrap {
//  set -o xtrace
//
//  if [ "$EUID" -eq 0 ]; then
//    echo 'bootstrap function should not be called by root' >&2
//    exit 1
//  fi
//
//  LOCAL_TARBALL='arch-tarball.tar.gz'
//
//  if [ ! -f "$LOCAL_TARBALL" ]; then
//    curl -l "$REMOTE_BOOTSTRAP_TARBALL" -o "$LOCAL_TARBALL"
//  fi
//
//  if [ ! -d "$LOCAL_DIR" ]; then
//    # --numeric-owner since host might not use the same user id's as arch
//    # Must have root permission as there are some UID 0 files
//    sudo tar xzf "$LOCAL_TARBALL" --numeric-owner
//
//    # Enable berkeley mirror
//    # -E means extended regex
//    # -i means update file in place
//    sudo sed -E -i 's/^#(.*berkeley)/\1/' "$LOCAL_DIR/etc/pacman.d/mirrorlist"
//    # disable CheckSpace setting
//    sudo sed -E -i 's/^CheckSpace/#CheckSpace/' "$LOCAL_DIR/etc/pacman.conf"
//  fi
//}

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
