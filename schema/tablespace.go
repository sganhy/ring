package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
)

type Tablespace struct {
	id          int32
	name        string
	description string
	filName     string
	tableName   string
	table       bool
	index       bool
}

func (tablespace *Tablespace) Init(id int32, name string, description string, tableName string, table bool, index bool) {
	tablespace.id = id
	tablespace.name = name
	tablespace.description = description
	tablespace.tableName = tableName
	tablespace.table = table
	tablespace.index = index
}

//******************************
// getters
//******************************

func (tablespace *Tablespace) GetId() int32 {
	return tablespace.id
}

func (tablespace *Tablespace) GetName() string {
	return tablespace.name
}

func (tablespace *Tablespace) GetDescription() string {
	return tablespace.description
}

func (tablespace *Tablespace) GetEntityType() entitytype.EntityType {
	return entitytype.Tablespace
}

//******************************
// public methods
//******************************
func (tablespace *Tablespace) Clone() *Tablespace {
	newTablespace := new(Tablespace)
	newTablespace.Init(tablespace.id, tablespace.name, tablespace.description,
		tablespace.tableName, tablespace.table, tablespace.index)
	return newTablespace
}

func (tablespace *Tablespace) GetDdl(statement ddlstatement.DdlStatement, provider databaseprovider.DatabaseProvider) string {
	switch statement {
	case ddlstatement.NotDefined:
		return "TABLESPACE " + tablespace.name
	}
	return ""
}
