package mithril

import (
	"net/http"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (me *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
