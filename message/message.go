package message

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	MessageId       string
	CorrelationId   string
	Timestamp       time.Time
	AppId           string
	ContentType     string
	ContentEncoding string
	Exchange        string
	RoutingKey      string
	Mandatory       bool
	Immediate       bool
	BodyBytes       []byte

	*http.Request
}

func NewMessage(req *http.Request) (*Message, error) {
	defer req.Body.Close()

	var (
		body      []byte
		err       error
		mandatory bool
		immediate bool
	)

	if body, err = ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	}

	reqPath := req.URL.Path
	pathParts := strings.Split(strings.TrimLeft(reqPath, "/"), "/")
	if len(pathParts) < 2 || len(pathParts[0]) == 0 || len(pathParts[1]) == 0 {
		return nil, fmt.Errorf("Missing required exchange and/or routing key "+
			"in PATH_INFO: %+v", reqPath)
	}

	reqQuery := req.URL.Query()
	if m := reqQuery.Get("m"); m == "1" {
		mandatory = true
	}

	if i := reqQuery.Get("i"); i == "1" {
		immediate = true
	}

	return &Message{
		req.Header.Get("Message-ID"),     // MessageID
		req.Header.Get("Correlation-ID"), //CorrelationID
		// FIXME parse "Date" header?
		time.Now().UTC(),                   // Timestamp
		req.Header.Get("From"),             // AppId
		req.Header.Get("Content-Type"),     // ContentType
		req.Header.Get("Content-Encoding"), //ContentEncoding
		pathParts[0],                       // Exchange
		pathParts[1],                       // RoutingKey
		mandatory,                          // Mandatory
		immediate,                          // Immediate
		body,                               // BodyString
		req,                                // *http.Request
	}, nil
}
