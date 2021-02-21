package schema

import (
	"database/sql"
	"time"
)

type connection struct {
	id           int
	creation     time.Time
	lastGet      time.Time
	lastPing     time.Time
	dbConnection *sql.DB
}

func newConnection(id int, connectionString string, provider string) (*connection, error) {
	db, err := sql.Open(provider, connectionString)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	newConnection := new(connection)
	newConnection.id = id
	newConnection.creation = time.Now()
	newConnection.lastGet = time.Now()
	newConnection.lastPing = time.Now()
	newConnection.dbConnection = db
	return newConnection, nil
}

func (conn *connection) destroy() {
	conn.dbConnection.Close()
}
