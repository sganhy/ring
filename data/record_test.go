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
	rcd.SetField("reference_id", "7585454")

	if rcd.GetField("description") != "758645454" {
		t.Errorf("Record.SetField() ==> 'description' is different of 758645454")
	}

	if rcd.GetField("reference_id") != "7585454" {
		t.Errorf("Record.SetField() ==> 'reference_id' is different of 758645454")
	}

}
