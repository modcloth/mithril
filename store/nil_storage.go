package store

import (
	"github.com/modcloth/mithril/log"
	"github.com/modcloth/mithril/message"
)

type nilstore struct {
}

func init() {
	register("", &nilstore{})
}
func (me *nilstore) Init(url string) error {
	log.Println("Using the nil logger, no messages will be stored")
	return nil
}
func (me *nilstore) UriFormat() string {
	return "Nil Logger"
}
func (me *nilstore) Store(msg *message.Message) error {
	return nil
}
