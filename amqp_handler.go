package mithril

import (
	"github.com/streadway/amqp"
)

var (
	pingMsg = &amqp.Publishing{Body: []byte("z")}
)

type AMQPHandler struct {
	amqpUri     string
	amqpConn    *amqp.Connection
	pingChannel *amqp.Channel
}

func NewAMQPHandler(amqpUri string) *AMQPHandler {
	return &AMQPHandler{amqpUri: amqpUri}
}

func (me *AMQPHandler) HandleRequest(req *Request) error {
	if err := me.ensureConnected(); err != nil {
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
	return me.pingChannel.Publish("", "mithril_pings", false, false, *pingMsg)
}
