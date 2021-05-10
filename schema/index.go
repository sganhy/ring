package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
	"ring/schema/sqlfmt"
	"ring/schema/tabletype"
	"strconv"
	"strings"
)

type Index struct {
	id          int32
	name        string
	description string
	tableId     int32
	fields      []string
	bitmap      bool
	unique      bool
	baseline    bool
	active      bool
}

const (
	createIndexPostGreSql string = "%s%s %s %s ON %s USING btree (%s) %s"
	physicalIndexPrefix   string = "idx_"
	indexToStringFormat   string = "name=%s; description=%s; bitmap=%t; unique=%t; baseline=%t; fields=%s"
)

func (index *Index) Init(id int32, name string, description string, fields []string, tableId int32, bitmap bool,
	unique bool, baseline bool, active bool) {
	index.id = id
	index.name = name
	index.description = description
	index.loadFields(fields)
	index.bitmap = bitmap
	index.unique = unique
	index.baseline = baseline
	index.active = active
	index.tableId = tableId
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

func (index *Index) GetTableId() int32 {
	return index.tableId
}

func (index *Index) GetEntityType() entitytype.EntityType {
	return entitytype.Index

}

//******************************
// public methods
//******************************

func (index *Index) Clone() *Index {
	newIndex := new(Index)

	fields := make([]string, len(index.fields))
	copy(fields, index.fields)
	newIndex.Init(index.id, index.name, index.description,
		fields, index.tableId, index.bitmap, index.unique, index.baseline, index.active)

	return newIndex
}

func (index *Index) GetDdl(statment ddlstatement.DdlStatement, table *Table, tablespace *Tablespace) string {
	switch statment {
	case ddlstatement.Create:
		return index.getDdlCreate(table, tablespace)
	case ddlstatement.Drop:
		return index.getDdlDrop(table)
	}
	return ""
}

func (index *Index) GetPhysicalName(table *Table) string {
	var result strings.Builder
	result.Grow(30)
	result.WriteString(physicalIndexPrefix)

	switch table.tableType {
	case tabletype.Business:
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(table.id)), "0", 4))
		result.WriteString("_")
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(index.id)), "0", 4))
		break
	case tabletype.Mtm:
		result.WriteString("mtm_")
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(index.id)), "0", 4))
		break
	default:
		result.WriteString(table.name[1:])
		result.WriteString("_")
		result.WriteString(sqlfmt.PadLeft(strconv.Itoa(int(index.id)), "0", 3))
	}
	return result.String()
}

func (index *Index) String() string {
	//	indexToStringFormat   string = "name=%s; description=%s; bitmap=%s; unique=%s; baseline=%s; fields=%s"
	return fmt.Sprintf(indexToStringFormat, index.name, index.description, index.bitmap, index.unique, index.baseline,
		strings.Join(index.fields, metaIndexSeparator))
}

//******************************
// private methods
//******************************
func (index *Index) toMeta() *Meta {
	// we cannot have error here
	var result = new(Meta)

	// key
	result.id = index.id
	result.refId = index.tableId
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
		index.fields = make([]string, 0, 1)
	}
}

func (index *Index) getDdlDrop(table *Table) string {
	if table == nil || table.provider == databaseprovider.NotDefined {
		return ""
	}
	var query strings.Builder
	query.WriteString(ddlstatement.Drop.String())
	query.WriteString(ddlSpace)
	query.WriteString(entitytype.Index.String())
	query.WriteString(ddlSpace)
	query.WriteString(table.getSchemaName())
	query.WriteString(".")
	query.WriteString(index.GetPhysicalName(table))
	return query.String()
}

func (index *Index) getDdlCreate(table *Table, tablespace *Tablespace) string {
	var sqlUnique = ""
	var sqlTablespace = ""
	var fields strings.Builder

	if index.unique == true {
		sqlUnique = " UNIQUE"
	}
	if tablespace != nil {
		sqlTablespace = tablespace.GetDdl(ddlstatement.NotDefined, table.provider)
	}
	if len(index.fields) > 0 {
		for i := 0; i < len(index.fields); i++ {
			fields.WriteString(sqlfmt.FormatEntityName(table.provider, index.fields[i]))
			if i < len(index.fields)-1 {
				fields.WriteString(",")
			}
		}
	}
	switch table.provider {
	case databaseprovider.PostgreSql:
		return strings.Trim(fmt.Sprintf(createIndexPostGreSql, ddlstatement.Create.String(), sqlUnique, entitytype.Index.String(),
			index.GetPhysicalName(table), table.GetPhysicalName(), fields.String(), sqlTablespace), ddlSpace)
	}
	return ""
}

func (index *Index) create(schema *Schema) error {
	var metaQuery = metaQuery{}
	var table = schema.GetTableById(index.tableId)

	metaQuery.query = index.GetDdl(ddlstatement.Create, table, schema.findTablespace(table, nil, nil))
	metaQuery.setSchema(schema.GetName())
	metaQuery.setTable(table.GetName())

	return metaQuery.create()
}
