// +build pg full

package main

import (
	"flag"

	"github.com/modcloth-labs/mithril"
)

var (
	enablePgFlag = flag.Bool("pg", false, "Enable PostgreSQL handler")
	pgUriFlag    = flag.String("pg.uri",
		"postgres://localhost?sslmode=disable", "PostgreSQL Server URI")
)

func init() {
	pipelineCallbacks["pg"] = func(pipeline mithril.Handler) mithril.Handler {
		if *enablePgFlag {
			mithril.Debugf("  --> pg enabled, so adding postgresql handler")
			pipeline = mithril.NewPostgreSQLHandler(*pgUriFlag, pipeline)
		} else {
			mithril.Debugf("  --> pg not enabled, so leaving pipeline unaltered")
		}
		return pipeline
	}
}
