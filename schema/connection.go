package schema

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"
)

type connection struct {
	id           int
	creation     time.Time
	lastGet      time.Time
	lastPing     time.Time
	closed       bool
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
	newConnection.closed = false
	newConnection.dbConnection = db

	runtime.SetFinalizer(newConnection, finalizer)
	return newConnection, nil
}

func (conn *connection) close() {
	if conn.dbConnection != nil && conn.closed == false {
		// close properly connection
		err := conn.dbConnection.Close()
		fmt.Println(err)
		conn.closed = true
	}

}

func finalizer(conn *connection) {
	if conn != nil && conn.dbConnection != nil && conn.closed == false {
		// close properly connection
		conn.dbConnection.Close()
		conn.closed = true
	}
}
