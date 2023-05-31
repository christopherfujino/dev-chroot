package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	bootstrapCmd := flag.NewFlagSet(
		"bootstrap",
		flag.PanicOnError,
	)
	var uid = bootstrapCmd.Int(
		"uid",
		-1,
		"user ID (required)",
	)
	var homeDir = bootstrapCmd.String(
		"home-dir",
		"",
		"path to home directory (required)",
	)
	attachCmd := flag.NewFlagSet(
		"attach",
		flag.PanicOnError,
	)

	if len(os.Args) < 2 {
		bootstrapCmd.Usage()
		attachCmd.Usage()
		os.Exit(1)
	}

	// TODO validate cwd has valid config file, then read config from that file
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var config = defaultConfig

	switch os.Args[1] {
	case "bootstrap":
		err := bootstrapCmd.Parse(os.Args[2:])
		check(
			err,
			fmt.Sprintf("parsing arguments to bootstrap: %v", os.Args[2:]),
		)
		if *uid == -1 {
			panic("Must provide the --uid option")
		}
		if *homeDir == "" {
			panic("Must provide the --home-dir option")
		}
		if err != nil {
			panic(err)
		}

		bootstrap(config, http.Get, *uid, cwd, *homeDir)
	case "attach":
		err := attachCmd.Parse(os.Args[2:])
		check(
			err,
			fmt.Sprintf("parsing arguments to attach: %v", os.Args[2:]),
		)
		attach(config, cwd)
	default:
		panic(fmt.Errorf("unrecognized sub-command %s", os.Args[1]))
	}
}
