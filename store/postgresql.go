package store

import (
	"database/sql"
	"github.com/lib/pq"
	"mithril/log"
	"mithril/message"
	"sync"
)

type psql struct {
	sync.Mutex
	db   *sql.DB
	stmt *sql.Stmt
	conn string
}

func init() {
	register("postgresql", &psql{})
}

func (me *psql) Init(url string) error {
	var err error

	if me.conn, err = pq.ParseURL(url); err != nil {
		return err
	}

	if err = me.establishConnection(); err != nil {
		return err
	}

	if err = NewPGSchemaEnsurer(me.db).EnsureSchema(); err != nil {
		return err
	}

	return me.PrepareStatement()
}

func (me *psql) UriFormat() string {
	return "postgres://username:password@hostname:port/database?sslmode=<verify-full|require|disable>"
}

func (me *psql) Store(msg *message.Message) error {
	var err error

	if err = me.establishConnection(); err == nil {
		_, err = me.stmt.Exec(
			msg.MessageId,
			msg.CorrelationId,
			msg.Timestamp,
			msg.AppId,
			msg.ContentType,
			msg.Exchange,
			msg.RoutingKey,
			msg.Mandatory,
			msg.Immediate,
			msg.BodyBytes)
	}

	if err != nil {
		log.Printf("Failed to store message: %+v", err)
		me.stmt = nil
		me.db = nil
	}
	return err
}

func (me *psql) PrepareStatement() error {
	var err error
	me.stmt, err = me.db.Prepare(`INSERT INTO mithril_requests (
								  message_id,
								  correlation_id,
								  created_at,
								  app_id,
								  content_type,
								  exchange,
								  routing_key,
								  mandatory,
								  immediate,
								  body_bytes)
								  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`)
	return err
}

func (me *psql) establishConnection() error {
	var err error

	me.Lock()
	defer me.Unlock()
	if me.db == nil {
		if me.db, err = sql.Open("postgres", me.conn); err != nil {
			return err
		}
	}

	return me.PrepareStatement()
}
