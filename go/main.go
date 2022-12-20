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
		if err != nil {
			panic(err)
		}

		bootstrap(config, http.Get, cwd)
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
