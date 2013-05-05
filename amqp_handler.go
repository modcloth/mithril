package mithril

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

type amqpAdaptedRequest struct {
	Publishing *amqp.Publishing
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
}

type AMQPHandler struct {
	amqpUri         string
	amqpConn        *amqp.Connection
	handlingChannel *amqp.Channel
	nextHandler     Handler
}

func NewAMQPHandler(amqpUri string, next Handler) *AMQPHandler {
	amqpHandler := &AMQPHandler{
		amqpUri: amqpUri,
	}

	amqpHandler.SetNextHandler(next)
	return amqpHandler
}

func (me *AMQPHandler) Init() error {
	var err error

	if err = me.establishConnection(); err != nil {
		return err
	}

	log.Println("AMQP handler initialized")

	if me.nextHandler != nil {
		return me.nextHandler.Init()
	}

	return nil
}

func (me *AMQPHandler) SetNextHandler(handler Handler) {
	me.nextHandler = handler
}

func (me *AMQPHandler) HandleRequest(req *http.Request) error {
	var (
		amqpReq *amqpAdaptedRequest
		err     error
	)

	if err = me.establishConnection(); err != nil {
		return err
	}

	if amqpReq, err = me.adaptHttpRequest(req); err != nil {
		return err
	}

	if err = me.publishAdaptedRequest(amqpReq); err != nil {
		log.Println("Failed to publish request:", err)
		return err
	}

	defer me.disconnect()

	return nil
}

func (me *AMQPHandler) establishConnection() error {
	conn, err := amqp.Dial(me.amqpUri)
	if err != nil {
		return err
	}

	me.amqpConn = conn

	handlingChannel, err := me.amqpConn.Channel()
	if err != nil {
		return err
	}

	me.handlingChannel = handlingChannel
	return nil
}

func (me *AMQPHandler) disconnect() {
	if me.handlingChannel != nil {
		me.handlingChannel.Close()
		me.handlingChannel = nil
	}

	if me.amqpConn != nil {
		me.amqpConn.Close()
		me.amqpConn = nil
	}
}

func (me *AMQPHandler) adaptHttpRequest(req *http.Request) (*amqpAdaptedRequest, error) {
	var (
		body      []byte
		err       error
		mandatory bool
		immediate bool
	)

	log.Printf("Adapting HTTP request %q", req)

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

	adaptedReq := &amqpAdaptedRequest{
		Publishing: &amqp.Publishing{
			MessageId:   req.Header.Get("Message-ID"),
			Timestamp:   time.Now().UTC(), // FIXME parse "Date" header?
			AppId:       req.Header.Get("From"),
			ContentType: req.Header.Get("Content-Type"),
			Body:        body,
		},
		Exchange:   pathParts[0],
		RoutingKey: pathParts[1],
		Mandatory:  mandatory,
		Immediate:  immediate,
	}

	return adaptedReq, nil
}

func (me *AMQPHandler) publishAdaptedRequest(amqpReq *amqpAdaptedRequest) error {
	log.Printf("Publishing adapted HTTP request %+v", amqpReq)

	return me.handlingChannel.Publish(amqpReq.Exchange,
		amqpReq.RoutingKey, amqpReq.Mandatory,
		amqpReq.Immediate, *amqpReq.Publishing)
}
