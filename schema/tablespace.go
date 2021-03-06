package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
)

type tablespace struct {
	id          int32
	name        string
	description string
	filName     string
	tableName   string
	table       bool
	index       bool
	constraint  bool
}

const (
	tablespaceToStringFormat   string = "name=%s; description=%s; filename=%s; table=%t; index=%t"
	createTablespacePostGreSql string = "%s %s %s LOCATION '%s'"
)

func (tableSpace *tablespace) Init(id int32, name string, description string, fileName string, table bool, index bool) {
	tableSpace.id = id
	tableSpace.name = name
	tableSpace.description = description
	tableSpace.filName = fileName
	tableSpace.table = table
	tableSpace.index = index
}

//******************************
// getters and setters
//******************************
func (tableSpace *tablespace) GetId() int32 {
	return tableSpace.id
}

func (tableSpace *tablespace) GetName() string {
	return tableSpace.name
}

func (tableSpace *tablespace) GetPhysicalName() string {
	return tableSpace.GetName()
}

func (tableSpace *tablespace) GetDescription() string {
	return tableSpace.description
}

func (tableSpace *tablespace) GetEntityType() entitytype.EntityType {
	return entitytype.Tablespace
}

func (tableSpace *tablespace) setName(name string) {
	tableSpace.name = name
}

//******************************
// public methods
//******************************
func (tableSpace *tablespace) Clone() *tablespace {
	newTablespace := new(tablespace)
	newTablespace.Init(tableSpace.id, tableSpace.name, tableSpace.description,
		tableSpace.tableName, tableSpace.table, tableSpace.index)
	return newTablespace
}

func (tableSpace *tablespace) GetDdl(statement ddlstatement.DdlStatement, provider databaseprovider.DatabaseProvider) string {
	switch statement {
	case ddlstatement.NotDefined:
		return tableSpace.GetEntityType().String() + " " + tableSpace.name
	case ddlstatement.Create:
		return tableSpace.getDdlCreate(provider)
	}
	return ""
}

func (tableSpace *tablespace) String() string {
	// tablespaceToStringFormat string = "name=%s; description=%s; filename=%s; table=%t; index=%t"
	return fmt.Sprintf(tablespaceToStringFormat, tableSpace.name, tableSpace.description, tableSpace.filName,
		tableSpace.table, tableSpace.index)
}

//******************************
// private methods
//******************************
func (tableSpace *tablespace) exists(schema *Schema) bool {
	cata := new(catalogue)
	return cata.exists(schema, tableSpace)
}

func (tableSpace *tablespace) create(schema *Schema) {
}

func (tableSpace *tablespace) getDdlCreate(provider databaseprovider.DatabaseProvider) string {
	return ""
}
