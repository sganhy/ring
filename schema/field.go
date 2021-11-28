package schema

import (
	"errors"
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/searchabletype"
	"ring/schema/sqlfmt"
	"ring/schema/tabletype"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

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

const (
	defaultNumberValue         string = "0"
	defaultBooleanValue        string = "false"
	defaultDateTimeValue       string = "0001-01-01T00:00:00.000"
	defaultShortDateTimeValue  string = "0001-01-01"
	defaultLongDateTimeValue   string = "0001-01-01T00:00:00Z"
	postgreSqlByteConstraint   string = " CHECK (%s >= '-128' AND %s <= 127)"
	maxInt08                   string = "127"
	maxInt16                   string = "32767"
	maxInt32                   string = "2147483647"
	maxInt64                   string = "9223372036854775807"
	minInt08                   string = "-128"
	minInt16                   string = "-32768"
	minInt32                   string = "-2147483648"
	minInt64                   string = "-9223372036854775808"
	defaultTimeFormat          string = "2006-01-02T15:04:05.000" // rfc3339
	defaultShortTimeFormat     string = "2006-01-02"              // rfc3339
	unknownFieldDataType       string = ""
	invalidValue               string = ""
	primaryKeyFieldName        string = "id"
	primaryKeyDesc             string = "Internal record number"
	errorInvalidValueType      string = "Invalid value type"
	errorInvalidDateTimeFormat string = "Invalid Date/Time format"
	postgreVarcharMaxSize      uint16 = 65535
	mySqlVarcharMaxSize        uint16 = 65535
	sqliteVarcharMaxSize       int64  = 1000000000
	searchableFieldPrefix      string = "s_"
	fieldToStringFormat        string = "id=%d; name=%s; description=%s; type=%s; defaultValue=%s; baseline=%t; notNull=%t; caseSensitive=%t; multilingual=%t; active=%t"
)

var (
	defaultPrimaryKeyInt64 *Field = nil
	defaultPrimaryKeyInt32 *Field = nil
	defaultPrimaryKeyInt16 *Field = nil
	postgreDataType               = map[fieldtype.FieldType]string{
		fieldtype.String:        "varchar",
		fieldtype.LongString:    "text",
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
	mysqlDataType = map[fieldtype.FieldType]string{
		fieldtype.String:        "VARCHAR",
		fieldtype.LongString:    "LONGTEXT",
		fieldtype.Double:        "DOUBLE",
		fieldtype.Float:         "FLOAT",
		fieldtype.Long:          "BIGINT(20)",
		fieldtype.Int:           "INT(11)",
		fieldtype.Short:         "SMALLINT(6)",
		fieldtype.Byte:          "TINYINT(4)",
		fieldtype.Boolean:       "BOOLEAN",
		fieldtype.ShortDateTime: "DATE",
		fieldtype.DateTime:      "TIMESTAMP",
		fieldtype.LongDateTime:  "TIMESTAMP"}
)

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
	field.defaultValue = field.getDefaultValue(defaultValue)
}

//******************************
// getters and setters
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

func (field *Field) GetEntityType() entitytype.EntityType {
	return entitytype.Field

}

func (field *Field) setType(fieldType fieldtype.FieldType) {
	field.fieldType = fieldType
}

func (field *Field) setName(name string) {
	field.name = name
}

func (field *Field) setSize(size uint16) {
	field.size = size
}

func (field *Field) setCaseSensitive(caseSensitive bool) {
	field.caseSensitive = caseSensitive
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

func (field *Field) GetDdl(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) string {
	datatype := field.getSqlDataType(provider)
	if datatype == unknownFieldDataType {
		return unknownFieldDataType
	}
	return strings.TrimSpace(field.GetPhysicalName(provider) + " " + field.getSqlDataType(provider))
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
		if field.isValidInteger(value) == true {
			return value, nil
		}
		break
	case fieldtype.DateTime, fieldtype.ShortDateTime, fieldtype.LongDateTime:
		// must support iso-8601
		t, e := field.getDateTimeIso8601(value)
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
	return invalidValue, errors.New(errorInvalidValueType)
}

func (field *Field) IsValueValid(value string) bool {
	// nullable field?
	if field.notNull == false && value == "" {
		return true
	}
	switch field.fieldType {
	case fieldtype.Long, fieldtype.Int, fieldtype.Short, fieldtype.Byte:
		return field.isValidInteger(value)
	case fieldtype.Double:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case fieldtype.Float:
		_, err := strconv.ParseFloat(value, 32)
		return err == nil
	case fieldtype.String:
		return true
	case fieldtype.DateTime, fieldtype.LongDateTime, fieldtype.ShortDateTime:
		_, err := field.getDateTimeIso8601(value)
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
	return fmt.Sprintf(fieldToStringFormat, field.id, field.name, field.description, fieldTyp, field.defaultValue, field.baseline,
		field.notNull, field.caseSensitive, field.multilingual, field.active)
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

func (field *Field) GetPhysicalName(provider databaseprovider.DatabaseProvider) string {
	return field.getPhysicalName(provider, field.name)
}

//******************************
// private methods
//******************************
func (field *Field) getSearchableDdl(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) string {
	datatype := field.getSqlDataType(provider)
	if datatype == unknownFieldDataType {
		return unknownFieldDataType
	}
	return strings.TrimSpace(field.getPhysicalName(provider, field.getSearchableFieldName()+" "+
		field.getSqlDataType(provider)))
}

// compare if the physical fields are equal
func (fieldA *Field) equal(fieldB *Field) bool {
	return strings.EqualFold(fieldA.name, fieldB.name) &&
		fieldA.size == fieldB.size &&
		fieldA.fieldType == fieldB.fieldType &&
		fieldA.caseSensitive == fieldB.caseSensitive
}

func (field *Field) getSearchableFieldName() string {
	return searchableFieldPrefix + field.name
}

func (field *Field) toMeta(tableId int32) *meta {
	// we cannot have error here
	var result = new(meta)

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

func (field *Field) isValidInteger(value string) bool {
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
	switch field.fieldType {
	case fieldtype.Long:
		return field.int64Condition(value, size, sign)
	case fieldtype.Int:
		return field.int32Condition(value, size, sign)
	case fieldtype.Short:
		return field.int16Condition(value, size, sign)
	case fieldtype.Byte:
		return field.int08Condition(value, size, sign)
	}
	return false
}

func (field *Field) int08Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 3) ||
		(size == 3 && value <= maxInt08 && sign > 0) ||
		(size == 3 && sign == -1) ||
		(size == 4 && value <= minInt08 && sign == -1)
}

func (field *Field) int16Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 5) ||
		(size == 5 && value <= maxInt16 && sign > 0) ||
		(size == 5 && sign == -1) ||
		(size == 6 && value <= minInt16 && sign == -1)
}
func (field *Field) int32Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 10) ||
		(size == 10 && value <= maxInt32 && sign > 0) ||
		(size == 10 && sign == -1) ||
		(size == 11 && value <= minInt32 && sign == -1)
}
func (field *Field) int64Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 19) ||
		(size == 19 && value <= maxInt64 && sign > 0) ||
		(size == 19 && sign == -1) ||
		(size == 20 && value <= minInt64 && sign == -1)
}

func (field *Field) getDefaultValue(defaultValue string) string {
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

func (field *Field) getDefaultPrimaryKey() *Field {
	switch field.fieldType {
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
	var mapper map[fieldtype.FieldType]string
	var maxStringSize uint16 = 0

	switch provider {
	case databaseprovider.MySql:
		mapper = mysqlDataType
		maxStringSize = mySqlVarcharMaxSize
		break
	case databaseprovider.PostgreSql:
		mapper = postgreDataType
		maxStringSize = postgreVarcharMaxSize
		break
	}

	if val, ok := mapper[field.fieldType]; ok {
		result = val
		if field.fieldType == fieldtype.String {
			if field.size <= 0 || field.size > maxStringSize {
				result = mapper[fieldtype.LongString]
			} else {
				// long text in postgresql
				result += fmt.Sprintf("(%d)", field.size)
			}
		}
	}

	return result
}

func (field *Field) getDateTimeIso8601(value string) (*time.Time, error) {
	var t time.Time
	var e error

	//fmt.Println("getDateTimeIso8601()")
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
			//TODO WRONG way!!!
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

func (field *Field) getPhysicalName(provider databaseprovider.DatabaseProvider, name string) string {
	return sqlfmt.FormatEntityName(provider, name)
}
