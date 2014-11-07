package mithril

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/modcloth/mithril/message"
	"github.com/modcloth/mithril/store"
)

type Server struct {
	amqp    *AMQPPublisher
	storage *store.Storage
}

func NewServer(storer *store.Storage, amqp *AMQPPublisher) *Server {
	return &Server{
		storage: storer,
		amqp:    amqp,
	}
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "PUT" {
		me.processMessage(w, r)
	} else {
		me.respondErr(
			fmt.Errorf(`Only "POST" and "PUT" are accepted, not %q`, r.Method),
			http.StatusMethodNotAllowed,
			w)
	}
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
