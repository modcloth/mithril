package mithril

import (
	"fmt"
	"os"
	"path"
)

var (
	// Version is the `git describe` string embedded via ldflags
	Version string
	// Rev is the `git rev-parse` string embedded via ldflags
	Rev      string
	progName string
)

func init() {
	progName = path.Base(os.Args[0])
	if Version == "" {
		Version = "<unknown>"
	}
	if Rev == "" {
		Rev = "<unknown>"
	}
}

func ProgVersion() string {
	return fmt.Sprintf("%s %s", progName, Version)
}
