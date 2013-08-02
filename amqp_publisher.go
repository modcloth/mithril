package mithril

import (
	"fmt"
	"mithril/log"
	"mithril/message"

	"github.com/streadway/amqp"
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
	notifyClose     chan *amqp.Error
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
		log.Println("amqp - Failed to publish request:", err)
		return err
	}
	return nil
}

func (me *AMQPPublisher) establishConnection() (err error) {

	if me.amqpConn != nil {
		return
	}

	log.Printf("amqp - connecting to rabbitmq...")
	me.amqpConn, err = amqp.Dial(me.amqpUri)
	if err != nil {
		return err
	}
	log.Printf("amqp - connected to rabbitmq")

	log.Printf("amqp - creating channel...")
	me.handlingChannel, err = me.amqpConn.Channel()
	if err != nil {
		return err
	}
	log.Printf("amqp - channel created")

	log.Printf("amqp - setting confirm mode...")
	if err = me.handlingChannel.Confirm(false); err != nil {
		return err
	}
	log.Printf("amqp - confirm mode set")

	me.confirmAck, me.confirmNack = me.handlingChannel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))
	log.Printf("amqp - notify confirm channels created.")
	go func() {
		closeChan := me.amqpConn.NotifyClose(make(chan *amqp.Error))
		select {
		case e := <-closeChan:
			log.Printf("amqp - The connection to rabbitmq has been closed. %d: %s", e.Code, e.Reason)
			me.disconnect()
		}
	}()

	log.Printf("amqp - Ready to publish messages!")
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

func (me *AMQPPublisher) publishAdaptedRequest(amqpReq *amqpAdaptedRequest) (err error) {
	err = me.establishConnection()
	if err != nil {
		return err
	}
	err = me.handlingChannel.Publish(amqpReq.Exchange,
		amqpReq.RoutingKey, amqpReq.Mandatory,
		amqpReq.Immediate, *amqpReq.Publishing)

	if err != nil {
		return err
	}

	select {
	case _ = <-me.confirmAck:
		return nil
	case _ = <-me.confirmNack:
		return fmt.Errorf("amqp - RabbitMQ nack'd message")
	}

	panic("amqp - I shouldn't be here")
}