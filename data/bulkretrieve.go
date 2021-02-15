package data

import (
	"errors"
	"fmt"
	"ring/schema"
)

type BulkRetrieve struct {
	data            []*bulkRetrieveQuery
	currentSchema   *schema.Schema
	currentLanguage *schema.Language
}

const errorInvalidIndex = "This BulkRetrieve does not have a level #{0} to retrieve results for."

func (bulkRetrieve *BulkRetrieve) setSchema(schema *schema.Schema) {
	bulkRetrieve.currentSchema = schema
	if schema != nil {
		bulkRetrieve.currentLanguage = schema.GetLanguage()
	}
}

func (bulkRetrieve *BulkRetrieve) SetSchema(schemaName string) error {
	var newSchema = schema.GetSchemaByName(schemaName)
	if newSchema != nil {
		bulkRetrieve.currentSchema = newSchema
		return nil
	}
	return errors.New(fmt.Sprintf("Invalid, name schema name: %s", schemaName))
}

func (bulkRetrieve *BulkRetrieve) setLanguage(language *schema.Language) {
	if language != nil {
		bulkRetrieve.currentLanguage = language
	}
}

func (bulkRetrieve *BulkRetrieve) SimpleQuery(entryIndex uint16, objectName string) error {
	return nil
}
