package mithril

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

var (
	pingMsg = &amqp.Publishing{Body: []byte("z")}
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
	pingChannel     *amqp.Channel
	handlingChannel *amqp.Channel
}

func NewAMQPHandler(amqpUri string) *AMQPHandler {
	return &AMQPHandler{amqpUri: amqpUri}
}

func (me *AMQPHandler) HandleRequest(req Request) error {
	var (
		amqpReq *amqpAdaptedRequest
		err     error
	)

	if err = me.ensureConnected(); err != nil {
		return err
	}

	if amqpReq, err = me.adaptHttpRequest(req); err != nil {
		return err
	}

	if err = me.publishAdaptedRequest(amqpReq); err != nil {
		return err
	}

	return nil
}

func (me *AMQPHandler) ensureConnected() error {
	if me.isConnected() {
		return nil
	}

	return me.establishConnection()
}

func (me *AMQPHandler) establishConnection() error {
	if me.amqpConn != nil {
		me.amqpConn.Close()
	}

	conn, err := amqp.Dial(me.amqpUri)
	if err != nil {
		return err
	}

	me.amqpConn = conn

	pingChannel, err := me.amqpConn.Channel()
	if err != nil {
		return err
	}

	me.pingChannel = pingChannel

	handlingChannel, err := me.amqpConn.Channel()
	if err != nil {
		return err
	}

	me.handlingChannel = handlingChannel
	return nil
}

func (me *AMQPHandler) isConnected() bool {
	if me.amqpConn == nil {
		return false
	}

	if me.pingChannel == nil {
		return false
	}

	return me.publishPing() == nil
}

func (me *AMQPHandler) publishPing() error {
	log.Println("Pinging RabbitMQ")
	return me.pingChannel.Publish("", "mithril_pings", false, false, *pingMsg)
}

func (me *AMQPHandler) adaptHttpRequest(req Request) (*amqpAdaptedRequest, error) {
	var (
		body      []byte
		err       error
		mandatory bool
		immediate bool
	)

	log.Printf("Adapting HTTP request %q", req)

	if body, err = ioutil.ReadAll(req.Body()); err != nil {
		return nil, err
	}

	reqPath := req.Path()
	pathParts := strings.Split(strings.TrimLeft(reqPath, "/"), "/")
	if len(pathParts) < 2 || len(pathParts[0]) == 0 || len(pathParts[1]) == 0 {
		return nil, fmt.Errorf("Missing required exchange and/or routing key "+
			"in PATH_INFO: %+v", reqPath)
	}

	reqQuery := req.Query()
	if m := reqQuery.Get("m"); m == "1" {
		mandatory = true
	}

	if i := reqQuery.Get("i"); i == "1" {
		immediate = true
	}

	adaptedReq := &amqpAdaptedRequest{
		Publishing: &amqp.Publishing{
			MessageId:   req.Headers().Get("Message-ID"),
			Timestamp:   time.Now().UTC(), // FIXME parse "Date" header?
			AppId:       req.Headers().Get("From"),
			ContentType: req.Headers().Get("Content-Type"),
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
