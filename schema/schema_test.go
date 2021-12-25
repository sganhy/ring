package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"testing"
)

// INIT
func Test__Schema__Init(t *testing.T) {
	var tables = []*Table{}
	var tablespaces = []tablespace{}
	var sequences = []Sequence{}
	var param = parameter{}
	var parameters = []parameter{}
	var schema = Schema{}

	parameters = append(parameters, *param.getLanguageParameter(212, "en"))
	//t.Errorf(parameters[0].String())
	schema.Init(212, "test name", "test physical name", "test desc", "test connectionString", tables,
		tablespaces, sequences, parameters, databaseprovider.Influx, 0, 0, true, true, true)

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
	if schema.GetPhysicalName() != "test physical name" {
		t.Errorf("Schema.Init() ==> physicalName() <> 'test physical name'")
	}
	if schema.GetDatabaseProvider() != databaseprovider.Influx {
		t.Errorf("Schema.Init() ==> provider() <> databaseprovider.Influx")
	}
	if schema.GetEntityType() != entitytype.Schema {
		t.Errorf("Schema.Init() ==> GetEntityType() <> entitytype.Schema")
	}
	if schema.IsEmpty() == true {
		t.Errorf("Schema.Init() ==> IsEmpty() <> false")
	}
	if schema.GetLanguage().GetNativeName() != "English" {
		t.Errorf("Schema.Init() ==> Native language <> 'English'")
	}
	if schema.GetTableCount() != 0 {
		t.Errorf("Schema.Init() ==> GetTableCount() <> 0")
	}
	if schema.logStatement(ddlstatement.Create) != true {
		t.Errorf("Schema.Init() ==> logStatment() <> true")
	}
	//
	var schema2 *Schema
	if schema2.GetTableCount() != 0 {
		t.Errorf("Schema.GetTableCount() ==> GetTableCount() <> 0")
	}
}

// test GetTableByName, GetTableById
func Test__Schema__GetTableByName(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}
	var tables = []*Table{}
	var tablespaces = []tablespace{}
	var sequences = []Sequence{}
	var parameters = []parameter{}
	var tableSpace = tablespace{}

	var schema = Schema{}
	const TABLE_COUNT = 20000

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)

	tableSpace.Init(55, "test", "test", "azezae", true, true)
	tablespaces = append(tablespaces, tableSpace)

	elemt := Table{}
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes,
		physicaltype.Table, 64, "", tabletype.Lexicon, databaseprovider.Undefined,
		"subject test", true, false, true, false)
	for i := -100; i <= TABLE_COUNT; i++ {
		table := elemt.Clone()
		table.setId(int32(i))
		nameLenght := (abs(i) % 30) + 2
		// fixture
		table.setName(randStringBytes(nameLenght))
		tables = append(tables, table)
	}

	/*
		id int32, name string, description string, connectionString string, language Language, tables []Table,
		tablespaces []Tablespace, provider databaseprovider.DatabaseProvider, baseline bool, active bool
	*/
	schema.Init(213, "test", "test", "phys test", "", tables, tablespaces, sequences, parameters,
		databaseprovider.Influx, 0, 0, true, true, true)
	// GetTableByName()
	for i := 0; i < len(tables); i++ {
		tableName := tables[i].GetName()
		val := schema.GetTableByName(tableName)
		if val == nil {
			t.Errorf("Schema.GetTableByName() ==> tables[i].name; i=%d, name=%s, id=%d", i, tableName, tables[i].GetId())
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
		tableId := tables[i].GetId()
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

func Test__Schema__GetSequenceByName(t *testing.T) {
	var schema = new(Schema)
	schema = schema.getMetaSchema(databaseprovider.MySql, "", 0, 0, true)

	for i := 0; i < len(schema.sequences); i++ {
		sequence := schema.GetSequenceByName(schema.sequences[i].GetName())
		if sequence == nil {
			t.Errorf("Schema.GetSequenceByName() ==> sequence name '%s' cannot be found", schema.sequences[i].GetName())
		}
	}

	sequence := schema.GetSequenceByName("????????????")
	if sequence != nil {
		t.Errorf("Schema.GetSequenceByName() ==> sequence name '????????????' must be null")
	}

}

func Test__Schema__getParameterByName(t *testing.T) {
	var schema = new(Schema)
	schema = schema.getMetaSchema(databaseprovider.MySql, "", 0, 0, true)

	for i := 0; i < len(schema.parameters); i++ {
		param := schema.getParameterByName(schema.parameters[i].GetName())
		if param == nil {
			t.Errorf("Schema.getParameterByName() ==> parameter name '%s' cannot be found", schema.sequences[i].GetName())
		}
	}

	param := schema.getParameterByName("????????????")
	if param != nil {
		t.Errorf("Schema.getParameterByName() ==> parameter name '????????????' must be null")
	}

}

func Test__Schema__GetDdl(t *testing.T) {
	var schema = new(Schema)
	schema = schema.getMetaSchema(databaseprovider.MySql, "rpg_test", 0, 0, true)

	if schema.GetDdl(ddlstatement.Create) != "CREATE SCHEMA rpg_sheet_test" {
		t.Errorf("Schema.GetDdl() ==> is not equal to 'CREATE SCHEMA rpg_sheet_test'")
	}
}

func Test__Schema__getPhysicalName(t *testing.T) {
	schema := new(Schema)

	if schema.getPhysicalName(databaseprovider.PostgreSql, metaSchemaName) != postgreSqlSchema {
		t.Errorf("Schema.getPhysicalName() ==> is not equal to '%s'", postgreSqlSchema)
	}
	validPhysicalName := "rpg_sheet"
	if schema.getPhysicalName(databaseprovider.PostgreSql, " Rpg Sheet ") != validPhysicalName {
		t.Errorf("Schema.getPhysicalName(' Rpg Sheet ') ==> is not equal to '%s'", validPhysicalName)
	}
	if schema.getPhysicalName(databaseprovider.PostgreSql, "Rpg     Sheet") != validPhysicalName {
		t.Errorf("Schema.getPhysicalName(' Rpg Sheet ') ==> is not equal to '%s'", validPhysicalName)
	}
	if schema.getPhysicalName(databaseprovider.PostgreSql, "Rpg_Sheet") != validPhysicalName {
		t.Errorf("Schema.getPhysicalName(' Rpg Sheet ') ==> is not equal to '%s'", validPhysicalName)
	}
	if schema.getPhysicalName(databaseprovider.PostgreSql, "Rpg sheet") != validPhysicalName {
		t.Errorf("Schema.getPhysicalName(' Rpg Sheet ') ==> is not equal to '%s'", validPhysicalName)
	}
	validPhysicalName = "" // take default schema
	if schema.getPhysicalName(databaseprovider.PostgreSql, "") != validPhysicalName {
		t.Errorf("Schema.getPhysicalName('') ==> should be equal to null")
	}

}

func Test__Schema__findTablespace(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}
	var tables = []*Table{}
	var tablespaces = []tablespace{}
	var sequences = []Sequence{}
	var parameters = []parameter{}
	var schema = Schema{}
	var uk Index = Index{}

	// creating tablepaces
	tbl01 := new(tablespace)
	tbl01.Init(3333, "indexspace 3", "ATable Test", "/data/data", false, false)
	tbl02 := new(tablespace)
	tbl02.Init(3334, "indexspace 4", "ATable Test", "/data/indexes", true, false)
	tbl03 := new(tablespace)
	tbl03.Init(3335, "indexspace 5", "ATable Test", "/data/indexes", false, true)
	tablespaces = append(tablespaces, *tbl01)
	tablespaces = append(tablespaces, *tbl02)
	tablespaces = append(tablespaces, *tbl03)

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)

	var indexedFields = []string{"Zorba"}
	uk.Init(1, "uk_test", "ATable Test", indexedFields, false, false, true, true)
	indexes = append(indexes, uk)

	elemt := Table{}
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes,
		physicaltype.Table, 64, "", tabletype.Lexicon, databaseprovider.Undefined,
		"subject test", true, false, true, false)
	tables = append(tables, &elemt)

	schema.Init(214, "test", "test", "phys test", "", tables, tablespaces, sequences, parameters,
		databaseprovider.Influx, 0, 0, true, true, true)

	tblspcResult := schema.findTablespace(&elemt, nil, nil)
	if tblspcResult.name != "indexspace 4" {
		t.Errorf("Schema.findTablespace() ==> should be equal to 'indexspace 4'")
	}
	tblspcResult = schema.findTablespace(nil, &uk, nil)
	if tblspcResult.name != "indexspace 5" {
		t.Errorf("Schema.findTablespace() ==> should be equal to 'indexspace 5'")
	}

}
