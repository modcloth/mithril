package mithril

import (
	"encoding/base64"
	"fmt"
	"mithril/log"
	"mithril/message"
	"mithril/store"
	"net/http"
	"time"

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

var faviconBytes []byte

type Server struct {
	amqp    *AMQPPublisher
	storage *store.Storage
	address string
}

func init() {
	faviconBytes, _ = base64.StdEncoding.DecodeString(faviconBase64)
}

func NewServer(configuration *Configuration) (*Server, error) {
	var (
		storer *store.Storage
		amqp   *AMQPPublisher
		err    error
	)

	if storer, err = store.Open(configuration.Storage, configuration.StorageUri); err != nil {
		return nil, err
	}

	if amqp, err = NewAMQPPublisher(configuration.AmqpUri); err != nil {
		return nil, err
	}

	return &Server{
		storage: storer,
		amqp:    amqp,
		address: configuration.ServerAddress,
	}, nil
}

func (me *Server) Serve() {
	http.Handle("/", me)
	log.Println("Serving on", me.address)
	log.Fatal(http.ListenAndServe(me.address, nil))
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	//r.Close = true

	log.Printf("%s %s Headers: %+v", r.Method, r.URL.Path, r.Header)

	if r.Method == "POST" || r.Method == "PUT" {
		me.processMessage(w, r)
	} else if r.Method == "GET" && r.URL.Path == "/favicon.ico" {
		me.respondFavicon(w)
	} else {
		me.respondErr(
			fmt.Errorf(`Only "POST" and "PUT" are accepted, not %q`, r.Method),
			http.StatusMethodNotAllowed,
			w)
	}

	log.Printf("Request handled in %s.\n", time.Now().Sub(startTime))
	return
}

func (me *Server) processMessage(w http.ResponseWriter, r *http.Request) {
	var (
		msg *message.Message
		err error
	)
	if msg, err = message.NewMessage(r); err != nil {
		me.respondErr(err, http.StatusBadRequest, w)
		return
	}

	log.Println("Processing message: ", msg.MessageId)
	if err = me.storage.Store(msg); err != nil {
		me.respondErr(err, http.StatusBadRequest, w)
		return
	}

	if err = me.amqp.Publish(msg); err != nil {
		me.respondErr(err, http.StatusBadRequest, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (me *Server) respondErr(err error, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "WOMP WOMP: %v\n", err)
}

func (me *Server) respondFavicon(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
	w.WriteHeader(http.StatusOK)
	w.Write(faviconBytes)
}
