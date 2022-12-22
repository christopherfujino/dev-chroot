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

type HttpGetter func(url string) (*http.Response, error)

// Should be run as root on the host
func bootstrap(
	config Config,
	httpGetter HttpGetter,
	cwd string,
) {
	// Get effective UID in case they are using sudo
	uid := os.Geteuid()
	if uid != 0 {
		panic(fmt.Errorf("bootstrap should be called as root from the host"))
	}

	fmt.Println("Bootstrapping chroot locally")
	var localPath = filepath.Join(cwd, config.LocalBootstrapTarball)
	if _, err := os.Stat(localPath); !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			panic(err)
		}
		// err might be nil here, in which case the file exists, in which case
		// we are happy and don't need to download it
		fmt.Printf("File %s already exists, skipping download step.\n", localPath)
		return
	} else {
		downloadTarball(
			httpGetter,
			config.RemoteBootstrapTarball,
			localPath,
		)
	}

	localRoot := extractTarball(config.LocalBootstrapTarball, cwd)
	processFile(
		filepath.Join(localRoot, "etc/pacman.d/mirrorlist"),
		"^#(.*berkeley)",
		"$1",
	)
	processFile(
		filepath.Join(localRoot, "etc/pacman.conf"),
		"^CheckSpace",
		"#CheckSpace",
	)
}

// Download remote tarball to local disk, unless the localPath already exists.
func downloadTarball(httpGetter HttpGetter, remotePath string, localPath string) {
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

func extractTarball(localTarball string, cwd string) string {
	var absoluteTarballPath = filepath.Join(cwd, localTarball)
	file, err := os.Open(absoluteTarballPath)
	if err != nil {
		panic(err)
	}
	tarRaw, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	tarReader := tar.NewReader(tarRaw)
	var linkQueue = LinkQueue{}
	var localRoot = ""

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		check(
			err,
			fmt.Sprintf("reading %s", absoluteTarballPath),
		)
		newpath := filepath.Join(cwd, header.Name)
		switch header.Typeflag {
		case tar.TypeDir: // Directory
			if localRoot == "" {
				// getRoot() will return root relative to CWD
				localRoot = filepath.Join(cwd, getRoot(header.Name))
			}
			os.Mkdir(newpath, fs.FileMode(header.Mode))
			err := os.Chown(newpath, header.Uid, header.Gid)
			check(err, fmt.Sprintf("chown of %s to %d:%d", newpath, header.Uid, header.Gid))
			fmt.Printf("Created dir %s\n", header.Name)
		case tar.TypeReg: // File
			file := create(newpath)
			_ = copyBuffer(
				file,
				tarReader,
				fmt.Sprintf("extracting %s from %s", newpath, absoluteTarballPath),
			)
			chown(newpath, header.Uid, header.Gid)
			chmod(newpath, header.Mode)
			fmt.Printf("Created file %s\n", header.Name)
			file.Close()
		case tar.TypeLink: // Hard link
			var destination = newpath
			var source = filepath.Join(cwd, header.Linkname)
			linkQueue.TryHardlink(
				source,
				destination,
				func() {
					chown(destination, header.Uid, header.Gid)
					chmod(destination, header.Mode)
					fmt.Printf("Created hard link %s -> %s\n", header.Name, header.Linkname)
				},
			)
			//log.Fatalf("Unimplemented type hard link")
		case tar.TypeSymlink: // Symlink
			var destination = newpath
			// This might be a relative path, so don't concatenate with CWD
			var source string
			if filepath.IsAbs(header.Linkname) {
				var rootPath = getRoot(header.Name)
				source = filepath.Join(cwd, rootPath, header.Linkname)
			} else {
				source = header.Linkname
			}
			linkQueue.TrySymlink(
				source,
				destination,
				func() {
					chown(destination, header.Uid, header.Gid)
					// TODO chmod symlinks on macOS
					fmt.Printf("Created symlink %s -> %s\n", header.Name, header.Linkname)
				},
			)
		default:
			log.Fatalf("%s has an unknown Tar type '%c'\n", header.Name, header.Typeflag)
		}
	}

	if localRoot == "" {
		panic("Whoops!")
	}

	// Run link jobs last
	for _, task := range linkQueue.queue {
		task.Callback()
	}

	fmt.Printf("Finished extracting %s to %s\n", localTarball, localRoot)

	return localRoot
}

// Create a file.
func create(path string) *os.File {
	file, err := os.Create(path)
	check(err, fmt.Sprintf("creating %s", path))
	return file
}

// Change owner of file.
//
// If the given file is a link, change the link itself.
func chown(path string, uid int, gid int) {
	err := os.Lchown(path, uid, gid)
	check(err, fmt.Sprintf("chown of %s to %d:%d", path, uid, gid))
}

// Change mode of file.
//
// No-op if the file is a symlink (TODO: fix for MacOS
// https://unix.stackexchange.com/questions/87200/change-permissions-for-a-symbolic-link)
func chmod(path string, mode int64) {
	fileInfo, err := os.Lstat(path)
	check(
		err,
		fmt.Sprintf("checking Lstat of %s", path),
	)
	if fileInfo.Mode()&os.ModeSymlink > 0 {
		log.Fatalf("skipping chmod on %s as it is a symlink", path)
		return
	}
	err = os.Chmod(path, fs.FileMode(mode))
	check(err, fmt.Sprintf("chmod of %s to %o", path, mode))
}

func copyBuffer(destination io.Writer, source io.Reader, errMessage string) int64 {
	length, err := io.Copy(destination, source)
	check(err, errMessage)
	return length
}

func getRoot(path string) string {
	var lastDir = filepath.Dir(path)
	var currentDir = filepath.Dir(lastDir)
	// If this was a relative path, we want the lastDir before the dot.
	// If this was an absolute path, we want top-level slash.
	for currentDir != "." && lastDir != currentDir {
		lastDir = currentDir
		currentDir = filepath.Dir(lastDir)
	}
	return lastDir
}
