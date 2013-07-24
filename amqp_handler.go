package mithril

import (
	"fmt"

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
	confirmAck      chan uint64
	confirmNack     chan uint64
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

	Debugln("AMQP handler initialized")

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

	defer me.disconnect()

	amqpReq = me.adaptHttpRequest(req)

	if err = me.publishAdaptedRequest(amqpReq); err != nil {
		Debugln("Failed to publish request:", err)
		return err
	}

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

	if err = handlingChannel.Confirm(false); err != nil {
		return err
	}

	me.confirmAck, me.confirmNack = handlingChannel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

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
			MessageId:     req.MessageId,
			CorrelationId: req.CorrelationId,
			Timestamp:     req.Timestamp,
			AppId:         req.AppId,
			ContentType:   req.ContentType,
			Body:          req.BodyBytes,
		},
		Exchange:   req.Exchange,
		RoutingKey: req.RoutingKey,
		Mandatory:  req.Mandatory,
		Immediate:  req.Immediate,
	}
}

func (me *AMQPHandler) publishAdaptedRequest(amqpReq *amqpAdaptedRequest) error {
	Debugf("Publishing adapted HTTP request %+v\n", amqpReq)

	err := me.handlingChannel.Publish(amqpReq.Exchange,
		amqpReq.RoutingKey, amqpReq.Mandatory,
		amqpReq.Immediate, *amqpReq.Publishing)

	if err != nil {
		return err
	}

	select {
	case _ = <-me.confirmAck:
		return nil
	case _ = <-me.confirmNack:
		return fmt.Errorf("RabbitMQ nack'd message")
	}

	panic("I shouldn't be here")
}
