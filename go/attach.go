package main

import (
	"fmt"
	"os"
	// Consider adding https://pkg.go.dev/golang.org/x/sys/unix
)

// function attach {
//   # -ot is older than
//   if [ ! -f "$LOCAL_DIR/$RUNNER" ] || [ "$LOCAL_DIR/$RUNNER" -ot "$RUNNER" ]; then
//     sudo cp --preserve=mode "$RUNNER" "$LOCAL_DIR/$RUNNER"
//     # We want global read/execute permissions, which git doesn't track
//     sudo chmod 755 "$LOCAL_DIR/$RUNNER"
//   fi
//
//   if [ ! -f "$HOME/.ssh/id_rsa" ]; then
//     echo "Expected $HOME/.ssh/id_rsa to exist to copy to chroot" >&2
//     exit 1
//   fi
//
//   if [ ! -d "$LOCAL_DIR/.ssh" ]; then
//     sudo cp -r "$HOME/.ssh/" "$LOCAL_DIR"
//     # ensure our coder user can copy this
//     sudo chmod +r "$LOCAL_DIR/.ssh"
//   fi
//
//   # chroot
//   sudo "$LOCAL_DIR/bin/arch-chroot" "$LOCAL_DIR/" "/$RUNNER" 'initialize-root'
// }

func attach(config Config) {
	// Get effective UID in case they are using sudo
	uid := os.Geteuid()
	if uid != 0 {
		panic(fmt.Errorf("attach should be called as root from the host"))
	}

	panic(fmt.Errorf("TODO implement attach\n"))
}
