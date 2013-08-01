package mithril

import (
	"fmt"
	"mithril/log"
	"mithril/message"

	"github.com/streadway/amqp" // explicitly cloned into place
)

type amqpAdaptedRequest struct {
	Publishing *amqp.Publishing
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
}

type AMQPPublisher struct {
	amqpUri         string
	amqpConn        *amqp.Connection
	handlingChannel *amqp.Channel
	confirmAck      chan uint64
	confirmNack     chan uint64
}

func NewAMQPPublisher(amqpUri string) (*AMQPPublisher, error) {
	publisher := &AMQPPublisher{
		amqpUri: amqpUri,
	}

	if err := publisher.establishConnection(); err != nil {
		return nil, err
	}
	return publisher, nil
}

func (me *AMQPPublisher) Publish(req *message.Message) error {
	var (
		amqpReq *amqpAdaptedRequest
		err     error
	)

	amqpReq = me.adaptHttpRequest(req)
	if err = me.publishAdaptedRequest(amqpReq); err != nil {
		log.Println("Failed to publish request:", err)
		return err
	}
	return nil
}

func (me *AMQPPublisher) establishConnection() error {
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

func (me *AMQPPublisher) disconnect() {
	if me.handlingChannel != nil {
		me.handlingChannel.Close()
		me.handlingChannel = nil
	}

	if me.amqpConn != nil {
		me.amqpConn.Close()
		me.amqpConn = nil
	}
}

func (me *AMQPPublisher) adaptHttpRequest(req *message.Message) *amqpAdaptedRequest {
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

func (me *AMQPPublisher) publishAdaptedRequest(amqpReq *amqpAdaptedRequest) error {
	log.Println("Publishing adapted HTTP request %+v\n", amqpReq)

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
