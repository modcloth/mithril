package mithril

import (
	"github.com/codegangsta/cli"
)

type Configuration struct {
	Storage    string
	StorageUri string
	AmqpUri    string
}

func NewConfigurationFromContext(c *cli.Context) *Configuration {
	return &Configuration{
		Storage:    c.String("storage"),
		StorageUri: c.String("storage-uri"),
		AmqpUri:    c.String("amqp-uri"),
	}
}
