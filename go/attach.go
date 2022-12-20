package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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

	var initializationScriptPath = ""
	if config.Provision != "" {
		ensureInitScript(localRoot, config)
	}

	// TODO copy .ssh folder?

	var archChroot = filepath.Join(localRoot, "bin", "arch-chroot")
	if initializationScriptPath != "" {
		// Note this will now be relative to chroot
		var cmd = exec.Command(archChroot, localRoot, fmt.Sprintf("/%s", initScriptName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		check(
			err,
			"running arch-chroot",
		)
	}
	Exec(archChroot, []string{localRoot, "su", "--login", config.UserName})
}

func createInitScriptFile(filePath string, contents string) {
	// Only root needs to execute, others can read
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0744)
	check(
		err,
		fmt.Sprintf("opening %s in write-only", filePath),
	)

	_, err = file.WriteString(contents)
	check(
		err,
		fmt.Sprintf("writing init script to %s", filePath),
	)
	file.Close()
}

func interpolateInitScript(templateString string, config Config) string {
	configTemplate, err := template.New("config.Provision").Parse(config.Provision)
	check(
		err,
		fmt.Sprintf("trying to create template"),
	)
	var buffer = strings.Builder{}
	err = configTemplate.Execute(&buffer, config)
	check(
		err,
		"interpolating template",
	)
	return buffer.String()
}

func createHashFile(hash string, filePath string) {
	var err = os.WriteFile(filePath, []byte(hash), 0644)
	check(
		err,
		fmt.Sprintf("writing sha256sum to %s", filePath),
	)
}

const initScriptName = "root_init_script.sh"

// Ensure the chroot dir has the init script installed.
//
// localRoot is the absolute path to the chroot dir.
func ensureInitScript(localRoot string, config Config) {
	var initScriptPath = filepath.Join(localRoot, initScriptName)
	var hashFilePath = fmt.Sprintf("%s.sha256", initScriptPath)

	var initScriptContents = interpolateInitScript(config.Provision, config)
	var hashBuffer = sha256.Sum256([]byte(initScriptContents))
	var hashString = fmt.Sprintf("%x\n", hashBuffer)
	fileBytes, err := os.ReadFile(hashFilePath)
	if err == nil {
		// hash file exists
		var fileString = string(fileBytes)
		if fileString != hashString {
			// invalidate hash file
			fmt.Printf("Invalidation of hash file, re-copying init script...\n")
			createInitScriptFile(initScriptPath, initScriptContents)
			createHashFile(hashString, hashFilePath)
		} else {
			// cache hit, nothing else to do
			return
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("Error reading %s: %s", hashFilePath, err.Error()))
	} else {
		// no hash file, should create it
		createInitScriptFile(initScriptPath, initScriptContents)
		createHashFile(hashString, hashFilePath)
	}
}
