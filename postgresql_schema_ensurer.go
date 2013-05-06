package mithril

import (
	"database/sql"
	"fmt"
	"log"
)

type pgSchemaEnsurer struct {
	db          *sql.DB
	schemaTable string
}

func newPgSchemaEnsurer(db *sql.DB, schemaTable string) *pgSchemaEnsurer {
	return &pgSchemaEnsurer{
		db:          db,
		schemaTable: schemaTable,
	}
}

func (me *pgSchemaEnsurer) Init() error {
	r, err := me.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version character varying(255) NOT NULL
		);
    `, me.schemaTable))
	log.Printf("Ensuring schema versions table exists "+
		"results=%+v, error=%+v", r, err)
	return err
}

func (me *pgSchemaEnsurer) Migrate(migrations map[string][]string) error {
	for schemaVersion, sqls := range migrations {
		if me.containsMigration(schemaVersion) {
			continue
		}

		if err := me.migrateTo(schemaVersion, sqls); err != nil {
			return err
		}
	}
	return nil
}

func (me *pgSchemaEnsurer) containsMigration(schemaVersion string) bool {
	var count int

	q := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s WHERE version = $1`, me.schemaTable)
	if err := me.db.QueryRow(q, schemaVersion).Scan(&count); err != nil {
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
			defer tx.Rollback()
			return err
		}
	}

	q := fmt.Sprintf(`INSERT INTO %s VALUES ($1)`, me.schemaTable)
	if _, err = tx.Exec(q, schemaVersion); err != nil {

		defer tx.Rollback()
		return err
	}

	return tx.Commit()
}
