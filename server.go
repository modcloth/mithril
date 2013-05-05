package mithril

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
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
	w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
	w.WriteHeader(status)
	w.Write(faviconBytes)
}
