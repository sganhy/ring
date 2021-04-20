package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"runtime"
	"strings"
)

type database struct {
	name  string
	index int
}

const (
	schemaSeparator   string = "."
	initialSliceCount int    = 16
)

// sorted by id
var (
	schemaById          *[]*Schema   // assign firstly --> sorted by Id
	schemaIndexByName   *[]*database // assign secondly  --> sorted by name
	defaultSchemaName   string
	databaseInitialized = false
)

func init() {
	lstById := make([]*Schema, 0, initialSliceCount)
	schemaById = &lstById
	lstByName := make([]*database, 0, initialSliceCount)
	schemaIndexByName = &lstByName
}

//******************************
// getters
//******************************

//******************************
// public methods
//******************************

func Init(provider databaseprovider.DatabaseProvider, connectionString string, minConnection uint16, maxConnection uint16) {
	// perform just once
	schema := new(Schema)

	if databaseInitialized == false {
		databaseInitialized = true
		var disableConnectionPool = false
		// disable connection pool for unit testing ??
		if minConnection == 0 && maxConnection == 0 && connectionString == "" {
			disableConnectionPool = true
		}
		// 1> instanciate meta schema
		var metaSchema = schema.getMetaSchema(provider, connectionString, minConnection, maxConnection, disableConnectionPool)
		defaultSchemaName = metaSchema.name

		// 2> add meta schema to collection
		addSchema(metaSchema)

		// 3> initialize logger
		initLogger(metaSchema, metaSchema.GetTableByName(metaLogTableName))

		// 4> initialize cache Id
		InitCacheId(metaSchema, metaSchema.GetTableByName(metaIdTableName))

		// 5> load other schemas if connection pool is not disable
		if disableConnectionPool == false {
			// generate meta tables first before getSchemaIdList()
			createMetaTables(metaSchema)

			// generate meta sequences
			createMetaSequences(metaSchema)

			var schemas = getSchemaIdList()
			for i := 0; i < len(schemas); i++ {
				loadSchemaById(schemas[i])
			}
		}
		// call garbage collector
		runtime.GC()
	}
}

func GetDefaultSchema() *Schema {
	return GetSchemaByName(metaSchemaName)
}

func GetTableBySchemaName(recordType string) *Table {
	var index = strings.Index(recordType, schemaSeparator)
	var schemaName = defaultSchemaName
	var tableName = recordType
	if index >= 0 {
		schemaName = strings.ToUpper(recordType[:index])
		tableName = recordType[index+1:]
	}
	var schema = GetSchemaByName(schemaName)
	if schema != nil {
		return schema.GetTableByName(tableName)
	}
	return nil
}

func GetSchemaByName(name string) *Schema {
	var currentSchemaById = schemaById
	var schemaName = strings.ToUpper(name)
	var currentSchemaByName = schemaIndexByName
	var index = -1
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(*currentSchemaByName) - 1, 0, 0

	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(schemaName, (*currentSchemaByName)[indexerMiddle].name)
		if indexerCompare == 0 {
			index = indexerMiddle
			break
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	if index >= 0 && index < len(*currentSchemaById) {
		return (*currentSchemaById)[index]
	}
	return nil
}

func GetSchemaById(id int32) *Schema {
	var currentSchemaById = schemaById
	var indexerLeft, indexerRight, indexerMiddle = 0, len(*currentSchemaById) - 1, 0
	var indexerCompare int32 = 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = id - (*currentSchemaById)[indexerMiddle].id
		if indexerCompare == 0 {
			return (*currentSchemaById)[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return nil
}

//******************************
// private methods
//******************************
// not thread safe !! slow!! Used only during initialization
func addSchema(schema *Schema) {
	metaDb := new(database)
	metaDb.index = len(*schemaById)
	metaDb.name = strings.ToUpper(schema.name)
	*schemaIndexByName = append(*schemaIndexByName, metaDb)
	*schemaById = append(*schemaById, schema)
	var currentDb *database
	var currentSch *Schema
	// insert to the right place
	for i := len(*schemaIndexByName) - 1; i > 0; i-- {
		currentDb = (*schemaIndexByName)[i]
		if currentDb.name < (*schemaIndexByName)[i-1].name {
			// then swap
			(*schemaIndexByName)[i] = (*schemaIndexByName)[i-1]
			(*schemaIndexByName)[i-1] = currentDb
		} else {
			break
		}
	}
	for i := len(*schemaById) - 1; i > 0; i-- {
		currentSch = (*schemaById)[i]
		if currentSch.id < (*schemaById)[i-1].id {
			// then swap
			(*schemaById)[i] = (*schemaById)[i-1]
			(*schemaById)[i-1] = currentSch
		} else {
			break
		}
	}
}

func createMetaTables(schema *Schema) {
	// first create log table
	logTable := schema.GetTableByName(metaLogTableName)
	if logTable.exists(schema) == false {
		err := logTable.create(schema)
		if err != nil {
			panic(err)
		}
	}
	for _, table := range schema.tables {
		if table.id != logTable.id && table.exists(schema) == false {
			err := table.create(schema)
			if err != nil {
				panic(err)
			}
		}
	}
}

func createMetaSequences(schema *Schema) {
	for _, sequence := range schema.sequences {
		if sequence.exists(schema) == false {
			sequence.create(schema)
		}
		if sequence.value.exists(entitytype.Sequence, sequence.id, schema.id) == false {
			sequence.value.create(entitytype.Sequence, sequence.id, schema.id)
		}
	}
}

// get schema list from @meta table
func getSchemaIdList() []Schema {
	var query = metaQuery{}
	var result []Schema

	// generate meta query
	query.setTable(metaTableName)
	query.addFilter(metaObjectType, operatorEqual, int8(entitytype.Schema))
	err := query.run()

	if err != nil {
		panic(err)
	}
	var metaList = query.getMetaList()
	var count = len(metaList)
	result = make([]Schema, count, count)

	// O(log n)
	for i := 0; i < count; i++ {
		result[i] = *metaList[i].ToSchema()
	}
	return result
}

// load schema from @meta table
func loadSchemaById(schema Schema) {
	var metaList = getMetaList(schema.id)
	var metaIdList = getMetaIdList(schema.id)
	var tables = getTables(schema, metaList) //from: meta.go

	//var schema = new(Schema)

	//schema.Init(212, "test", "test", "", language, tables, tablespaces, databaseprovider.Influx, 0, 0, true, true, true)
	fmt.Println(len(metaList))
	fmt.Println(len(metaIdList))
	fmt.Println(len(tables))
}

// load meta from db @meta table
func getMetaList(schemaId int32) []Meta {
	var query = metaQuery{}

	query.setTable(metaTableName)
	query.addFilter(metaSchemaId, operatorEqual, schemaId)
	err := query.run()
	if err != nil {
		panic(err)
	}
	return query.getMetaList()
}

// load metaI from db @meta_id table
func getMetaIdList(schemaId int32) []MetaId {
	var query = metaQuery{}

	query.setTable(metaIdTableName)
	query.addFilter(metaSchemaId, operatorEqual, schemaId)
	err := query.run()
	if err != nil {
		panic(err)
	}
	return query.getMetaIdList()
}
