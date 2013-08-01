package mithril

import (
	"flag"
	"fmt"
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
	pidFlag     = flag.String("p", "", "PID file (only written if provided)")
	versionFlag = flag.Bool("version", false, "Print version and exit")
	revFlag     = flag.Bool("rev", false, "Print git revision and exit")
	debug       = flag.Bool("d", false, "Enable Debug logging")
	showStorage = flag.Bool("s", false, "show the list of compiled in storage drivers.")
	storage     = flag.String("storage", "postgresql", "The storage type to presist messages to.")
	storageUri  = flag.String("storage.uri", "postgres://localhost/mithril_test?sslmode=disable", "The connection uri used by the storage engine.")
	addrFlag    = flag.String("a", ":8371", "Mithril server address")
	amqpUriFlag = flag.String("amqp.uri", "amqp://guest:guest@localhost:5672", "AMQP Server URI")
)

func NewConfigurationFromFlags() *Configuration {

	flag.Usage = func() {
		fmt.Println("Usage: mithril-server [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	return &Configuration{
		DisplayVersion: *versionFlag,
		DisplayRev:     *revFlag,
		PidFile:        *pidFlag,
		EnableDebug:    *debug,
		ServerAddress:  *addrFlag,
		Storage:        *storage,
		StorageUri:     *storageUri,
		AmqpUri:        *amqpUriFlag,
		ShowStorage:    *showStorage,
	}

}
