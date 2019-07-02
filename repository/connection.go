package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	unknownID       = -1
	maxPingAttempts = 5
)

type Queryer interface {
	Query(query string, args ...interface{}) (result *sql.Rows, err error)
}

type Execer interface {
	Exec(query string, args ...interface{}) (lastInsertedID int64, err error)
}

type QueryExecer interface {
	Queryer
	Execer
}

type DB struct {
	Conn    *sql.DB
	connStr string
	Logger  logrus.Logger
}

func New(connStr string, logger logrus.Logger) (*DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %s", err)
	}
	if conn == nil {
		return nil, errors.New("connection failed: unknown error")
	}

	db := &DB{
		Conn:    conn,
		connStr: connStr,
		Logger:  logger,
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping: %s", err)
	}

	return db, nil
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %s", err)
	}

	return rows, nil
}

func (db *DB) Exec(query string, args ...interface{}) (int64, error) {
	result, err := db.Conn.Exec(query, args...)
	if err != nil {
		return unknownID, fmt.Errorf("exec: %s", err)
	}
	if result == nil {
		return unknownID, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return unknownID, nil
	}
	return id, nil
}

func (db DB) Ping() error {
	return db.pingNTimes(maxPingAttempts)
}

func (db DB) Close() error {
	return db.Conn.Close()
}

func (db DB) pingNTimes(maxAttempts int) error {
	attempts := 0

	for {
		_, err := db.Conn.Exec("SELECT 1;")
		attempts++

		if err == nil {
			return nil
		}
		if attempts >= maxAttempts {
			return fmt.Errorf("giving up after %d retries: %s", attempts, err)
		}
		db.Logger.Infof("database ping fail count %d: %s", attempts, err)
	}
}
