// +build pg

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
			pipeline = mithril.NewPostgreSQLHandler(*pgUriFlag, pipeline)
		}
		return pipeline
	}
}
