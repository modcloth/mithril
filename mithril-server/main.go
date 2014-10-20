package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/modcloth/mithril"
	"github.com/modcloth/mithril/log"
	"github.com/modcloth/mithril/store"
)

func main() {
	app := cli.NewApp()
	app.Usage = "HTTP -> AMQP proxy"
	app.Version = fmt.Sprintf("%s (%s)", mithril.Version, mithril.Rev)
	app.Commands = []cli.Command{
		{
			Name:        "server",
			ShortName:   "s",
			Usage:       "start server",
			Description: "Start the AMQP -> HTTP proxy server",
			Action: func(c *cli.Context) {
				config := mithril.NewConfigurationFromContext(c)

				log.Initialize(config.EnableDebug)
				log.Println("Initializing Mithril...")
				if server, err := mithril.NewServer(config); err != nil {
					log.Fatal(err)
				} else {
					server.Serve()
				}
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "debug, d",
					Usage:  "Enable debug logging.",
					EnvVar: "MITHRIL_DEBUG",
				},
				cli.StringFlag{
					Name:   "storage-uri, s",
					Usage:  "The url used by the storage driver.",
					Value:  "",
					EnvVar: "MITHRIL_STORAGE_URI",
				},
				cli.StringFlag{
					Name:   "amqp-uri, a",
					Usage:  "The url of the AMQP server",
					Value:  "amqp://localhost:5672",
					EnvVar: "MITHRIL_AMQP_URI",
				},
				cli.StringFlag{
					Name:   "bind, b",
					Usage:  "The address to bind to",
					Value:  ":8371",
					EnvVar: "MITHRIL_AMQP_URI",
				},
			},
		},
		{
			Name:        "list-storage",
			ShortName:   "l",
			Usage:       "list storage backends",
			Description: "List the avaliable storage backends for Mithril",
			Action: func(c *cli.Context) {
				store.ShowStorage()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
