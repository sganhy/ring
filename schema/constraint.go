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
	relation  *Relation
}

const (
	constraintPkPrefix      string = "pk_"
	constraintPkShortPrefix string = "pk"
	constraintCkPrefix      string = "ck_"
	constraintFkPrefix      string = "fk_"
	constraintSeparator     string = "_"
	createFkPostGreSql      string = "%s %s %s ADD CONSTRAINT %s %s (%s) REFERENCES %s (%s)"
	createPkPostGreSql      string = "%s %s %s ADD CONSTRAINT %s %s (%s) %s"
	createNnPostGreSql      string = "%s %s %s ALTER COLUMN %s SET NOT NULL"
	createNnMySql           string = "%s %s %s MODIFY %s NOT NULL"
	createCkPostGreSql      string = "%s %s %s ADD CONSTRAINT %s CHECK (%s BETWEEN %d AND 127)"
)

func (constr *constraint) Init(consttype constrainttype.ConstraintType, table *Table) {
	constr.constType = consttype
	constr.table = table
}

//******************************
// getters and setters
//******************************
func (constr *constraint) GetId() int32 {
	return int32(0)
}
func (constr *constraint) GetName() string {
	return ""
}
func (constr *constraint) GetPhysicalName() string {
	return constr.getPhysicalName()
}
func (constr *constraint) GetEntityType() entitytype.EntityType {
	return entitytype.Constraint
}
func (constr *constraint) setField(field *Field) {
	constr.field = field
}
func (constr *constraint) setRelation(relation *Relation) {
	constr.relation = relation
}

func (constr *constraint) logStatement(statment ddlstatement.DdlStatement) bool {
	return constr.constType == constrainttype.ForeignKey || constr.constType == constrainttype.UniqueKey
}

//******************************
// public methods
//******************************
func (constr *constraint) GetDdl(statment ddlstatement.DdlStatement, tableSpace *tablespace) string {
	if statment == ddlstatement.Create {
		switch constr.constType {
		case constrainttype.PrimaryKey:
			return constr.getDdlCreatePrimaryKey(tableSpace)
		case constrainttype.Check:
			return constr.getDdlCheck()
		case constrainttype.ForeignKey:
			return constr.getDdlCreateForeignKey()
		case constrainttype.NotNull:
			return constr.getDdlNotNull()
		}
	}
	return ""
}

//******************************
// private methods
//******************************
func (constr *constraint) create(jobId int64, schema *Schema) error {
	if constr.constType == constrainttype.NotNull && constr.table.GetType() == tabletype.Business {
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
		var eventId int32 = 23
		metaQuery.query = query
		metaQuery.Init(schema, constr.table)
		// create table
		return metaQuery.create(eventId, jobId, constr)
	}
	return nil
}

func (constr *constraint) getPhysicalName() string {
	result := ""
	switch constr.constType {
	case constrainttype.PrimaryKey:
		//name:  pk_{table_name}
		var tableName = constr.table.GetName()
		if constr.table.GetType() == tabletype.Mtm {
			tableName = strings.Replace(tableName, mtmTableNamePrefix+mtmSeperator, "@", 1)
		}
		// keep constraint lenght less or equal to 30
		if len(tableName) > 27 {
			result = constraintPkShortPrefix + tableName
		} else {
			result = constraintPkPrefix + tableName
		}
		break
	case constrainttype.Check:
		//name:  idx_{table_id}_{index_id}
		result = constr.getCheckName()
		break
	case constrainttype.ForeignKey:
		if constr.table.GetType() == tabletype.Mtm {

			//TODO define fk name for MTM tables
			result = constraintFkPrefix + strconv.FormatInt(int64(constr.relation.GetId()), 10)
			result += strings.ReplaceAll(constr.table.GetName(), mtmTableNamePrefix, "")

		} else {

			//==> OK
			//name:  fk_{table_id}_{relation_id}
			tableId := sqlfmt.PadLeft(strconv.FormatInt(int64(constr.table.GetId()), 10), "0", 5)
			result = constraintFkPrefix + tableId + constraintSeparator +
				sqlfmt.PadLeft(strconv.FormatInt(int64(constr.relation.GetId()), 10), "0", 5)
		}
		break
	}
	return sqlfmt.FormatEntityName(constr.table.GetDatabaseProvider(), result)
}

func (constr *constraint) getCheckName() string {
	result := ""
	switch constr.table.GetType() {
	case tabletype.Business:
		//name: ck_{table_id}_{field_id}
		result = constraintCkPrefix + strconv.Itoa(int(constr.table.GetId())) + "_" +
			sqlfmt.PadLeft(strconv.Itoa(int(constr.field.GetId())), "0", 4)
		break
	case tabletype.Meta, tabletype.MetaId, tabletype.Log:
		//name: ck_{table_name}_{field_id}
		result = constraintCkPrefix + constr.table.GetName() + "_" +
			sqlfmt.PadLeft(strconv.Itoa(int(constr.field.GetId())), "0", 4)
		break

	}
	return result
}

func (constr *constraint) getDdlCreatePrimaryKey(tableSpace *tablespace) string {
	var sqlTablespace = ""
	var fields = constr.table.getUniqueFieldList()
	var provider = constr.table.GetDatabaseProvider()

	// unique fields ?
	if fields != "" {
		if tableSpace != nil && provider == databaseprovider.PostgreSql {
			// postgresql only ==>
			sqlTablespace = "USING INDEX " + tableSpace.GetDdl(ddlstatement.NotDefined, constr.table.GetDatabaseProvider())
		}

		switch provider {
		case databaseprovider.PostgreSql, databaseprovider.MySql:
			return strings.Trim(fmt.Sprintf(createPkPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
				constr.table.GetPhysicalName(), constr.getPhysicalName(), constrainttype.PrimaryKey.String(),
				fields, sqlTablespace), " ")
		}
	}
	return ""
}

func (constr *constraint) getDdlCreateForeignKey() string {
	if constr.relation != nil && constr.relation.HasConstraint() == true {

		var provider = constr.table.GetDatabaseProvider()
		switch provider {
		case databaseprovider.PostgreSql, databaseprovider.MySql:
			return strings.Trim(fmt.Sprintf(createFkPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
				constr.table.GetPhysicalName(), constr.getPhysicalName(), constrainttype.ForeignKey.String(),
				constr.relation.GetPhysicalName(provider), constr.relation.GetToTable().GetPhysicalName(),
				constr.relation.GetToTable().GetPrimaryKey().GetPhysicalName(provider)), " ")
		}
	}
	return ""
}

func (constr *constraint) getDdlNotNull() string {
	var fieldName string
	var provider = constr.table.GetDatabaseProvider()
	var dataType string

	if constr.field != nil && constr.field.IsNotNull() {
		fieldName = constr.field.GetPhysicalName(provider)
		dataType = constr.field.getSqlDataType(provider)
	}
	if constr.relation != nil && constr.relation.IsNotNull() {
		fieldName = constr.relation.GetPhysicalName(provider)
		//dataType = constr.relation.getSqlDataType(provider)
	}
	switch provider {
	case databaseprovider.PostgreSql:
		return strings.Trim(fmt.Sprintf(createNnPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
			constr.table.GetPhysicalName(), fieldName), ddlSpace)
	case databaseprovider.MySql:
		return strings.Trim(fmt.Sprintf(createNnMySql, ddlstatement.Alter.String(), entitytype.Table.String(),
			constr.table.GetPhysicalName(), fieldName+ddlSpace+dataType), ddlSpace)
		//constr.field.GetDdl(provider, constr.table.GetType()))
	}

	return ""
}

func (constr *constraint) getDdlCheck() string {
	var provider = constr.table.GetDatabaseProvider()

	// %s %s %s ADD CONSTRAINT %s CHECK (%s BETWEEN %d AND 127)
	if provider == databaseprovider.PostgreSql && constr.field != nil && constr.field.GetType() == fieldtype.Byte {
		var minValue = -128

		if constr.table.tableType == tabletype.Meta && constr.field.GetName() == metaObjectType {
			minValue = 0
		}

		return strings.Trim(fmt.Sprintf(createCkPostGreSql, ddlstatement.Alter.String(), entitytype.Table.String(),
			constr.table.GetPhysicalName(), constr.getPhysicalName(), constr.field.GetPhysicalName(constr.table.GetDatabaseProvider()),
			minValue), ddlSpace)

	}
	return ""
}
