package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/modcloth-labs/mithril"
)

var (
	addrFlag = flag.String("a", ":8371", "Mithril server address")

	amqpUriFlag = flag.String("amqp.uri",
		"amqp://guest:guest@localhost:5672", "AMQP Server URI")

	enablePgFlag = flag.Bool("pg", false, "Enable PostgreSQL handler")
	pgUriFlag    = flag.String("pg.uri",
		"postgres://localhost?sslmode=disable", "PostgreSQL Server URI")
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: mithril-server [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var pipeline mithril.Handler

	server := mithril.NewServer()
	pipeline = mithril.NewAMQPHandler(*amqpUriFlag, nil)

	if *enablePgFlag {
		pipeline = mithril.NewPostgreSQLHandler(*pgUriFlag, pipeline)
	}

	if err := pipeline.Init(); err != nil {
		log.Fatalf("Failed to initialize handler pipeline: %q", err)
	}

	server.SetHandlerPipeline(pipeline)
	http.Handle("/", server)

	log.Println("Serving on", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}
