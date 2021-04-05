package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"runtime"
	"strconv"
	"strings"
)

// sorted by id
var schemaById *[]*Schema          // assign firstly --> sorted by Id
var schemaIndexByName *[]*database // assign secondly  --> sorted by name
var defaultSchemaName string
var databaseInitialized = false

const schemaSeparator = "."
const initialSliceCount = 16

type database struct {
	name  string
	index int
}

func init() {
	lstById := make([]*Schema, 0, initialSliceCount)
	schemaById = &lstById
	lstByName := make([]*database, 0, initialSliceCount)
	schemaIndexByName = &lstByName
}

//******************************
// public methods
//******************************

func Init(provider databaseprovider.DatabaseProvider, connectionString string, minConnection uint16, maxConnection uint16) {
	// perform just once
	if databaseInitialized == false {
		var disableConnectionPool = false
		// disable connection pool for unit testing ??
		if minConnection == 0 && maxConnection == 0 && connectionString == "" {
			disableConnectionPool = true
		}
		var metaSchema = getMetaSchema(provider, connectionString, minConnection, maxConnection, disableConnectionPool)
		defaultSchemaName = metaSchema.name

		//
		addSchema(metaSchema)
		databaseInitialized = true

		// unitests context ?
		if disableConnectionPool == false {
			// generate meta tables first before getSchemaIdList()
			generateMetaTables(metaSchema)

			var schemas = getSchemaIdList()
			fmt.Println("schema id ==> ")

			for i := 0; i < len(schemas); i++ {
				loadSchemaById(schemas[i])
			}
			// call garbage collector
		}
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

//******************************
// private methods
//******************************
func generateMetaTables(schema *Schema) {
	for _, table := range schema.tables {
		if table.exists(schema) == false {
			table.create(schema)
		}
	}
}

// get schema list from @meta table
func getSchemaIdList() []Schema {
	var query = metaQuery{}
	var result []Schema

	// generate meta query
	query.setTable(metaTableName)
	query.addFilter(metaObjectType, "=", entitytype.Schema.GetId())
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
	query.addFilter(metaSchemaId, "=", strconv.Itoa(int(schemaId)))
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
	query.addFilter(metaSchemaId, "=", strconv.Itoa(int(schemaId)))
	err := query.run()
	if err != nil {
		panic(err)
	}
	return query.getMetaIdList()
}
