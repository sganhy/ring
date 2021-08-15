package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"testing"
)

func Test__Tablespace__Init(t *testing.T) {
	tbl01 := new(tablespace)
	tbl01.Init(3333, "indexspace", "ATable Test", "/data/indexes", false, false)

	if tbl01.GetId() != 3333 {
		t.Errorf("Tablespace.Init() ==> id <> GetId()")
	}
	if tbl01.GetName() != "indexspace" {
		t.Errorf("Tablespace.Init() ==> name <> GetName()")
	}
	if tbl01.GetDescription() != "ATable Test" {
		t.Errorf("Tablespace.Init() ==> description <> GetDescription()")
	}
	if tbl01.GetPhysicalName() != "indexspace" {
		t.Errorf("Tablespace.Init() ==> physical name <> GetPhysicalName()")
	}
	if tbl01.GetEntityType() != entitytype.Tablespace {
		t.Errorf("Tablespace.Init() ==> entitytype.Tablespace <> GetEntityType()")
	}
	tbl01.setName("tablespace")
	if tbl01.GetName() != "tablespace" {
		t.Errorf("Tablespace.Init() ==> name <> GetName()")
	}
}

func Test__Tablespace__getDdlCreate(t *testing.T) {
	tbl01 := new(tablespace)
	expectedDll := "CREATE TABLESPACE indexspace LOCATION '/data/indexes'"
	tbl01.Init(111, "indexspace", "", "/data/indexes", false, false)

	if tbl01.getDdlCreate(databaseprovider.PostgreSql) != expectedDll {
		t.Errorf("Tablespace.getDdlCreate() ==> must be to %s", expectedDll)
	}

}