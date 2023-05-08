package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func ConnectDB() (*DB, error) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=SergioAn200313272006 dbname=crudDB sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *DB) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	return d.db.QueryRow(query, args...), nil
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}
