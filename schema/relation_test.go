package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/relationtype"
	"testing"
)

// INIT
func Test__Relation__Init(t *testing.T) {
	elemr0 := Relation{}
	elemt := GetMetaTable(databaseprovider.PostgreSql)
	elemr0.Init(-23, "arel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)

	if elemr0.GetName() != "arel test" {
		t.Errorf("Relation.Init() ==> name <> GetName()")
	}
	if elemr0.GetId() != -23 {
		t.Errorf("Relation.Init() ==> id <> GetId()")
	}
	if elemr0.GetDescription() != "hellkzae" {
		t.Errorf("Relation.Init() ==> description <> GetDescription()")
	}
	if elemr0.GetInverseRelationName() != "hell1" {
		t.Errorf("Relationeld.Init() ==> inverseRelationNam <> GetInverseRelationName()")
	}
	if elemr0.GetMtmTableName() != "52" {
		t.Errorf("Relationeld.Init() ==> GetMtmTable() <> mtm table")
	}
	if elemr0.GetType() != relationtype.Mto {
		t.Errorf("Relationeld.Init() ==> type <> GetType()")
	}
	if elemr0.IsNotNull() != false {
		t.Errorf("Relationeld.Init() ==> IsNotNull() <> false")
	}
	if elemr0.IsBaseline() != true {
		t.Errorf("Relationeld.Init() ==> IsBaseline() <> true")
	}
	if elemr0.IsActive() != false {
		t.Errorf("Relationeld.Init() ==> IsActive() <> false")
	}
	if elemr0.GetToTable() != elemt {
		t.Errorf("Relationeld.Init() ==> GetToTable() <> table pointer")
	}

}
