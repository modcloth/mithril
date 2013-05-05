package mithril

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

var (
	dbIsNil = fmt.Errorf("PostgreSQL handler database is nil!")
)

type PostgreSQLHandler struct {
	connUri     string
	db          *sql.DB
	nextHandler Handler
}

func NewPostgreSQLHandler(connUri string, next Handler) *PostgreSQLHandler {
	pgHandler := &PostgreSQLHandler{connUri: connUri}
	pgHandler.SetNextHandler(next)
	return pgHandler
}

func (me *PostgreSQLHandler) SetNextHandler(handler Handler) {
	me.nextHandler = handler
}

func (me *PostgreSQLHandler) Init() error {
	var err error

	if err = me.ensureConnected(); err != nil {
		return err
	}

	if err = me.ensureSchemaPresent(); err != nil {
		return err
	}

	log.Println("PostgreSQL handler initialized")

	if me.nextHandler != nil {
		return me.nextHandler.Init()
	}

	return nil
}

func (me *PostgreSQLHandler) HandleRequest(req *FancyRequest) error {
	if err := me.insertRequest(req); err != nil {
		return err
	}

	if me.nextHandler != nil {
		return me.nextHandler.HandleRequest(req)
	}

	return nil
}

func (me *PostgreSQLHandler) insertRequest(req *FancyRequest) error {
	if me.db == nil {
		return dbIsNil
	}

	r, err := me.db.Exec(`
		INSERT INTO mithril_requests (
			message_id,
			created_at,
			app_id,
			content_type,
			exchange,
			routing_key,
			mandatory,
			immediate,
			body_bytes
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`,
		req.MessageId,
		req.Timestamp,
		req.AppId,
		req.ContentType,
		req.Exchange,
		req.RoutingKey,
		req.Mandatory,
		req.Immediate,
		req.BodyBytes,
	)

	log.Printf("Insert returned result=%+v, err=%+v", r, err)
	return err
}

func (me *PostgreSQLHandler) ensureConnected() error {
	if me.isConnected() {
		return nil
	}

	return me.establishConnection()
}

func (me *PostgreSQLHandler) isConnected() bool {
	if me.db == nil {
		return false
	}

	if me.selectNow() != nil {
		return false
	}

	return true
}

func (me *PostgreSQLHandler) selectNow() error {
	_, err := me.db.Exec(`SELECT now() "mithril ping test";`)

	if err != nil {
		log.Println("PostgreSQL failed to execute 'SELECT now()':", err)
		log.Println("Is PostgreSQL running?")
	}

	return err
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
	return me.selectNow()
}

func (me *PostgreSQLHandler) ensureSchemaPresent() error {
	// TODO this should probably delegate to some kind
	// of schema-checking thingydoo, maybe even with fancy
	// pants migration gadgetry.

	if me.db == nil {
		return dbIsNil
	}

	var (
		r   sql.Result
		err error
	)

	r, err = me.db.Exec(`
		CREATE TABLE IF NOT EXISTS mithril_requests (
			id serial PRIMARY KEY,
			message_id character varying(128) NOT NULL,
			created_at timestamp without time zone NOT NULL,
			app_id character varying(128) NOT NULL,
			content_type character varying(64) NOT NULL,
			exchange character varying(256) NOT NULL,
			routing_key character varying(256) NOT NULL,
			mandatory boolean NOT NULL,
			immediate boolean NOT NULL,
			body_bytes text NOT NULL
		);
	`)

	log.Printf("Table ensuring query result=%+v, err=%+v", r, err)
	if err != nil {
		return err
	}

	// XXX kinda don't care if these fail, given the hackiness of the
	// situation.  See above comment regarding using a real thing to do real
	// things.
	me.db.Exec(`
		CREATE INDEX mr_app_ids
		ON mithril_requests (app_id);
	`)
	me.db.Exec(`
		CREATE INDEX mr_exchanges
		ON mithril_requests (exchange);
	`)
	me.db.Exec(`
		CREATE INDEX mr_routing_keys
		ON mithril_requests (routing_key);
	`)

	return nil
}
