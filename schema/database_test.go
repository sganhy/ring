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
	var tables = []Table{}
	var tablespaces = []Tablespace{}
	var tablespace = Tablespace{}
	var schema = Schema{}
	var schemas = []*Schema{}
	var language = Language{}
	Init(databaseprovider.MySql, "zorba")

	const SCHEMA_COUNT = 2000 // 20k is too slow (~50 seconds)

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)

	tablespace.Init(55, "test", "test", "azezae", true, true)
	tablespaces = append(tablespaces, tablespace)

	elemt := Table{}
	elemt.Init(22, "zorro", "hellkzae", fields, relations, indexes,
		"schema.t_site", physicaltype.Table, 64, tabletype.Lexicon, "subject test", true, false, true, false)
	tables = append(tables, elemt)
	schema.Init(212, "test", "test", "test", language, tables, tablespaces, databaseprovider.Influx, true, true)
	for i := -100; i < SCHEMA_COUNT; i++ {
		var newSchema = schema.Clone()
		nameLenght := (abs(i) % 30) + 2
		schemaName := randStringBytes(nameLenght)
		newSchema.id = int32(i)
		newSchema.name = schemaName
		addSchema(newSchema)
		schemas = append(schemas, newSchema)
	}

	schemaName := schemas[SCHEMA_COUNT>>1].name
	if GetSchemaByName(schemaName) == nil {
		t.Errorf("Database.GetSchemaByName() ==> schema name %s not found", schemaName)
	}
	schemaName = schemas[SCHEMA_COUNT>>2].name
	if GetSchemaByName(schemaName) == nil {
		t.Errorf("Database.GetSchemaByName() ==> schema name %s not found", schemaName)
	}
	if GetSchemaByName("11111111111111111111") != nil {
		t.Errorf("Database.GetSchemaByName() ==> schema name %s cannot be not found", "11111111111111111111")
	}
	schemaName += ".zorro"
	if GetTableBySchemaName(schemaName) == nil {
		t.Errorf("Database.GetTableBySchemaName() ==> recordtype %s not found", schemaName)
	}
	schemaId := schemas[SCHEMA_COUNT>>1].id
	if GetSchemaById(schemaId) == nil {
		t.Errorf("Database.GetSchemaById() ==> schema id %d not found", schemaId)
	}
	schemaId = schemas[SCHEMA_COUNT>>2].id
	if GetSchemaById(schemaId) == nil {
		t.Errorf("Database.GetSchemaById() ==> schema id %d not found", schemaId)
	}
	if GetSchemaById(-1455555555) != nil {
		t.Errorf("Database.GetSchemaById() ==> schema id %d cannot be found", -1455555555)
	}

}
