package mithril

import (
	"sync"
	"time"

	"github.com/modcloth/mithril/message"

	log "github.com/Sirupsen/logrus"
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
	sync.Mutex
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
		log.Warnf("amqp - Failed to publish request: %v", err)
		return err
	}
	return nil
}

func (me *AMQPPublisher) establishConnection() (err error) {
	me.Lock()
	defer me.Unlock()

	if me.amqpConn != nil {
		return
	}

	log.Info("amqp - no RabbitMQ connection found, establishing new connection...")
	me.amqpConn, err = amqp.Dial(me.amqpUri)
	if err != nil {
		return err
	}
	log.Debug("amqp - connected to RabbitMQ")

	log.Debug("amqp - creating channel...")
	me.handlingChannel, err = me.amqpConn.Channel()
	if err != nil {
		return err
	}
	log.Debug("amqp - channel created")

	log.Debug("amqp - setting confirm mode...")
	if err = me.handlingChannel.Confirm(false); err != nil {
		return err
	}
	log.Debug("amqp - confirm mode set")

	me.confirmAck, me.confirmNack = me.handlingChannel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))
	log.Debug("amqp - notify confirm channels created.")

	go func() {
		closeChan := me.handlingChannel.NotifyClose(make(chan *amqp.Error))

		select {
		case e := <-closeChan:
			log.Printf("amqp - The channel opened with RabbitMQ has been closed. %d: %s", e.Code, e.Reason)
			me.disconnect()
		}
	}()

	log.Info("amqp - Ready to publish messages!")
	return nil
}

func (me *AMQPPublisher) disconnect() {
	me.Lock()
	defer me.Unlock()

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
			MessageId:       req.MessageId,
			CorrelationId:   req.CorrelationId,
			Timestamp:       req.Timestamp,
			AppId:           req.AppId,
			ContentType:     req.ContentType,
			ContentEncoding: req.ContentEncoding,
			Body:            req.BodyBytes,
			DeliveryMode:    amqp.Persistent,
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

	me.Lock()
	err = me.handlingChannel.Publish(amqpReq.Exchange,
		amqpReq.RoutingKey, amqpReq.Mandatory,
		amqpReq.Immediate, *amqpReq.Publishing)
	me.Unlock()

	if err != nil {
		return err
	}

	select {
	case _ = <-me.confirmAck:
		return nil
	case _ = <-me.confirmNack:
		log.Printf("amqp - RabbitMQ nack'd a message at %s", time.Now().UTC())
		return nil
	}
}
