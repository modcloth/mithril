package log

import (
	"flag"
	"log"
	"os"
)

type Log interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

var logger Log

func init() {
	var debug bool
	flag.BoolVar(&debug, "d", false, "Enable Debugging handler")

	if debug {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		logger = &nullLogger{}
	}
}

func Print(v ...interface{}) {
	logger.Print(v)
}
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v)
}
func Println(v ...interface{}) {
	logger.Println(v)
}
func Fatal(v ...interface{}) {
	logger.Fatal(v)
}
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v)
}
