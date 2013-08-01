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
		return
	}

	if config.DisplayRev {
		fmt.Println(mithril.Rev)
		return
	}

	if config.ShowStorage {
		store.ShowStorage()
		return
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
		panic(err)
	} else {
		server.Serve()
	}
}
