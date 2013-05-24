// +build debug full

package main

import (
	"github.com/modcloth-labs/mithril"
)

func init() {
	pipelineCallbacks["debug"] = func(pipeline mithril.Handler) mithril.Handler {
		if mithril.DebugEnabled {
			mithril.Debugf("  --> debug enabled, so adding debugging handler")
			pipeline = mithril.NewDebuggingHandler(pipeline)
		} else {
			mithril.Debugf("  --> debug not enabled, so leaving pipeline unaltered")
		}
		return pipeline
	}
}
