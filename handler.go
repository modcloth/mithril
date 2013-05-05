package mithril

type Handler interface {
	HandleRequest(*FancyRequest) error
	Init() error
	SetNextHandler(Handler)
}
