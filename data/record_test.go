package data

import (
	"ring/schema"
	"ring/schema/databaseprovider"
	"testing"
)

func Test__Record__SetField(t *testing.T) {
	var rcd = new(Record)

	schema.Init(databaseprovider.MySql, "zorba")
	rcd.SetRecordType("@meta")
	rcd.SetField("description", "758645454")
	rcd.SetField("reference_id", 7585454)

	if rcd.GetField("description") != "758645454" {
		t.Errorf("Record.SetField() ==> 'description' is not equal to 758645454")
	}
	if rcd.GetField("reference_id") != "7585454" {
		t.Errorf("Record.SetField() ==> 'reference_id' is not equal to 758645454")
		t.Errorf(rcd.GetField("reference_id"))
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
		t.Errorf("Record.SetField() ==> 'active' is not equal to false")
		t.Errorf(rcd.GetField("value"))
	}

}
