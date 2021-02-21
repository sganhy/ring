package schema

import (
	"ring/schema/databaseprovider"
	"strings"
)

// sorted by id
var schemaById *[]*Schema          // assign firstly --> sorted by Id
var schemaIndexByName *[]*database // assign secondly  --> sorted by name
var metaSchemaName string
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
		metaSchemaName = metaSchema.name
		defaultSchemaName = metaSchema.name
		//
		addSchema(metaSchema)
		databaseInitialized = true
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
