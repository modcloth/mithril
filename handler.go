package mithril

type Handler interface {
	HandleRequest(Request) error
	Init() error
	SetNextHandler(Handler)
}
