package data

import (
	"errors"
	"ring/data/bulksavetype"
	"ring/schema"
)

type BulkSave struct {
	data []*schema.Query
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (bulkSave *BulkSave) InsertRecord(record *Record) error {
	if record == nil || record.recordType == nil {
		return errors.New(errorUnknownRecordType)
	}
	// allow object
	if bulkSave.data == nil {
		data := make([]*schema.Query, 0, initialSliceCount)
		bulkSave.data = data
	}
	bulkSave.data = append(bulkSave.data, bulkSave.getQuery(record, bulksavetype.InsertRecord))
	return nil
}

func (bulkSave *BulkSave) RelateRecords(sourceRecord *Record, targetRecord *Record, relationName string) {
}

func (bulkSave *BulkSave) Save() error {
	return nil //bulkSave.currentSchema.Execute(bulkSave.data)
}

func (bulkSave *BulkSave) Clear() {
	if bulkSave.data != nil {
		if len(bulkSave.data) < 1000 {
			// re-slicing
			bulkSave.data = bulkSave.data[:0]
		} else {
			bulkSave.data = nil
		}
	}
}

//******************************
// private methods
//******************************
func (bulkSave *BulkSave) getQuery(record *Record, bulkSaveType bulksavetype.BulkSaveType) *schema.Query {
	var query = bulkSaveQuery{}

	query.Init(record, bulkSaveType)
	var result schema.Query = query

	return &result
}
