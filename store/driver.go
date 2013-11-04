package store

import "github.com/modcloth-labs/mithril/message"

type Driver interface {
	Init(uri string) error
	Store(msg *message.Message) error
	UriFormat() string
}
