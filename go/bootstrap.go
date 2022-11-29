package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

type HttpGetter func(url string) (*http.Response, error)

// Should be run as root on the host
func bootstrap(
	config Config,
	httpGetter HttpGetter,
) {
	// Get effective UID in case they are using sudo
	uid := os.Geteuid()
	if uid == 0 {
		panic(fmt.Errorf("bootstrap should be called as a regular user from the host"))
	}
	fmt.Println("Bootstrapping chroot locally")
	downloadTarball(httpGetter, config.remoteBootstrapTarball, config.localBootstrapTarball)
	extractTarball(config.localBootstrapTarball)
}

// Download remote tarball to local disk, unless the localPath already exists.
func downloadTarball(httpGetter HttpGetter, remotePath string, localPath string) {
	if _, err := os.Stat(localPath); !errors.Is(err, os.ErrNotExist) {
		// err might be nil here, in which case the file exists, in which case
		// we are happy and don't need to download it
		fmt.Printf("File %s already exists, skipping download step.\n", localPath)
		if err != nil {
			panic(err)
		}
		return
	}
	localFile, err := os.Create(localPath) // TODO append to working dir
	if err != nil {
		panic(err)
	}
	defer localFile.Close()
	fmt.Printf("Downloading %s to %s...\n", remotePath, localFile.Name())
	resp, err := httpGetter(remotePath)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if _, err = io.Copy(localFile, resp.Body); err != nil {
		panic(err)
	}
}

func extractTarball(localTarball string) {}
