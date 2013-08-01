package log

import "fmt"

type nullLogger struct{}

func (me *nullLogger) Print(v ...interface{}) {
}
func (me *nullLogger) Printf(format string, v ...interface{}) {
}
func (me *nullLogger) Println(v ...interface{}) {
}
func (me *nullLogger) Fatal(v ...interface{}) {
	panic(v)
}
func (me *nullLogger) Fatalf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v))
}
