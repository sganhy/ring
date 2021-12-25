package data

import (
	"database/sql"
	"ring/data/bulksavetype"
	"ring/schema/dmlstatement"
)

type bulkSaveQuery struct {
	currentRecord *Record
	bulkSaveType  bulksavetype.BulkSaveType
	cancelled     bool
}

func (bulkSQ *bulkSaveQuery) Init(record *Record, bulkSaveType bulksavetype.BulkSaveType) {
	bulkSQ.currentRecord = record
	bulkSQ.bulkSaveType = bulkSaveType
	bulkSQ.cancelled = false
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************
func (query bulkSaveQuery) Execute(dbConnection *sql.DB) error {
	// execute without transaction!
	//var provider = query.targetObject.GetDatabaseProvider()
	return nil
}

//******************************
// private methods
//******************************
func (query *bulkSaveQuery) getDmlStatement(dbConnection *sql.DB) dmlstatement.DmlStatement {
	// execute without transaction!
	switch query.bulkSaveType {
	case bulksavetype.InsertMtm, bulksavetype.InsertRecord, bulksavetype.InsertMtmIfNotExist:
	case bulksavetype.UpdateRecord:
	}
	return dmlstatement.Undefined
}
