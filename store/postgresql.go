// +build pg full

package store

import (
	"database/sql"
	"mithril/log"
	"mithril/message"
	"sync"

	"github.com/lib/pq"
)

type psql struct {
	sync.RWMutex
	db   *sql.DB
	stmt *sql.Stmt
	conn string
}

func init() {
	register("postgresql", &psql{})
}

func (me *psql) Init(url string) (err error) {

	log.Println("pg - Parsing url")
	if me.conn, err = pq.ParseURL(url); err != nil {
		log.Printf("The postgresql url specified could not be parsed.\nurl: %s\nerr: %s\n", url, err)
		return err
	}
	log.Println("pg - Parsed url")

	log.Println("pg - Establishing connection to postgresql server")
	if err = me.establishConnection(); err != nil {
		return err
	}
	log.Println("pg - Connection established")

	log.Println("pg - Verifying database schema and performing migrations if necessary")
	if err = newPGSchemaEnsurer(me.db).EnsureSchema(); err != nil {
		return err
	}
	log.Println("pg - Schema verification complete")

	log.Println("pg - Preparing mithril request statement")
	if err = me.PrepareStatement(); err != nil {
		return err
	}

	return nil
}

func (me *psql) UriFormat() string {
	return "postgres://username:password@hostname:port/database?sslmode=<verify-full|require|disable>"
}

func (me *psql) Store(msg *message.Message) (err error) {
	if err = me.establishConnection(); err != nil {
		return err
	}
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

	if err != nil {
		me.db.Close()
		me.db = nil
		log.Printf("pg - Failed to store message: %+v", err)
	}
	return err
}

func (me *psql) PrepareStatement() (err error) {
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
	if err != nil {
		log.Printf("pg - An error occurred while preparing the insert statement: %s\n", err)
	}
	return err
}

func (me *psql) establishConnection() (err error) {
	if me.db != nil {
		return nil
	}

	me.db, err = sql.Open("postgres", me.conn)
	if err != nil {
		log.Printf("pg - An error occurred while preparing the insert statement: %s\n", err)
		return err
	}

	return nil
}
