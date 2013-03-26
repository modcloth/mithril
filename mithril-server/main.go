package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/modcloth-labs/mithril"
)

var (
	addr    = flag.String("a", ":8371", "Mithril server address")
	amqpUri = flag.String("u", "amqp://guest:guest@localhost:5672", "AMQP Server URI")
)

func main() {
	server := mithril.NewServer()
	server.AddHandler(mithril.NewAMQPHandler(*amqpUri))
	http.Handle("/", server)
	log.Println("Serving on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
