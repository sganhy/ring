package data

import (
	"fmt"
	"ring/schema"
)

type BulkSave struct {
	data          *[]schema.Query
	currentSchema *schema.Schema
	language      *schema.Language
}

//******************************
// getters and setters
//******************************
func (bulkSave *BulkSave) SetSchema(schemaName string) {
	var newSchema = schema.GetSchemaByName(schemaName)
	if newSchema == nil {
		panic(fmt.Sprintf(errorInvalidSchemaName, schemaName))
	}
	bulkSave.setSchema(newSchema)
}

func (bulkSave *BulkSave) SetLanguage(language *schema.Language) {
	bulkSave.language = language
}

func (bulkSave *BulkSave) setSchema(schema *schema.Schema) {
	// schema ==> cannot be null
	bulkSave.currentSchema = schema
	bulkSave.SetLanguage(schema.GetLanguage())
}

//******************************
// public methods
//******************************
func (bulkSave *BulkSave) InsertRecord(record *Record) {
	// schema ==> cannot be null
}

func (bulkSave *BulkSave) RelateRecords(sourceRecord *Record, targetRecord *Record, relationName string) {
}

func (bulkSave *BulkSave) Save() error {
	return bulkSave.currentSchema.Execute(*bulkSave.data)
}

func (bulkSave *BulkSave) Clear() {
	if bulkSave.data != nil {
		// re-slicing
		*bulkSave.data = (*bulkSave.data)[:0]
	}
}

//******************************
// private methods
//******************************
