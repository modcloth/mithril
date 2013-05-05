package mithril

import (
	"log"

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

func (me *AMQPHandler) HandleRequest(req *FancyRequest) error {
	var (
		amqpReq *amqpAdaptedRequest
		err     error
	)

	if err = me.establishConnection(); err != nil {
		return err
	}

	amqpReq = me.adaptHttpRequest(req)

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

func (me *AMQPHandler) adaptHttpRequest(req *FancyRequest) *amqpAdaptedRequest {
	return &amqpAdaptedRequest{
		Publishing: &amqp.Publishing{
			MessageId:   req.MessageId,
			Timestamp:   req.Timestamp,
			AppId:       req.AppId,
			ContentType: req.ContentType,
			Body:        req.BodyBytes,
		},
		Exchange:   req.Exchange,
		RoutingKey: req.RoutingKey,
		Mandatory:  req.Mandatory,
		Immediate:  req.Immediate,
	}
}

func (me *AMQPHandler) publishAdaptedRequest(amqpReq *amqpAdaptedRequest) error {
	log.Printf("Publishing adapted HTTP request %+v", amqpReq)

	return me.handlingChannel.Publish(amqpReq.Exchange,
		amqpReq.RoutingKey, amqpReq.Mandatory,
		amqpReq.Immediate, *amqpReq.Publishing)
}
