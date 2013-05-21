package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/modcloth-labs/mithril"
)

var (
	addrFlag = flag.String("a", ":8371", "Mithril server address")

	amqpUriFlag = flag.String("amqp.uri",
		"amqp://guest:guest@localhost:5672", "AMQP Server URI")
	pipelineCallbacks = map[string]func(mithril.Handler) mithril.Handler{}
	pipelineOrder     = []string{"pg"}

	pidFlag = flag.String("p", "", "PID file (only written if provided)")
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: mithril-server [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(*pidFlag) > 0 {
		var (
			f   *os.File
			err error
		)

		if f, err = os.Create(*pidFlag); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "%d\n", os.Getpid())
	}

	var pipeline mithril.Handler

	server := mithril.NewServer()
	pipeline = mithril.NewAMQPHandler(*amqpUriFlag, nil)

	for _, name := range pipelineOrder {
		callback := pipelineCallbacks[name]
		log.Printf("Calling %q pipeline callback", name)
		pipeline = callback(pipeline)
	}

	if err := pipeline.Init(); err != nil {
		log.Fatalf("Failed to initialize handler pipeline: %q", err)
	}

	server.SetHandlerPipeline(pipeline)
	http.Handle("/", server)

	log.Println("Serving on", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}
