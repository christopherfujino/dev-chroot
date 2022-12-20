package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"os"
	"os/exec"
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

	const initScriptName = "root_init_script.sh"
	var initializationScriptPath = ""
	if config.Provision != "" {
		initializationScriptPath = filepath.Join(localRoot, initScriptName)
		_, err := os.Stat(initializationScriptPath)
		if err == nil {
			// TODO hash and check if it should be overwritten
		} else if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		// Only root needs to execute
		file, err := os.OpenFile(initializationScriptPath, os.O_CREATE | os.O_WRONLY, 0700)
		check(
			err,
			fmt.Sprintf("opening %s in write-only", initializationScriptPath),
		)

		template, err := template.New("config.Provision").Parse(config.Provision)
		check(
			err,
			fmt.Sprintf("trying to create template"),
		)
		err = template.Execute(file, config)
		check(
			err,
			"interpolating template",
		)
		file.Close()
	}

		// TODO copy .ssh folder?

	var archChroot = filepath.Join(localRoot, "bin", "arch-chroot")
	if initializationScriptPath != "" {
		// Note this will now be relative to chroot
		var cmd = exec.Command(archChroot, localRoot, fmt.Sprintf("/%s", initScriptName))
		fmt.Println("about to run arch-chroot")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		check(
			err,
			"running arch-chroot",
		)
		fmt.Println("back")
		Exec(archChroot, []string{localRoot, "su", "--login", config.UserName})
	}
}
