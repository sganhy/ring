package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/sourcetype"
	"sort"
	"strings"
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

const (
	metaSchemaName        string = "@meta"
	metaSchemaDescription string = "@meta"
	postgreSqlSchema      string = "rpg_sheet_test"
)

var (
	createSchemaSql string = "%s %s %s"
)

func (schema *Schema) Init(id int32, name string, physicalName string, description string, connectionString string, language Language, tables []Table,
	tableSpaces []Tablespace, provider databaseprovider.DatabaseProvider, minConnection uint16, maxConnection uint16, baseline bool,
	active bool, disablePool bool) {

	logger := new(log)
	schema.id = id // first assign the id !!!!
	schema.poolInitialized = false
	schema.name = name
	schema.physicalName = physicalName
	schema.description = description
	schema.source = sourcetype.NativeDataBase // default value
	schema.logger = logger
	schema.sequences = make([]*Sequence, 0, 1)
	schema.connections = new(connectionPool)

	if disablePool == false {
		// load sequences ==> before log init
		logger.Init(id, false)
		schema.connections.Init(id, connectionString, provider, minConnection, maxConnection)
	} else {
		// load sequences ==> before log init
		logger.Init(id, true)
		// instanciate connectionPool without database connections
		schema.connections.connectionString = connectionString
		schema.connections.provider = provider
		schema.connections.poolId = -1 // disable connection pool
	}
	schema.poolInitialized = true
	schema.loadTables(tables)
	schema.loadTablespaces(tableSpaces)
	schema.loadSequences()
	schema.baseline = baseline
	schema.active = active
}

//******************************
// getters and setters
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

func (schema *Schema) GetPhysicalName() string {
	return schema.physicalName
}

func (schema *Schema) GetLanguage() *Language {
	return &schema.language
}

func (schema *Schema) GetDatabaseProvider() databaseprovider.DatabaseProvider {
	return schema.connections.provider
}

func (schema *Schema) GetEntityType() entitytype.EntityType {
	return entitytype.Schema
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
func (schema *Schema) GetSequenceByName(name string) *Sequence {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(schema.sequences) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, schema.sequences[indexerMiddle].name)
		if indexerCompare == 0 {
			return schema.sequences[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
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

func (schema *Schema) LogWarn(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.warn(id, 0, messages...)
	}
}
func (schema *Schema) LogInfo(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.info(id, 0, messages...)
	}
}
func (schema *Schema) LogError(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.error(id, 0, messages...)
	}
}
func (schema *Schema) LogDebug(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.debug(id, 0, messages...)
	}
}
func (schema *Schema) LogFatal(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.fatal(id, 0, messages...)
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

func (schema *Schema) GetDdl(statement ddlstatement.DdlStatement) string {
	var query string
	if statement == ddlstatement.Create {
		query = fmt.Sprintf(createSchemaSql, ddlstatement.Create.String(), entitytype.Schema.String(), schema.GetPhysicalName())
	}
	return query
}

//******************************
// private methods
//******************************
func (schema *Schema) create() error {
	var metaQuery = metaQuery{}
	var creationTime = time.Now()
	var logger = schema.logger
	var err error

	metaQuery.Init(schema, nil)
	metaQuery.query = schema.GetDdl(ddlstatement.Create)
	err = metaQuery.create()

	if err != nil {
		logger.error(-1, 0, err)
		panic(err)
	}

	duration := time.Now().Sub(creationTime)
	logger.info(16, 0, "Create "+entitytype.Schema.String(), fmt.Sprintf("name=%s; execution_time=%d (ms)",
		schema.physicalName, int(duration.Seconds()*1000)))
	return nil
}

func (schema *Schema) exists() bool {
	cata := new(catalogue)
	return cata.exists(schema, schema)
}

func (schema *Schema) findTablespace(table *Table, index *Index, constr *constraint) *Tablespace {
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
		schema.tables[table.GetName()] = &table
		schema.tablesById[table.id] = &table
	}
}

func (schema *Schema) loadTablespaces(tablespaces []Tablespace) {
	schema.tableSpaces = make([]*Tablespace, 0, len(tablespaces))
	for i := 0; i < len(tablespaces); i++ {
		schema.tableSpaces = append(schema.tableSpaces, &tablespaces[i])
	}
}

func (schema *Schema) loadSequences() {
	var seq = Sequence{}
	if schema.name == metaSchemaName {
		schema.sequences = append(schema.sequences, seq.getLexiconId(schema.id))
		schema.sequences = append(schema.sequences, seq.getLanguageId(schema.id))
		schema.sequences = append(schema.sequences, seq.getUserId(schema.id))
		schema.sequences = append(schema.sequences, seq.getIndexId(schema.id))
		schema.sequences = append(schema.sequences, seq.getEventId(schema.id))
	}
	schema.sequences = append(schema.sequences, seq.getJobId(schema.id))
	// sort sequences
	sort.Slice(schema.sequences, func(i, j int) bool {
		return schema.sequences[i].name < schema.sequences[j].name
	})
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

	return result
}

func (schema *Schema) getJobIdValue() int64 {
	var jobIdSequence = schema.GetSequenceByName(sequenceJobIdName)
	if jobIdSequence != nil {
		return jobIdSequence.value.CurrentId
	}
	return -1
}
