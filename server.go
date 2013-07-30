package mithril

import (
	"encoding/base64"
	"flag"
	"fmt"
	"mithril/log"
	"net/http"

	_ "net/http/pprof" // hey, why not
)

const faviconBase64 = `
AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCw
AAAAAAAAAAAAD//////////zMna/8zJ2v/Mydr/zMna/8zJ2v/////////////////////
/////////////////////////////////zMna/8zJ2v/Mydr/zMna/8zJ2v/Mydr/zMna/
///////////////////////////////////////////zMna/8zJ2v/Mydr////////////
/////zMna/8zJ2v/Mydr//////////////////////////////////////8zJ2v/Mydr//
//////////////////////////Mydr/zMna///////////////////////////////////
////////////////////Mydr/zMna////////////zMna/8zJ2v///////////8zJ2v/My
dr//////////////////////////////////////8zJ2v/Mydr//////8zJ2v/Mydr////
//8zJ2v/Mydr/////////////////////////////////////////////////zMna/8zJ2
v/Mydr/zMna/8zJ2v/Mydr////////////////////////////////////////////////
//////8zJ2v/Mydr/zMna/8zJ2v/Mydr/zMna/////////////////////////////////
////////////////8zJ2v/Mydr//////8zJ2v/Mydr//////8zJ2v/Mydr////////////
////////////////////////////////Mydr////////////Mydr/zMna////////////z
Mna////////////////////////////////////////////zMna////////////zMna/8z
J2v///////////8zJ2v/////////////////////////////////////////////////My
dr/zMna/8zJ2v/Mydr/zMna/8zJ2v/////////////////////////////////////////
////////////////////////Mydr/zMna/////////////////////////////////////
//Mydr/zMna/8zJ2v//////////////////////zMna/8zJ2v//////zMna/8zJ2v/Mydr
/zMna/8zJ2v///////////8zJ2v/Mydr/zMna/8zJ2v/Mydr/zMna/8zJ2v/Mydr/zMna/
8zJ2v/Mydr/zMna/8zJ2v/Mydr/zMna////////////zMna/8zJ2v/Mydr/zMna/8zJ2v/
//////////////////////////////////////////8zJ2v/AAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==
`

var (
	faviconBytes []byte

	addrFlag    = flag.String("a", ":8371", "Mithril server address")

	amqpUriFlag       = flag.String("amqp.uri", "amqp://guest:guest@localhost:5672", "AMQP Server URI")
	pipelineCallbacks = map[string]func(Handler) Handler{}
	pipelineOrder     = []string{"debug", "pg"}

)

func init() {
	faviconBytes, _ = base64.StdEncoding.DecodeString(faviconBase64)
}

// ServerMain is the entry point used by `mithril-server`
func ServerMain() {

	var pipeline Handler

	server := newServer()
	pipeline = NewAMQPHandler(*amqpUriFlag, nil)

	for _, name := range pipelineOrder {
		if callback, ok := pipelineCallbacks[name]; ok {
			log.Println("Calling %q pipeline callback\n", name)
			pipeline = callback(pipeline)
		}
	}

	if err := pipeline.Init(); err != nil {
		log.Fatalf("Failed to initialize handler pipeline: %q", err)
	}

	server.SetHandlerPipeline(pipeline)
	http.Handle("/", server)

	log.Println("Serving on", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}

type server struct {
	handlerPipeline Handler
}

func newServer() *server {
	return &server{}
}

func (me *server) SetHandlerPipeline(handler Handler) {
	me.handlerPipeline = handler
}

func (me *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)

	defer func() { // defer to capture status on the method exit
		log.Println("\"%v %v %v\" %v -\n", r.Method, r.URL, r.Proto, status)
	}()

	if r.Method == "GET" && r.URL.Path == "/favicon.ico" {
		status = http.StatusOK
		me.respondFavicon(status, w)
		return
	}

	if r.Method != "POST" && r.Method != "PUT" {
		status = http.StatusMethodNotAllowed
		err = fmt.Errorf(`Only "POST" and "PUT" are accepted, not %q`, r.Method)
		me.respondErr(err, status, w)
		return
	}

	fReq, err := NewFancyRequest(r)
	if err != nil {
		status = http.StatusBadRequest
		me.respondErr(err, status, w)
		return
	}

	if err = me.handlerPipeline.HandleRequest(fReq); err != nil {
		status = http.StatusInternalServerError
		me.respondErr(err, status, w)
		return
	}

	status = http.StatusNoContent
	me.respond(status, []byte(""), w)
}

func (me *server) respondErr(err error, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "WOMP WOMP: %v\n", err)
}

func (me *server) respond(status int, body []byte, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write(body)
}

func (me *server) respondFavicon(status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
	w.WriteHeader(status)
	w.Write(faviconBytes)
}
