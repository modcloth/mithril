package log

import (
	stdlog "log"
	"os"
	"sync"
)

type Log interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type indirectLogger struct {
	Log
	sync.Mutex
}

var (
	logger Log
	mu     = new(sync.Mutex)
)

func Initialize(debug bool) {
	mu.Lock()
	defer mu.Unlock()
	logger = NewLogger(debug)
}

func NewLogger(debug bool) Log {
	if debug {
		return stdlog.New(os.Stderr, "[mithril] ", stdlog.LstdFlags)
	} else {
		return &nullLogger{}
	}
}

func Print(v ...interface{}) {
	logger.Print(v...)
}
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
func Println(v ...interface{}) {
	logger.Println(v...)
}
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}
