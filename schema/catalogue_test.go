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
	expectedSql = "SELECT 1 FROM information_schema.tables WHERE upper(table_schema)=? AND upper(table_name)=?"
	if cata.GetDql(databaseprovider.MySql, entitytype.Table) != expectedSql {
		t.Errorf("Catalogue.GetDql() ==> query must be equal to " + expectedSql)
	}

}
