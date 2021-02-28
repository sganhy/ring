package schema

import (
	"database/sql"
)

type Query interface {
	Execute(dbConnection *sql.DB) error
}
