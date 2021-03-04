package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
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

	// computed default value
	elemf1 := Field{}
	elemf1.Init(0, "aName Test", "AField Test", fieldtype.Long, 5, "", true, true, false, true, true)
	if elemf1.GetDefaultValue() != "0" {
		t.Errorf("Field.Init() ==> defaultValue <> GetDefaultValue()")
	}
	elemf2 := Field{}
	elemf2.Init(0, "aName Test", "AField Test", fieldtype.Float, 5, "4154", true, true, false, true, true)
	if elemf2.GetDefaultValue() != "4154" {
		t.Errorf("Field.Init() ==> defaultValue <> GetDefaultValue()")
	}
	elemf2.Init(0, "aName Test", "AField Test", fieldtype.Boolean, 5, "", true, true, false, true, true)
	if strings.ToLower(elemf2.GetDefaultValue()) != "false" {
		t.Errorf("Field.Init() ==> defaultValue <> GetDefaultValue()")
	}

}

// getDefaultPrimaryKey
func Test__Field__getDefaultPrimaryKey(t *testing.T) {
	elemf0 := getDefaultPrimaryKey(fieldtype.Short)

	if elemf0.IsValid() != true {
		t.Errorf("Field.getDefaultPrimaryKey() ==> IsValid() <> true")
	}
	if elemf0.IsPrimaryKey() != true {
		t.Errorf("Field.getDefaultPrimaryKey() ==> IsPrimaryKey() <> true")
	}
	if elemf0.fieldType != fieldtype.Short {
		t.Errorf("Field.getDefaultPrimaryKey() ==> fieldType <> fieldtype.Short")
	}

	elemf1 := getDefaultPrimaryKey(fieldtype.Int)
	if elemf1.fieldType != fieldtype.Int {
		t.Errorf("Field.getDefaultPrimaryKey() ==> fieldType <> fieldtype.Int")
	}

	elemf2 := getDefaultPrimaryKey(fieldtype.Double)
	if elemf2.fieldType != fieldtype.Long {
		t.Errorf("Field.getDefaultPrimaryKey() ==> fieldType <> fieldtype.Long")
	}
}

// GetDdlSql
func Test__Field__GetDdlSql(t *testing.T) {
	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 5, "test default", true, true, false, true, true)

	var sql = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Business)
	if strings.ToUpper(sql) != "ANAME FLOAT4 NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME FLOAT4 NULL")
	}

	elemf0.Init(11, "aName", "AField Test", fieldtype.Long, 5, "test default", true, true, false, true, true)
	sql = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME INT8 NOT NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME INT8 NOT NULL")
	}

	elemf0.Init(11, "aName", "AField Test", fieldtype.String, 5, "test default", true, true, false, true, true)
	sql = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME VARCHAR(5) NOT NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME VARCHAR(5)")
	}

	// generated long text datatype
	elemf0.Init(11, "aName", "AField Test", fieldtype.String, 50000000, "test default", true, true, false, true, true)
	sql = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME TEXT NOT NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME TEXT")
	}

}

// GetSearchableValue
func Test__Field__GetSearchableValue(t *testing.T) {
	var lang = Language{}
	lang.Init("FR")

	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 5, "test default", true, false, false, true, true)

	if elemf0.GetSearchableValue("žůžo", lang) == "ZUZO" {
		//t.Errorf("elemf0.GetSearchableValue(\"\") should be equal to \"\" instead of '"+elemf0.GetSearchableValue("")) + "'"
		//t.Errorf("Field.GetSearchableValue(\"žůžo\") ==> should be equal to ZUZO")
	}

	/*
		if len(elemf0.GetSearchableValue("")) == 0 {
			t.Errorf("len(elemf0.GetSearchableValue(\"\")) should be equal to 0")
		}

		if elemf0.GetSearchableValue("a") == "A" {
			t.Errorf("elemf0.GetSearchableValue(\"a\") should be equal to \"A\"")
		}
	*/
}

// ToMeta - test mapping - ToMeta(), ToField
func Test__Field__ToMeta(t *testing.T) {

	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 15, "test default", true, false, false, true, true)
	meta := elemf0.ToMeta(777)
	elemf1 := meta.ToField()

	if elemf0.GetId() != elemf1.GetId() {
		t.Errorf("Field.ToMeta() ==> f0.GetId() must be equal to f1.GetId()")
	}
	if elemf0.GetName() != elemf1.GetName() {
		t.Errorf("Field.ToMeta() ==> f0.GetName() must be equal to f1.GetName()")
	}
	if elemf0.GetDescription() != elemf1.GetDescription() {
		t.Errorf("Field.ToMeta() ==> f0.GetDescription() must be equal to f1.GetDescription()")
	}
	if elemf0.GetType() != elemf1.GetType() {
		t.Errorf("Field.ToMeta() ==> f0.GetType() must be equal to f1.GetType()")
	}
	if elemf0.GetSize() != elemf1.GetSize() {
		t.Errorf("Field.ToMeta() ==> f0.GetSize() must be equal to f1.GetSize()")
	}
	if elemf0.GetDefaultValue() != elemf1.GetDefaultValue() {
		t.Errorf("Field.ToMeta() ==> f0.GetDefaultValue() must be equal to f1.GetDefaultValue()")
	}
	if elemf0.IsBaseline() != elemf1.IsBaseline() {
		t.Errorf("Field.ToMeta() ==> f0.IsBaseline() must be equal to f1.IsBaseline()")
	}
	if elemf0.IsNotNull() != elemf1.IsNotNull() {
		t.Errorf("Field.ToMeta() ==> f0.IsNotNull() must be equal to f1.IsNotNull()")
	}
	if elemf0.IsCaseSensitive() != elemf1.IsCaseSensitive() {
		t.Errorf("Field.ToMeta() ==> f0.IsCaseSensitive() must be equal to f1.IsCaseSensitive()")
	}
	if elemf0.IsMultilingual() != elemf1.IsMultilingual() {
		t.Errorf("Field.ToMeta() ==> f0.IsMultilingual() must be equal to f1.IsMultilingual()")
	}
	if elemf0.IsActive() != elemf1.IsActive() {
		t.Errorf("Field.ToMeta() ==> f0.IsActive() must be equal to f1.IsActive()")
	}

}

//Test isValidInteger
func Test__Field__isValidInteger(t *testing.T) {

	// ========== POSITIVE tests
	// INT
	if isValidInteger("55451", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> 55451 is a valid integer (32 bits)")
	}
	if isValidInteger("-55451", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> -55451 is a valid integer (32 bits)")
	}
	if isValidInteger("0", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> 0 is a valid integer (32 bits)")
	}
	if isValidInteger("2147483647", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> 2147483647 is a valid integer (32 bits)")
	}
	if isValidInteger("2147483646", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> 2147483646 is a valid integer (32 bits)")
	}
	if isValidInteger("-7483646", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> -7483646 is a valid integer (32 bits)")
	}
	if isValidInteger("-2147483648", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> -2147483648 is a valid integer (32 bits)")
	}
	if isValidInteger("-2147483647", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> -2147483647 is a valid integer (32 bits)")
	}
	if isValidInteger("-214748364", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> -214748364 is a valid integer (32 bits)")
	}
	if isValidInteger("758645454", fieldtype.Int) == false {
		t.Errorf("Field.isValidInteger() ==> 758645454 is a valid integer (32 bits)")
	}
	// LONG
	if isValidInteger("9223372036854775807", fieldtype.Long) == false {
		t.Errorf("Field.isValidInteger() ==> 9223372036854775807 string is not a valid integer (64 bits)")
	}
	if isValidInteger("0", fieldtype.Long) == false {
		t.Errorf("Field.isValidInteger() ==> 0 string is not a valid integer (64 bits)")
	}
	//SHORT
	if isValidInteger("-32768", fieldtype.Short) == false {
		t.Errorf("Field.isValidInteger() ==> -32768 string is not a valid integer (16 bits)")
	}
	if isValidInteger("32767", fieldtype.Short) == false {
		t.Errorf("Field.isValidInteger() ==> 32767 string is not a valid integer (16 bits)")
	}
	if isValidInteger("-564", fieldtype.Short) == false {
		t.Errorf("Field.isValidInteger() ==> -564 string is not a valid integer (16 bits)")
	}
	if isValidInteger("0", fieldtype.Short) == false {
		t.Errorf("Field.isValidInteger() ==> 0 string is not a valid integer (16 bits)")
	}

	//BYTE
	if isValidInteger("127", fieldtype.Byte) == false {
		t.Errorf("Field.isValidInteger() ==> 127 string is not a valid integer (8 bits)")
	}
	if isValidInteger("11", fieldtype.Byte) == false {
		t.Errorf("Field.isValidInteger() ==> 127 string is not a valid integer (8 bits)")
	}
	if isValidInteger("-128", fieldtype.Byte) == false {
		t.Errorf("Field.isValidInteger() ==> -128 string is not a valid integer (8 bits)")
	}

	// ========== NEGATIVE tests
	//INT
	if isValidInteger("", fieldtype.Long) != false {
		t.Errorf("Field.isValidInteger() ==> empty string is not a valid integer (32 bits)")
	}
	if isValidInteger("21474836490", fieldtype.Int) != false {
		t.Errorf("Field.isValidInteger() ==> 21474836490 is not a valid integer (32 bits)")
	}
	if isValidInteger("2-14836490", fieldtype.Int) != false {
		t.Errorf("Field.isValidInteger() ==> 2-14836490 is not a valid integer (32 bits)")
	}
	if isValidInteger("2147483648", fieldtype.Int) != false {
		t.Errorf("Field.isValidInteger() ==> 2147483648 is not a valid integer (32 bits)")
	}
	if isValidInteger("-2147483651", fieldtype.Int) != false {
		t.Errorf("Field.isValidInteger() ==> -2147483651 is not a valid integer (32 bits)")
	}
	if isValidInteger("", fieldtype.Int) != false {
		t.Errorf("Field.isValidInteger() ==> empty string is not a valid integer (32 bits)")
	}
	if isValidInteger("9223372036854775808", fieldtype.Long) != false {
		t.Errorf("Field.isValidInteger() ==> 9223372036854775808 string is not a valid integer (64 bits)")
	}
	//SHORT
	if isValidInteger("40000", fieldtype.Short) != false {
		t.Errorf("Field.isValidInteger() ==> 40000 string is not a valid integer (16 bits)")
	}
	if isValidInteger("-40000", fieldtype.Short) != false {
		t.Errorf("Field.isValidInteger() ==> -40000 string is not a valid integer (64 bits)")
	}
	//BYTE
	if isValidInteger("256", fieldtype.Byte) != false {
		t.Errorf("Field.isValidInteger() ==> 256 string is not a valid integer (8 bits)")
	}
	if isValidInteger("128", fieldtype.Byte) != false {
		t.Errorf("Field.isValidInteger() ==> 128 string is not a valid integer (8 bits)")
	}
	if isValidInteger("-129", fieldtype.Byte) != false {
		t.Errorf("Field.isValidInteger() ==> -129 string is not a valid integer (8 bits)")
	}

}

//Test getDateTimeIso8601, GetDateTimeString
func Test__Field__getDateTimeIso8601(t *testing.T) {

	elemf0 := new(Field)
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.DateTime, 5, "test default", true, false, false, true, true)
	var ti *time.Time
	_, offset := time.Now().Zone()

	// test ==> fieldtype.DateTime
	// sample #1
	ti, _ = getDateTimeIso8601("2016-12-12")
	*ti = ti.Add(time.Second * time.Duration(offset)) // adapt to local time for date time
	if elemf0.GetDateTimeString(*ti) != "2016-12-12T00:00:00.000" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-12-12') must return '2016-12-12T00:00:00.000'")
	}
	// sample #2
	ti, _ = getDateTimeIso8601("2016-11-12 12:13:14")
	*ti = ti.Add(time.Second * time.Duration(offset))
	if elemf0.GetDateTimeString(*ti) != "2016-11-12T12:13:14.000" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-11-12 12:13:14') must return '2016-11-12T12:13:14.000'")
	}
	// sample #3
	ti, _ = getDateTimeIso8601("2016-11-12T12:13:14")
	*ti = ti.Add(time.Second * time.Duration(offset))
	if elemf0.GetDateTimeString(*ti) != "2016-11-12T12:13:14.000" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-11-12 12:13:14') must return '2016-11-12T12:13:14.000'")
	}

	// test ==> fieldtype.DateTime
	// sample #1
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.ShortDateTime, 5, "test default", true, false, false, true, true)
	ti, _ = getDateTimeIso8601("2016-11-12T12:13:14")
	if elemf0.GetDateTimeString(*ti) != "2016-11-12" {
		t.Errorf("Field.getDateTimeIso8601() ==> getDateTimeIso8601('2016-11-12 12:13:14') must return '2016-11-12'")
	}

}
