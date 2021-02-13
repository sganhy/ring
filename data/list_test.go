package data

import (
	"ring/schema"
	"ring/schema/databaseprovider"
	"testing"
)

func Test__List__AppendItem(t *testing.T) {
	var rcd = new(Record)
	var lst = new(List)
	schema.Init(databaseprovider.MySql, "zorba")
	rcd.SetRecordType("@meta")
	rcd.SetField("description", "758645454")
	rcd.SetField("reference_id", "7585454")
	rcd.SetField("id", "1")
	lst.appendItem(rcd.Copy())
	rcd.SetField("id", "2")
	lst.appendItem(rcd.Copy())
	rcd.SetField("id", "3")
	lst.appendItem(rcd.Copy())
	rcd.SetField("id", "4")
	lst.appendItem(rcd.Copy())

	if lst.Count() != 4 {
		t.Errorf("List.AppendItem() ==> Count() must be equal to 4")
	}
	rcd = lst.ItemByIndex(0)
	if rcd.GetField("id") != "1" {
		t.Errorf("List.AppendItem() ==> ItemByIndex(1) must be equal to 2")
	}
	rcd = lst.ItemByIndex(1)
	if rcd.GetField("id") != "2" {
		t.Errorf("List.AppendItem() ==> ItemByIndex(1) must be equal to 2")
	}
	rcd = lst.ItemByIndex(2)
	if rcd.GetField("id") != "3" {
		t.Errorf("List.AppendItem() ==> ItemByIndex(1) must be equal to 3")
	}
	rcd = lst.ItemByIndex(3)
	if rcd.GetField("id") != "4" {
		t.Errorf("List.AppendItem() ==> ItemByIndex(1) must be equal to 4")
	}
}
