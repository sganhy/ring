package schema

import (
	"ring/schema/constrainttype"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/tabletype"
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

func Test__Constraint__getDdlCreatePrimaryKey(t *testing.T) {
	table := new(Table)
	//======================
	//==== testing PostgreSql
	//======================
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	constr := new(constraint)
	constr.Init(constrainttype.PrimaryKey, metaTable)
	expectedSql := "ALTER TABLE information_schema.\"@meta\" ADD CONSTRAINT \"pk_@meta\" PRIMARY KEY (id,schema_id,object_type,reference_id)"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlCreatePrimaryKey() ==> query must be equal to " + expectedSql)
	}

	//======================
	//==== testing MySql
	//======================
	metaTable = table.getMetaTable(databaseprovider.MySql, "information_schema")
	constr = new(constraint)
	constr.Init(constrainttype.PrimaryKey, metaTable)
	expectedSql = "ALTER TABLE information_schema.`@meta` ADD CONSTRAINT `pk_@meta` PRIMARY KEY (id,schema_id,object_type,reference_id)"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlCreatePrimaryKey() ==> query must be equal to " + expectedSql)
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

func Test__Constraint__getDdlCreateForeignKey(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}
	constr := new(constraint)

	//TABLE_A
	relations = make([]Relation, 1, 1)
	elemr3 := Relation{}
	elemr3.Init(25, "relation_to_B", "[description]", nil, relationtype.Mtm, true, false, true, false)
	relations[0] = elemr3

	elemt01 := Table{}
	elemt01.Init(22, "table_a", "[description]", fields, relations, indexes, physicaltype.Table, 64, "information_schema",
		tabletype.Business, databaseprovider.Oracle, "subject test", true, false, true, false)

	//TABLE 2
	relations = make([]Relation, 1, 1)
	elemr4 := Relation{}
	elemr4.Init(24, "relation_to_A", "[description]", nil, relationtype.Mtm, true, false, true, false)
	relations[0] = elemr4
	elemt02 := Table{}
	elemt02.Init(23, "table_b", "[description]", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business,
		databaseprovider.Oracle, "subject test", true, false, true, false)

	elemt01.relations[0].setToTable(&elemt02)
	elemt02.relations[0].setToTable(&elemt01)
	elemt01.relations[0].setInverseRelation(elemt02.relations[0])
	elemt02.relations[0].setInverseRelation(elemt01.relations[0])

	//======================
	//==== testing PostgreSql
	//======================
	elemt01.setDatabaseProvider(databaseprovider.PostgreSql)
	constr.Init(constrainttype.ForeignKey, &elemt01)
	constr.setRelation(elemt01.relations[0])
	expectedSql := "ALTER TABLE information_schema.t_table_a ADD CONSTRAINT fk_00022_00025 FOREIGN KEY (relation_to_B) REFERENCES t_table_b (id)"
	if constr.GetDdl(ddlstatement.Create, nil) != expectedSql {
		t.Errorf("Constraint.getDdlCreateForeignKey() ==> query must be equal to " + expectedSql)
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
