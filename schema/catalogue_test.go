package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"testing"
)

func Test__Catalogue__GetDql(t *testing.T) {
	cata := new(catalogue)
	//======================
	//==== testing PostgreSql
	//======================
	expectedSql := "SELECT 1 FROM pg_tables WHERE upper(tablename)=$1 AND upper(schemaname)=$2"
	if cata.GetDql(databaseprovider.PostgreSql, entitytype.Table) != expectedSql {
		t.Errorf("Catalogue.GetDql() ==> query must be equal to " + expectedSql)
	}

	//======================
	//==== testing MySql
	//======================
	expectedSql = "SELECT 1 FROM information_schema.tables WHERE upper(table_name)=? AND upper(table_schema)=?"
	if cata.GetDql(databaseprovider.MySql, entitytype.Table) != expectedSql {
		t.Errorf("Catalogue.GetDql() ==> query must be equal to " + expectedSql)
	}

}

func Test__Catalogue__getEntityName(t *testing.T) {
	schema := new(Schema)
	table := new(Table)
	cata := new(catalogue)
	tbl01 := new(tablespace)

	//======================
	//==== testing PostgreSql
	//======================
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	table = table.getMetaTable(databaseprovider.PostgreSql, "zorba")
	if cata.getEntityName(schema, table) != "@meta" {
		t.Errorf("Catalogue.getEntityName() ==> entity name must be equal to '@meta'")
	}
	table = table.getMetaTable(databaseprovider.PostgreSql, "")
	if cata.getEntityName(schema, table) != "@meta" {
		t.Errorf("Catalogue.getEntityName() ==> entity name must be equal to '@meta'")
	}
	tbl01.Init(111, "indexspace", "", "/Temp", false, false)
	if cata.getEntityName(schema, tbl01) != "indexspace" {
		t.Errorf("Catalogue.getEntityName() ==> entity name must be equal to 'indexspace'")
	}
}
