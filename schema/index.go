package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/sqlfmt"
	"ring/schema/tabletype"
	"strconv"
	"strings"
)

type Index struct {
	id          int32
	name        string
	description string
	fields      []string
	bitmap      bool
	unique      bool
	baseline    bool
	active      bool
}

const (
	createIndexPostGreSql string = "%s%s %s %s ON %s USING btree (%s) %s"
	createIndexMysql      string = "%s%s %s %s ON %s (%s) USING btree %s"
	physicalIndexPrefix   string = "idx_"
	indexToStringFormat   string = "name=%s; description=%s; bitmap=%t; unique=%t; baseline=%t; fields=%s"
)

func (index *Index) Init(id int32, name string, description string, fields []string, bitmap bool,
	unique bool, baseline bool, active bool) {
	index.id = id
	index.name = name
	index.description = description
	index.loadFields(fields)
	index.bitmap = bitmap
	index.unique = unique
	index.baseline = baseline
	index.active = active
}

//******************************
// getters and setters
//******************************
func (index *Index) GetId() int32 {
	return index.id
}

func (index *Index) GetName() string {
	return index.name
}

func (index *Index) GetPhysicalName() string {
	return index.GetName()
}

func (index *Index) GetDescription() string {
	return index.description
}

func (index *Index) GetFields() []string {
	return index.fields
}

func (index *Index) IsUnique() bool {
	return index.unique
}

func (index *Index) IsBitmap() bool {
	return index.bitmap
}

func (index *Index) IsBaseline() bool {
	return index.baseline
}

func (index *Index) IsActive() bool {
	return index.active
}

func (index *Index) GetEntityType() entitytype.EntityType {
	return entitytype.Index
}

func (index *Index) logStatement(statment ddlstatement.DdlStatement) bool {
	return true
}
func (index *Index) setUnique(unique bool) {
	index.unique = unique
}

//******************************
// public methods
//******************************
func (index *Index) Clone() *Index {
	newIndex := new(Index)

	fields := make([]string, len(index.fields))
	copy(fields, index.fields)
	newIndex.Init(index.id, index.name, index.description,
		fields, index.bitmap, index.unique, index.baseline, index.active)

	return newIndex
}

func (index *Index) GetDdl(statment ddlstatement.DdlStatement, table *Table, tableSpace *tablespace) string {
	switch statment {
	case ddlstatement.Create:
		return index.getDdlCreate(table, tableSpace)
	case ddlstatement.Drop:
		return index.getDdlDrop(table)
	}
	return ""
}

func (index *Index) String() string {
	//	indexToStringFormat   string = "name=%s; description=%s; bitmap=%s; unique=%s; baseline=%s; fields=%s"
	return fmt.Sprintf(indexToStringFormat, index.name, index.description, index.bitmap, index.unique, index.baseline,
		strings.Join(index.fields, metaIndexSeparator))
}

//******************************
// private methods
//******************************

// compare if the physical indexes are equal
func (indexA *Index) equal(indexB *Index) bool {
	//fmt.Println(indexA.String())
	//fmt.Println(indexB.String())
	return strings.EqualFold(strings.Join(indexA.fields, metaIndexSeparator),
		strings.Join(indexB.fields, metaIndexSeparator)) &&
		indexA.bitmap == indexB.bitmap &&
		indexA.unique == indexB.unique
}

func (index *Index) getPhysicalName(table *Table) string {
	var result strings.Builder
	result.Grow(30)
	result.WriteString(physicalIndexPrefix)

	switch table.GetType() {
	case tabletype.Business:
		//name:  idx_{table_id}_{index_id}
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(table.GetId())), "0", 4))
		result.WriteString("_")
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(index.id)), "0", 4))
		break
	case tabletype.Mtm:
		//name:  idx_{from_table_id}_{to_table_id}_{from_relation_id}
		var indexBody = strings.Replace(table.GetName(), mtmTableNamePrefix, "", 1)
		result.WriteString(indexBody[1:])
		break
	default:
		//name:  idx_{table_name}_{index_id}
		result.WriteString(table.GetName()[1:])
		result.WriteString("_")
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(index.id)), "0", 3))
	}
	return result.String()
}

func (index *Index) toMeta(tableId int32) *meta {
	// we cannot have error here
	var result = new(meta)

	// key
	result.id = index.id
	result.refId = tableId
	result.objectType = int8(entitytype.Index)

	// others
	result.dataType = 0
	result.name = index.name // max length 30 !! must be validated before
	result.description = index.description
	result.value = strings.Join(index.fields, metaIndexSeparator)

	// flags
	result.flags = 0
	result.setIndexBitmap(index.bitmap)
	result.setIndexUnique(index.unique)
	result.setEntityBaseline(index.baseline)
	result.enabled = index.active
	return result
}

func (index *Index) loadFields(fields []string) {
	// copy slice -- func make([]T, len, cap) []T
	if fields != nil {
		index.fields = make([]string, len(fields))
		copy(index.fields, fields)
	} else {
		index.fields = make([]string, 0, 0)
	}
}

func (index *Index) getDdlDrop(table *Table) string {
	if table == nil || table.GetDatabaseProvider() == databaseprovider.Undefined {
		return ""
	}
	var query strings.Builder
	var schema = index.getSchema(table.GetSchemaId())

	// generate error later !!
	if schema != nil {
		query.WriteString(ddlstatement.Drop.String())
		query.WriteString(ddlSpace)
		query.WriteString(entitytype.Index.String())
		query.WriteString(ddlSpace)
		query.WriteString(schema.GetPhysicalName())
		query.WriteString(".")
		query.WriteString(index.getPhysicalName(table))
	}

	return query.String()
}

func (index *Index) getSchema(schemaId int32) *Schema {
	var schema = GetSchemaById(schemaId)
	if schema == nil {
		schema = getUpgradingSchema()
	}
	return schema
}

func (index *Index) getDdlCreate(table *Table, tableSpace *tablespace) string {
	var sqlUnique = ""
	var sqlTablespace = ""
	var fields string
	var provider = table.GetDatabaseProvider()

	if index.unique == true {
		sqlUnique = " UNIQUE"
	}
	if tableSpace != nil {
		sqlTablespace = tableSpace.GetDdl(ddlstatement.Undefined, provider)
	}

	fields = index.getFieldList(table)

	switch table.GetDatabaseProvider() {
	case databaseprovider.PostgreSql:
		return strings.Trim(fmt.Sprintf(createIndexPostGreSql, ddlstatement.Create.String(), sqlUnique, entitytype.Index.String(),
			index.getPhysicalName(table), table.GetPhysicalName(), fields, sqlTablespace), ddlSpace)
	case databaseprovider.MySql:
		return strings.Trim(fmt.Sprintf(createIndexMysql, ddlstatement.Create.String(), sqlUnique, entitytype.Index.String(),
			index.getPhysicalName(table), table.GetPhysicalName(), fields, sqlTablespace), ddlSpace)
	}
	return ""
}

func (index *Index) getFieldList(table *Table) string {
	var result strings.Builder
	var provider = table.GetDatabaseProvider()

	if len(index.fields) > 0 {
		for i := 0; i < len(index.fields); i++ {
			field := table.GetFieldByName(index.fields[i])

			// manage searchable fields
			if field != nil && field.IsCaseSensitive() == false && field.GetType() == fieldtype.String {
				result.WriteString(sqlfmt.FormatEntityName(provider, field.getSearchableFieldName()))
			} else {
				result.WriteString(sqlfmt.FormatEntityName(provider, index.fields[i]))
			}

			if i < len(index.fields)-1 {
				result.WriteString(",")
			}
		}
	}

	return result.String()
}

func (index *Index) create(jobId int64, schema *Schema, table *Table) error {
	var metaQuery = metaQuery{}
	var eventId int32 = 88

	metaQuery.query = index.GetDdl(ddlstatement.Create, table, schema.findTablespace(nil, index, nil))
	metaQuery.Init(schema, table)

	err := metaQuery.create(eventId, jobId, index)

	return err
}

func (index *Index) drop(jobId int64, schema *Schema, table *Table) error {
	var metaQuery = metaQuery{}
	var eventId int32 = 73

	metaQuery.query = index.GetDdl(ddlstatement.Drop, table, nil)
	metaQuery.Init(schema, table)
	return metaQuery.drop(eventId, jobId, index)
}
