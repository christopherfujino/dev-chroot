package main

import (
	"flag"
	"fmt"
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

	switch os.Args[1] {
	case "bootstrap":
		err := bootstrapCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		err = bootstrap(defaultConfig)
		if err != nil {
			panic(err)
		}
	case "attach":
		err := attachCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		err = attach(defaultConfig)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("Do not recognize sub-command %s", os.Args[1]))
	}
}
