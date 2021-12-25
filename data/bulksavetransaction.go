package data

import (
	"database/sql"
	"ring/schema"
)

type bulkSaveTransaction struct {
	data []*schema.Query
}

func (bulkST *bulkSaveTransaction) Init(data []*schema.Query) {
	bulkST.data = data
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************
func (query bulkSaveTransaction) Execute(dbConnection *sql.DB) error {
	// execute without transaction!
	return nil
}
