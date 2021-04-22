package schema

import (
	"fmt"
	"ring/schema/constrainttype"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/sqlfmt"
	"strings"
)

type constraint struct {
	constType constrainttype.ConstraintType
	table     *Table
	fields    []*Field
}

const (
	constraintPkPrefix string = "pk_"
	createPkPostGreSql string = "%s %s %s ADD CONSTRAINT %s PRIMARY KEY (%s) %s"
)

func (constr *constraint) Init(consttype constrainttype.ConstraintType, table *Table) {
	constr.constType = consttype
	constr.table = table
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (constr *constraint) GetDdl(statment ddlstatement.DdlStatement, tablespace *Tablespace) string {
	switch statment {
	case ddlstatement.Create:
		switch constr.constType {
		case constrainttype.PrimaryKey:
			return constr.getDdlPrimaryKey(constr.table.provider, tablespace)
			break
		}
	}
	return ""
}

//******************************
// private methods
//******************************
func (constr *constraint) create(schema *Schema) error {
	var metaQuery = metaQuery{}
	var query = constr.GetDdl(ddlstatement.Create, schema.findTablespace(nil, nil, constr))
	//	var firstUniqueIndex = true
	if query != "" {
		metaQuery.query = query
		metaQuery.Init(schema, constr.table)
		// create table
		return metaQuery.create()
	}
	return nil
}

func (constr *constraint) getPhysicalName() string {
	result := ""
	switch constr.constType {
	case constrainttype.PrimaryKey:
		result = constraintPkPrefix + constr.table.name
		break
	case constrainttype.Check:
		result = constraintPkPrefix + constr.table.name
		break

	}
	return sqlfmt.FormatEntityName(constr.table.provider, result)
}

func (constr *constraint) getDdlPrimaryKey(provider databaseprovider.DatabaseProvider, tablespace *Tablespace) string {
	var sqlTablespace = ""
	var fields = constr.table.getUniqueFieldList()

	// unique fields ?
	if fields != "" {
		if tablespace != nil {
			// postgresql only ==>
			sqlTablespace = "USING INDEX " + tablespace.GetDdl(ddlstatement.NotDefined, constr.table.provider)
		}

		switch provider {
		case databaseprovider.PostgreSql:
			return strings.Trim(fmt.Sprintf(createPkPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
				constr.table.GetPhysicalName(), constr.getPhysicalName(), fields, sqlTablespace), " ")
		}
	}
	return ""
}
