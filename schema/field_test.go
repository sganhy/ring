package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/tabletype"
	"strings"
	"testing"
)

// INIT
func Test__Field__Init(t *testing.T) {
	elemf0 := Field{}
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.DateTime, 5, "test default", true, false, false, true, true)

	if elemf0.GetName() != "aName Test" {
		t.Errorf("Field.Init() ==> name <> GetName()")
	}
	if elemf0.GetId() != 11 {
		t.Errorf("Field.Init() ==> id <> GetId()")
	}
	if elemf0.GetDescription() != "AField Test" {
		t.Errorf("Field.Init() ==> description <> GetDescription()")
	}
	if elemf0.GetType() != fieldtype.DateTime {
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
	if elemf0.IsDateTime() != true {
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

	elemf1 := Field{}
	elemf1.Init(0, "aName Test", "AField Test", fieldtype.DateTime, 5, "test default", true, false, false, true, true)
	if elemf1.IsValid() != false {
		t.Errorf("Field.Init() ==> IsValid() <> false")
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

	var sql, _ = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Business)
	if strings.ToUpper(sql) != "ANAME FLOAT4 NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME FLOAT4 NULL")
	}

	elemf0.Init(11, "aName", "AField Test", fieldtype.Long, 5, "test default", true, true, false, true, true)
	sql, _ = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME INT8 NOT NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME INT8 NOT NULL")
	}

	elemf0.Init(11, "aName", "AField Test", fieldtype.String, 5, "test default", true, true, false, true, true)
	sql, _ = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME VARCHAR(5) NOT NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME VARCHAR(5)")
	}

	// generated long text datatype
	elemf0.Init(11, "aName", "AField Test", fieldtype.String, 50000000, "test default", true, true, false, true, true)
	sql, _ = elemf0.GetDdlSql(databaseprovider.PostgreSql, tabletype.Meta)
	if strings.ToUpper(sql) != "ANAME TEXT NOT NULL" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to ANAME TEXT")
	}

}

// GetSearchableValue
func Test__Field__GetSearchableValue(t *testing.T) {
	var lang = NewLanguage("Fr")

	elemf0 := Field{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemf0.Init(11, "aName", "AField Test", fieldtype.Float, 5, "test default", true, false, false, true, true)

	if elemf0.GetSearchableValue("žůžo", *lang) == "ZUZO" {
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
