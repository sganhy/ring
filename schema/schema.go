package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/sourcetype"
	"time"
)

type Schema struct {
	id              int32
	name            string
	physicalName    string
	description     string
	language        Language
	tables          map[string]*Table
	tablesById      map[int32]*Table
	tableSpaces     []*Tablespace
	sequences       []*Sequence
	connections     *connectionPool
	source          sourcetype.SourceType
	logger          *log
	poolInitialized bool // connection pool initialized
	baseline        bool
	active          bool
}

const metaSchemaName string = "@meta"
const metaSchemaDescription string = "@meta"
const postgreSqlSchema string = "rpg_sheet"

var currentJobId int64 = 101007 // min job id

func (schema *Schema) Init(id int32, name string, physicalName string, description string, connectionString string, language Language, tables []Table,
	tableSpaces []Tablespace, provider databaseprovider.DatabaseProvider, minConnection uint16, maxConnection uint16, baseline bool,
	active bool, disablePool bool) {

	logger := new(log)
	schema.id = id
	schema.poolInitialized = false
	schema.name = name
	schema.physicalName = physicalName
	schema.description = description
	schema.source = sourcetype.NativeDataBase // default value
	schema.logger = logger
	schema.sequences = make([]*Sequence, 0, 1)

	if disablePool == false {
		logger.Init(id, false)
		connectionPool, err := newConnectionPool(id, connectionString, provider, minConnection, maxConnection)
		if connectionPool != nil {
			schema.connections = connectionPool
		} else {
			panic(err)
		}
	} else {
		logger.Init(id, true)
		schema.connections = new(connectionPool)
		schema.connections.connectionString = connectionString
		schema.connections.provider = provider
		schema.connections.poolId = -1 // disable connection pool
	}
	schema.poolInitialized = true
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
	newSchema.Init(schema.id, schema.name, schema.physicalName, schema.description, schema.GetConnectionString(), schema.language, tables, tableSpaces,
		schema.connections.provider, uint16(schema.connections.minConnection), uint16(schema.connections.maxConnection),
		schema.baseline, schema.active, disabledPool)
	return newSchema
}

func (schema *Schema) LogWarn(id int32, jobId int64, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.warn(id, jobId, messages...)
	}
}
func (schema *Schema) LogInfo(id int32, jobId int64, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.info(id, jobId, messages...)
	}
}
func (schema *Schema) LogError(id int32, jobId int64, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.error(id, jobId, messages...)
	}
}
func (schema *Schema) LogDebug(id int32, jobId int64, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.debug(id, jobId, messages...)
	}
}
func (schema *Schema) LogFatal(id int32, jobId int64, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.fatal(id, jobId, messages...)
	}
}

func (schema *Schema) Execute(queries []Query) error {
	var connection = schema.connections.get()
	var err error

	connection.lastGet = time.Now()

	for i := 0; i < len(queries); i++ {
		err = queries[i].Execute(connection.dbConnection)
		if err != nil {
			schema.connections.put(connection)
			return err
		}
	}
	schema.connections.put(connection)
	duration := time.Now().Sub(connection.lastGet)
	fmt.Println("Execution Time:")
	fmt.Println(duration.Milliseconds())
	return nil
}

//******************************
// private methods
//******************************
func (schema *Schema) findTablespace(table *Table, index *Index) *Tablespace {
	result := new(Tablespace)
	result.name = "rpg_data"
	return result
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

func (schema *Schema) getPhysicalName(provider databaseprovider.DatabaseProvider, name string) string {
	//
	if name == metaSchemaName {
		return postgreSqlSchema
	}
	return name
}

func (schema *Schema) getMetaSchema(provider databaseprovider.DatabaseProvider, connectionstring string, minConnection uint16, maxConnection uint16, disablePool bool) *Schema {
	var table = new(Table)
	var tables []Table
	var tablespaces []Tablespace
	var result = new(Schema)
	var language = Language{}
	var sequence = new(Sequence)
	var physicalName = result.getPhysicalName(provider, metaSchemaName)

	var metaTable = table.getMetaTable(provider, physicalName)
	var metaIdTable = table.getMetaIdTable(provider, physicalName)
	var metaLogTable = table.getLogTable(provider, physicalName)

	language.Init("EN")

	//TODO meta schema name hardcoded ("information_schema")
	tables = append(tables, *metaTable)
	tables = append(tables, *metaIdTable)
	tables = append(tables, *metaLogTable)

	// schema.Init(212, "test", "test", "test", language, tables, tablespaces, databaseprovider.Influx, true, true)
	result.Init(0, metaSchemaName, physicalName, metaSchemaDescription, connectionstring, language, tables, tablespaces, provider, minConnection, maxConnection, true, true,
		disablePool)

	// add sequences
	result.sequences = append(result.sequences, sequence.getJobId())

	return result
}

func getCurrentJobId() int64 {
	return currentJobId
}
