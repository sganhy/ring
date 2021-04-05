package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"testing"
)

// INIT
func Test__Schema__Init(t *testing.T) {
	var tables = []Table{}
	var language = Language{}
	var tablespaces = []Tablespace{}
	var schema = Schema{}

	schema.Init(212, "test name", "test physical name", "test desc", "test connectionString", language, tables, tablespaces, databaseprovider.Influx, 0, 0, true, true, true)

	if schema.GetName() != "test name" {
		t.Errorf("Schema.Init() ==> name <> GetName()")
	}
	if schema.GetId() != 212 {
		t.Errorf("Schema.Init() ==> id <> GetId()")
	}
	if schema.GetDescription() != "test desc" {
		t.Errorf("Schema.Init() ==> description <> GetDescription()")
	}
	if schema.GetConnectionString() != "test connectionString" {
		t.Errorf("Schema.Init() ==> connection string <> GetConnectionString()")
	}
	if schema.IsActive() != true {
		t.Errorf("Schema.Init() ==> IsActive() <> true")
	}
	if schema.IsBaseline() != true {
		t.Errorf("Schema.Init() ==> IsBaseline() <> true")
	}

}

// test GetTableByName, GetTableById
func Test__Schema__GetTableByName(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}
	var tables = []Table{}
	var tablespaces = []Tablespace{}
	var tablespace = Tablespace{}
	var schema = Schema{}
	var language = Language{}
	const TABLE_COUNT = 20000

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)

	tablespace.Init(55, "test", "test", "azezae", true, true)
	tablespaces = append(tablespaces, tablespace)

	elemt := Table{}
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes,
		physicaltype.Table, 64, "", tabletype.Lexicon, databaseprovider.NotDefined,
		"subject test", true, false, true, false)
	for i := -100; i <= TABLE_COUNT; i++ {
		table := elemt.Clone()
		table.id = int32(i)
		nameLenght := (abs(i) % 30) + 2
		// fixture
		table.name = randStringBytes(nameLenght)
		tables = append(tables, *table)
	}

	/*
		id int32, name string, description string, connectionString string, language Language, tables []Table,
		tablespaces []Tablespace, provider databaseprovider.DatabaseProvider, baseline bool, active bool
	*/
	schema.Init(212, "test", "test", "phys test", "", language, tables, tablespaces, databaseprovider.Influx, 0, 0, true, true, true)
	// GetTableByName()
	for i := 0; i < len(tables); i++ {
		tableName := tables[i].name
		val := schema.GetTableByName(tableName)
		if val == nil {
			t.Errorf("Schema.GetTableByName() ==> tables[i].name; i=%d, name=%s, id=%d", i, tableName, tables[i].id)
			break
		} else {
			if tableName != val.GetName() {
				t.Errorf("Schema.GetTableByName() ==> tables[i].name; i=%d, name=%s, found=%s", i, tableName, val.GetName())
				break
			}
		}
	}

	// GetTableById()
	for i := 0; i < len(tables); i++ {
		tableId := tables[i].id
		val := schema.GetTableById(tableId)

		if val == nil {
			t.Errorf("Schema.GetTableById() ==> tables[i].name; i=%d, id=%d", i, int(tableId))
			break
		}
	}

	// test nil
	val := schema.GetTableById(-444444444)
	if val != nil {
		t.Errorf("Schema.GetTableById() ==> table id -444444444 cannot be found")
	}
	// test nil
	val = schema.GetTableByName("77777777")
	if val != nil {
		t.Errorf("Schema.GetTableById() ==> table name '77777777' cannot be found")
	}
}
