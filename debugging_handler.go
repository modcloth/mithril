// +build debug full

package mithril

type DebuggingHandler struct {
	nextHandler Handler
}

func NewDebuggingHandler(next Handler) *DebuggingHandler {
	debuggingHandler := &DebuggingHandler{}
	debuggingHandler.SetNextHandler(next)
	return debuggingHandler
}

func (me *DebuggingHandler) SetNextHandler(handler Handler) {
	me.nextHandler = handler
}

func (me *DebuggingHandler) Init() error {
	return nil
}

func (me *DebuggingHandler) HandleRequest(req *FancyRequest) error {
	Debugf("Handling request -> %+v\n", req)

	if me.nextHandler == nil {
		return nil
	}

	if result := me.nextHandler.HandleRequest(req); result != nil {
		Debugf("ERROR: %+v\n", result)
		return result
	}

	return nil
}
