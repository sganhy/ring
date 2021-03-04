package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/sourcetype"
	"time"
)

type Schema struct {
	id          int32
	name        string
	description string
	language    Language
	tables      map[string]*Table
	tablesById  map[int32]*Table
	tableSpaces []*Tablespace
	connections *connectionPool
	source      sourcetype.SourceType
	baseline    bool
	active      bool
}

const metaSchemaName string = "@meta"
const metaSchemaDescription string = "@meta"
const postgreSqlSchema string = "information_schema"

func (schema *Schema) Init(id int32, name string, description string, connectionString string, language Language, tables []Table,
	tableSpaces []Tablespace, provider databaseprovider.DatabaseProvider, minConnection uint16, maxConnection uint16, baseline bool,
	active bool, disablePool bool) {
	schema.id = id
	schema.name = name
	schema.description = description
	schema.source = sourcetype.NativeDataBase // default value
	if disablePool == false {
		connectionPool, err := newConnectionPool(connectionString, provider, minConnection, maxConnection)
		if connectionPool != nil {
			schema.connections = connectionPool
		} else {
			panic(err)
		}
	} else {
		schema.connections = new(connectionPool)
		schema.connections.connectionString = connectionString
		schema.connections.provider = provider
		schema.connections.poolId = -1 // disable connection pool
	}
	schema.loadTables(tables)
	schema.loadTablespaces(tableSpaces)
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
	return schema.connections.connectionString
}

func (schema *Schema) IsBaseline() bool {
	return schema.baseline
}

func (schema *Schema) IsActive() bool {
	return schema.active
}

func (schema *Schema) GetLanguage() *Language {
	return &schema.language
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
func (schema *Schema) GetTableCount() int {
	if schema != nil {
		return len(schema.tables)
	} else {
		return 0
	}
}

func (schema *Schema) Clone() *Schema {
	newSchema := new(Schema)
	var tables []Table
	var tableSpaces []Tablespace
	var disabledPool = false
	for _, v := range schema.tables {
		var table = (*v).Clone()
		tables = append(tables, *table)
	}
	for i := 0; i < len(schema.tableSpaces); i++ {
		var tablespace = *schema.tableSpaces[i]
		tableSpaces = append(tableSpaces, *tablespace.Clone())
	}
	if schema.connections.poolId == -1 {
		disabledPool = true
	}
	newSchema.Init(schema.id, schema.name, schema.description, schema.GetConnectionString(), schema.language, tables, tableSpaces,
		schema.connections.provider, uint16(schema.connections.minConnection), uint16(schema.connections.maxConnection),
		schema.baseline, schema.active, disabledPool)
	return newSchema
}

func (schema *Schema) Execute(queries []Query) error {
	var connection = schema.connections.get()
	var err error
	connection.lastGet = time.Now()

	for i := 0; i < len(queries); i++ {
		err = queries[i].Execute(schema.connections.provider, connection.dbConnection)
		if err != nil {
			schema.connections.put(connection)
			return err
		}
	}

	schema.connections.put(connection)
	return nil
}

//******************************
// private methods
//******************************

func (schema *Schema) setSourceType(source sourcetype.SourceType) {
	schema.source = source
}

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
	schema.tableSpaces = make([]*Tablespace, 0, len(tablespaces))
	for i := 0; i < len(tablespaces); i++ {
		schema.tableSpaces = append(schema.tableSpaces, &tablespaces[i])
	}
}

func getPhysicalName(provider databaseprovider.DatabaseProvider, name string) string {
	//
	if name == metaSchemaName {
		return postgreSqlSchema
	}
	return name
}

func getMetaSchema(provider databaseprovider.DatabaseProvider, connectionstring string, minConnection uint16, maxConnection uint16, disablePool bool) *Schema {
	var tables []Table
	var tablespaces []Tablespace
	var schema = Schema{}
	var language = Language{}

	language.Init("EN")

	tables = append(tables, *getMetaTable(provider, getPhysicalName(provider, metaSchemaName)))
	tables = append(tables, *getMetaIdTable(provider, getPhysicalName(provider, metaSchemaName)))
	tables = append(tables, *getLogTable(provider, getPhysicalName(provider, metaSchemaName)))

	// schema.Init(212, "test", "test", "test", language, tables, tablespaces, databaseprovider.Influx, true, true)
	schema.Init(0, metaSchemaName, metaSchemaDescription, connectionstring, language, tables, tablespaces, provider, minConnection, maxConnection, true, true,
		disablePool)
	return &schema
}
