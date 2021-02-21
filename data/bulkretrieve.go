package data

import (
	"errors"
	"fmt"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
)

type BulkRetrieve struct {
	data            *[]*bulkRetrieveQuery
	currentSchema   *schema.Schema
	currentLanguage *schema.Language
}

const errorInvalidIndex = "This BulkRetrieve does not have a level #%d to retrieve results for."
const errorIndexAlreadyExist = "ThisIndex %d already exist."
const errorUnknownSchema = "Unknown schema."
const errorInvalidObject = "Object type '%s' is not valid."
const initialSliceCount = 4

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

func (bulkRetrieve *BulkRetrieve) SimpleQuery(entryIndex int, objectName string) error {
	if bulkRetrieve.data == nil {
		data := make([]*bulkRetrieveQuery, 0, initialSliceCount)
		bulkRetrieve.data = &data
		// get default schema
		bulkRetrieve.currentSchema = schema.GetDefaultSchema()
	}
	queryCount := len(*bulkRetrieve.data)
	if entryIndex > queryCount {
		return errors.New(fmt.Sprintf(errorInvalidIndex, queryCount))
	}
	if entryIndex < queryCount {
		return errors.New(fmt.Sprintf(errorIndexAlreadyExist, queryCount))
	}
	if bulkRetrieve.currentSchema == nil {
		return errors.New(errorUnknownSchema)
	}
	var table = bulkRetrieve.currentSchema.GetTableByName(objectName)
	if table == nil {
		return errors.New(fmt.Sprintf(errorInvalidObject, objectName))
	}
	*bulkRetrieve.data = append(*bulkRetrieve.data, newSimpleQuery(table))
	return nil
}

func (bulkRetrieve *BulkRetrieve) AppendFilter(entryIndex int, fieldName string, operation operationtype.OperationType,
	operand interface{}) error {
	queryCount := len(*bulkRetrieve.data)

	if entryIndex < 0 || entryIndex >= queryCount {
		return errors.New(fmt.Sprintf(errorInvalidIndex, queryCount))
	}
	var field = (*bulkRetrieve.data)[entryIndex].targetObject.GetFieldByName(fieldName)
	if field == nil {
		return errors.New(fmt.Sprintf(errorUnknownFieldName, fieldName, (*bulkRetrieve.data)[entryIndex].targetObject.GetName()))
	}
	//TODO type validations
	filter, err := newQueryFilter(field, operation, operand)
	if err == nil {
		(*bulkRetrieve.data)[entryIndex].addFilter(filter)
	}
	return err
}

func (bulkRetrieve *BulkRetrieve) AppendSort(entryIndex int, fieldName string, sortType sortordertype.SortOrderType) error {
	queryCount := len(*bulkRetrieve.data)

	if entryIndex < 0 || entryIndex >= queryCount {
		return errors.New(fmt.Sprintf(errorInvalidIndex, queryCount))
	}
	var field = (*bulkRetrieve.data)[entryIndex].targetObject.GetFieldByName(fieldName)
	if field == nil {
		return errors.New(fmt.Sprintf(errorUnknownFieldName, fieldName, (*bulkRetrieve.data)[entryIndex].targetObject.GetName()))
	}

	// check if field is not already sorted
	return nil
}

func (bulkRetrieve *BulkRetrieve) RetrieveRecords() error {

	return nil
}