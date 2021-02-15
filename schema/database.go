package schema

import (
	"ring/schema/databaseprovider"
	"strings"
)

//******************************
// getters
//******************************

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************

// sorted by id
var schemaCollection *[]*Schema // assign firstly --> not sorted
//TODO replace by slice
var schemaNameCollection *map[string]int // assign secondly structure -->
var metaSchemaName string
var defaultSchemaName string

const schemaSeparator = "."

func init() {
	schemaCollection = &([]*Schema{})
	var schemasByName = make(map[string]int)
	schemaNameCollection = &schemasByName
}

func Init(provider databaseprovider.DatabaseProvider, connectionstring string) {
	var schemas = []*Schema{}
	var metaschema = GetMetaSchema(provider, connectionstring)
	var schemasByName = make(map[string]int)
	metaSchemaName = metaschema.name

	schemas = append(schemas, metaschema)
	schemaCollection = &schemas
	schemasByName[strings.ToUpper(metaSchemaName)] = 0
	schemaNameCollection = &schemasByName
	defaultSchemaName = metaschema.name
}

func GetMetaSchemaName() string {
	return metaSchemaName
}

func GetDefaultSchema() *Schema {
	return GetSchemaByName(metaSchemaName)
}

func SetDefaultSchema(name string) {
	var currentSchemaCollection = schemaCollection
	var currentSchemaNameCollection = schemaNameCollection

	// thread safe?
	// check if schema name exist
	if schemaId, ok := (*currentSchemaNameCollection)[name]; ok {
		currSchema := (*currentSchemaCollection)[schemaId]
		defaultSchemaName = currSchema.name
	}
}

func GetTableBySchemaName(recordType string) *Table {
	var index = strings.Index(recordType, schemaSeparator)
	if index >= 0 {
		var schemaName = strings.ToUpper(recordType[:index])
		var tableName = recordType[index+1:]
		var currentSchemaCollection = schemaCollection
		var currentSchemaNameCollection = schemaNameCollection

		// find schema id
		if schemaId, ok := (*currentSchemaNameCollection)[schemaName]; ok {
			currSchema := (*currentSchemaCollection)[schemaId]
			return currSchema.GetTableByName(tableName)
		}
		//TODO Init launched before ??
		return nil
	}
	var currentSchemaCollection = schemaCollection
	var currentSchemaNameCollection = schemaNameCollection
	// no schema separator, use default schema
	if schemaId, ok := (*currentSchemaNameCollection)[strings.ToUpper(defaultSchemaName)]; ok {
		currSchema := (*currentSchemaCollection)[schemaId]
		return currSchema.GetTableByName(recordType)
	}
	return nil
}

func GetSchemaByName(name string) *Schema {
	var currentSchemaCollection = schemaCollection
	var currentSchemaNameCollection = schemaNameCollection
	if schemaId, ok := (*currentSchemaNameCollection)[name]; ok {
		return (*currentSchemaCollection)[schemaId]
	}
	return nil
}
