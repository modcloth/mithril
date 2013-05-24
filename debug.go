// +build debug full

package mithril

import (
	"flag"
	"log"
)

var DebugEnabled = false

func init() {
	flag.BoolVar(&DebugEnabled, "d", false, "Enable Debugging handler")
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
