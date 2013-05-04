package mithril

import (
	"fmt"
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

type HTTPRequestWrapper struct {
	Req *http.Request
}

func (me *HTTPRequestWrapper) Body() io.Reader {
	return me.Req.Body
}

func (me *HTTPRequestWrapper) Headers() *http.Header {
	return &me.Req.Header
}

func (me *HTTPRequestWrapper) Path() string {
	return me.Req.URL.Path
}

func (me *HTTPRequestWrapper) String() string {
	return fmt.Sprintf("%s %s %s", me.Req.Method, me.Req.URL.Path, me.Req.Proto)
}

func (me *HTTPRequestWrapper) Query() url.Values {
	return me.Req.URL.Query()
}
