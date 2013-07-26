// +build pg full

package mithril

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/lib/pq"
)

var (
	enablePgFlag = flag.Bool("pg", false, "Enable PostgreSQL handler")
	pgUriFlag    = flag.String("pg.uri", "postgres://localhost?sslmode=disable", "PostgreSQL Server URI")

	dbIsNil      = fmt.Errorf("PostgreSQL handler database is nil!")
	pgMigrations = map[string][]string{
		"20120505000000": {`
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
		`,
			`CREATE INDEX mithril_app_ids ON mithril_requests (app_id);`,
			`CREATE INDEX mithril_exchanges ON mithril_requests (exchange);`,
			`CREATE INDEX mithril_routing_keys ON mithril_requests (routing_key);`,
		},
	}
)

func init() {
	pipelineCallbacks["pg"] = func(pipeline Handler) Handler {
		if *enablePgFlag {
			pipeline = NewPostgreSQLHandler(*pgUriFlag, pipeline)
		}
		return pipeline
	}
}

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

	_, err := me.db.Exec(`
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
	if me.db == nil {
		return dbIsNil
	}

	ensurer := newPgSchemaEnsurer(me.db, "mithril_schema_migrations")
	if err := ensurer.Init(); err != nil {
		return err
	}

	return ensurer.Migrate(pgMigrations)
}
