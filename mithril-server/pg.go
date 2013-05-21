// +build pg full

package main

import (
	"flag"
	"log"

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
			log.Printf("  --> pg enabled, so adding postgresql handler")
			pipeline = mithril.NewPostgreSQLHandler(*pgUriFlag, pipeline)
		} else {
			log.Printf("  --> pg not enabled, so leaving pipeline unaltered")
		}
		return pipeline
	}
}
