package data

import (
	"errors"
	"ring/data/bulksavetype"
	"ring/schema"
)

type BulkSave struct {
	data map[int32][]schema.Query
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (bulkSave *BulkSave) InsertRecord(record *Record) error {
	return bulkSave.addRecord(record, bulksavetype.InsertRecord)
}

func (bulkSave *BulkSave) UpdatRecord(record *Record) error {
	return bulkSave.addRecord(record, bulksavetype.UpdateRecord)
}

func (bulkSave *BulkSave) RelateRecords(sourceRecord *Record, targetRecord *Record, relationName string) {
}

func (bulkSave *BulkSave) Save() error {
	if bulkSave.data != nil {
		// without transaction
		for key, element := range bulkSave.data {
			if element != nil && len(element) > 0 {
				var sch = schema.GetSchemaById(key)
				sch.Execute(element, true)
			}
		}

	}
	return nil //bulkSave.currentSchema.Execute(bulkSave.data)
}

func (bulkSave *BulkSave) Clear() {
	if bulkSave.data != nil {
		// re-slicing
		for key, element := range bulkSave.data {
			if element != nil {
				if len(element) > 1000 {
					bulkSave.data[key] = nil
				} else {
					bulkSave.data[key] = bulkSave.data[key][:0]
				}
			}
		}
	}
}

//******************************
// private methods
//******************************
func (bulkSave *BulkSave) getQuery(record *Record, bulkSaveType bulksavetype.BulkSaveType) schema.Query {
	var query = bulkSaveQuery{}
	query.Init(record, bulkSaveType)
	var result schema.Query = query
	return result
}

func (bulkSave *BulkSave) addRecord(record *Record, bulkSaveType bulksavetype.BulkSaveType) error {
	if record == nil || record.recordType == nil {
		return errors.New(errorUnknownRecordType)
	}
	// allow object
	if bulkSave.data == nil {
		bulkSave.initializeData()
	}
	var schemaId = record.recordType.GetSchemaId()
	var data []schema.Query
	if _, ok := bulkSave.data[schemaId]; ok {
		data = bulkSave.data[schemaId]
	}
	if data == nil {
		data = make([]schema.Query, 0, initialSliceCount)
	}
	data = append(data, bulkSave.getQuery(record, bulkSaveType))
	bulkSave.data[schemaId] = data
	return nil
}

func (bulkSave *BulkSave) initializeData() {
	bulkSave.data = make(map[int32][]schema.Query, schema.GetSchemaCount()*2)
}
