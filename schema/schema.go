package schema

import (
	"ring/schema/databaseprovider"
)

type Schema struct {
	id               int32
	name             string
	description      string
	connectionString string
	language         Language
	tables           map[string]*Table
	tablesById       map[int32]*Table
	tablespaces      []*Tablespace
	provider         databaseprovider.DatabaseProvider
	baseline         bool
	active           bool
}

func (schema *Schema) Init(id int32, name string, description string, connectionString string, language Language, tables []Table,
	tablespaces []Tablespace, provider databaseprovider.DatabaseProvider, baseline bool, active bool) {
	schema.id = id
	schema.name = name
	schema.description = description
	schema.connectionString = connectionString
	schema.provider = provider
	schema.loadTables(tables)
	schema.loadTablespaces(tablespaces)
	schema.baseline = baseline
	schema.active = active
}

//******************************
// getters
//******************************
func (schema *Schema) GetId() int32 {
	return schema.id
}

func (schema *Schema) GetName() string {
	return schema.name
}

func (schema *Schema) GetDescription() string {
	return schema.description
}

func (schema *Schema) GetConnectionString() string {
	return schema.connectionString
}

func (schema *Schema) IsBaseline() bool {
	return schema.baseline
}

func (schema *Schema) IsActive() bool {
	return schema.active
}

//******************************
// public methods
//******************************
func (schema *Schema) GetTableByName(name string) *Table {
	if val, ok := schema.tables[name]; ok {
		return val
	}
	return nil
}
func (schema *Schema) GetTableById(id int32) *Table {
	if val, ok := schema.tablesById[id]; ok {
		return val
	}
	return nil
}

//******************************
// private methods
//******************************
func (schema *Schema) loadTables(tables []Table) {
	tableCount := len(tables) + 1
	// reducing collision then *2
	capacity := tableCount * 2
	// *4 if small number of tables
	if tableCount < 1000 {
		capacity *= 2
	}
	schema.tables = make(map[string]*Table, capacity)
	schema.tablesById = make(map[int32]*Table, capacity)

	for i := 0; i < len(tables); i++ {
		table := tables[i]
		schema.tables[table.name] = &table
		schema.tablesById[table.id] = &table
	}
}

func (schema *Schema) loadTablespaces(tablespaces []Tablespace) {
	schema.tablespaces = make([]*Tablespace, 0, len(tablespaces))
	for i := 0; i < len(tablespaces); i++ {
		schema.tablespaces = append(schema.tablespaces, &tablespaces[i])
	}
}
