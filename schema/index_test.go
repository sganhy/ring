package schema

import (
	"testing"
)

func Test__Index__Init(t *testing.T) {
	var aarr = []string{"Gga", "Zorba"}
	elemi := Index{}
	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool, actif bool
	elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)

	if elemi.GetName() != "rel test" {
		t.Errorf("Index.Init() ==> name <> GetName()")
	}
	if elemi.GetId() != 21 {
		t.Errorf("Index.Init() ==> id <> GetId()")
	}
	if elemi.GetDescription() != "hellkzae" {
		t.Errorf("Index.Init() ==> description <> GetDescription()")
	}
	if elemi.GetTableId() != 52 {
		t.Errorf("Index.Init() ==> tableId <> GetTableId()")
	}
	if elemi.IsBaseline() != true {
		t.Errorf("Index.Init() ==> IsBaseline() <> true")
	}
	if elemi.IsActive() != true {
		t.Errorf("Index.Init() ==> IsActive() <> true")
	}
	if elemi.IsUnique() != true {
		t.Errorf("Index.Init() ==> IsUnique() <> true")
	}
	if elemi.IsBitmap() != false {
		t.Errorf("Index.Init() ==> IsBitmap() <> false")
	}

}
