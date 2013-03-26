package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/modcloth-labs/mithril"
)

var (
	addr = flag.String("a", ":8371", "Server address")
)

func main() {
	http.Handle("/", mithril.NewServer())
	log.Println("Serving on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
