package main

import (
	"flag"
	"fmt"
	"mithril"
	"mithril/log"
	"os"
)

var (
	pidFlag     = flag.String("p", "", "PID file (only written if provided)")
	versionFlag = flag.Bool("version", false, "Print version and exit")
	revFlag     = flag.Bool("rev", false, "Print git revision and exit")
	debug       = flag.Bool("d", false, "Enable Debugging handler")
	pidFile     *os.File
)

func main() {
	var err error

	flag.Usage = func() {
		fmt.Println("Usage: mithril-server [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println(mithril.ProgVersion())
		return
	}

	if *revFlag {
		fmt.Println(mithril.Rev)
		return
	}

	if len(*pidFlag) > 0 {

		if pidFile, err = os.Create(*pidFlag); err != nil {
			log.Fatal(err)
		}
		defer func() { os.Remove(*pidFlag) }()
		fmt.Fprintf(pidFile, "%d\n", os.Getpid())
	}

	log.Initialize(*debug)
	mithril.ServerMain()
}
