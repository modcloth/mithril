package mithril

import (
	"io"
	"net/http"
	"net/url"
)

type Request interface {
	Body() io.Reader
	Path() string
	Query() url.Values
	Headers() *http.Header
	String() string
}
