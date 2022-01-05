package data

import (
	"errors"
	"ring/data/bulksavetype"
	"ring/schema"
)

type BulkSave struct {
	data        map[int32][]schema.Query
	insertCount int
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (bulkSave *BulkSave) InsertRecord(record *Record) error {
	bulkSave.insertCount++
	return bulkSave.addQuery(record, bulksavetype.InsertRecord)
}

func (bulkSave *BulkSave) UpdatRecord(record *Record) error {
	return bulkSave.addQuery(record, bulksavetype.UpdateRecord)
}

func (bulkSave *BulkSave) RelateRecords(sourceRecord *Record, targetRecord *Record, relationName string) {
}

func (bulkSave *BulkSave) DeleteRecord(record *Record) error {
	return bulkSave.addQuery(record, bulksavetype.DeleteRecord)
}

func (bulkSave *BulkSave) DeleteRecordById(recordType string, id int64) error {
	record := new(Record)
	err := record.SetRecordType(recordType)
	if err != nil {
		return err
	}
	record.setField(id)
	return bulkSave.addQuery(record, bulksavetype.DeleteRecord)
}

func (bulkSave *BulkSave) Save() error {
	if bulkSave.data != nil {
		// generate Ids
		if bulkSave.insertCount > 0 {
			bulkSave.loadObjectId()
		}
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
	bulkSave.insertCount = 0

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
func (bulkSave *BulkSave) addQuery(record *Record, bulkSaveType bulksavetype.BulkSaveType) error {
	if record == nil || record.recordType == nil {
		return errors.New(errorUnknownRecordType)
	}
	// allow object
	if bulkSave.data == nil {
		bulkSave.data = make(map[int32][]schema.Query, schema.GetSchemaCount()*2)
	}
	var schemaId = record.recordType.GetSchemaId()
	var data []schema.Query
	if _, ok := bulkSave.data[schemaId]; ok {
		data = bulkSave.data[schemaId]
	}
	if data == nil {
		data = make([]schema.Query, 0, initialSliceCount)
	}

	var query = bulkSaveQuery{}
	query.Init(record, bulkSaveType)
	bulkSave.data[schemaId] = append(data, query)

	return nil
}

func (bulkSave *BulkSave) loadObjectId() error {
	var dico map[int32]int64
	var bsQuery bulkSaveQuery
	var currentId int32 = 0
	var schem *schema.Schema
	var table *schema.Table
	var schemaIndex = 0

	// [schemaId, number of insert]
	dico = make(map[int32]int64, bulkSave.insertCount)

	// build dictionary
	for schemaId, operations := range bulkSave.data {
		schem = schema.GetSchemaById(schemaId)
		// generate dico
		for i := 0; i < len(operations); i++ {
			bsQuery = operations[i].(bulkSaveQuery)
			if bsQuery.bulkSaveType == bulksavetype.InsertRecord {
				currentId = bsQuery.targetObject.GetId()
				dico[currentId] = dico[currentId] + 1
			}
		}
		// generate id
		for tableId, count := range dico {
			table = schem.GetTableById(tableId)
			dico[tableId] = table.GetNewObjid(uint32(count))
		}
		// set bulksave queries
		for i := len(operations) - 1; i >= 0; i-- {
			bsQuery = operations[i].(bulkSaveQuery)
			if bsQuery.bulkSaveType == bulksavetype.InsertRecord {
				currentId = bsQuery.targetObject.GetId()
				bsQuery.parameters[0] = dico[currentId]
				//fmt.Printf("==> objid=%d \n", dico[currentId])
				dico[currentId] = dico[currentId] - 1
			}
		}
		schemaIndex++
		for key, _ := range dico {
			delete(dico, key)
		}
		/*
			// clear dictionary
			if schemaIndex < len(bulkSave.data) {

			}
		*/
	}

	return nil
}

func (bulkSave *BulkSave) clearData() {

}
