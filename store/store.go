package store

import (
	"fmt"

	"github.com/modcloth/mithril/message"

	log "github.com/Sirupsen/logrus"
)

type Storage struct {
	driver Driver
	uri    string
}

var drivers = make(map[string]Driver)

func register(name string, driver Driver) {
	if driver == nil {
		panic("Cannot register a null storage driver.")
	}
	if _, dup := drivers[name]; dup {
		panic(fmt.Sprintf("The driver, %s, has already been registered", name))
	}
	drivers[name] = driver
}

func ShowStorage() {
	fmt.Println("Available Storage Drivers:")
	for k, v := range drivers {
		fmt.Printf("\t%s: %s\n", k, v.UriFormat())
	}
}

func Open(name string, uri string) (*Storage, error) {
	driver, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("Unknown storage driver %q, did you forget to build it?", name)
	}

	if name != "" {
		log.Infof("Persisting messages to: %s.\n", name)
	}
	if err := driver.Init(uri); err != nil {
		return nil, err
	}
	return &Storage{driver, uri}, nil
}

func (me *Storage) Store(message *message.Message) error {
	return me.driver.Store(message)
}
