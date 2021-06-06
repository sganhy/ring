package schema

import (
	"ring/schema/constrainttype"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"testing"
)

func Test__Constraint__Init(t *testing.T) {
	table := new(Table)
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	//======================
	//==== testing PostgreSql
	//======================
	constr := new(constraint)
	constr.Init(constrainttype.PrimaryKey, metaTable)

	if constr.GetEntityType() != entitytype.Constraint {
		t.Errorf("Constraint.GetEntityType() ==>  must be equal to " + entitytype.Constraint.String())
	}

}

func Test__Constraint__getDdlPrimaryKey(t *testing.T) {
	table := new(Table)
	//======================
	//==== testing PostgreSql
	//======================
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	constr := new(constraint)
	constr.Init(constrainttype.PrimaryKey, metaTable)
	expectedSql := "ALTER TABLE information_schema.\"@meta\" ADD CONSTRAINT \"pk_@meta\" PRIMARY KEY (id,schema_id,object_type,reference_id)"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlPrimaryKey() ==> query must be equal to " + expectedSql)
	}

	//======================
	//==== testing MySql
	//======================
	metaTable = table.getMetaTable(databaseprovider.MySql, "information_schema")
	constr = new(constraint)
	constr.Init(constrainttype.PrimaryKey, metaTable)
	expectedSql = "ALTER TABLE information_schema.`@meta` ADD CONSTRAINT `pk_@meta` PRIMARY KEY (id,schema_id,object_type,reference_id)"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlPrimaryKey() ==> query must be equal to " + expectedSql)
	}
}

func Test__Constraint__getDdlNotNull(t *testing.T) {
	table := new(Table)
	//======================
	//==== testing PostgreSql
	//======================
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	constr := new(constraint)
	constr.Init(constrainttype.NotNull, metaTable)
	constr.setField(metaTable.GetFieldByName("reference_id"))
	expectedSql := "ALTER TABLE information_schema.\"@meta\" ALTER COLUMN reference_id SET NOT NULL"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlNotNull() ==> query must be equal to " + expectedSql)
	}

	//======================
	//==== testing MySql
	//======================
	metaTable = table.getMetaTable(databaseprovider.MySql, "information_schema")
	constr.Init(constrainttype.NotNull, metaTable)
	constr.setField(metaTable.GetFieldByName("reference_id"))
	expectedSql = "ALTER TABLE information_schema.`@meta` MODIFY reference_id INT(11) NOT NULL"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlNotNull() ==> query must be equal to " + expectedSql)
	}
}

func Test__Constraint__getDdlCheck(t *testing.T) {
	table := new(Table)
	//======================
	//==== testing PostgreSql
	//======================
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	constr := new(constraint)
	constr.Init(constrainttype.Check, metaTable)
	constr.setField(metaTable.GetFieldByName("object_type"))
	expectedSql := "ALTER TABLE information_schema.\"@meta\" ADD CONSTRAINT \"ck_@meta_1019\" CHECK (object_type BETWEEN -128 AND 127)"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlCheck() ==> query must be equal to " + expectedSql)
	}
}

func Test__Constraint__getPhysicalName(t *testing.T) {

}
