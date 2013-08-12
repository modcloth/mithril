package main

import (
	"fmt"
	"mithril"
	"mithril/log"
	"mithril/store"
	"os"
)

func main() {
	config := mithril.NewConfigurationFromFlags()

	if config.DisplayVersion {
		fmt.Println(mithril.ProgVersion())
	}

	if config.DisplayRev {
		fmt.Println(mithril.Rev)
	}

	if config.ShowStorage {
		store.ShowStorage()
	}

	if config.ExitImmediate {
		os.Exit(1)
	}

	if len(config.PidFile) > 0 {
		if pidFile, err := os.Create(config.PidFile); err != nil {
			log.Fatal(err)
		} else {
			defer func() { os.Remove(config.PidFile) }()
			fmt.Fprintf(pidFile, "%d\n", os.Getpid())
		}
	}

	log.Initialize(config.EnableDebug)
	log.Println("Initializing Mithril...")
	if server, err := mithril.NewServer(config); err != nil {
		log.Fatal(err)
	} else {
		server.Serve()
	}
}
