package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/arken/arken/config"
)

//GlobalFlags contains the flags for commands.
type GlobalFlags struct {
	Config  string `short:"c" long:"config" desc:"Specify a custom config path."`
	Verbose bool   `short:"v" long:"verbose" desc:"Show More Information"`
}

var Root = &cmd.Root{
	Name:    "arken",
	Short:   "Arken is a distributed archive system.",
	Version: config.Version,
	License: "Licensed under the Apache License, Version 2.0",
	Flags:   &GlobalFlags{},
}

func checkError(flags *GlobalFlags, err error) {
	if err != nil {
		if flags.Verbose {
			log.Fatal(err)
		}
		fmt.Println(err)
		os.Exit(1)
	}
}
