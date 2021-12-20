package data

import (
	"database/sql"
	"ring/schema"
)

type bulkSaveQuery struct {
	targetObject *schema.Table
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************
func (query bulkSaveQuery) Execute(dbConnection *sql.DB) error {
	return nil
}
