package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/sqlfmt"
	"strings"
)

type catalogue struct {
	fieldEntityName string
	fieldSchemaName string
	viewName        string
}

var (
	postgreSqlCatalogue = map[entitytype.EntityType]catalogue{
		entitytype.Table:      {fieldEntityName: "tablename", fieldSchemaName: "schemaname", viewName: "pg_tables"},
		entitytype.Schema:     {fieldEntityName: "", fieldSchemaName: "nspname", viewName: "pg_catalog.pg_namespace"},
		entitytype.Tablespace: {fieldEntityName: "spcname", fieldSchemaName: "", viewName: "pg_tablespace"},
	}
	mySqlCatalogue = map[entitytype.EntityType]catalogue{
		entitytype.Table:      {fieldEntityName: "table_schema", fieldSchemaName: "table_name", viewName: "information_schema.tables"},
		entitytype.Schema:     {fieldEntityName: "", fieldSchemaName: "schema_name", viewName: "information_schema.schemata"},
		entitytype.Tablespace: {fieldEntityName: "", fieldSchemaName: "", viewName: ""},
	}
	catalogTable *Table
)

func init() {
	// take random table from metaSchema
	table := new(Table)
	//TODO create catalogue table
	catalogTable = table.getMetaTable(databaseprovider.NotDefined, "")
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (cata *catalogue) GetDql(provider databaseprovider.DatabaseProvider, entityType entitytype.EntityType) string {
	var result strings.Builder
	var mapper map[entitytype.EntityType]catalogue
	var variableIndex = 0

	table := new(Table)
	table.setDatabaseProvider(provider) // SET provider

	result.WriteString(dqlSelect)
	result.WriteString("1")
	result.WriteString(dqlFrom)

	switch provider {
	case databaseprovider.PostgreSql:
		mapper = postgreSqlCatalogue
		break
	case databaseprovider.MySql:
		mapper = mySqlCatalogue
		break
	}

	var objectName = mapper[entityType].fieldEntityName
	result.WriteString(mapper[entityType].viewName)
	result.WriteString(dqlWhere)

	if objectName != "" {
		result.WriteString("upper(")
		result.WriteString(mapper[entityType].fieldEntityName)
		result.WriteString(")=")
		result.WriteString(table.getVariableName(variableIndex))
		if entityType != entitytype.Tablespace {
			result.WriteString(filterSeparator)
		}
		variableIndex++
	}
	if entityType != entitytype.Tablespace {
		result.WriteString("upper(")
		result.WriteString(mapper[entityType].fieldSchemaName)
		result.WriteString(")=")
		result.WriteString(table.getVariableName(variableIndex))
	}

	return result.String()
}

//******************************
// private methods
//******************************
func (cata *catalogue) exists(schema *Schema, ent entity) bool {
	query := cata.GetDql(schema.GetDatabaseProvider(), ent.GetEntityType())
	if query != "" {
		var metaQuery = metaQuery{}
		metaQuery.query = query
		metaQuery.Init(schema, catalogTable)
		if ent.GetEntityType() != entitytype.Schema {
			metaQuery.addParam(strings.ToUpper(cata.getEntityName(schema, ent)))
		}
		if ent.GetEntityType() != entitytype.Tablespace {
			metaQuery.addParam(strings.ToUpper(schema.GetPhysicalName()))
		}
		result, err := metaQuery.exists()
		if err != nil {
			panic(err)
		}
		return result
	}
	panic("Unable to query from catalog")
}

func (cata *catalogue) getEntityName(schema *Schema, ent entity) string {

	// table physical name contains schema information ==> <schema_name>.t_<table_name>
	if ent.GetEntityType() == entitytype.Table {
		physicalName := ent.GetPhysicalName()
		if strings.Index(physicalName, ".") > 0 {
			arr := strings.Split(physicalName, ".")
			return sqlfmt.UnFormatEntityName(schema.GetDatabaseProvider(), arr[1])
		} else {
			return sqlfmt.UnFormatEntityName(schema.GetDatabaseProvider(), physicalName)
		}
	} else {
		return ent.GetName()
	}
}
