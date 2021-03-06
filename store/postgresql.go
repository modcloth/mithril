// +build !nil

package store

import (
	"database/sql"
	"sync"

	"github.com/modcloth/mithril/message"

	log "github.com/Sirupsen/logrus"
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
	log.Debug("pg - Parsing url")
	if me.conn, err = pq.ParseURL(url); err != nil {
		log.Errorf("The postgresql url specified could not be parsed.\nurl: %s\nerr: %s\n", url, err)
		return err
	}
	log.Debug("pg - Parsed url")

	if err = me.establishConnection(); err != nil {
		return err
	}

	log.Debug("pg - Verifying database schema and performing migrations if necessary")
	if err = newPGSchemaEnsurer(me.db).EnsureSchema(); err != nil {
		return err
	}
	log.Debug("pg - Schema verification complete")

	if err = me.PrepareStatement(); err != nil {
		return err
	}

	return nil
}

func (me *psql) UriFormat() string {
	return "postgres://username:password@hostname:port/database?sslmode=<verify-full|require|disable>"
}

func (me *psql) Store(msg *message.Message) (err error) {
	if me.db == nil {
		if err = me.restablishConnection(); err != nil {
			return err
		}
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
		log.Warnf("pg - Failed to store message: %+v", err)
	}
	return err
}

func (me *psql) PrepareStatement() (err error) {
	log.Debug("pg - Preparing mithril request statement")
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
		log.Errorf("pg - An error occurred while preparing the insert statement: %s\n", err)
	}
	return err
}

func (me *psql) establishConnection() (err error) {
	log.Info("pg - Establishing connection to postgresql server")
	me.db, err = sql.Open("postgres", me.conn)
	if err != nil {
		log.Errorf("pg - An error occurred while preparing the insert statement: %s\n", err)
		return err
	}

	log.Info("pg - Connection established")
	return nil
}

func (me *psql) restablishConnection() (err error) {
	log.Info("pg - Detected reconnection to postgresql server is needed")
	if err = me.establishConnection(); err != nil {
		return err
	}

	if err = me.PrepareStatement(); err != nil {
		return err
	}
	log.Info("pg - Reconnection successful")
	return nil
}
