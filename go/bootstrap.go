package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	if uid != 0 {
		panic(fmt.Errorf("bootstrap should be called as root from the host"))
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if cwd != "/home/fujino/git/dev-chroot/go" {
		log.Fatalf("CWD = %s", cwd)
	}
	fmt.Println("Bootstrapping chroot locally")
	downloadTarball(
		httpGetter,
		config.remoteBootstrapTarball,
		filepath.Join(cwd, config.localBootstrapTarball),
	)
	extractTarball(config.localBootstrapTarball, cwd)
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
	localFile, err := os.Create(localPath)
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

func extractTarball(localTarball string, cwd string) {
	file, err := os.Open(filepath.Join(cwd, localTarball))
	if err != nil {
		panic(err)
	}
	tarRaw, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	tarReader := tar.NewReader(tarRaw)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error %s reading %s", err.Error(), localTarball)
		}

		extractEntity(header, cwd)
	}
}

func extractEntity(header *tar.Header, cwd string) {
	switch header.Typeflag {
	case tar.TypeDir:
		//fmt.Printf("%s\tdir\t0%o\n", header.Name, header.Mode)
		newpath := filepath.Join(cwd, header.Name)
		os.Mkdir(newpath, fs.FileMode(header.Mode))
		err := os.Chown(newpath, header.Uid, header.Gid)
		if err != nil {
			panic(err)
		}
	case tar.TypeReg:
		fmt.Printf("%s\tfile\t0%o\n", header.Name, header.Mode)
		//fmt.Printf("Making %s", filepath.Join(cwd, header.Name))
	case tar.TypeLink:
		fmt.Printf("%s -> %s\thard link\towned by %d\n", header.Name, header.Linkname, header.Uid)
	case tar.TypeSymlink:
		fmt.Printf("%s -> %s\tsymlink\towned by %d\n", header.Name, header.Linkname, header.Uid)
	default:
		log.Fatalf("%s has an unknown Tar type '%c'\n", header.Name, header.Typeflag)
	}
}
