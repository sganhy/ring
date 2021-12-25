package data

import (
	"errors"
	"fmt"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
)

type BulkRetrieve struct {
	data          []schema.Query
	currentSchema *schema.Schema
	language      *schema.Language
}

const (
	errorInvalidIndex      = "This BulkRetrieve does not have a level #%d to retrieve results for."
	errorInvalidPageSize   = "An invalid page size was detected (page size should be greater than 0)."
	errorInvalidPageNumber = "An invalid page number was detected (page number should be greater or equal than 1)."
	errorIndexAlreadyExist = "ThisIndex %d already exist."
	errorUnknownSchema     = "Unknown schema."
	errorInvalidObject     = "Object type '%s' is not valid."
	errorInvalidSchemaName = "Invalid schema name: %s."
	initialSliceCount      = 4
)

//******************************
// getters and setters
//******************************
func (bulkRetrieve *BulkRetrieve) SetSchema(schemaName string) {
	var newSchema = schema.GetSchemaByName(schemaName)
	if newSchema == nil {
		panic(fmt.Sprintf(errorInvalidSchemaName, schemaName))
	}
	bulkRetrieve.setSchema(newSchema)
}

func (bulkRetrieve *BulkRetrieve) SetLanguage(language *schema.Language) {
	bulkRetrieve.language = language
}

func (bulkRetrieve *BulkRetrieve) setSchema(schema *schema.Schema) {
	// schema ==> cannot be null
	bulkRetrieve.currentSchema = schema
	bulkRetrieve.SetLanguage(schema.GetLanguage())
}

//******************************
// public methods
//******************************
func (bulkRetrieve *BulkRetrieve) SimpleQuery(entryIndex int, objectName string) error {
	if bulkRetrieve.data == nil {
		bulkRetrieve.data = make([]schema.Query, 0, initialSliceCount)
		// get default schema
		bulkRetrieve.currentSchema = schema.GetDefaultSchema()
	}
	queryCount := len(bulkRetrieve.data)
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
	bulkRetrieve.data = append(bulkRetrieve.data, newSimpleQuery(table))
	return nil
}

func (bulkRetrieve *BulkRetrieve) AppendFilter(entryIndex int, fieldName string, operation operationtype.OperationType,
	operand interface{}) error {
	queryCount := len(bulkRetrieve.data)

	if entryIndex < 0 || entryIndex >= queryCount {
		return errors.New(fmt.Sprintf(errorInvalidIndex, queryCount))
	}

	// cast interface schema.Query
	var query = bulkRetrieve.data[entryIndex].(bulkRetrieveQuery)
	var field = query.targetObject.GetFieldByName(fieldName)

	if field == nil {
		return errors.New(fmt.Sprintf(errorUnknownFieldName, fieldName, query.targetObject.GetName()))
	}
	//TODO type validations
	item, err := newQueryFilter(field, operation, operand)
	if err == nil {
		query.addFilter(item)
	}
	return err
}

func (bulkRetrieve *BulkRetrieve) AppendSort(entryIndex int, fieldName string, sortType sortordertype.SortOrderType) error {
	queryCount := len(bulkRetrieve.data)

	if entryIndex < 0 || entryIndex >= queryCount {
		return errors.New(fmt.Sprintf(errorInvalidIndex, queryCount))
	}

	// cast interface schema.Query
	var query = bulkRetrieve.data[entryIndex].(bulkRetrieveQuery)
	var field = query.targetObject.GetFieldByName(fieldName)

	if field == nil {
		return errors.New(fmt.Sprintf(errorUnknownFieldName, fieldName, query.targetObject.GetName()))
	}
	sort := newQuerySort(field, sortType)
	query.addSort(sort)
	// check if field is not already sorted
	return nil
}

func (bulkRetrieve *BulkRetrieve) RetrieveRecords() error {
	return bulkRetrieve.currentSchema.Execute(bulkRetrieve.data, false)
}

func (bulkRetrieve *BulkRetrieve) GetRecordList(entryIndex int) List {
	queryCount := len(bulkRetrieve.data)
	if entryIndex >= queryCount {
		var emptyResult = List{}
		return emptyResult
	}
	var query = bulkRetrieve.data[entryIndex].(bulkRetrieveQuery)
	return *query.result
}

func (bulkRetrieve *BulkRetrieve) Clear() {
	if bulkRetrieve.data != nil {
		// re-slicing
		bulkRetrieve.data = bulkRetrieve.data[:0]
	}
}

func (bulkRetrieve *BulkRetrieve) SetPage(entryIndex int, pageNumber int, pageSize int) error {
	queryCount := len(bulkRetrieve.data)
	if entryIndex < 0 || entryIndex >= queryCount {
		return errors.New(fmt.Sprintf(errorInvalidIndex, queryCount))
	}
	if pageSize <= 0 {
		return errors.New(errorInvalidPageSize)
	}
	if pageNumber <= 0 {
		return errors.New(errorInvalidPageNumber)
	}
	return nil
}

func (bulkRetrieve *BulkRetrieve) SetRootById(objectName string, id int64) error {
	table := bulkRetrieve.currentSchema.GetTableByName(objectName)
	if table == nil {
		return errors.New(fmt.Sprintf(errorInvalidObject, objectName))
	}
	rcd := new(Record)
	rcd.setTable(table)
	rcd.setField(id)
	return bulkRetrieve.SetRoot(rcd)
}

func (bulkRetrieve *BulkRetrieve) SetRoot(rootRecord *Record) error {
	table := rootRecord.getTable()
	if table == nil {
		return errors.New(errorUnknownRecordType)
	}
	return nil
}
