package data

import (
	"errors"
	"fmt"
	"ring/schema"
	"strconv"
	"time"
)

const emptyField string = ""
const errorInvalidObjectType = "Object type '%s' is not valid."
const errorUnknownRecordType = "This Record object has an unknown RecordType.  The RecordType property must be set before performing this operation."
const errorUnknownFieldName = "Field name '%s' does not exist for object type '%s'."
const errorInvalidNumber = "Invalid '%s' value %s."

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

func (record *Record) SetField(name string, value interface{}) error {
	if record.recordType != nil {
		var fieldId = record.recordType.GetFieldIndexByName(name)
		if fieldId >= 0 {
			var field = record.recordType.GetFieldByIndex(fieldId)
			var val string
			switch value.(type) {
			case string:
				val = value.(string)
				break
			case float32:
				// conversion issues
				val = fmt.Sprintf("%g", value.(float32))
				break
			case float64:
				val = strconv.FormatFloat(value.(float64), 'f', -1, 64)
				break
			case int:
				val = strconv.Itoa(value.(int))
				break
			case bool:
				val = strconv.FormatBool(value.(bool))
				break
			case int8:
				val = strconv.Itoa(int(value.(int8)))
				break
			case int16:
				val = strconv.Itoa(int(value.(int16)))
				break
			case int32:
				val = strconv.Itoa(int(value.(int32)))
				break
			case int64:
				val = strconv.FormatInt(value.(int64), 10)
				break
			case uint8:
				val = strconv.FormatUint(uint64(value.(uint8)), 10)
				break
			case uint16:
				val = strconv.FormatUint(uint64(value.(uint16)), 10)
				break
			case uint32:
				val = strconv.FormatUint(uint64(value.(uint32)), 10)
				break
			case uint64:
				val = strconv.FormatUint(value.(uint64), 10)
				break
			case time.Time:
				val = field.GetDateTimeString(value.(time.Time))
				if field.IsDateTime() {
					// avoid dateTime revalidation
					record.data[fieldId] = val
					return nil
				}
				break
			default:
				return errors.New("Unsupported type")
			}
			var err error
			val, err = field.GetValue(val)
			if err == nil {
				record.data[fieldId] = val
			} else {
				var fieltyp = field.GetType()
				return errors.New(fmt.Sprintf(errorInvalidNumber, fieltyp.ToString(), val))
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

func (record *Record) Copy() *Record {
	var result = new(Record)
	if record.recordType != nil {
		result.recordType = record.recordType
		result.data = make([]string, len(record.data), cap(record.data))
		for i := 0; i < len(record.data); i++ {
			result.data[i] = record.data[i]
		}
	}
	return result
}

//******************************
// private methods
//******************************

func (record *Record) setRecordType(recordType *schema.Table) {
	// is it the same ?
	if recordType != nil {
		capacity := recordType.GetFieldCount()
		if capacity > 0 {
			// maye ba we can avoid to allocate if the capacity is smaller
			if capacity != len(record.data) {
				record.data = make([]string, capacity, capacity)
			} else {
				// reset all values if recordType changed?? is there
				// good idea may be
				//if record.recordType != recordType {
				for i := 0; i < len(record.data); i++ {
					record.data[i] = emptyField
				}
				//}
			}
		}
	}
	record.recordType = recordType
}
