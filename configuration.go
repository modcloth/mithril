package mithril

import (
	"flag"
	"fmt"
	"os"
)

type Configuration struct {
	DisplayVersion bool
	DisplayRev     bool
	PidFile        string
	EnableDebug    bool
	ServerAddress  string
	Storage        string
	StorageUri     string
	AmqpUri        string
	ShowStorage    bool
}

var (
	pidFlag     = flag.String("p", "", "-p PID\tCreate a pid file. If the pid file already exits, the application will exit immediately.")
	debug       = flag.Bool("d", false, "-d\tEnable debug logging.")
	showStorage = flag.Bool("l", false, "-l\tList the available, compiled storage drivers.")
	storage     = flag.String("s", "", "-s DRIVER  Which storage driver to use.  Messages will not be persisted if unset.")
	storageUri  = flag.String("u", "", "-u URL\tThe url used by the storage driver.")
	revFlag     = flag.Bool("r", false, "-r\tPrint git revision and exit.")
	versionFlag = flag.Bool("v", false, "-v\tPrint version and exit.")
)

func NewConfigurationFromFlags() *Configuration {

	flag.Usage = func() {
		fmt.Println("Usage: mithril-server [options] <hosting-address> <amqp uri>")
		printOptions()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	return &Configuration{
		DisplayVersion: *versionFlag,
		DisplayRev:     *revFlag,
		PidFile:        *pidFlag,
		EnableDebug:    *debug,
		Storage:        *storage,
		StorageUri:     *storageUri,
		ShowStorage:    *showStorage,
		ServerAddress:  flag.Arg(0),
		AmqpUri:        flag.Arg(1),
	}

}

func printOptions() {
	fmt.Println("Options:")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Println(flag.Usage)
	})
}
