package mithril

import (
	"github.com/codegangsta/cli"
)

type Configuration struct {
	EnableDebug   bool
	ServerAddress string
	Storage       string
	StorageUri    string
	AmqpUri       string
}

func NewConfigurationFromContext(c *cli.Context) *Configuration {
	return &Configuration{
		EnableDebug:   c.Bool("debug"),
		Storage:       c.String("storage"),
		StorageUri:    c.String("storage-uri"),
		ServerAddress: c.String("bind"),
		AmqpUri:       c.String("amqp-uri"),
	}
}
