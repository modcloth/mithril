package mithril

type Handler interface {
	HandleRequest(Request) error
}
