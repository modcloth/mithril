// +build debug full

package mithril

import (
	"flag"
	"log"

	_ "net/http/pprof"
)

var DebugEnabled = false

func init() {
	flag.BoolVar(&DebugEnabled, "d", false, "Enable Debugging handler")

	pipelineCallbacks["debug"] = func(pipeline Handler) Handler {
		if DebugEnabled {
			Debugf("  --> debug enabled, so adding debugging handler")
			pipeline = NewDebuggingHandler(pipeline)
		} else {
			Debugf("  --> debug not enabled, so leaving pipeline unaltered")
		}
		return pipeline
	}
}

func Debugf(format string, args ...interface{}) {
	if DebugEnabled {
		log.Printf(format, args...)
	}
}

func Debugln(args ...interface{}) {
	if DebugEnabled {
		log.Println(args...)
	}
}
