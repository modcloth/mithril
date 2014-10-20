package store

import "github.com/modcloth/mithril/message"

type Driver interface {
	Init(uri string) error
	Store(msg *message.Message) error
	UriFormat() string
}
