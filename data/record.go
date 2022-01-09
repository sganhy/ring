package data

import (
	"errors"
	"fmt"
	"ring/schema"
	"ring/schema/fieldtype"
	"strconv"
	"strings"
	"time"
)

const (
	emptyField             string = ""
	errorInvalidObjectType string = "Object type '%s' is not valid."
	errorUnknownRecordType string = "This Record object has an unknown RecordType.  The RecordType property must be set before performing this operation."
	errorUnknownFieldName  string = "Field name '%s' does not exist for object type '%s'."
	errorUnknownRelName    string = "Relation name %s does not exist for object type %s."
	errorInvalidNumber     string = "Invalid '%s' value %s."
	recordIdNotDefined     int64  = -1
)

type Record struct {
	data         []string      // values from rows
	recordType   *schema.Table // table definition
	stateChanged *node         // store information value changes
}

//******************************
// getters / setters
//******************************
func (record *Record) setTable(table *schema.Table) {
	record.recordType = table
}
func (record *Record) getTable() *schema.Table {
	return record.recordType
}

//******************************
// public methods
//******************************
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
			return record.setFieldByIndex(fieldId, value)
		}
		return errors.New(fmt.Sprintf(errorUnknownFieldName, name, record.recordType.GetName()))
	}
	return errors.New(errorUnknownRecordType)
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

func (record Record) IsDirty() bool {
	return record.stateChanged != nil
}

func (record *Record) String() string {
	if record.recordType == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(record.recordType.GetName())
	sb.WriteString(" object:\n")
	for i := 0; i < record.recordType.GetFieldCount(); i++ {
		var currentField = record.recordType.GetFieldIdByIndex(i)

		// field indent
		sb.WriteString("   ")
		sb.WriteString(currentField.GetName())
		sb.WriteString(" ")
		// display by type
		if currentField.GetType() == fieldtype.String {
			sb.WriteString("'")
			sb.WriteString(record.GetField(currentField.GetName()))
			sb.WriteString("'")
		} else {
			sb.WriteString(record.GetField(currentField.GetName()))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

//******************************
// private methods
//******************************
func (record *Record) setField(id int64) {
	if record.recordType != nil {
		var index = record.recordType.GetPrimaryKeyIndex()
		if index >= 0 {
			record.data[index] = strconv.FormatInt(id, 10)
		}
	}
}

func (record *Record) getField() int64 {
	if record.recordType != nil {
		var index = record.recordType.GetPrimaryKeyIndex()
		if index >= 0 {
			id, err := strconv.ParseInt(record.data[index], 10, 64)
			if err == nil {
				return id
			}
		}
	}
	return recordIdNotDefined
}

func (record *Record) setFieldByIndex(index int, value interface{}) error {
	var field = record.recordType.GetFieldByIndex(index)
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
	case uint:
		val = strconv.FormatUint(uint64(value.(uint)), 10)
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
			record.data[index] = val
			return nil
		}
		break
	default:
		return errors.New("Unsupported type.")
	}
	var err error
	val, err = field.GetValue(val)
	if err == nil {
		if val != record.data[index] {
			record.data[index] = val
			record.updateState(index)
		}
	} else {
		var fieldType = field.GetType()
		return errors.New(fmt.Sprintf(errorInvalidNumber, fieldType.String(), value))
	}
	return nil
}

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
			}
		}
	}
	record.recordType = recordType
}

func (record *Record) getUpdatedFields() []*schema.Field {
	if record.stateChanged != nil {
		count := record.stateChanged.CountSetBits()
		if count > 0 {
			result := make([]*schema.Field, count, count)
			for i := 0; i < record.recordType.GetFieldCount() && count >= 0; i++ {
				if record.stateChanged.GetValue(uint8(i)) == true {
					count--
					result[count] = record.recordType.GetFieldByIndex(i)
				}
			}
			return result
		}
	}
	return nil
}

func (record *Record) updateState(index int) {
	if record.stateChanged == nil {
		record.stateChanged = new(node)
	}
	if index == record.recordType.GetPrimaryKeyIndex() {
		//record.stateChanged.ResetAll(uint8(record.recordType.GetFieldCount()), true)
		//do nothing
	} else {
		// change state
		record.stateChanged.SetValue(uint8(index), true)
	}
}
