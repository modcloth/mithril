package log

type nullLogger struct{}

func (me *nullLogger) Print(v ...interface{}) {
}
func (me *nullLogger) Printf(format string, v ...interface{}) {
}
func (me *nullLogger) Println(v ...interface{}) {
}
func (me *nullLogger) Fatal(v ...interface{}) {
}
func (me *nullLogger) Fatalf(format string, v ...interface{}) {
}
