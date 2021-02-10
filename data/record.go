package data

import (
	"errors"
	"fmt"
	"ring/schema"
	"ring/schema/fieldtype"
)

const unknowFieldDataType string = ""
const emptyField string = ""
const maxint32 string = "2147483647"
const minint32 string = "-2147483648"

const errorInvalidObjectType = "Object type '%s' is not valid."
const errorUnknownRecordType = "This Record object has an unknown RecordType.  The RecordType property must be set before performing this operation."
const errorUnknownFieldName = "Field name '%s' does not exist for object type '%s'."

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
		record.data = make([]string, 0, capacity)
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
			var fieldType = record.recordType.GetFieldTypeByIndex(fieldId)
			switch fieldType {
			case fieldtype.Int:
			case fieldtype.Long:
			case fieldtype.Short:
			case fieldtype.Byte:
				if isValidInteger(name, fieldType) {
					record.data[fieldId] = value
				}
				break
			case fieldtype.Double:
			case fieldtype.Float:
				break
			case fieldtype.String:
				record.data[fieldId] = value
				break
			}
			return nil
		}
		return errors.New(fmt.Sprintf(errorUnknownFieldName, name, record.recordType.GetName()))
	}
	return errors.New(errorUnknownRecordType)
}

//******************************
// private methods
//******************************

func IsValidInteger(value string, fieldtyp fieldtype.FieldType) bool {
	return isValidInteger(value, fieldtyp)
}

func isValidInteger(value string, fieldtyp fieldtype.FieldType) bool {
	var i = 0
	for _, v := range value {
		if v >= '0' && v <= '9' {
			i++
			continue
		} else if v == '-' && i == 0 {
			i -= len(value)
			continue
		} else {
			return false
		}
	}

	// it's a digit
	switch fieldtyp {
	case fieldtype.Int:
		return (len(value) > 0 && len(value) < 10) ||
			(len(value) == 10 && value <= maxint32 && i > 0) ||
			(len(value) == 10 && i == -1) ||
			(len(value) == 11 && value >= minint32 && i == -1)
	}
	return false
}
