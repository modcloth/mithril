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
	addr = flag.String("a", ":8371", "Mithril server address")

	amqpUri = flag.String("amqp.uri",
		"amqp://guest:guest@localhost:5672", "AMQP Server URI")

	enablePg = flag.Bool("pg", false, "Enable PostgreSQL handler")
	pgUri    = flag.String("pg.uri",
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
	pipeline = mithril.NewAMQPHandler(*amqpUri, nil)

	if *enablePg {
		pipeline = mithril.NewPostgreSQLHandler(*pgUri, pipeline)
	}

	if err := pipeline.Init(); err != nil {
		log.Fatalf("Failed to initialize handler pipeline: %q", err)
	}

	server.SetHandlerPipeline(pipeline)
	http.Handle("/", server)

	log.Println("Serving on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
