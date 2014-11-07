package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/jingweno/negroni-gorelic"

	"github.com/meatballhat/negroni-logrus"
	"github.com/modcloth/mithril"
	"github.com/modcloth/mithril/store"
)

var (
	logLevels = map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}

	logFormats = map[string]logrus.Formatter{
		"text": new(logrus.TextFormatter),
		"json": new(logrus.JSONFormatter),
	}
)

func main() {
	var (
		logLevelOptions  []string
		logFormatOptions []string
	)

	for s := range logLevels {
		logLevelOptions = append(logLevelOptions, s)
	}

	for s := range logFormats {
		logFormatOptions = append(logFormatOptions, s)
	}

	app := cli.NewApp()
	app.Usage = "HTTP -> AMQP proxy"
	app.Version = fmt.Sprintf("%s (%s)", mithril.Version, mithril.Rev)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: fmt.Sprintf("Log level (options: %s)", strings.Join(logLevelOptions, ",")),
		},
		cli.StringFlag{
			Name:  "log-format, f",
			Value: "text",
			Usage: fmt.Sprintf("Log format (options: %s)", strings.Join(logFormatOptions, ",")),
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "serve",
			ShortName:   "s",
			Usage:       "start server",
			Description: "Start the AMQP -> HTTP proxy server",
			Action: func(c *cli.Context) {
				level, err := getLogLevel(c)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				formatter, err := getLogFormatter(c)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				storer, err := store.Open(c.String("storage"), c.String("storage-uri"))
				if err != nil {
					log.Fatal(err)
				}

				amqp, err := mithril.NewAMQPPublisher(c.String("amqp-uri"))
				if err != nil {
					log.Fatal(err)
				}

				n := negroni.New(negroni.NewRecovery(), negronilogrus.NewCustomMiddleware(level, formatter, "mithril"))

				if c.Bool("new-relic-agent-enabled") {
					if c.String("new-relic-license-key") == "" {
						log.Warn("newrelic license key is absent")
					} else {
						n.Use(negronigorelic.New(
							c.String("new-relic-license-key"),
							fmt.Sprintf("[%s] %s", c.String("new-relic-env"), "Mithril"),
							c.Bool("new-relic-verbose")))
					}
				}

				n.UseHandler(mithril.NewServer(storer, amqp))
				n.Run(c.String("bind"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "storage, s",
					Usage:  "Which storage driver to use (see `list-storage` command).",
					Value:  "",
					EnvVar: "MITHRIL_STORAGE",
				},
				cli.StringFlag{
					Name:   "storage-uri, u",
					Usage:  "The url used by the storage driver.",
					Value:  "",
					EnvVar: "MITHRIL_STORAGE_URI",
				},
				cli.StringFlag{
					Name:   "amqp-uri, a",
					Usage:  "The url of the AMQP server",
					Value:  "amqp://guest:guest@localhost:5672",
					EnvVar: "MITHRIL_AMQP_URI",
				},
				cli.StringFlag{
					Name:   "bind, b",
					Usage:  "The address to bind to",
					Value:  ":8371",
					EnvVar: "MITHRIL_BIND",
				},
				cli.BoolTFlag{
					Name:   "new-relic-agent-enabled, E",
					Usage:  "Enable the NewRelic agent",
					EnvVar: "NEW_RELIC_AGENT_ENABLED",
				},
				cli.StringFlag{
					Name:   "new-relic-license-key, n",
					Usage:  "New Relic License Key",
					EnvVar: "NEW_RELIC_LICENSE_KEY",
				},
				cli.StringFlag{
					Name:   "new-relic-env, e",
					Usage:  "New Relic Env",
					EnvVar: "NEW_RELIC_ENV",
				},
				cli.BoolFlag{
					Name:   "new-relic-verbose, V",
					Usage:  "New Relic Logging Level",
					EnvVar: "NEW_RELIC_VERBOSE",
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

func getLogLevel(c *cli.Context) (level logrus.Level, err error) {
	level, ok := logLevels[c.GlobalString("log-level")]
	if !ok {
		return 0, fmt.Errorf("invalid log level %s", c.GlobalString("log-level"))
	}

	return level, nil
}

func getLogFormatter(c *cli.Context) (formatter logrus.Formatter, err error) {
	formatter, ok := logFormats[c.GlobalString("log-format")]
	if !ok {
		return nil, fmt.Errorf("invalid log format %s", c.GlobalString("log-format"))
	}

	return formatter, nil
}
