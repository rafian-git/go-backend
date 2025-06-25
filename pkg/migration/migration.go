package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"gitlab.techetronventures.com/core/backend/pkg/sqlxdb"
)

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown Direction = "down"
)

// Check checks Direction values
func (m Direction) Check() (err error) {
	if m != DirectionUp && m != DirectionDown && m != "" {
		return fmt.Errorf("migration flag is not up or down: %s", m)
	}

	return
}

// SQLFromUrl creates sql.DB instance from url
func SQLFromUrl(url string) (*sql.DB, error) {
	cfg := &sqlxdb.Config{URL: url}
	db, err := sqlxdb.New(cfg)
	if err != nil {
		return nil, err
	}

	return db.DB.DB, nil
}

// MigrateFromFS performs migration from fs.FS source, wraps MigrateFromSource
func MigrateFromFS(db *sql.DB, direction Direction, database string,
	files fs.FS) (err error) {
	src, err := httpfs.New(http.FS(files), "migrations")
	if err != nil {
		return err
	}

	return MigrateFromSource(db, direction, database, "httpfs", src)
}

// MigrateFromSource performs migration from source.Driver
func MigrateFromSource(db *sql.DB, direction Direction, database string,
	source string, files source.Driver) (err error) {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithInstance(source, files, database, driver)
	if err != nil {
		return err
	}

	switch direction {
	case DirectionUp:
		err = m.Up()
	case DirectionDown:
		err = m.Down()
	}

	// don't emit errors on no changes made
	if err == migrate.ErrNoChange || err == migrate.ErrNilVersion {
		return nil
	}

	return err
}
