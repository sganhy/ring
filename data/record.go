package data

import (
	"errors"
	"fmt"
	"ring/schema"
	"ring/schema/fieldtype"
)

const unknowFieldDataType string = ""
const emptyField string = ""

const errorInvalidObjectType = "Object type '%s' is not valid."
const errorUnknownRecordType = "This Record object has an unknown RecordType.  The RecordType property must be set before performing this operation."
const errorUnknownFieldName = "Field name '%s' does not exist for object type '%s'."
const errorInvalidNumber = "Invalid %s value %s."

type Record struct {
	data       []string
	recordType *schema.Table
}

// schemaName.TableName
func (record *Record) SetRecordType(recordType string) error {
	record.recordType = schema.GetTableBySchemaName(recordType)
	if record.recordType != nil {
		record.setRecordType(record.recordType)
	} else {
		return errors.New(fmt.Sprintf(errorInvalidObjectType, recordType))
	}
	return nil
}

func (record *Record) setRecordType(recordType *schema.Table) {
	record.recordType = recordType
	capacity := recordType.GetFieldCount()
	if capacity > 0 {
		record.data = make([]string, capacity, capacity)
	}
}

func (record *Record) GetField(name string) string {
	if record.recordType != nil {
		var fieldId = record.recordType.GetFieldIndexByName(name)
		if fieldId >= 0 {
			var value = record.data[fieldId]
			if value != emptyField {
				return value
			} else {

				return record.recordType.GetFieldByIndex(fieldId).GetDefaultValue()
			}
		}
	}
	return emptyField
}

func (record *Record) SetField(name string, value string) error {
	if record.recordType != nil {
		var fieldId = record.recordType.GetFieldIndexByName(name)
		if fieldId >= 0 {
			var field = record.recordType.GetFieldByIndex(fieldId)
			if field.IsValueValid(value) {
				record.data[fieldId] = value
			} else {
				return errors.New(getMsgInvalidValue(value, field.GetType()))
			}
			return nil
		}
		return errors.New(fmt.Sprintf(errorUnknownFieldName, name, record.recordType.GetName()))
	}
	return errors.New(errorUnknownRecordType)
}

// set field without validation
func (record *Record) setField(name string, value string) {
	record.data[record.recordType.GetFieldIndexByName(name)] = value
}

//******************************
// private methods
//******************************

func getMsgInvalidValue(value string, fieldtyp fieldtype.FieldType) string {
	var msg = "Invalid value"
	switch fieldtyp {
	case fieldtype.Long:
		msg = fmt.Sprintf(errorInvalidNumber, "Long", value)
		break
	case fieldtype.Int:
		msg = fmt.Sprintf(errorInvalidNumber, "Int", value)
		break
	case fieldtype.Short:
		msg = fmt.Sprintf(errorInvalidNumber, "Short", value)
		break
	case fieldtype.Byte:
		msg = fmt.Sprintf(errorInvalidNumber, "Byte", value)
		break
	case fieldtype.Boolean:
		msg = fmt.Sprintf(errorInvalidNumber, "Boolean", value)
		break
	}
	return msg
}
