package schema

import (
	"ring/schema/constrainttype"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"testing"
)

// INIT

func Test__Constraint__getPrimaryKey(t *testing.T) {
	table := new(Table)
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	constr := new(constraint)
	constr.Init(constrainttype.PrimaryKey, metaTable)
	t.Errorf(constr.GetDdl(ddlstatement.Create, nil))

}
