package schema

import (
	"strings"
	"testing"
)

func Test__Index__Init(t *testing.T) {
	var aarr = []string{"Gga", "Zorba"}
	elemi := Index{}
	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool, actif bool
	elemi.Init(21, "rel test", "hellkzae", aarr, false, true, true, true)

	if elemi.GetName() != "rel test" {
		t.Errorf("Index.Init() ==> name <> GetName()")
	}
	if elemi.GetId() != 21 {
		t.Errorf("Index.Init() ==> id <> GetId()")
	}
	if elemi.GetDescription() != "hellkzae" {
		t.Errorf("Index.Init() ==> description <> GetDescription()")
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

func Test__Index__ToMeta(t *testing.T) {

	elemi0 := Index{}
	aarr := []string{"Gga", "Zorba", "testllk", "testllk22"}

	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemi0.Init(21, "rel test", "hellkzae", aarr, false, true, false, true)

	meta := elemi0.ToMeta(777)
	elemi1 := meta.ToIndex()

	if elemi0.GetId() != elemi1.GetId() {
		t.Errorf("Index.ToMeta() ==> i0.GetId() must be equal to i1.GetId()")
	}
	if elemi0.GetName() != elemi1.GetName() {
		t.Errorf("Index.ToMeta() ==> i0.GetName() must be equal to i1.GetName()")
	}
	if elemi0.GetDescription() != elemi1.GetDescription() {
		t.Errorf("Index.ToMeta() ==> i0.GetDescription() must be equal to i1.GetDescription()")
	}
	if elemi0.IsBitmap() != elemi1.IsBitmap() {
		t.Errorf("Index.ToMeta() ==> i0.IsBitmap() must be equal to i1.IsBitmap()")
	}
	if elemi0.IsUnique() != elemi1.IsUnique() {
		t.Errorf("Index.ToMeta() ==> i0.IsUnique() must be equal to i1.IsUnique()")
	}
	if elemi0.IsBaseline() != elemi1.IsBaseline() {
		t.Errorf("Index.ToMeta() ==> i0.IsBaseline() must be equal to i1.IsBaseline()")
	}
	if elemi0.IsActive() != elemi1.IsActive() {
		t.Errorf("Index.ToMeta() ==> i0.IsActive() must be equal to i1.IsActive()")
	}
	// test fields
	if elemi1.fields == nil {
		t.Errorf("Index.ToMeta() ==> i1.fields cannot be nil")
	} else {
		// keep ";" hardcoded to detectec metaIndexSeparator constant change
		arr0str := strings.Join(elemi0.fields[:], ";")
		arr1str := strings.Join(elemi1.fields[:], ";")
		if arr0str != arr1str {
			t.Errorf("Index.ToMeta() ==> elemi0.fields is not equal to elemi1.fields")
		}
	}
}
