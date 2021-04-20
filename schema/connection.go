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

func (conn *connection) Init(id int, connectionString string, provider string) error {
	db, err := sql.Open(provider, connectionString)

	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	conn.id = id
	conn.creation = time.Now()
	conn.lastGet = time.Now()
	conn.lastPing = time.Now()
	conn.closed = false
	conn.dbConnection = db

	runtime.SetFinalizer(conn, finalizer)
	return nil
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
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
