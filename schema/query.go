package schema

import (
	"database/sql"
	"ring/schema/databaseprovider"
)

type Query interface {
	Execute(provider databaseprovider.DatabaseProvider, dbConnection *sql.DB) error
}
