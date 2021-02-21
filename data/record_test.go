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

//test SetField(), GetField()
func Test__Record__SetField(t *testing.T) {
	var rcd = new(Record)

	// disable connection pool empty connection, string min & max == 0
	schema.Init(databaseprovider.MySql, "", 0, 0)
	rcd.SetRecordType("@meta")

	rcd.SetField("description", "758645454")
	rcd.SetField("reference_id", 7585454)

	if rcd.GetField("description2") != "" {
		t.Errorf("Record.GetFiteld() ==> 'description2' is not empty")
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

	rcd.SetRecordType("@log")
	dt := time.Now()
	rcd.SetField("entry_time", dt)
	if rcd.GetField("entry_time") != dt.UTC().Format(defaultTimeFormat) {
		t.Errorf("Record.SetField() ==> 'entry_time' is not equal to %s", dt.UTC().Format(time.RFC3339))
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

	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes, "schema.@meta",
		physicaltype.Table, -111, tabletype.Business, "", true, false, true, true)
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
