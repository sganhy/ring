package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"testing"
)

//test: GetTableBySchemaName, GetSchemaByName, GetSchemaById
func Test__Database__GetTableBySchemaName(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}
	var tables = []*Table{}
	var tablespaces = []tablespace{}
	var sequences = []Sequence{}
	var parameters = []parameter{}
	var tableSpace = tablespace{}
	var schema = Schema{}
	var schemas = []*Schema{}

	// disable connection pool empty connection, string min & max == 0
	Init(databaseprovider.MySql, "", 0, 0)

	const SCHEMA_COUNT = 8000 // 20k is too slow (~50 seconds)

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)

	tableSpace.Init(55, "test", "test", "azezae", true, true)
	tablespaces = append(tablespaces, tableSpace)

	elemt := Table{}
	elemt.Init(22, "zorro", "hellkzae", fields, relations, indexes,
		physicaltype.Table, 64, "", tabletype.Lexicon, databaseprovider.Undefined, "subject test", true, false, true, false)
	tables = append(tables, &elemt)
	schema.Init(211, "test", "test", "test", "test", tables, tablespaces, sequences, parameters,
		databaseprovider.Influx, 0, 0, true, true, true)
	for i := -100; i < SCHEMA_COUNT; i++ {
		var newSchema = schema.Clone()
		nameLenght := (abs(i) % 30) + 2
		schemaName := randStringBytes(nameLenght)
		newSchema.setId(int32(i))
		newSchema.setName(schemaName)
		addSchema(newSchema)
		schemas = append(schemas, newSchema)
	}
	// search all schema
	for i := 0; i < len(schemas); i++ {
		var schemaName = schemas[i].GetName()
		if GetSchemaByName(schemaName) == nil {
			t.Errorf("Database.GetSchemaByName() ==> schema name %s is missing", schemaName)
		}
	}
	if GetSchemaByName("11111111111111111111") != nil {
		t.Errorf("Database.GetSchemaByName() ==> schema name %s cannot be not found", "11111111111111111111")
	}
	schemaName := schemas[len(schemas)>>1].GetName() + ".zorro"
	if GetTableBySchemaName(schemaName) == nil {
		t.Errorf("Database.GetTableBySchemaName() ==> recordtype %s not found", schemaName)
	}
	if GetTableBySchemaName("?????.?????") != nil {
		t.Errorf("Database.GetTableBySchemaName() ==> recordtype %s cannot be found", "?????.?????")
	}
	schemaId := schemas[SCHEMA_COUNT>>1].GetId()
	if GetSchemaById(schemaId) == nil {
		t.Errorf("Database.GetSchemaById() ==> schema id %d not found", schemaId)
	}
	schemaId = schemas[SCHEMA_COUNT>>2].GetId()
	if GetSchemaById(schemaId) == nil {
		t.Errorf("Database.GetSchemaById() ==> schema id %d not found", schemaId)
	}
	if GetSchemaById(-1455555555) != nil {
		t.Errorf("Database.GetSchemaById() ==> schema id %d cannot be found", -1455555555)
	}

}
