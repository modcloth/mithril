package mithril

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type PostgreSQLHandler struct {
	connUri string
	db      *sql.DB
}

func NewPostgreSQLHandler(connUri string) *PostgreSQLHandler {
	return &PostgreSQLHandler{connUri: connUri}
}

func (me *PostgreSQLHandler) HandleRequest(req Request) error {
	var err error

	if err = me.ensureConnected(); err != nil {
		return err
	}

	log.Println("PostgreSQLHandler not really handling request")
	if _, err = me.db.Query("SELECT now()"); err != nil {
		log.Println(err)
	}

	return nil
}

func (me *PostgreSQLHandler) ensureConnected() error {
	if me.isConnected() {
		return nil
	}

	return me.establishConnection()
}

func (me *PostgreSQLHandler) isConnected() bool {
	return me.db != nil
}

func (me *PostgreSQLHandler) establishConnection() error {
	var (
		connStr string
		err     error
		db      *sql.DB
	)

	if connStr, err = pq.ParseURL(me.connUri); err != nil {
		return err
	}

	if db, err = sql.Open("postgres", connStr); err != nil {
		return err
	}

	me.db = db
	return nil
}
