package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/searchabletype"
	"ring/schema/tabletype"
	"strings"
	"testing"
	"time"
)

// INIT
func Test__Field__Init(t *testing.T) {
	elemf0 := Field{}
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.String, 5, "test default", true, false, false, true, true)

	if elemf0.GetName() != "aName Test" {
		t.Errorf("Field.Init() ==> name <> GetName()")
	}
	if elemf0.GetId() != 11 {
		t.Errorf("Field.Init() ==> id <> GetId()")
	}
	if elemf0.GetDescription() != "AField Test" {
		t.Errorf("Field.Init() ==> description <> GetDescription()")
	}
	if elemf0.GetType() != fieldtype.String {
		t.Errorf("Field.Init() ==> type <> GetType()")
	}
	if elemf0.GetSize() != 5 {
		t.Errorf("Field.Init() ==> size <> GetSize()")
	}
	if elemf0.GetDefaultValue() != "test default" {
		t.Errorf("Field.Init() ==> defaultValue <> GetDefaultValue()")
	}
	if elemf0.IsBaseline() != true {
		t.Errorf("Field.Init() ==> baseline <> IsBaseline()")
	}
	if elemf0.IsNotNull() != false {
		t.Errorf("Field.Init() ==> notnull <> IsNotNull()")
	}
	if elemf0.IsCaseSensitive() != false {
		t.Errorf("Field.Init() ==> caseSensitive <> IsCaseSensitive()")
	}
	if elemf0.IsMultilingual() != true {
		t.Errorf("Field.Init() ==> multilingual <> IsMultilingual()")
	}
	if elemf0.IsActive() != true {
		t.Errorf("Field.Init() ==> active <> IsActive()")
	}
	if elemf0.IsDateTime() != false {
		t.Errorf("Field.Init() ==> IsDateTime() <> true")
	}
	if elemf0.IsNumeric() != false {
		t.Errorf("Field.Init() ==> IsNumeric() <> false")
	}
	if elemf0.IsValid() != true {
		t.Errorf("Field.Init() ==> IsValid() <> true")
	}
	if elemf0.IsPrimaryKey() != false {
		t.Errorf("Field.Init() ==> IsPrimaryKey() <> false")
	}
	if elemf0.GetEntityType() != entitytype.Field {
		t.Errorf("Field.Init() ==>  GetEntityType() <> entitytype.Field")
	}

	// computed default value
	elemf1 := Field{}
	elemf1.Init(0, "aName Test", "AField Test", fieldtype.Long, 5, "", true, true, false, true, true)
	if elemf1.GetDefaultValue() != "0" {
		t.Errorf("Field.GetDefaultValue() ==> defaultValue <> GetDefaultValue()")
	}
	elemf2 := Field{}
	elemf2.Init(0, "aName Test", "AField Test", fieldtype.Float, 5, "4154", true, true, false, true, true)
	if elemf2.GetDefaultValue() != "4154" {
		t.Errorf("Field.GetDefaultValue() ==> defaultValue <> GetDefaultValue()")
	}
	elemf2.Init(0, "aName Test", "AField Test", fieldtype.Boolean, 5, "", true, true, false, true, true)
	if strings.ToLower(elemf2.GetDefaultValue()) != "false" {
		t.Errorf("Field.GetDefaultValue() ==> defaultValue <> GetDefaultValue()")
	}
	elemf3 := Field{}
	elemf3.Init(0, "aName Test", "AField Test", fieldtype.Double, 5, "4154", true, true, false, true, true)
	if elemf3.GetDefaultValue() != "4154" {
		t.Errorf("Field.GetDefaultValue() ==> defaultValue <> GetDefaultValue()")
	}

}

// getDefaultPrimaryKey
func Test__Field__getDefaultPrimaryKey(t *testing.T) {
	field := new(Field)
	field.setType(fieldtype.Short)
	elemf0 := field.getDefaultPrimaryKey()

	if elemf0.IsValid() != true {
		t.Errorf("Field.getDefaultPrimaryKey() ==> IsValid() <> true")
	}
	if elemf0.IsPrimaryKey() != true {
		t.Errorf("Field.getDefaultPrimaryKey() ==> IsPrimaryKey() <> true")
	}
	if elemf0.GetType() != fieldtype.Short {
		t.Errorf("Field.getDefaultPrimaryKey() ==> fieldType <> fieldtype.Short")
	}

	field.setType(fieldtype.Int)
	elemf1 := field.getDefaultPrimaryKey()
	if elemf1.GetType() != fieldtype.Int {
		t.Errorf("Field.getDefaultPrimaryKey() ==> fieldType <> fieldtype.Int")
	}

	field.setType(fieldtype.Double)
	elemf2 := field.getDefaultPrimaryKey()
	if elemf2.GetType() != fieldtype.Long {
		t.Errorf("Field.getDefaultPrimaryKey() ==> fieldType <> fieldtype.Long")
	}
}

// GetDdlSql
func Test__Field__GetDdl(t *testing.T) {
	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 5, "test default", true, true, false, true, true)

	var sql = elemf0.GetDdl(databaseprovider.PostgreSql, tabletype.Business)
	if strings.ToUpper(sql) != "ANAME FLOAT4" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME FLOAT4")
	}

	elemf0.Init(11, "aName", "AField Test", fieldtype.Long, 5, "test default", true, true, false, true, true)
	sql = elemf0.GetDdl(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME INT8" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME INT8")
	}

	elemf0.Init(11, "aName", "AField Test", fieldtype.String, 5, "test default", true, true, false, true, true)
	sql = elemf0.GetDdl(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME VARCHAR(5)" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME VARCHAR(5)")
	}

	// generated long text datatype
	elemf0.Init(11, "aName", "AField Test", fieldtype.String, 50000000, "test default", true, true, false, true, true)
	sql = elemf0.GetDdl(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME TEXT" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME TEXT")
	}

}

// GetSearchableValue
func Test__Field__GetSearchableValue(t *testing.T) {
	var lang = Language{}
	lang.Init(1, "FR")

	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 5, "test default", true, false, false, true, true)

	if strings.Compare(elemf0.GetSearchableValue("žůžo", searchabletype.None), "ZUZO") != 0 {
		t.Errorf("Field.GetSearchableValue(\"žůžo\") ==> should be equal to ZUZO")
	}
	if elemf0.GetSearchableValue("žůžo", searchabletype.None) != "ZUZO" {
		t.Errorf("Field.GetSearchableValue(\"žůžo\") ==> should be equal to ZUZO")
	}
	if len(elemf0.GetSearchableValue("", searchabletype.None)) != 0 {
		t.Errorf("len(elemf0.GetSearchableValue(\"\")) should be equal to 0")
	}
	if elemf0.GetSearchableValue("a", searchabletype.None) != "A" {
		t.Errorf("elemf0.GetSearchableValue(\"a\") should be equal to \"A\"")
	}
}

//test mappers Meta to Field, and Index to Meta
func Test__Field__toMeta(t *testing.T) {

	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 15, "test default", true, false, false, true, true)
	metaData := elemf0.toMeta(777)
	elemf1 := metaData.toField()

	if elemf0.GetId() != elemf1.GetId() {
		t.Errorf("Field.toMeta() ==> f0.GetId() must be equal to f1.GetId()")
	}
	if elemf0.GetName() != elemf1.GetName() {
		t.Errorf("Field.toMeta() ==> f0.GetName() must be equal to f1.GetName()")
	}
	if elemf0.GetDescription() != elemf1.GetDescription() {
		t.Errorf("Field.toMeta() ==> f0.GetDescription() must be equal to f1.GetDescription()")
	}
	if elemf0.GetType() != elemf1.GetType() {
		t.Errorf("Field.toMeta() ==> f0.GetType() must be equal to f1.GetType()")
	}
	if elemf0.GetSize() != elemf1.GetSize() {
		t.Errorf("Field.toMeta() ==> f0.GetSize() must be equal to f1.GetSize()")
	}
	if elemf0.GetDefaultValue() != elemf1.GetDefaultValue() {
		t.Errorf("Field.toMeta() ==> f0.GetDefaultValue() must be equal to f1.GetDefaultValue()")
	}
	if elemf0.IsBaseline() != elemf1.IsBaseline() {
		t.Errorf("Field.toMeta() ==> f0.IsBaseline() must be equal to f1.IsBaseline()")
	}
	if elemf0.IsNotNull() != elemf1.IsNotNull() {
		t.Errorf("Field.toMeta() ==> f0.IsNotNull() must be equal to f1.IsNotNull()")
	}
	if elemf0.IsCaseSensitive() != elemf1.IsCaseSensitive() {
		t.Errorf("Field.toMeta() ==> f0.IsCaseSensitive() must be equal to f1.IsCaseSensitive()")
	}
	if elemf0.IsMultilingual() != elemf1.IsMultilingual() {
		t.Errorf("Field.toMeta() ==> f0.IsMultilingual() must be equal to f1.IsMultilingual()")
	}
	if elemf0.IsActive() != elemf1.IsActive() {
		t.Errorf("Field.toMeta() ==> f0.IsActive() must be equal to f1.IsActive()")
	}
	if elemf0.String() != elemf1.String() {
		t.Errorf("Field.toMeta() ==> f0.String() must be equal to f1.String()")
	}

}

//Test isValidInteger
func Test__Field__isValidInteger(t *testing.T) {
	var field = new(Field)
	// ========== POSITIVE tests
	// INT
	field.setType(fieldtype.Int)
	if field.isValidInteger("55451") == false {
		t.Errorf("Field.isValidInteger() ==> 55451 is a valid integer (32 bits)")
	}
	if field.isValidInteger("-55451") == false {
		t.Errorf("Field.isValidInteger() ==> -55451 is a valid integer (32 bits)")
	}
	if field.isValidInteger("0") == false {
		t.Errorf("Field.isValidInteger() ==> 0 is a valid integer (32 bits)")
	}
	if field.isValidInteger("2147483647") == false {
		t.Errorf("Field.isValidInteger() ==> 2147483647 is a valid integer (32 bits)")
	}
	if field.isValidInteger("2147483646") == false {
		t.Errorf("Field.isValidInteger() ==> 2147483646 is a valid integer (32 bits)")
	}
	if field.isValidInteger("-7483646") == false {
		t.Errorf("Field.isValidInteger() ==> -7483646 is a valid integer (32 bits)")
	}
	if field.isValidInteger("-2147483648") == false {
		t.Errorf("Field.isValidInteger() ==> -2147483648 is a valid integer (32 bits)")
	}
	if field.isValidInteger("-2147483647") == false {
		t.Errorf("Field.isValidInteger() ==> -2147483647 is a valid integer (32 bits)")
	}
	if field.isValidInteger("-214748364") == false {
		t.Errorf("Field.isValidInteger() ==> -214748364 is a valid integer (32 bits)")
	}
	if field.isValidInteger("758645454") == false {
		t.Errorf("Field.isValidInteger() ==> 758645454 is a valid integer (32 bits)")
	}
	// LONG
	field.setType(fieldtype.Long)
	if field.isValidInteger("9223372036854775807") == false {
		t.Errorf("Field.isValidInteger() ==> 9223372036854775807 string is not a valid integer (64 bits)")
	}
	if field.isValidInteger("0") == false {
		t.Errorf("Field.isValidInteger() ==> 0 string is not a valid integer (64 bits)")
	}
	//SHORT
	field.setType(fieldtype.Short)
	if field.isValidInteger("-32768") == false {
		t.Errorf("Field.isValidInteger() ==> -32768 string is not a valid integer (16 bits)")
	}
	if field.isValidInteger("32767") == false {
		t.Errorf("Field.isValidInteger() ==> 32767 string is not a valid integer (16 bits)")
	}
	if field.isValidInteger("-564") == false {
		t.Errorf("Field.isValidInteger() ==> -564 string is not a valid integer (16 bits)")
	}
	if field.isValidInteger("0") == false {
		t.Errorf("Field.isValidInteger() ==> 0 string is not a valid integer (16 bits)")
	}

	//BYTE
	field.setType(fieldtype.Byte)
	if field.isValidInteger("127") == false {
		t.Errorf("Field.isValidInteger() ==> 127 string is not a valid integer (8 bits)")
	}
	if field.isValidInteger("11") == false {
		t.Errorf("Field.isValidInteger() ==> 127 string is not a valid integer (8 bits)")
	}
	if field.isValidInteger("-128") == false {
		t.Errorf("Field.isValidInteger() ==> -128 string is not a valid integer (8 bits)")
	}

	// ========== NEGATIVE tests
	//INT
	field.setType(fieldtype.Int)
	if field.isValidInteger("21474836490") != false {
		t.Errorf("Field.isValidInteger() ==> 21474836490 is not a valid integer (32 bits)")
	}
	if field.isValidInteger("2-14836490") != false {
		t.Errorf("Field.isValidInteger() ==> 2-14836490 is not a valid integer (32 bits)")
	}
	if field.isValidInteger("2147483648") != false {
		t.Errorf("Field.isValidInteger() ==> 2147483648 is not a valid integer (32 bits)")
	}
	if field.isValidInteger("-2147483651") != false {
		t.Errorf("Field.isValidInteger() ==> -2147483651 is not a valid integer (32 bits)")
	}
	if field.isValidInteger("") != false {
		t.Errorf("Field.isValidInteger() ==> empty string is not a valid integer (32 bits)")
	}
	field.setType(fieldtype.Long)
	if field.isValidInteger("9223372036854775808") != false {
		t.Errorf("Field.isValidInteger() ==> 9223372036854775808 string is not a valid integer (64 bits)")
	}
	//SHORT
	field.setType(fieldtype.Short)
	if field.isValidInteger("40000") != false {
		t.Errorf("Field.isValidInteger() ==> 40000 string is not a valid integer (16 bits)")
	}
	if field.isValidInteger("-40000") != false {
		t.Errorf("Field.isValidInteger() ==> -40000 string is not a valid integer (64 bits)")
	}
	//BYTE
	field.setType(fieldtype.Byte)
	if field.isValidInteger("256") != false {
		t.Errorf("Field.isValidInteger() ==> 256 string is not a valid integer (8 bits)")
	}
	if field.isValidInteger("128") != false {
		t.Errorf("Field.isValidInteger() ==> 128 string is not a valid integer (8 bits)")
	}
	if field.isValidInteger("-129") != false {
		t.Errorf("Field.isValidInteger() ==> -129 string is not a valid integer (8 bits)")
	}
	// force false last return ... test
	field.setType(fieldtype.Boolean)
	if field.isValidInteger("11111") != false {
		t.Errorf("Field.isValidInteger() ==> 11111 string is not a valid Boolean")
	}

}

//Test getDateTimeIso8601, GetDateTimeString
func Test__Field__getDateTimeIso8601(t *testing.T) {

	elemf0 := new(Field)
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.DateTime, 5, "test default", true, false, false, true, true)
	var ti *time.Time
	_, offset := time.Now().Zone()

	//*** test ==> fieldtype.DateTime
	// sample #1
	ti, _ = elemf0.getDateTimeIso8601("2016-12-12")
	*ti = ti.Add(time.Second * time.Duration(offset)) // adapt to local time for date time
	if elemf0.GetDateTimeString(*ti) != "2016-12-12T00:00:00.000" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-12-12') must return '2016-12-12T00:00:00.000'")
	}
	// sample #2
	ti, _ = elemf0.getDateTimeIso8601("2016-11-12 12:13:14")
	*ti = ti.Add(time.Second * time.Duration(offset))
	if elemf0.GetDateTimeString(*ti) != "2016-11-12T12:13:14.000" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-11-12 12:13:14') must return '2016-11-12T12:13:14.000'")
	}
	// sample #3
	ti, _ = elemf0.getDateTimeIso8601("2016-11-12T12:13:14")
	*ti = ti.Add(time.Second * time.Duration(offset))
	if elemf0.GetDateTimeString(*ti) != "2016-11-12T12:13:14.000" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-11-12 12:13:14') must return '2016-11-12T12:13:14.000'")
	}

	// test ==> fieldtype.ShortDateTime
	// sample #1
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.ShortDateTime, 5, "test default", true, false, false, true, true)
	ti, _ = elemf0.getDateTimeIso8601("2016-11-12T12:13:14")
	if elemf0.GetDateTimeString(*ti) != "2016-11-12" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-11-12 12:13:14') must return '2016-11-12'")
	}

	//*** test ==> fieldtype.ShortDateTime

}

//Test GetParameterValue
func Test__Field__GetParameterValue(t *testing.T) {
	elemf0 := new(Field)
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.Boolean, 5, "test default", true, false, false, true, true)

	// boolean
	if elemf0.GetParameterValue("true") != true {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('true') must be equal to true")
	}
	if elemf0.GetParameterValue("false") != false {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('false') must be equal to false")
	}
	// integer
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.Short, 5, "test default", true, false, false, true, true)
	elemf0.GetParameterValue("444")
	var valueInt = elemf0.GetParameterValue("444").(int64)
	if valueInt != 444 {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('444') must be equal to 444")
	}
	valueInt = elemf0.GetParameterValue("-32000").(int64)
	if valueInt != -32000 {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('-32000') must be equal to -32000")
	}
	// float
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.Float, 5, "test default", true, false, false, true, true)
	var valueFloat = elemf0.GetParameterValue("-32.10001").(float64)
	if valueFloat != -32.10001 {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('-32.10001') must be equal to -32.10001")
	}
	// string
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.String, 5, "test default", true, false, false, true, true)
	var valueStr = elemf0.GetParameterValue("Hello world").(string)
	if valueStr != "Hello world" {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('Hello world') must be equal to 'Hello world'")
	}
	// dateTime
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.DateTime, 5, "", true, false, false, true, true)
	var valueDt = elemf0.GetParameterValue("2016-11-12T12:13:14.071").(time.Time)
	if valueDt.String() != "2016-11-12 12:13:14.071 +0000 UTC" {
		t.Errorf("Field.GetParameterValue() ==> GetParameterValue('2016-11-12T12:13:14.071') must be equal to '2016-11-12 12:13:14.071 +0000 UTC'")
	}

}

//Test GetValue
func Test__Field__GetValue(t *testing.T) {
	elemf0 := new(Field)
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.Int, 5, "test default", true, false, false, true, true)

	//*** test 1 - int32
	value, err := elemf0.GetValue("2222")
	if value != "2222" {
		t.Errorf("Field.GetValue() ==> GetValue('2222') must be equal to '2222'")
	}
	if err != nil {
		t.Errorf("Field.GetValue() ==> GetValue('2222') error must be equal to null")
	}
	//*** test 2 - int32
	value, err = elemf0.GetValue("22222222222222")
	if value != invalidValue {
		t.Errorf("Field.GetValue() ==> GetValue('22222222222222') must be equal to invalidValue")
	}
	if err.Error() != errorInvalidValueType {
		t.Errorf("Field.GetValue() ==> GetValue('22222222222222') error must be equal to null to '%s'", errorInvalidValueType)
	}
	//*** test 3 - int64
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.Long, 5, "test default", true, false, false, true, true)
	value, err = elemf0.GetValue("-9223372036854775808")
	if value != "-9223372036854775808" {
		t.Errorf("Field.GetValue() ==> GetValue('-9223372036854775808') must be equal to '-9223372036854775808'")
	}
	//*** test 4 - bool
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.Boolean, 5, "test default", true, false, false, true, true)
	value, err = elemf0.GetValue("true")
	if value != "true" {
		t.Errorf("Field.GetValue() ==> GetValue('true') must be equal to 'true'")
	}
	//*** test 5 - bool
	value, err = elemf0.GetValue("22222")
	if value != invalidValue {
		t.Errorf("Field.GetValue() ==> GetValue('22222') must be equal to invalidValue")
	}

	//*** test 6 - string
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.String, 10, "test default", true, false, false, true, true)
	value, err = elemf0.GetValue("01234567890")
	if value != "0123456789" {
		t.Errorf("Field.GetValue() ==> GetValue('01234567890') must be equal to '0123456789'")
	}
	value, err = elemf0.GetValue("0123456789")
	if value != "0123456789" {
		t.Errorf("Field.GetValue() ==> GetValue('0123456789') must be equal to '0123456789'")
	}

	//*** test 7 - dateTime
	/*
		elemf0.Init(11, "aName Test", "AField Test", fieldtype.DateTime, 10, "test default", true, false, false, true, true)
		value, err = elemf0.GetValue("2020-04-30T00:00:00.000Z")
		if value != "2020-04-29T22:00:00.000" {
			t.Errorf("Field.GetValue() ==> GetValue('2020-04-30T00:00:00.000Z') must be equal to '2020-04-29T22:00:00.000'")
		}
	*/
}

func Test__Field__IsValueValid(t *testing.T) {
	elemf0 := new(Field)
	elemf0.Init(11, "a_name", "AField Test", fieldtype.NotDefined, 10, "test default", true, false, false, true, true)

	// get default validation should be equal to false
	if elemf0.IsValueValid("00") != false {
		t.Errorf("Field.IsValueValid() ==> IsValueValid('00') must be equal to false")
	}
}

func Test__Field__getSearchableDdl(t *testing.T) {
	elemf0 := new(Field)
	elemf0.Init(11, "a_name", "AField Test", fieldtype.String, 10, "test default", true, false, false, true, true)

	var provider = databaseprovider.PostgreSql
	var tableType = tabletype.Business

	if elemf0.getSearchableDdl(provider, tableType) != "s_a_name varchar(10)" {
		t.Errorf("Field.getSearchableDdl() ==> getSearchableDdl('a_name') must be equal to 's_a_name varchar(10)'")
	}

	if elemf0.getSearchableDdl(databaseprovider.NotDefined, tableType) != unknownFieldDataType {
		t.Errorf("Field.getSearchableDdl() ==> getSearchableDdl(databaseprovider.NotDefined, 'a_name') must be empty")
	}

}

func Test__Field__equal(t *testing.T) {
	elemf0 := new(Field)
	elemf0.Init(11, "a_name", "AField Test", fieldtype.String, 10, "test default", true, false, false, true, true)
	elemf1 := elemf0.Clone()

	// positive test
	if elemf0.equal(elemf1) == false {
		t.Errorf("Field.equal() ==> f0.equal(f0.Clone()) must be equal to true")
	}
	// positive test - name comparison is not case sensitive
	elemf0.setName("test")
	elemf1 = elemf0.Clone()
	elemf1.setName("tEst")
	if elemf0.equal(elemf1) == false {
		t.Errorf("Field.equal() ==> {0} f0.equal(f1) must be equal to true")
	}

	// negative test - size changed
	elemf1 = elemf0.Clone()
	elemf1.setSize(11)
	if elemf0.equal(elemf1) == true {
		t.Errorf("Field.equal() ==> {1} f0.equal(f1) must be equal to false")
	}
	// negative test - casesensitive
	elemf0.setCaseSensitive(false)
	elemf1 = elemf0.Clone()
	elemf1.setCaseSensitive(true)
	if elemf0.equal(elemf1) == true {
		t.Errorf("Field.equal() ==> {2} f0.equal(f1) must be equal to false")
	}

}
