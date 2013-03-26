package mithril

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	handlers []Handler
}

func NewServer() *Server {
	return &Server{}
}

func (me *Server) AddHandler(handler Handler) {
	me.handlers = append(me.handlers, handler)
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
		req    *Request
	)

	defer func() {
		log.Printf("\"%v %v %v\" %v -", r.Method, r.URL, r.Proto, status)
	}()

	if req, err = me.parseRequest(r); err != nil {
		status = http.StatusBadRequest
		me.respondErr(err, status, w)
		return
	}

	for _, handler := range me.handlers {
		err = handler.HandleRequest(req)
		if err != nil {
			status = http.StatusInternalServerError
			me.respondErr(err, status, w)
			return
		}
	}

	status = http.StatusNoContent
	me.respond(status, []byte(""), w)
}

func (me *Server) parseRequest(r *http.Request) (*Request, error) {
	return &Request{}, nil
}

func (me *Server) respondErr(err error, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "WOMP WOMP: %v\n", err)
}

func (me *Server) respond(status int, body []byte, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write(body)
}
