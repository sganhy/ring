package schema

import (
	"ring/schema/databaseprovider"
	"sort"
	"strings"
)

// sorted by id
var schemaById *[]*Schema          // assign firstly --> sorted by Id
var schemaIndexByName *[]*database // assign secondly  --> sorted by name
var metaSchemaName string
var defaultSchemaName string

const schemaSeparator = "."
const initialSliceCount = 16
const schemaNotFound = 16

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

func Init(provider databaseprovider.DatabaseProvider, connectionstring string) {
	var metaschema = getMetaSchema(provider, connectionstring)
	metaSchemaName = metaschema.name
	defaultSchemaName = metaschema.name
	//
	addSchema(metaschema)
}

func GetMetaSchemaName() string {
	return metaSchemaName
}

func GetDefaultSchema() *Schema {
	return GetSchemaByName(metaSchemaName)
}

func setDefaultSchema(name string) {
	var schema = GetSchemaByName(name)
	// thread safe?
	if schema != nil {
		defaultSchemaName = schema.name
	}
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
	var indexerLeft, indexerRigth, indexerMiddle, indexerCompare int = 0, len(*currentSchemaByName) - 1, 0, 0

	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(schemaName, (*currentSchemaByName)[indexerMiddle].name)
		if indexerCompare == 0 {
			index = indexerMiddle
			break
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	if index >= 0 && index < len(*currentSchemaById) {
		return (*currentSchemaById)[index]
	}
	return nil
}

func GetSchemaById(id int32) *Schema {
	var currentSchemaById = schemaById
	var indexerLeft, indexerRigth, indexerMiddle int = 0, len(*currentSchemaById) - 1, 0
	var indexerCompare int32 = 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = id - (*currentSchemaById)[indexerMiddle].id
		if indexerCompare == 0 {
			return (*currentSchemaById)[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	return nil
}

// not thread safe !! very slow
func addSchema(schema *Schema) {
	metaDb := new(database)
	metaDb.index = len(*schemaById)
	metaDb.name = strings.ToUpper(schema.name)
	*schemaIndexByName = append(*schemaIndexByName, metaDb)
	*schemaById = append(*schemaById, schema)
	// sort
	sort.Slice(*schemaIndexByName, func(i, j int) bool {
		return (*schemaIndexByName)[i].name < (*schemaIndexByName)[j].name
	})
	sort.Slice(*schemaById, func(i, j int) bool {
		return (*schemaById)[i].id < (*schemaById)[j].id
	})
}

//******************************
// private methods
//******************************
