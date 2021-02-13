package schema

import (
	"errors"
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/tabletype"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var defaultPrimaryKeyInt64 *Field = nil
var defaultPrimaryKeyInt32 *Field = nil
var defaultPrimaryKeyInt16 *Field = nil
var defaultNumberValue = "0"
var defaultBooleanValue = "false"
var maxint08 string = "127"
var maxint16 string = "32767"
var maxint32 string = "2147483647"
var maxint64 string = "9223372036854775807"
var minint08 string = "-128"
var minint16 string = "-32768"
var minint32 string = "-2147483648"
var minint64 string = "-9223372036854775808"

const unknowFieldDataType string = ""

var postgreDataType = map[fieldtype.FieldType]string{
	fieldtype.String:        "varchar",
	fieldtype.Double:        "float8",
	fieldtype.Float:         "float4",
	fieldtype.Long:          "int8",
	fieldtype.Int:           "int4",
	fieldtype.Short:         "int2",
	fieldtype.Byte:          "int2",
	fieldtype.Boolean:       "bool",
	fieldtype.ShortDateTime: "date",
	fieldtype.DateTime:      "timestamp without time zone",
	fieldtype.LongDateTime:  "timestamp with time zone"}

const primaryKeyFielName = "id"
const primaryKeyDesc = "Internal record number"
const fieldFlags = "id"

// max lenght for a varchar
const longTextDefaultSize = 0
const postgreVarcharMaxSize = 65535
const mySqlVarcharMaxSize = 65535
const sqliteVarcharMaxSize = 1000000000

type Field struct {
	id            int32
	name          string
	description   string
	fieldType     fieldtype.FieldType
	size          uint16
	defaultValue  string
	baseline      bool
	notNull       bool
	caseSensitive bool
	multilingual  bool
	active        bool
}

func init() {

	// elemf3.Init(7, "Gga", "Gga", fieldtype.Double, 5, "", true, false, true, true)
	// --
	// id , name , description, fieldType fieldtype.FieldType, size uint32,
	// 		defaultValue string, baseline bool, notNull bool, multilingual bool, active bool
	// --
	// id, name, description, FieldType.Long, 0, null, false, true,true, true, false
	// --

	//64
	defaultPrimaryKeyInt64 = new(Field)
	defaultPrimaryKeyInt64.Init(0, primaryKeyFielName, primaryKeyDesc, fieldtype.Long, 0, "", false, true, true, false, true)
	//32
	defaultPrimaryKeyInt32 = new(Field)
	defaultPrimaryKeyInt32.Init(0, primaryKeyFielName, primaryKeyDesc, fieldtype.Int, 0, "", false, true, true, false, true)
	//16
	defaultPrimaryKeyInt16 = new(Field)
	defaultPrimaryKeyInt16.Init(0, primaryKeyFielName, primaryKeyDesc, fieldtype.Short, 0, "", false, true, true, false, true)

}

// call exemple elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true)
func (field *Field) Init(id int32, name string, description string, fieldType fieldtype.FieldType, size uint32,
	defaultValue string, baseline bool, notNull bool, casesensitive bool, multilingual bool, active bool) {
	field.id = id
	field.name = name
	field.description = description
	field.fieldType = fieldType
	if size > 65535 {
		field.size = 0
	} else {
		field.size = uint16(size)
	}

	field.baseline = baseline
	field.notNull = notNull
	field.active = active
	field.multilingual = multilingual
	field.caseSensitive = casesensitive
	//!!! at the end only
	field.defaultValue = getDefaultValue(defaultValue, field)
}

//******************************
// getters
//******************************
func (field *Field) GetId() int32 {
	return field.id
}

func (field *Field) GetName() string {
	return field.name
}

func (field *Field) GetDescription() string {
	return field.description
}

func (field *Field) GetType() fieldtype.FieldType {
	return field.fieldType
}

func (field *Field) GetSize() uint16 {
	return field.size
}

func (field *Field) GetDefaultValue() string {
	return field.defaultValue
}

func (field *Field) IsBaseline() bool {
	return field.baseline
}

func (field *Field) IsNotNull() bool {
	return field.notNull
}

func (field *Field) IsCaseSensitive() bool {
	return field.caseSensitive
}

func (field *Field) IsActive() bool {
	return field.active
}

func (field *Field) IsMultilingual() bool {
	return field.multilingual
}

//******************************
// public methods
//******************************
func (field *Field) IsValid() bool {
	// compare addresses
	if defaultPrimaryKeyInt64 == field || defaultPrimaryKeyInt32 == field || defaultPrimaryKeyInt16 == field {
		return true
	}
	return field.id > 0
}

func (field *Field) IsPrimaryKey() bool {
	// compare addresses
	return defaultPrimaryKeyInt64 == field || defaultPrimaryKeyInt32 == field || defaultPrimaryKeyInt16 == field
}

func (field *Field) IsNumeric() bool {
	return field.fieldType == fieldtype.Long || field.fieldType == fieldtype.Int ||
		field.fieldType == fieldtype.Short || field.fieldType == fieldtype.Byte ||
		field.fieldType == fieldtype.Float || field.fieldType == fieldtype.Double

}

func (field *Field) IsDateTime() bool {
	return field.fieldType == fieldtype.DateTime || field.fieldType == fieldtype.LongDateTime ||
		field.fieldType == fieldtype.ShortDateTime
}

///
/// Calculate searchable field value (remove diacritic characters and value.ToUpper())
///
func (field *Field) GetSearchableValue(value string, language Language) string {
	//TODO specific treatmen by language
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, value)

	// log here
	if err != nil {
		return strings.ToUpper(value)
	} else {
		return strings.ToUpper(output)
	}
}

func (field *Field) ToMeta(tableId int32) *Meta {
	// we cannot have error here
	var result = new(Meta)

	// key
	result.id = field.id
	result.refId = tableId
	result.objectType = int8(entitytype.Field)

	// others
	result.dataType = int32(field.fieldType)
	result.name = field.name // max lenght 30 !! must be valided before
	result.description = field.description
	result.value = field.defaultValue

	// flags
	result.flags = 0
	result.setFieldNotNull(field.notNull)
	result.setFieldCaseSensitive(field.caseSensitive)
	result.setFieldMultilingual(field.multilingual)
	result.setEntityBaseline(field.baseline)
	result.setFieldSize(uint32(field.size))

	result.enabled = field.active
	return result
}

func (field *Field) GetDdlSql(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) (string, error) {
	datatype := field.getSqlDataType(provider)
	if datatype == unknowFieldDataType {
		return unknowFieldDataType, errors.New(fmt.Sprintf("Unknow datatype {provider: %s, dataTypeId: %d, fieldName: %s}",
			provider, field.fieldType, field.name))
	}
	return strings.TrimSpace(field.getSqlFieldName(provider) + " " + field.getSqlDataType(provider) + " " +
		field.getSqlConstraint(provider, tableType)), nil
}

func (field *Field) IsValueValid(value string) bool {
	switch field.fieldType {
	case fieldtype.Long, fieldtype.Int, fieldtype.Short, fieldtype.Byte:
		return isValidInteger(value, field.fieldType)
	case fieldtype.Double:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case fieldtype.Float:
		_, err := strconv.ParseFloat(value, 32)
		return err == nil
	case fieldtype.String:
		return true
	case fieldtype.DateTime, fieldtype.LongDateTime, fieldtype.ShortDateTime:
		_, err := time.Parse("2006-01-02 15:04", "2011-01-19 22:15")
		return err == nil
	case fieldtype.Boolean:
		return strings.ToLower(value) == "true" || strings.ToLower(value) == "false"
	}
	return false
}

func (field *Field) Clone() *Field {
	newField := new(Field)
	/*
		id int32, name string, description string, fieldType fieldtype.FieldType, size uint32,
		defaultValue string, baseline bool, notNull bool, casesensitive bool, multilingual bool, active bool
	*/
	newField.Init(field.id, field.name, field.description, field.fieldType, uint32(field.size), field.defaultValue, field.baseline,
		field.notNull, field.caseSensitive, field.multilingual, field.active)
	return newField
}

//******************************
// private methods
//******************************
func isValidInteger(value string, fieldtyp fieldtype.FieldType) bool {
	var sign, size int = 0, len(value)
	for _, v := range value {
		if v >= '0' && v <= '9' {
			sign++
			continue
		} else if v == '-' && sign == 0 {
			sign -= size
			continue
		} else {
			return false
		}
	}
	// it's a digit
	switch fieldtyp {
	case fieldtype.Byte:
		return int08Condition(value, size, sign)
	case fieldtype.Short:
		return int16Condition(value, size, sign)
	case fieldtype.Int:
		return int32Condition(value, size, sign)
	case fieldtype.Long:
		return int64Condition(value, size, sign)
	}
	return false
}

func int08Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 3) ||
		(size == 3 && value <= maxint08 && sign > 0) ||
		(size == 3 && sign == -1) ||
		(size == 4 && value <= minint08 && sign == -1)
}

func int16Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 5) ||
		(size == 5 && value <= maxint16 && sign > 0) ||
		(size == 5 && sign == -1) ||
		(size == 6 && value <= minint16 && sign == -1)
}
func int32Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 10) ||
		(size == 10 && value <= maxint32 && sign > 0) ||
		(size == 10 && sign == -1) ||
		(size == 11 && value <= minint32 && sign == -1)
}
func int64Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 19) ||
		(size == 19 && value <= maxint64 && sign > 0) ||
		(size == 19 && sign == -1) ||
		(size == 20 && value <= minint64 && sign == -1)
}

func getDefaultValue(defaultValue string, field *Field) string {
	if defaultValue == "" {
		if field.IsNumeric() && field.notNull {
			return defaultNumberValue
		}
		if field.fieldType == fieldtype.Boolean && field.notNull {
			return defaultBooleanValue
		}
	}
	return defaultValue
}

func getDefaultPrimaryKey(fldtype fieldtype.FieldType) *Field {
	switch fldtype {
	case fieldtype.Int:
		return defaultPrimaryKeyInt32
	case fieldtype.Long:
		return defaultPrimaryKeyInt64
	case fieldtype.Short:
		return defaultPrimaryKeyInt16
	}
	return defaultPrimaryKeyInt64
}

func (field *Field) getSqlDataType(provider databaseprovider.DatabaseProvider) string {
	var result = unknowFieldDataType
	switch provider {
	case databaseprovider.MySql:
		if val, ok := postgreDataType[field.fieldType]; ok {
			result = val
		}
	case databaseprovider.Oracle:
	case databaseprovider.PostgreSql:
		if val, ok := postgreDataType[field.fieldType]; ok {
			result = val
			if field.fieldType == fieldtype.String {
				if field.size <= 0 || field.size > postgreVarcharMaxSize {
					result = "text"
				} else {
					// long text in postgresql
					result += fmt.Sprintf("(%d)", field.size)
				}
			}
		}
	case databaseprovider.Sqlite3:
	case databaseprovider.Influx:
	}
	return result
}

func (field *Field) getSqlFieldName(provider databaseprovider.DatabaseProvider) string {
	return field.name
}

func (field *Field) getSqlConstraint(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) string {
	// postgresql
	var result = "NULL"
	if tableType != tabletype.Business && field.notNull == true {
		result = "NOT " + result
	}
	return result
}
