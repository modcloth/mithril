package mithril

import (
	"fmt"
	"os"
	"path"
)

var (
	VersionString string
	RevString     string
)

func init() {
	progName := path.Base(os.Args[0])
	for _, arg := range os.Args {
		if arg == "--version" {
			if VersionString == "" {
				VersionString = "<unknown>"
			}
			fmt.Printf("%s %s\n", progName, VersionString)
			os.Exit(0)
		} else if arg == "--rev" {
			if RevString == "" {
				RevString = "<unknown>"
			}
			fmt.Printf("%s\n", RevString)
			os.Exit(0)
		}
	}
}
