package schema

import (
	"errors"
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/searchabletype"
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

const defaultNumberValue = "0"
const defaultBooleanValue = "false"
const defaultDateTimeValue = "0001-01-01T00:00:00.000"
const defaultShortDateTimeValue = "0001-01-01"
const defaultLongDateTimeValue = "0001-01-01T00:00:00Z"
const maxInt08 = "127"
const maxInt16 = "32767"
const maxInt32 = "2147483647"
const maxInt64 = "9223372036854775807"
const minInt08 = "-128"
const minInt16 = "-32768"
const minInt32 = "-2147483648"
const minInt64 = "-9223372036854775808"
const defaultTimeFormat = "2006-01-02T15:04:05.000" // rfc3339
const defaultShortTimeFormat = "2006-01-02"         // rfc3339
const unknownFieldDataType = ""

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

const primaryKeyFieldName = "id"
const primaryKeyDesc = "Internal record number"
const errorInvalidValueType = "Invalid value type"
const errorInvalidDateTimeFormat = "Invalid Date/Time format"

// max length for a varchar
const postgreVarcharMaxSize = 65535
const mySqlVarcharMaxSize = 65535
const sqliteVarcharMaxSize = 1000000000
const fieldToStringFormat = "name=%s; description=%s; type=%s; defaultValue=%s; baseline=%t; notNull=%t; caseSensitive=%t; active=%t"

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
	//64
	defaultPrimaryKeyInt64 = new(Field)
	defaultPrimaryKeyInt64.Init(0, primaryKeyFieldName, primaryKeyDesc, fieldtype.Long, 0, "", false, true, true, false, true)
	//32
	defaultPrimaryKeyInt32 = new(Field)
	defaultPrimaryKeyInt32.Init(0, primaryKeyFieldName, primaryKeyDesc, fieldtype.Int, 0, "", false, true, true, false, true)
	//16
	defaultPrimaryKeyInt16 = new(Field)
	defaultPrimaryKeyInt16.Init(0, primaryKeyFieldName, primaryKeyDesc, fieldtype.Short, 0, "", false, true, true, false, true)

}

func (field *Field) Init(id int32, name string, description string, fieldType fieldtype.FieldType, size uint32,
	defaultValue string, baseline bool, notNull bool, caseSensitive bool, multilingual bool, active bool) {
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
	field.caseSensitive = caseSensitive
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

func (field *Field) GetSize() uint32 {
	return uint32(field.size)
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
	return field.fieldType == fieldtype.DateTime || field.fieldType == fieldtype.LongDateTime || field.fieldType == fieldtype.ShortDateTime
}

///
/// Calculate searchable field value (remove diacritic characters and value.ToUpper())
///
func (field *Field) GetSearchableValue(value string, searchableType searchabletype.SearchableType) string {
	//TODO specific treatment by language
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, value)
	//output = fmt.Sprint(output)
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
	result.name = field.name // max length 30 !! must be validated before
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

func (field *Field) GetDdl(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) string {
	datatype := field.getSqlDataType(provider)
	if datatype == unknownFieldDataType {
		return unknownFieldDataType
	}
	return strings.TrimSpace(field.GetPhysicalName(provider) + " " + field.getSqlDataType(provider) + " " +
		field.getSqlConstraint(provider, tableType))
}

// reformat value for records
func (field *Field) GetValue(value string) (string, error) {
	switch field.fieldType {
	case fieldtype.String:
		if field.size > 0 && len(value) > int(field.size) {
			//truncate or error
			return value[0:field.size], nil
		}
		return value, nil
	case fieldtype.Long, fieldtype.Int, fieldtype.Short, fieldtype.Byte:
		if isValidInteger(value, field.fieldType) == true {
			return value, nil
		}
		break
	case fieldtype.DateTime, fieldtype.ShortDateTime, fieldtype.LongDateTime:
		// must support iso-8601
		t, e := getDateTimeIso8601(value)
		if e == nil {
			return field.GetDateTimeString(*t), nil
		}
		return value, e
	default:
		if field.IsValueValid(value) {
			return value, nil
		}
		break
	}
	return value, errors.New(errorInvalidValueType)
}

func (field *Field) IsValueValid(value string) bool {
	// nullable field?
	if field.notNull == false && value == "" {
		return true
	}
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
		_, err := getDateTimeIso8601(value)
		return err == nil
	case fieldtype.Boolean:
		return strings.ToLower(value) == "true" || strings.ToLower(value) == "false"
	}
	return false
}

func (field *Field) Clone() *Field {
	newField := new(Field)
	newField.Init(field.id, field.name, field.description, field.fieldType, uint32(field.size), field.defaultValue, field.baseline,
		field.notNull, field.caseSensitive, field.multilingual, field.active)
	return newField
}

func (field *Field) String() string {
	var fieldTyp = field.fieldType.String()
	if field.fieldType == fieldtype.String && field.size > 0 {
		fieldTyp += fmt.Sprintf("(%d)", field.size)
	}
	return fmt.Sprintf(fieldToStringFormat, field.name, field.description, fieldTyp, field.defaultValue, field.baseline,
		field.notNull, field.caseSensitive, field.active)
}

func (field *Field) Equal(fieldB *Field) bool {
	if field == nil && fieldB == nil {
		return true
	}
	if (field == nil && fieldB != nil) || (field != nil && fieldB == nil) {
		return false
	}
	if field.String() == fieldB.String() {
		return true
	}
	return false
}

func (field *Field) GetDateTimeString(t time.Time) string {
	switch field.fieldType {
	case fieldtype.DateTime:
		return t.UTC().Format(defaultTimeFormat)
	case fieldtype.ShortDateTime:
		return t.Format(defaultShortTimeFormat)
	case fieldtype.LongDateTime:
		return t.Format(time.RFC3339Nano)
	case fieldtype.String:
		return t.UTC().Format(defaultTimeFormat)
	}
	return ""
}

// calculate
func (field *Field) GetParameterValue(value string) interface{} {
	switch field.fieldType {
	case fieldtype.Long, fieldtype.Int, fieldtype.Byte, fieldtype.Short:
		val, _ := strconv.ParseInt(value, 10, 64)
		return val
	case fieldtype.Double, fieldtype.Float:
		val, _ := strconv.ParseFloat(value, 64)
		return val
	case fieldtype.Boolean:
		if value == "true" {
			return true
		}
		return false
	case fieldtype.DateTime:
		val, _ := time.Parse(defaultTimeFormat, value)
		return val
	}
	return value
}

//******************************
// private methods
//******************************
func isValidInteger(value string, fieldType fieldtype.FieldType) bool {
	var sign, size = 0, len(value)
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
	switch fieldType {
	case fieldtype.Long:
		return int64Condition(value, size, sign)
	case fieldtype.Int:
		return int32Condition(value, size, sign)
	case fieldtype.Short:
		return int16Condition(value, size, sign)
	case fieldtype.Byte:
		return int08Condition(value, size, sign)
	}
	return false
}

func int08Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 3) ||
		(size == 3 && value <= maxInt08 && sign > 0) ||
		(size == 3 && sign == -1) ||
		(size == 4 && value <= minInt08 && sign == -1)
}

func int16Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 5) ||
		(size == 5 && value <= maxInt16 && sign > 0) ||
		(size == 5 && sign == -1) ||
		(size == 6 && value <= minInt16 && sign == -1)
}
func int32Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 10) ||
		(size == 10 && value <= maxInt32 && sign > 0) ||
		(size == 10 && sign == -1) ||
		(size == 11 && value <= minInt32 && sign == -1)
}
func int64Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 19) ||
		(size == 19 && value <= maxInt64 && sign > 0) ||
		(size == 19 && sign == -1) ||
		(size == 20 && value <= minInt64 && sign == -1)
}

func getDefaultValue(defaultValue string, field *Field) string {
	var newDefaultValue = ""
	if field.IsValueValid(defaultValue) == true {
		newDefaultValue = defaultValue
	}
	if newDefaultValue == "" {
		if field.notNull {

			if field.IsNumeric() {
				return defaultNumberValue
			}
			switch field.fieldType {
			case fieldtype.Boolean:
				return defaultBooleanValue
			case fieldtype.DateTime:
				return defaultDateTimeValue
			case fieldtype.ShortDateTime:
				return defaultShortDateTimeValue
			case fieldtype.LongDateTime:
				return defaultLongDateTimeValue
			}
		}
	}
	return newDefaultValue
}

func getDefaultPrimaryKey(fieldType fieldtype.FieldType) *Field {
	switch fieldType {
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
	var result = unknownFieldDataType
	switch provider {
	case databaseprovider.MySql:
		if val, ok := postgreDataType[field.fieldType]; ok {
			result = val
		}
		break
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
		break
	}
	return result
}

func (field *Field) GetPhysicalName(provider databaseprovider.DatabaseProvider) string {
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

func getDateTimeIso8601(value string) (*time.Time, error) {
	var t time.Time
	var e error

	//short format
	// YYYY-MM-DD
	if len(value) >= 10 {
		//currentTime := time.Now
		// 2019-10-12T14:20:50.52+07:00
		// get currentTimeZone
		t, e = time.Parse(time.RFC3339, value[0:10]+"T00:00:00Z")

		//no information about time zone, take the local one (can be cached)
		if e != nil {
			return nil, e
		}
		if strings.Contains(value, "Z") || strings.Contains(value, "-") || strings.Contains(value, "+") {
			_, offset := time.Now().Zone()
			t = t.Add(time.Second * time.Duration(offset*-1))
		}

		if len(value) >= 19 {
			// where is the timezone?
			// detect timezone
			hour, errH := strconv.Atoi(value[11:13])
			minute, errM := strconv.Atoi(value[14:16])
			second, errS := strconv.Atoi(value[17:19])
			if errH == nil && errM == nil && errS == nil {
				t = t.Add((time.Hour * time.Duration(hour)) + (time.Minute * time.Duration(minute)) +
					(time.Second * time.Duration(second)))
			}
		}
		return &t, e
	} else {
		return nil, errors.New(errorInvalidDateTimeFormat)
	}
}
