package schema

import (
	"context"
	"database/sql"
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/relationtype"
	"ring/schema/sourcetype"
	"ring/schema/sqlfmt"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Schema struct {
	id              int32
	name            string
	physicalName    string
	description     string
	language        Language // default language
	tables          map[string]*Table
	tablesById      map[int32]*Table
	tablespaces     []*tablespace
	sequences       []*Sequence
	parameters      []*parameter
	connections     *connectionPool
	source          sourcetype.SourceType
	logger          *log
	poolInitialized bool // connection pool initialized
	baseline        bool
	active          bool
}

const (
	metaSchemaName        string = "@meta"
	metaSchemaID          int32  = 0
	metaSchemaDescription string = "@meta"
	postgreSqlSchema      string = "rpg_sheet_test"
	emptySchemaName       string = ""
)

var (
	createSchemaSql string = "%s %s %s"
)

func (schema *Schema) Init(id int32, name string, physicalName string, description string, connectionString string,
	tables []*Table, tableSpaces []tablespace, sequences []Sequence, parameters []parameter, provider databaseprovider.DatabaseProvider,
	minConnection uint16, maxConnection uint16, baseline bool, active bool, disablePool bool) {

	logger := new(log)
	schema.id = id // first assign the id !!!!
	schema.poolInitialized = false
	schema.name = name
	schema.physicalName = physicalName
	schema.description = description
	schema.source = sourcetype.NativeDataBase // default value
	schema.logger = logger
	schema.sequences = make([]*Sequence, 0, 1)
	schema.parameters = make([]*parameter, 0, 1)
	schema.connections = new(connectionPool)

	if disablePool == false {
		// load sequences ==> before log init
		logger.Init(id, 0, false)
		schema.connections.Init(id, connectionString, provider, minConnection, maxConnection)
	} else {
		// load sequences ==> before log init
		logger.Init(id, 0, true)
		// instanciate connectionPool without database connections
		schema.connections.setConnectionString(connectionString)
		schema.connections.setDatabaseProvider(provider)
		schema.connections.setDisabled() // disable connection pool
	}

	schema.poolInitialized = true
	schema.loadTables(tables)
	schema.loadTablespaces(tableSpaces)
	schema.baseline = baseline
	schema.active = active

	// !!! at end only sequences !!!
	schema.loadSequences(sequences)
	schema.loadParameters(parameters)
	schema.loadMtmTables()

	// after connection pool initialization
	// postpone the cacheId for @meta schema, @meta_id may be doesn't exist yet
	if id != metaSchemaID {
		schema.loadCacheid()
	}

	schema.language = Language{}
	var paramLanguage = schema.getParameterByName(parameterDefaultLanguage)
	if paramLanguage != nil {
		//fmt.Println(paramLanguage.GetValue())
		schema.language.Init(paramLanguage.GetValue())
	}

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
	return schema.connections.GetConnectionString()
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
	return schema.connections.GetDatabaseProvider()
}
func (schema *Schema) GetEntityType() entitytype.EntityType {
	return entitytype.Schema
}
func (schema *Schema) setId(id int32) {
	schema.id = id
}
func (schema *Schema) setName(name string) {
	schema.name = name
}
func (schema *Schema) IsEmpty() bool {
	return schema.name == emptySchemaName
}
func (schema *Schema) getLogger() *log {
	return schema.logger
}
func (schema *Schema) logStatement(statment ddlstatement.DdlStatement) bool {
	return schema.id != metaSchemaID
}
func (schema *Schema) ToMeta() []*meta {
	return schema.toMeta()
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
	var tables []*Table
	var tableSpaces []tablespace
	var sequences []Sequence
	var parameters []parameter

	for _, v := range schema.tables {
		var table = (*v).Clone()
		tables = append(tables, table)
	}
	for i := 0; i < len(schema.tablespaces); i++ {
		var tablespace = *schema.tablespaces[i]
		tableSpaces = append(tableSpaces, *tablespace.Clone())
	}
	for i := 0; i < len(schema.sequences); i++ {
		var sequence = *schema.sequences[i]
		sequences = append(sequences, *sequence.Clone())
	}
	for i := 0; i < len(schema.parameters); i++ {
		var parameter = *schema.parameters[i]
		parameters = append(parameters, *parameter.Clone())
	}
	newSchema.Init(schema.id, schema.name, schema.physicalName, schema.description, schema.GetConnectionString(), tables,
		tableSpaces, sequences, parameters,
		schema.connections.GetDatabaseProvider(), schema.connections.GetMinConnection(),
		schema.connections.GetMaxConnection(), schema.baseline, schema.active, schema.connections.IsDisabled())

	return newSchema
}

//go:noinline
func (schema *Schema) LogWarn(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.writePartialLog(id, levelWarning, 0, messages...)
	}
}

//go:noinline
func (schema *Schema) LogInfo(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.writePartialLog(id, levelInfo, 0, messages...)
	}
}

//go:noinline
func (schema *Schema) LogError(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.writePartialLog(id, levelError, 0, messages...)
	}
}

//go:noinline
func (schema *Schema) LogQueryError(id int32, err error, query string, parameters []string) {
	if schema.logger != nil {
		var description = schema.logger.QueryError(schema.GetDatabaseProvider(), err, query, parameters)
		schema.logger.writePartialLog(id, levelError, 0, err.Error(), description)
	}
}

//go:noinline
func (schema *Schema) LogDebug(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.writePartialLog(id, levelDebug, 0, messages...)
	}
}

//go:noinline
func (schema *Schema) LogFatal(id int32, messages ...interface{}) {
	if schema.logger != nil {
		schema.logger.writePartialLog(id, levelFatal, 0, messages...)
	}
}

func (schema *Schema) Execute(queries []Query, transaction bool) error {
	var connection = schema.connections.get()
	var err error

	connection.lastGet = time.Now()

	if transaction {
		var trans *sql.Tx

		ctx := context.Background()
		trans, err = connection.dbConnection.BeginTx(ctx, nil)
		if err != nil {
			schema.LogError(101, err)
			return err
		}
		for j := 0; j < len(queries); j++ {
			err = queries[j].Execute(connection.dbConnection, trans)
			if err != nil {
				schema.connections.put(connection)
				trans.Rollback()
				return err
			}
		}
		err = trans.Commit()
	} else {
		for i := 0; i < len(queries); i++ {
			err = queries[i].Execute(connection.dbConnection, nil)
			if err != nil {
				schema.connections.put(connection)
				return err
			}
		}
	}

	schema.connections.put(connection)
	duration := time.Now().Sub(connection.lastGet)
	fmt.Println("Execution Time:")
	fmt.Println(duration.Milliseconds())

	return err
}

func (schema *Schema) ExecuteTransaction(queries []Query) error {

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
func (schema *Schema) getParameterByName(name string) *parameter {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(schema.parameters) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, schema.parameters[indexerMiddle].name)
		if indexerCompare == 0 {
			return schema.parameters[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return nil
}

func (schema *Schema) loadCacheid() {
	if schema.connections.IsDisabled() {
		return
	}
	var metaSchema *Schema
	if schema.id == metaSchemaID {
		metaSchema = schema
	} else {
		metaSchema = GetSchemaByName(metaSchemaName)
	}
	// sequences
	for i := 0; i < len(schema.sequences); i++ {
		sequence := schema.sequences[i]
		cacheId := sequence.getCacheId()
		cacheId.Init(metaSchema, schema.id, sequence)
		cacheId.setReservedRange(0)
	}
	// tables
	if schema.id != metaSchemaID {
		for _, table := range schema.tables {
			cacheId := table.getCacheId()
			cacheId.Init(metaSchema, schema.id, table)
			// by default cache id - add parameter
			cacheId.setReservedRange(1)
		}
	}
}

// copy of Execute() method (used only by metaQuery)
func (schema *Schema) execute(query *metaQuery) error {
	var conn = schema.connections.get()
	var err error

	conn.lastGet = time.Now()
	err = query.Execute(conn.dbConnection, nil)
	if err != nil {
		schema.connections.put(conn)
		return err
	}
	schema.connections.put(conn)
	duration := time.Now().Sub(conn.lastGet)
	fmt.Println("Execution Time:")
	fmt.Println(duration.Milliseconds())
	return nil
}

// create physical schema
func (schema *Schema) create(jobId int64) error {
	var metaQuery = metaQuery{}
	var eventId int32 = 16
	var err error

	metaQuery.Init(schema, nil)
	metaQuery.query = schema.GetDdl(ddlstatement.Create)
	err = metaQuery.create(eventId, jobId, schema)

	if err != nil {
		panic(err)
	}

	return nil
}

func (schema *Schema) exists() bool {
	cata := new(catalogue)
	return cata.exists(schema, schema)
}

func (currentSchema *Schema) upgrade(jobId int64, newSchema *Schema) {
	// borrow connection pool to meta schema
	metaSchema := GetSchemaByName(metaSchemaName)
	newSchema.connections = metaSchema.connections

	// create physical schema
	if !newSchema.exists() {
		// initialize logger
		newSchema.logger.Init(newSchema.id, jobId, false)
		newSchema.create(jobId)
	}

	newSchema.createTablespaces(jobId)

	newDico := newSchema.getTableDictionary()
	currDico := currentSchema.getTableDictionary()

	currentSchema.dropTables(jobId, currDico, newDico)
	newSchema.createTables(jobId, currDico, newDico)
	currentSchema.alterTables(jobId, newSchema, newDico)

	newSchema.connections = nil
}

func (schema *Schema) findTablespace(table *Table, index *Index, constr *constraint) *tablespace {
	for i := 0; i < len(schema.tablespaces); i++ {
		tablespace := schema.tablespaces[i]
		if table != nil && tablespace.table {
			return tablespace
		}
		if index != nil && tablespace.index {
			return tablespace
		}
		if constr != nil && tablespace.constraint {
			return tablespace
		}
	}
	// not found tablespace for constraints  use index one
	// recursive call
	if constr != nil && len(schema.tablespaces) > 0 {
		index := new(Index)
		return schema.findTablespace(nil, index, nil)
	}
	return nil
}

func (newSchema *Schema) createTables(jobId int64, prevDico map[string]string, newDico map[string]string) {
	// 1> create missing tables
	for tablePhysName, tableName := range newDico {
		if _, ok := prevDico[strings.ToUpper(tablePhysName)]; !ok {
			table := newSchema.GetTableByName(tableName)
			if table.exists() == false {
				table.create(jobId)
			}
		}
	}
	// 2> create constraints
	for tablePhysName, tableName := range newDico {
		if _, ok := prevDico[tablePhysName]; !ok {
			table := newSchema.GetTableByName(tableName)
			table.createConstraints(jobId, newSchema)
		}
	}
	newSchema.createMtmTables(jobId, prevDico, newDico)
}

func (currentSchema *Schema) dropTables(jobId int64, prevDico map[string]string, newDico map[string]string) {
	// 1> drop tables
	for tablePhysName, tableName := range prevDico {
		if _, ok := newDico[strings.ToUpper(tablePhysName)]; !ok {
			table := currentSchema.GetTableByName(tableName)
			if table.exists() == true {
				table.drop(jobId)
			}
		}
	}
}

func (currentSchema *Schema) alterTables(jobId int64, newSchema *Schema, newDico map[string]string) {
	for _, currTable := range currentSchema.tables {
		tablePhysName := currTable.GetPhysicalName()
		if name, ok := newDico[strings.ToUpper(tablePhysName)]; ok {
			newTable := newSchema.GetTableByName(name)
			if newTable.equal(currTable) == false {
				newTable.alter(jobId, currTable)
			}
		}
	}
}

func (schema *Schema) createTablespaces(jobId int64) error {
	for i := 0; i < len(schema.tablespaces); i++ {
		var tableSpace = schema.tablespaces[i]
		if tableSpace.exists(schema) == false {
			err := tableSpace.create(jobId, schema)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (schema *Schema) createMtmTables(jobId int64, prevDico map[string]string, newDico map[string]string) {
	// 3> create missing mtm relations
	for tablePhysName, tableName := range newDico {
		if _, ok := prevDico[tablePhysName]; !ok {
			table := schema.GetTableByName(tableName)
			for i := 0; i < len(table.relations); i++ {
				relation := table.relations[i]
				if relation.GetType() == relationtype.Mtm && relation.GetMtmTable().exists() == false {
					relation.GetMtmTable().create(jobId)
					relation.GetMtmTable().createConstraints(jobId, schema)
				}
			}
		}
	}
}

// get dictionary of table <Upper(physicalName), name>
func (schema *Schema) getTableDictionary() map[string]string {
	result := make(map[string]string, len(schema.tables))
	for _, tbl := range schema.tables {
		result[strings.ToUpper(tbl.GetPhysicalName())] = tbl.GetName()
	}
	return result
}

func (schema *Schema) loadTables(tables []*Table) {
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
		schema.tables[table.GetName()] = table
		schema.tablesById[table.GetId()] = table
	}
}

func (schema *Schema) loadTablespaces(tablespaces []tablespace) {
	schema.tablespaces = make([]*tablespace, 0, len(tablespaces))
	for i := 0; i < len(tablespaces); i++ {
		schema.tablespaces = append(schema.tablespaces, &tablespaces[i])
	}
	// sort tablespaces by name
	sort.Slice(schema.tablespaces, func(i, j int) bool {
		return schema.tablespaces[i].name < schema.tablespaces[j].name
	})
}

func (schema *Schema) loadSequences(sequences []Sequence) {
	schema.sequences = make([]*Sequence, 0, len(sequences))
	var seq = Sequence{}

	if schema.id == metaSchemaID {
		// initialize cache id before instance sequences
		// meta is ready finally, we need to initialize InitCacheId before sequence instanciation
		//initCacheId(schema, schema.GetTableByName(metaIdTableName), schema.GetTableByName(metaLongTableName))
		schema.sequences = append(schema.sequences, seq.getLexiconId(schema.id))
		schema.sequences = append(schema.sequences, seq.getUserId(schema.id))
		schema.sequences = append(schema.sequences, seq.getIndexId(schema.id))
		schema.sequences = append(schema.sequences, seq.getEventId(schema.id))
		schema.sequences = append(schema.sequences, seq.getJobId(schema.id))
	} else {
		for i := 0; i < len(sequences); i++ {
			var seq = sequences[i]
			schema.sequences = append(schema.sequences, &seq)
		}
	}

	// sort sequences by name
	sort.Slice(schema.sequences, func(i, j int) bool {
		return schema.sequences[i].name < schema.sequences[j].name
	})
}

func (schema *Schema) loadParameters(parameters []parameter) {
	schema.parameters = make([]*parameter, 0, len(parameters))

	for i := 0; i < len(parameters); i++ {
		schema.parameters = append(schema.parameters, &parameters[i])
	}

	// sort parameters
	sort.Slice(schema.parameters, func(i, j int) bool {
		return schema.parameters[i].name < schema.parameters[j].name
	})
}

func (schema *Schema) getPhysicalName(provider databaseprovider.DatabaseProvider, name string) string {
	//
	if len(name) <= 0 {
		return ""
	}
	if name == metaSchemaName {
		return postgreSqlSchema
	}
	var tempName = strings.Trim(name, " ")
	tempName = strings.ReplaceAll(tempName, "  ", " ")
	tempName = strings.ReplaceAll(tempName, "  ", " ")
	var runeArr = []rune(tempName)
	var result strings.Builder
	for i := 0; i < len(runeArr); i++ {
		var chr = runeArr[i]
		if unicode.IsLetter(chr) || unicode.IsDigit(chr) || unicode.IsSpace(chr) || chr == '_' {
			result.WriteRune(chr)
		}
	}
	return strings.ToLower(sqlfmt.ToSnakeCase(result.String()))
}

func (schema *Schema) getEmptySchema() *Schema {
	var tables = []*Table{}
	var tablespaces = []tablespace{}
	var minConnection uint16 = 5
	var maxConnection uint16 = 20
	var sequences = []Sequence{}
	var parameters = []parameter{}
	var result = new(Schema)
	var language = Language{}
	language.Init("en-US")

	var connectionstring string = ""
	var disablePool = true
	var provider = databaseprovider.Undefined

	result.Init(-1, emptySchemaName, "", "", connectionstring, tables,
		tablespaces, sequences, parameters, provider, minConnection, maxConnection, true, true,
		disablePool)

	return result
}

func (schema *Schema) getSchema(schemaId int32, metaList []meta, metaIdList []metaId, disablePool bool) *Schema {
	schemaName, schemaDescription, physicalName, provider := schema.getSchemaInfo(metaList)

	var tables = schema.getTables(provider, physicalName, schemaId, metaList, metaIdList)
	var tablespaces = schema.getTableSpaces(metaList)
	var minConnection uint16 = 5
	var maxConnection uint16 = 20
	var sequences = schema.getSequences(schemaId, metaList)
	var parameters = schema.getParameters(metaList)
	var result = new(Schema)

	// manage pool information
	var connectionstring = schema.getConnectionString(metaList, disablePool)

	// schema.Init(212, "test", "test", "test", language, tables, tablespaces, databaseprovider.Influx, true, true)
	result.Init(schemaId, schemaName, physicalName, schemaDescription, connectionstring, tables,
		tablespaces, sequences, parameters, provider, minConnection, maxConnection, true, true,
		disablePool)

	return result
}

func (schema *Schema) getTables(provider databaseprovider.DatabaseProvider, physicalSchemaName string,
	schemaId int32, metaList []meta, metaIdList []metaId) []*Table {
	var result []*Table
	// map[tableId] *table_meta
	var metaTables map[int32][]*meta
	var metaRefItemCount map[int32]int
	table := Table{}

	metaTables = make(map[int32][]*meta)
	metaRefItemCount = make(map[int32]int)

	// {1} init metaTablesItemCount map
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if val, ok := metaRefItemCount[metaData.refId]; ok {
			metaRefItemCount[metaData.refId] = val + 1
		} else {
			// new item + tableId
			metaRefItemCount[metaData.refId] = 2
		}
	}

	// {2} init metaTables map
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			val := make([]*meta, 0, metaRefItemCount[metaData.id])
			val = append(val, &metaData)
			metaTables[metaData.id] = val
		}
	}

	// {3} build metaTables
	result = make([]*Table, 0, len(metaTables))
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if val, ok := metaTables[metaData.refId]; ok {
			val = append(val, &metaData)
			metaTables[metaData.refId] = val
		}
	}

	// {4} build result
	for _, element := range metaTables {
		result = append(result, table.getTable(provider, physicalSchemaName, schemaId, element))
	}
	schema.loadRelations(result, metaList)
	return result
}

func (schema *Schema) getTableSpaces(metaList []meta) []tablespace {
	var result []tablespace
	result = make([]tablespace, 0, 3)
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if metaData.GetEntityType() == entitytype.Tablespace {
			result = append(result, *metaData.toTablespace())
		}
	}
	return result
}

func (schema *Schema) getDefaultLanguage(metaList []meta) Language {
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if metaData.GetEntityType() == entitytype.Language {
			return *metaData.toLanguage()
		}
	}

	// get language from @meta schema
	var metaSchema = GetSchemaByName(metaSchemaName)
	return metaSchema.language
}

func (schema *Schema) getSequences(schemaId int32, metaList []meta) []Sequence {
	var result []Sequence
	result = make([]Sequence, 0, 3)
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if metaData.GetEntityType() == entitytype.Sequence {
			result = append(result, *metaData.toSequence(schemaId))
		}
	}
	return result
}

func (schema *Schema) getParameters(metaList []meta) []parameter {
	var result []parameter
	result = make([]parameter, 0, 3)
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if metaData.GetEntityType() == entitytype.Parameter {
			result = append(result, *metaData.toParameter(schema.id))
		}
	}
	return result
}

func (schema *Schema) loadRelations(tables []*Table, metaList []meta) {
	// build map of table
	var tableDico map[int32]*Table

	tableDico = make(map[int32]*Table, len(tables))
	for i := 0; i < len(tables); i++ {
		table := tables[i]
		tableDico[table.id] = table
	}
	// load toTable, and inverseRelation
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Relation {
			var fromTable = tableDico[metaData.refId]
			var toTable = tableDico[metaData.dataType]
			var relation = fromTable.GetRelationByName(metaData.name)
			relation.setToTable(toTable)
			relation.setInverseRelation(toTable.GetRelationByName(metaData.value))
		}
	}
}

func (schema *Schema) getMetaSchema(provider databaseprovider.DatabaseProvider, connectionstring string, minConnection uint16, maxConnection uint16, disablePool bool) *Schema {
	const schemaId int32 = metaSchemaID
	var table = new(Table)
	var tables []*Table
	var tablespaces []tablespace
	var sequences []Sequence
	var parameters []parameter
	var result = new(Schema)
	var physicalName = result.getPhysicalName(provider, metaSchemaName)
	var metaTable = table.getMetaTable(provider, physicalName)
	var metaIdTable = table.getMetaIdTable(provider, physicalName)
	var metaLogTable = table.getLogTable(provider, physicalName)
	var metaLexiconTable = table.getLexiconTable(provider, physicalName)
	var metaLexiconItmTable = table.getLexiconItemTable(provider, physicalName)
	var metaLongTable = table.getLongTable()
	var param = new(parameter)
	var ver = new(version)

	tables = append(tables, metaTable)
	tables = append(tables, metaIdTable)
	tables = append(tables, metaLogTable)
	tables = append(tables, metaLongTable)
	tables = append(tables, metaLexiconTable)
	tables = append(tables, metaLexiconItmTable) // should be created after the lexicon

	parameters = append(parameters, *param.getCreationTimeParameter(schemaId, schemaId, entitytype.Schema))
	parameters = append(parameters, *param.getVersionParameter(schemaId, schemaId, entitytype.Schema, ver.String()))
	parameters = append(parameters, *param.getLastUpgradeParameter(schemaId, schemaId, entitytype.Schema))
	parameters = append(parameters, *param.getLanguageParameter(schemaId, "en-US"))

	// schema.Init(212, "test", "test", "test", language, tables, tablespaces, databaseprovider.Influx, true, true)
	result.Init(schemaId, metaSchemaName, physicalName, metaSchemaDescription, connectionstring, tables,
		tablespaces, sequences, parameters, provider, minConnection, maxConnection, true, true,
		disablePool)

	return result
}

func (schema *Schema) getJobIdNextValue() int64 {
	var jobIdSequence = schema.GetSequenceByName(sequenceJobIdName)
	if jobIdSequence != nil {
		return jobIdSequence.NextValue()
	}
	return -1
}

func (schema *Schema) getSchemaInfo(metaList []meta) (string, string, string, databaseprovider.DatabaseProvider) {
	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		if metaData.GetEntityType() == entitytype.Schema {
			return metaData.name, metaData.description, metaData.value,
				databaseprovider.GetDatabaseProviderById(int(metaData.flags))
		}
	}
	return "", "", "", databaseprovider.Undefined
}

func (schema *Schema) loadMtmTables() {
	mtmDico := schema.getMtmDictionary()

	// load tables
	for _, table := range schema.tables {
		for i := 0; i < len(table.relations); i++ {
			var relation = table.relations[i]
			if relation.GetType() == relationtype.Mtm {
				mtmName := relation.getMtmName(table.GetId())
				relation.setMtmTable(mtmDico[mtmName])
			}
		}
	}
}

func (schema *Schema) getMtmDictionary() map[string]*Table {
	result := make(map[string]*Table, len(schema.tables))

	for _, table := range schema.tables {
		for i := 0; i < len(table.relations); i++ {
			var relation = table.relations[i]
			if relation.GetType() == relationtype.Mtm {
				mtmName := relation.getMtmName(table.GetId())
				if _, ok := result[mtmName]; !ok {
					result[mtmName] = table.getMtmTable(schema, relation, mtmName)
				}
			}
		}
	}

	return result
}

func (schema *Schema) getConnectionString(metaList []meta, disablePool bool) string {
	var result string

	if disablePool == true {
		result = ""
	} else {
		//TODO find into metaList connection string
		// get connection string from metaSchema
		var metaSchema = GetSchemaByName(metaSchemaName)
		result = metaSchema.GetConnectionString()
	}

	return result
}

// save meta schema
func (schema *Schema) toMeta() []*meta {
	var metaList = make([]*meta, 0, 100)
	var language = Language{}

	for _, table := range schema.tables {
		metaList = append(metaList, table.toMeta())
		for _, field := range table.fields {
			metaList = append(metaList, field.toMeta(table.id))
		}
		for _, index := range table.indexes {
			metaList = append(metaList, index.toMeta(table.id))
		}
		for _, relation := range table.relations {
			metaList = append(metaList, relation.toMeta(table.id))
		}
	}

	if schema.id == metaSchemaID {
		for _, language := range language.getList() {
			metaList = append(metaList, language.toMeta())
		}
	}

	return metaList
}
