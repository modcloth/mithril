package mithril

import (
	"net/http"
)

type Handler interface {
	HandleRequest(*http.Request) error
	Init() error
	SetNextHandler(Handler)
}
