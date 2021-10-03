package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"testing"
	"time"
)

func Test__Metaquery__getLogAlterDescription(t *testing.T) {
	schema := new(Schema)
	var eventTime = time.Now().Add(10 * time.Minute) // be sure duration == 0
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	table := schema.GetTableByName(metaTableName)
	relation := Relation{}
	var relations = []*Relation{}
	var query = metaQuery{}
	query.Init(schema, table)
	relation.Init(-23, "zorba07", "hellkzae", table, relationtype.Mto, true, false, true, false)
	field := table.GetFieldByName(metaObjectType)
	relations = append(relations, &relation)
	table.relations = relations

	// alter drop field
	operator := field.GetName()
	query.query = table.getDdlAlter(field)
	expectedDescription := "rpg_sheet_test.\"@meta\" drop field object_type (done) | time=0ms"
	if query.getLogAlterDescription(table, eventTime, operator) != expectedDescription {
		t.Errorf("metaQuery.getLogAlterDescription() ==> description is different '%s'", expectedDescription)
	}
	// alter create field
	elemf := Field{}
	elemf.Init(21, "hello_world", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)
	query.query = table.getDdlAlter(&elemf)
	expectedDescription = "rpg_sheet_test.\"@meta\" add field object_type (done) | time=0ms"
	if query.getLogAlterDescription(table, eventTime, operator) != expectedDescription {
		t.Errorf("metaQuery.getLogAlterDescription() ==> description is different than '%s'", expectedDescription)
	}
	// alter create field : getLogDescription
	if query.getLogDescription("", table, eventTime, ddlstatement.Alter, operator) != expectedDescription {
		t.Errorf("metaQuery.getLogDescription() ==> description is different than '%s'", expectedDescription)
	}
	// alter create relation

	expectedDescription = "rpg_sheet_test.\"@meta\" drop relation zorba07 (done) | time=0ms"
	elemf2 := Field{}
	elemf2.setName(relation.GetName())
	query.query = table.getDdlAlter(&elemf2)
	operator = relation.GetName()
	if query.getLogAlterDescription(table, eventTime, operator) != expectedDescription {
		t.Errorf("metaQuery.getLogAlterDescription() ==> description is different than '%s'", expectedDescription)
	}
}

func Test__Metaquery__getLogCreateDescription(t *testing.T) {
	schema := new(Schema)
	var eventTime = time.Now().Add(10 * time.Minute) // be sure duration == 0
	//var eventTime = time.Now()
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	table := schema.GetTableByName(metaIdTableName)
	index := table.GetIndexByName(metaIdTableName)

	var query = metaQuery{}
	query.Init(schema, table)
	query.query = index.GetDdl(ddlstatement.Create, table, nil)
	entityName := index.getPhysicalName(table)
	expectedDescription := "idx_meta_id_001 on rpg_sheet_test.\"@meta_id\" (done) | time=0ms"
	if query.getLogCreateDescription(entityName, index, eventTime) != expectedDescription {
		t.Errorf("metaQuery.getLogCreateDescription() ==> description is different than '%s'", expectedDescription)
	}
	// test via getLogDescription()
	if query.getLogDescription(entityName, index, eventTime, ddlstatement.Create, "") != expectedDescription {
		t.Errorf("metaQuery.getLogDescription() ==> description is different than '%s'", expectedDescription)
	}
	if query.getLogDescription(entityName, index, eventTime, ddlstatement.Drop, "") != expectedDescription {
		t.Errorf("metaQuery.getLogDescription() ==> description is different than '%s'", expectedDescription)
	}

}
