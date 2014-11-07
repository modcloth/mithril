package mithril

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/modcloth/mithril/message"
	"github.com/modcloth/mithril/store"
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

func NewServer(storer *store.Storage, amqp *AMQPPublisher) *Server {
	return &Server{
		storage: storer,
		amqp:    amqp,
	}
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	r.Close = true

	log.Debugf("%s %s Headers: %+v", r.Method, r.URL.Path, r.Header)

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

	log.Infof("Request handled in %s.\n", time.Now().Sub(startTime))
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

	log.Infof("Processing message: ", msg.MessageId)
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
