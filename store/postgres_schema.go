package store

import (
	"database/sql"

	"github.com/modcloth/mithril/log"
)

type pgSchemaEnsurer struct {
	db *sql.DB
}

var migrations = map[string][]string{
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
	"20130725000000": {
		`ALTER TABLE mithril_requests ADD COLUMN correlation_id character varying(128)`,
	},
}

func newPGSchemaEnsurer(db *sql.DB) *pgSchemaEnsurer {
	return &pgSchemaEnsurer{db}
}

func (me *pgSchemaEnsurer) EnsureSchema() error {
	if err := me.ensureMigrationsTable(); err != nil {
		return err
	}
	return me.runMigrations()
}

func (me *pgSchemaEnsurer) ensureMigrationsTable() error {
	_, err := me.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version character varying(255) NOT NULL);`)
	return err
}

func (me *pgSchemaEnsurer) runMigrations() error {
	for schemaVersion, sqls := range migrations {
		if me.containsMigration(schemaVersion) {
			continue
		}

		log.Printf("pg - Executing postgresql migration %s\n", schemaVersion)
		if err := me.migrateTo(schemaVersion, sqls); err != nil {
			return err
		}
	}
	return nil
}

func (me *pgSchemaEnsurer) containsMigration(schemaVersion string) bool {
	var count int
	if err := me.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", schemaVersion).Scan(&count); err != nil {
		return false
	}

	return count == 1
}

func (me *pgSchemaEnsurer) migrateTo(schemaVersion string, sqls []string) error {
	var (
		tx  *sql.Tx
		err error
	)

	if tx, err = me.db.Begin(); err != nil {
		return err
	}

	for _, sql := range sqls {
		if _, err = tx.Exec(sql); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err = tx.Exec("INSERT INTO schema_migrations VALUES ($1)", schemaVersion); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
