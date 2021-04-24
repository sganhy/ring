package schema

import (
	"fmt"
	"ring/schema/constrainttype"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/sqlfmt"
	"ring/schema/tabletype"
	"strconv"
	"strings"
)

type constraint struct {
	constType constrainttype.ConstraintType
	table     *Table
	field     *Field
}

const (
	constraintPkPrefix string = "pk_"
	constraintCkPrefix string = "ck_"
	createPkPostGreSql string = "%s %s %s ADD CONSTRAINT %s PRIMARY KEY (%s) %s"
	createNnPostGreSql string = "%s %s %s ALTER COLUMN %s SET NOT NULL"
	createNnMySql      string = "%s %s %s MODIFY %s NOT NULL"
	createCkPostGreSql string = "%s %s %s ADD CONSTRAINT %s CHECK (%s BETWEEN -128 AND 127)"
)

func (constr *constraint) Init(consttype constrainttype.ConstraintType, table *Table) {
	constr.constType = consttype
	constr.table = table
}

//******************************
// getters and setters
//******************************
func (constr *constraint) GetEntityType() entitytype.EntityType {
	return entitytype.Constraint

}

//******************************
// public methods
//******************************
func (constr *constraint) GetDdl(statment ddlstatement.DdlStatement, tablespace *Tablespace) string {
	if statment == ddlstatement.Create {
		switch constr.constType {
		case constrainttype.PrimaryKey:
			return constr.getDdlPrimaryKey(tablespace)
		case constrainttype.Check:
			return constr.getDdlCheck()
		case constrainttype.NotNull:
			return constr.getDdlNotNull()
		}
	}
	return ""
}

//******************************
// private methods
//******************************
func (constr *constraint) create(schema *Schema) error {
	if constr.constType == constrainttype.NotNull && constr.table.tableType == tabletype.Business {
		// do nothing for busines tables
		return nil
	}
	if constr.constType == constrainttype.Check && schema.GetDatabaseProvider() == databaseprovider.MySql {
		// do nothing
		return nil
	}
	var query = constr.GetDdl(ddlstatement.Create, schema.findTablespace(nil, nil, constr))
	//	var firstUniqueIndex = true
	if query != "" {
		var metaQuery = metaQuery{}
		metaQuery.query = query
		metaQuery.Init(schema, constr.table)
		// create table
		return metaQuery.create()
	}
	return nil
}

func (constr *constraint) setField(field *Field) {
	constr.field = field
}

func (constr *constraint) getPhysicalName() string {
	result := ""
	switch constr.constType {
	case constrainttype.PrimaryKey:
		result = constraintPkPrefix + constr.table.name
		break
	case constrainttype.Check:
		result = constr.getCheckName()
		break

	}
	return sqlfmt.FormatEntityName(constr.table.provider, result)
}

func (constr *constraint) getCheckName() string {
	result := ""
	switch constr.table.tableType {
	case tabletype.Business:
		result = constraintCkPrefix + strconv.Itoa(int(constr.table.id)) + "_" +
			sqlfmt.PadLeft(strconv.Itoa(int(constr.field.id)), "0", 4)
		break
	case tabletype.Meta, tabletype.MetaId, tabletype.Log:
		result = constraintCkPrefix + constr.table.name + "_" +
			sqlfmt.PadLeft(strconv.Itoa(int(constr.field.id)), "0", 4)
		break

	}
	return result
}

func (constr *constraint) getDdlPrimaryKey(tablespace *Tablespace) string {
	var sqlTablespace = ""
	var fields = constr.table.getUniqueFieldList()
	var provider = constr.table.provider

	// unique fields ?
	if fields != "" {
		if tablespace != nil && provider == databaseprovider.PostgreSql {
			// postgresql only ==>
			sqlTablespace = "USING INDEX " + tablespace.GetDdl(ddlstatement.NotDefined, constr.table.provider)
		}

		switch provider {
		case databaseprovider.PostgreSql, databaseprovider.MySql:
			return strings.Trim(fmt.Sprintf(createPkPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
				constr.table.GetPhysicalName(), constr.getPhysicalName(), fields, sqlTablespace), " ")
		}
	}
	return ""
}

func (constr *constraint) getDdlNotNull() string {
	if constr.field != nil && constr.field.IsNotNull() {
		provider := constr.table.provider
		switch provider {
		case databaseprovider.PostgreSql:
			return strings.Trim(fmt.Sprintf(createNnPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
				constr.table.GetPhysicalName(), constr.field.GetPhysicalName(provider)), ddlSpace)
		case databaseprovider.MySql:
			return strings.Trim(fmt.Sprintf(createNnMySql, ddlstatement.Alter.String(), entitytype.Table.String(),
				constr.table.GetPhysicalName(), constr.field.GetDdl(provider, constr.table.tableType)), ddlSpace)
		}
	}
	return ""
}

func (constr *constraint) getDdlCheck() string {
	var provider = constr.table.provider
	// "%s %s %s ADD CONSTRAINT %s CHECK (%s>-129 AND %s<128)"
	if provider == databaseprovider.PostgreSql && constr.field != nil && constr.field.fieldType == fieldtype.Byte {
		return strings.Trim(fmt.Sprintf(createCkPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
			constr.table.GetPhysicalName(), constr.getPhysicalName(), constr.field.GetPhysicalName(constr.table.provider)),
			ddlSpace)
	}
	return ""
}
