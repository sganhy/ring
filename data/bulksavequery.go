package data

import (
	"database/sql"
	"fmt"
	"ring/data/bulksavetype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/dmlstatement"
)

type bulkSaveQuery struct {
	currentRecord *Record
	targetObject  *schema.Table
	bulkSaveType  bulksavetype.BulkSaveType
	cancelled     bool
}

var (
	dmlQuery *bulkQuery
)

func init() {
	dmlQuery = new(bulkQuery)
}

func (bulkSQ *bulkSaveQuery) Init(record *Record, bulkSaveType bulksavetype.BulkSaveType) {
	bulkSQ.currentRecord = record
	bulkSQ.targetObject = record.recordType
	bulkSQ.bulkSaveType = bulkSaveType
	bulkSQ.cancelled = false
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************
func (query bulkSaveQuery) Execute(dbConnection *sql.DB, transaction *sql.Tx) error {
	// execute without transaction!
	//var provider = query.targetObject.GetDatabaseProvider()
	var dmlStatement = query.getDmlStatement()
	//var parameters = query.getParameters(provider)
	var sqlQuery = query.targetObject.GetDml(dmlStatement, nil)

	//_, err := dmlQuery.Execute(dbConnection, sqlQuery, parameters)
	fmt.Println(sqlQuery)
	return nil
}

//******************************
// private methods
//******************************
func (query *bulkSaveQuery) getDmlStatement() dmlstatement.DmlStatement {
	// execute without transaction!
	switch query.bulkSaveType {
	case bulksavetype.InsertMtm, bulksavetype.InsertRecord, bulksavetype.InsertMtmIfNotExist:
		return dmlstatement.Insert
	case bulksavetype.UpdateRecord:
		return dmlstatement.Update
	case bulksavetype.DeleteRecord:
		return dmlstatement.Delete
	}
	return dmlstatement.Undefined
}

func (query *bulkSaveQuery) getParameters(provider databaseprovider.DatabaseProvider) []interface{} {
	var parameters []interface{}
	return parameters
}
