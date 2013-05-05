package mithril

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

const faviconBase64 = `
iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAgMAAABinRfyAAAABGdBTUEAALGPC/xhBQAAAA
FzUkdCAK7OHOkAAAAgY0hSTQAAeiYAAICEAAD6AAAAgOgAAHUwAADqYAAAOpgAABdwnLpR
PAAAAAlQTFRFAAAAaycz////cJGUEAAAAAF0Uk5TAEDm2GYAAAABYktHRAJmC3xkAAAAPE
lEQVQI1zWMsQ0AMAjDzNAT+AeG7q0E/79SqMRiJVYUViBoZhKyA8QArYD74O5xH13F3Tg9
jlXQLF9XPG8QCLmv6srMAAAAJXRFWHRkYXRlOmNyZWF0ZQAyMDEzLTA1LTAzVDIyOjI1Oj
AzLTA0OjAwA+emcwAAACV0RVh0ZGF0ZTptb2RpZnkAMjAxMy0wNS0wM1QyMjoyNTowMi0w
NDowMNTNFXsAAAAASUVORK5CYII=
`

var faviconBytes []byte

func init() {
	faviconBytes, _ = base64.StdEncoding.DecodeString(faviconBase64)
}

type Server struct {
	handlerPipeline Handler
}

func NewServer() *Server {
	return &Server{}
}

func (me *Server) SetHandlerPipeline(handler Handler) {
	me.handlerPipeline = handler
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
		req    Request
	)

	defer func() {
		log.Printf("\"%v %v %v\" %v -", r.Method, r.URL, r.Proto, status)
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

	if req, err = me.parseRequest(r); err != nil {
		status = http.StatusBadRequest
		me.respondErr(err, status, w)
		return
	}

	if err = me.handlerPipeline.HandleRequest(req); err != nil {
		status = http.StatusInternalServerError
		me.respondErr(err, status, w)
		return
	}

	status = http.StatusNoContent
	me.respond(status, []byte(""), w)
}

func (me *Server) parseRequest(r *http.Request) (Request, error) {
	return &HTTPRequestWrapper{Req: r}, nil
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

func (me *Server) respondFavicon(status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(status)
	w.Write(faviconBytes)
}
