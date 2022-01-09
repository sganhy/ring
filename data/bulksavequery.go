package data

import (
	"database/sql"
	"ring/data/bulksavetype"
	"ring/schema"
	"ring/schema/dmlstatement"
)

type bulkSaveQuery struct {
	record       *Record
	targetObject *schema.Table
	bulkSaveType bulksavetype.BulkSaveType
	parameters   []interface{}
	fields       []*schema.Field
}

func (bulkSQ *bulkSaveQuery) Init(record *Record, bulkSaveType bulksavetype.BulkSaveType) {
	bulkSQ.record = record
	bulkSQ.targetObject = record.recordType
	bulkSQ.bulkSaveType = bulkSaveType

	switch bulkSaveType {
	case bulksavetype.InsertRecord:
		bulkSQ.parameters = record.recordType.GetInsertParameters(bulkSQ.record.data)
		break
	case bulksavetype.DeleteRecord:
		bulkSQ.parameters = record.recordType.GetDeleteParameters(record.getField())
		break
	case bulksavetype.UpdateRecord:
		bulkSQ.fields = record.getUpdatedFields()
		bulkSQ.parameters = record.recordType.GetUpdateParameters(bulkSQ.fields, bulkSQ.record.data)
		break
	}
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************
func (query bulkSaveQuery) Execute(dbConnection *sql.DB, transaction *sql.Tx) error {
	// execute without transaction!
	var sqlQuery = query.targetObject.GetDml(query.getDmlStatement(), query.fields)
	var err error

	if transaction != nil {
		_, err = transaction.Exec(sqlQuery, query.parameters...)
	} else {
		var rows *sql.Rows
		rows, err = dbConnection.Query(sqlQuery, query.parameters...)
		//rows, err := dmlQuery.Execute(dbConnection, sqlQuery, parameters)
		if rows != nil {
			rows.Close() //WARN: don't forget rows.Close()
		}
	}

	if err != nil {
		var sch = schema.GetSchemaById(query.targetObject.GetSchemaId())
		//TODO cast query.parameters to string[]
		sch.LogQueryError(102, err, sqlQuery, nil)
	}

	//fmt.Println(err)
	return err
}
func (bulkSQ *bulkSaveQuery) Change(record *Record) {
	switch bulkSQ.bulkSaveType {
	case bulksavetype.InsertRecord:
		bulkSQ.parameters = record.recordType.GetInsertParameters(bulkSQ.record.data)
		break
	case bulksavetype.DeleteRecord:
		bulkSQ.parameters = record.recordType.GetDeleteParameters(record.getField())
		break
	case bulksavetype.UpdateRecord:
		bulkSQ.fields = record.getUpdatedFields()
		bulkSQ.parameters = record.recordType.GetUpdateParameters(bulkSQ.fields, bulkSQ.record.data)
		break
	}
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
