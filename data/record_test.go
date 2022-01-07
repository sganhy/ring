package data

import (
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"testing"
	"time"
)

const defaultTimeFormat = "2006-01-02T15:04:05.000" // rfc3339
const defaultShortTimeFormat = "2006-01-02"         // rfc3339

//test SetField(), GetField()
func Test__Record__SetField(t *testing.T) {
	var rcd = new(Record)

	// disable connection pool empty connection, string min & max == 0
	schema.Init(databaseprovider.MySql, "", 0, 0)
	rcd.SetRecordType("@meta")

	rcd.SetField("description", "758645454")
	rcd.SetField("reference_id", 7585454)
	rcd.SetField("name", "1234567890123456789012345678901")

	if len(rcd.GetField("name")) > 30 {
		t.Errorf("Record.GetField() ==> 'name' must be truncated to 30")
	}
	if rcd.GetField("description2") != "" {
		t.Errorf("Record.GetField() ==> 'description2' is not empty")
	}
	if rcd.GetField("description") != "758645454" {
		t.Errorf("Record.SetField() ==> 'description' is not equal to 758645454")
	}
	if rcd.GetField("reference_id") != "7585454" {
		t.Errorf("Record.SetField() ==> 'reference_id' is not equal to 758645454")
		t.Errorf(rcd.GetField("reference_id"))
	}
	rcd.SetRecordType("@meta") // no reset if recordtype changing
	if rcd.GetField("description") != "" {
		t.Errorf("Record.SetField() ==> 'description' is not empty")
	}
	if rcd.GetField("active") != "false" {
		t.Errorf("Record.SetField() ==> 'active' is not equal to false")
	}
	// get default
	if rcd.GetField("flags") != "0" {
		t.Errorf("Record.SetField() ==> 'flags' is not equal to 0")
	}
	rcd.SetField("active", true)
	if rcd.GetField("active") != "true" {
		t.Errorf("Record.SetField() ==> 'active' is not equal to true")
	}
	rcd.SetField("value", 40.45454545646)
	if rcd.GetField("value") != "40.45454545646" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to false")
	}
	var a uint32 = 444
	rcd.SetField("value", a)
	if rcd.GetField("value") != "444" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 444")
	}
	var b float64 = 4444.4444
	rcd.SetField("value", b)
	if rcd.GetField("value") != "4444.4444" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 444.444")
	}
	var c float32 = 7777.777
	rcd.SetField("value", c)
	if rcd.GetField("value") != "7777.777" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 444.444")
	}
	var d uint64 = 33333333333
	rcd.SetField("value", d)
	if rcd.GetField("value") != "33333333333" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 33333333333")
	}
	var e int8 = -124
	rcd.SetField("value", e)
	if rcd.GetField("value") != "-124" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to -124")
	}
	var f int16 = -12400
	rcd.SetField("value", f)
	if rcd.GetField("value") != "-12400" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to -12400")
	}
	var g int32 = -124000000
	rcd.SetField("value", g)
	if rcd.GetField("value") != "-124000000" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to -124000000")
	}
	var h int64 = -12400000000001
	rcd.SetField("value", h)
	if rcd.GetField("value") != "-12400000000001" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to -12400000000001")
	}
	var j uint8 = 253
	rcd.SetField("value", j)
	if rcd.GetField("value") != "253" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 253")
	}
	var k uint16 = 35300
	rcd.SetField("value", k)
	if rcd.GetField("value") != "35300" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 35300")
	}
	var p uint = 124000220
	rcd.SetField("value", p)
	if rcd.GetField("value") != "124000220" {
		t.Errorf("Record.SetField() ==> 'value' is not equal to 124000220")
	}

	rcd.SetRecordType("@log")
	dt := time.Now()
	rcd.SetField("entry_time", dt)
	if rcd.GetField("entry_time") != dt.UTC().Format(defaultTimeFormat) {
		t.Errorf("Record.SetField() ==> 'entry_time' is not equal to %s", dt.UTC().Format(time.RFC3339))
	}
	// a date time stored in string
	rcd.SetField("method", dt)
	if rcd.GetField("method") != dt.UTC().Format(defaultTimeFormat) {
		t.Errorf("Record.SetField() ==> 'method' is not equal to %s", dt.UTC().Format(time.RFC3339))
	}
	// return errors.New("Unsupported type.")
	if rcd.SetField("method", new(node)) == nil {
		t.Errorf("Record.SetField() ==> 'method' does not return an error")
	}
	// return errors.New("Unsupported type.")
	levelIdError := rcd.SetField("level_id", "hello")
	if levelIdError == nil {
		t.Errorf("Record.SetField() ==> 'level_id' does not return an error")
	}
	if rcd.SetField("222222222222", 12) == nil {
		t.Errorf("Record.SetField() ==> '222222222222' does not return an error")
	}
	rcd = new(Record)
	if rcd.SetField("222222222222", 12) == nil {
		t.Errorf("Record.SetField() ==> '222222222222' does not return an error")
	}
	if rcd.SetField("222222222222", 12) == nil {
		t.Errorf("Record.SetField() ==> '222222222222' does not return an error")
	}

	// TEST dateTime
	var fields = []schema.Field{}
	var relations = []schema.Relation{}
	var indexes = []schema.Index{}
	var table = new(schema.Table)
	var field1 = new(schema.Field)
	var field2 = new(schema.Field)
	var field3 = new(schema.Field)

	field1.Init(1, "test1", "", fieldtype.ShortDateTime, 0, "", true, true, true, false, true)
	fields = append(fields, *field1)
	field2.Init(2, "test2", "", fieldtype.LongDateTime, 0, "", true, true, true, false, true)
	fields = append(fields, *field2)
	field3.Init(3, "test3", "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	fields = append(fields, *field3)

	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "@meta", tabletype.Business, databaseprovider.Undefined, "", true, false, true, true)
	rcd.setRecordType(table)
	dt = time.Now()
	rcd.SetField("test1", dt)
	rcd.SetField("test2", dt)
	rcd.SetField("test3", dt)
	if rcd.GetField("test1") != dt.Format(defaultShortTimeFormat) {
		t.Errorf("Record.SetField() ==> 'test1' is not equal to %s", dt.Format(defaultShortTimeFormat))
	}
	if rcd.GetField("test2") != dt.Format(time.RFC3339Nano) {
		t.Errorf("Record.SetField() ==> 'test2' is not equal to %s", dt.Format(time.RFC3339Nano))
	}
	if rcd.GetField("test3") != dt.UTC().Format(defaultTimeFormat) {
		t.Errorf("Record.SetField() ==> 'test3' is not equal to %s", dt.UTC().Format(defaultTimeFormat))
	}

}

func Test__Record__getField(t *testing.T) {
	var rcd = new(Record)
	const testId = 171717

	schema.Init(databaseprovider.MySql, "", 0, 0)
	rcd.SetRecordType("@lexicon")
	if rcd.getField() != recordIdNotDefined {
		t.Errorf("Record.getField() ==> is not equal to recordIdNotDefined")
	}
	rcd.setField(testId)
	if rcd.getField() != testId {
		t.Errorf("Record.getField() ==> is not equal to %d", testId)
	}
}

func Test__Record__getUpdatedFields(t *testing.T) {
	var rcd = new(Record)
	const testId = 181818

	// disable connection pool empty connection, string min & max == 0
	schema.Init(databaseprovider.MySql, "", 0, 0)
	rcd.SetRecordType("@lexicon")

	// set id
	rcd.setField(testId)
	// detect no changes
	lst := rcd.getUpdatedFields()
	if lst != nil {
		t.Errorf("Record.getUpdatedFields() ==> is not equal to null")
	}

	rcd.SetField("description", "758645454")
	rcd.SetField("uuid", "554554")
	rcd.SetField("name", "1234567890123456789012345678901")
	rcd.SetField("table_id", 4)

	lst = rcd.getUpdatedFields()
	if len(lst) != 4 {
		t.Errorf("Record.getUpdatedFields() ==> len()is not equal to %d", 4)
	}
	// is there nil value
	dict := make(map[string]bool, 12)
	for i := 0; i < len(lst); i++ {
		field := lst[i]
		if field == nil {
			t.Errorf("Record.getUpdatedFields() ==> field cannot be null")
		} else {
			dict[field.GetName()] = true
		}
	}
	if ok := dict["description"]; !ok {
		t.Errorf("Record.getUpdatedFields() ==> field 'description' missing")
	}
	if ok := dict["uuid"]; !ok {
		t.Errorf("Record.getUpdatedFields() ==> field 'uuid' missing")
	}
	if ok := dict["name"]; !ok {
		t.Errorf("Record.getUpdatedFields() ==> field 'name' missing")
	}
	if ok := dict["table_id"]; !ok {
		t.Errorf("Record.getUpdatedFields() ==> field 'table_id' missing")
	}
}
