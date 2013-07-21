package main

import "github.com/modcloth-labs/mithril"
import "github.com/modcloth-labs/versioning"

func init() {
	versioning.Parse()
}

func main() {
	mithril.ServerMain()
}
