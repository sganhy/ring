package schema

import (
	"fmt"
	"ring/schema/constrainttype"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/searchabletype"
	"ring/schema/sqlfmt"
	"ring/schema/tabletype"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Table struct {
	id              int32
	name            string
	description     string
	fields          []*Field                          // sorted by name
	relations       []*Relation                       // sorted by name
	indexes         []*Index                          // sorted by name
	mapper          []uint16                          // mapping from .fieldsById to .fields ==> max number of column 65535!!
	physicalName    string                            //
	physicalType    physicaltype.PhysicalType         //
	schemaId        int32                             //
	tableType       tabletype.TableType               //
	fieldList       string                            // list of field without searchable fields
	sqlInsert       string                            // cache insert query
	subject         string                            //
	provider        databaseprovider.DatabaseProvider //
	cacheid         *cacheId                          // store id primary key information
	searchableCount uint8                             //
	cached          bool
	readonly        bool
	baseline        bool
	active          bool
}

const (
	createTableSql       string = "%s %s %s (\n%s\n)"
	ddlSpace             string = " "
	dmlInsertEnd         string = ")"
	dmlInsertStart       string = " ("
	dmlInsertValues      string = ") VALUES ("
	dmlSpace             string = " "
	dmlUpdateSet         string = " SET "
	dqlFrom              string = " FROM "
	dqlOrderBy           string = " ORDER BY "
	dqlSelect            string = "SELECT "
	dqlSpace             string = " "
	dqlWhere             string = " WHERE "
	fieldListSeparator   string = ","
	fieldNotFound        int    = -1
	maxNumberOfColumn    int    = 65535
	metaDataType         string = "data_type"
	metaDescription      string = "description"
	metaFlags            string = "flags"
	metaFieldId          string = "id"
	metaLogCallSite      string = "call_site"
	metaLogEntryTime     string = "entry_time"
	metaLogId            string = "id"
	metaLogJobId         string = "job_id"
	metaLogLevelId       string = "level_id"
	metaLoglineNumber    string = "line_number"
	metaLogMessage       string = "message"
	metaLogMethod        string = "method"
	metaLogTableName     string = "@log"
	metaLogThreadId      string = "thread_id"
	metaName             string = "name"
	metaUuid             string = "uuid"
	metaLexTableId       string = "table_id"
	metaLexFromFieldId   string = "source_field_id"
	metaLexToFieldId     string = "target_field_id"
	metaLexRelationId    string = "relation_id"
	metaLexRelationValue string = "relation_value"
	metaLexModifyStmp    string = "modify_stmp"
	metaLexItmLexId      string = "lexicon_id"
	metaObjectType       string = "object_type"
	metaReferenceId      string = "reference_id"
	metaLangId           string = "lang_id"
	metaSchemaId         string = "schema_id"
	metaTableName        string = "@meta"
	metaIdTableName      string = "@meta_id"
	metaLexTableName     string = "@lexicon"
	metaLexItmTableName  string = "@lexicon_itm"
	metaLongTableName    string = "@long"
	metaValue            string = "value"
	metaDataModifyStmp   string = "modify_stmp"
	metaActive           string = "active"
	metaCached           string = "cached"
	postGreAddColumn     string = "ADD"
	postGreDropColumn    string = "DROP"
	postGreColumn        string = "COLUMN "
	postGreCreateOptions string = " WITH (autovacuum_enabled=false) "
	postGreParameterName string = "$"
	postGreVacuum        string = "VACUUM %s"
	postGreFullVacuum    string = "FULL "
	postGreCascade       string = "CASCADE"
	postGreAnalyze       string = "ANALYZE %s"
	relationSeparator    string = "|"
	physicalTablePrefix  string = "t_"
	mysqlParameterName   string = "?"
	relationNotFound     int    = -1
	crlf                 string = "\n"
)

func (table *Table) Init(id int32, name string, description string, fields []Field, relations []Relation, indexes []Index,
	physicalType physicaltype.PhysicalType, schemaId int32, schemaPhysicalName string, tableType tabletype.TableType, provider databaseprovider.DatabaseProvider,
	subject string, cached bool, readonly bool, baseline bool, active bool) {

	table.id = id
	table.name = name
	table.description = description
	table.loadFields(fields, tableType)
	table.loadRelations(relations)
	table.loadMapper()         // !!! load after loadFields()
	table.loadIndexes(indexes) // !!! run at the end only
	table.tableType = tableType
	table.schemaId = schemaId
	table.tableType = tableType
	table.subject = subject
	table.provider = provider
	table.physicalType = physicalType
	table.cached = cached
	table.readonly = readonly
	table.baseline = baseline
	table.active = active
	table.searchableCount = table.getSearchableCount()
	if tableType == tabletype.Business {
		table.cacheid = new(cacheId)
	}
	if provider != databaseprovider.Undefined {
		table.physicalName = table.getPhysicalName(schemaPhysicalName)
		table.fieldList = table.getFieldList()
		table.sqlInsert = table.getInsertSql(provider)
	}

}

//******************************
// getters and setters
//******************************
func (table *Table) GetId() int32 {
	return table.id
}

func (table *Table) GetName() string {
	return table.name
}

func (table *Table) GetDescription() string {
	return table.description
}

func (table *Table) getCacheId() *cacheId {
	return table.cacheid
}

func (table *Table) GetPhysicalName() string {
	return table.physicalName
}

func (table *Table) GetSchemaId() int32 {
	return table.schemaId
}

func (table *Table) GetType() tabletype.TableType {
	return table.tableType
}

func (table *Table) GetSubject() string {
	return table.subject
}

func (table *Table) GetPhysicalType() physicaltype.PhysicalType {
	return table.physicalType
}

func (table *Table) GetDatabaseProvider() databaseprovider.DatabaseProvider {
	return table.provider
}

func (table *Table) IsCached() bool {
	return table.cached
}

func (table *Table) IsReadonly() bool {
	return table.readonly
}

func (table *Table) IsBaseline() bool {
	return table.baseline
}

func (table *Table) IsActive() bool {
	return table.active
}

func (table *Table) GetFieldCount() int {
	return len(table.fields)
}

func (table *Table) GetEntityType() entitytype.EntityType {
	return entitytype.Table
}

func (table *Table) GetFieldIdByIndex(index int) *Field {
	return table.fields[table.mapper[index]]
}

func (table *Table) GetFieldByIndex(index int) *Field {
	return table.fields[index]
}

func (table *Table) setId(id int32) {
	table.id = id
}

func (table *Table) setName(name string) {
	table.name = name
}

func (table *Table) setDatabaseProvider(provider databaseprovider.DatabaseProvider) {
	table.provider = provider
}

func (table *Table) setTableType(tableType tabletype.TableType) {
	table.tableType = tableType
}

func (table *Table) logStatement(statment ddlstatement.DdlStatement) bool {
	return true
}

//******************************
// public methods
//******************************
func (table *Table) GetFieldByName(name string) *Field {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(table.fields) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.fields[indexerMiddle].name)
		if indexerCompare == 0 {
			return table.fields[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return nil
}

// case insensitive search
// Get field by name ==> O(n) complexity
func (table *Table) GetFieldByNameI(name string) *Field {
	for i := len(table.fields) - 1; i >= 0; i-- {
		if strings.EqualFold(name, table.fields[i].name) {
			return table.fields[i]
		}
	}
	return nil
}

//return -1 if not found
func (table *Table) GetFieldIndexByName(name string) int {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(table.fields) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.fields[indexerMiddle].name)
		if indexerCompare == 0 {
			return indexerMiddle
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return fieldNotFound
}

//
func (table *Table) GetRelationByName(name string) *Relation {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(table.relations) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.relations[indexerMiddle].name)
		if indexerCompare == 0 {
			return table.relations[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return nil
}

func (table *Table) GetRelationByNameI(name string) *Relation {
	for i := len(table.relations) - 1; i >= 0; i-- {
		if strings.EqualFold(name, table.relations[i].name) {
			return table.relations[i]
		}
	}
	return nil
}

func (table *Table) GetNewObjid(reservedRange uint32) int64 {
	return table.cacheid.getNewRangeId(reservedRange)
}

func (table *Table) GetIndexByName(name string) *Index {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(table.indexes) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.indexes[indexerMiddle].name)
		if indexerCompare == 0 {
			return table.indexes[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return nil
}

//return -1 if not found
func (table *Table) GetRelationIndexByName(name string) int {
	var indexerLeft, indexerRight, indexerMiddle, indexerCompare = 0, len(table.relations) - 1, 0, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.relations[indexerMiddle].name)
		if indexerCompare == 0 {
			return indexerMiddle
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
		}
	}
	return relationNotFound
}

func (table *Table) GetPrimaryKey() *Field {
	if len(table.fields) > 0 && (table.tableType == tabletype.Business ||
		table.tableType == tabletype.Lexicon) {
		return table.fields[table.mapper[0]]
	}
	return nil
}

func (table *Table) GetPrimaryKeyIndex() int {
	return int(table.mapper[0])
}

func (table *Table) GetDdl(statement ddlstatement.DdlStatement, tableSpace *tablespace, field *Field) string {
	var query strings.Builder

	switch statement {
	case ddlstatement.Create:
		var fields []string
		for i := 0; i < len(table.fields); i++ {
			index := table.mapper[i]
			fieldSql := table.fields[index].GetDdl(table.provider, table.tableType)
			fields = append(fields, fieldSql)
			if table.fields[index].IsCaseSensitive() == false && table.fields[index].GetType() == fieldtype.String {
				fieldSql := table.fields[index].getSearchableDdl(table.provider, table.tableType)
				fields = append(fields, fieldSql)
			}
		}
		for i := 0; i < len(table.relations); i++ {
			if table.relations[i].GetType() == relationtype.Mto || table.relations[i].GetType() == relationtype.Otop {
				relationSql := table.relations[i].GetDdl(table.provider)
				fields = append(fields, relationSql)
			}
		}
		query.WriteString(fmt.Sprintf(createTableSql, ddlstatement.Create.String(), entitytype.Table.String(), table.physicalName,
			strings.Join(fields, ",\n")))
		query.WriteString(table.getOptions())
		if tableSpace != nil {
			query.WriteString(ddlSpace)
			query.WriteString(tableSpace.GetDdl(ddlstatement.Undefined, table.provider))
		}
		break
	case ddlstatement.Alter:
		return table.getDdlAlter(field)
	case ddlstatement.Drop:
		query.WriteString(ddlstatement.Drop.String())
		query.WriteString(ddlSpace)
		query.WriteString(entitytype.Table.String())
		query.WriteString(ddlSpace)
		query.WriteString(table.physicalName)
		query.WriteString(ddlSpace)
		query.WriteString(postGreCascade)
		break
	case ddlstatement.Truncate:
		query.WriteString(ddlstatement.Truncate.String())
		query.WriteString(ddlSpace)
		query.WriteString(entitytype.Table.String())
		query.WriteString(ddlSpace)
		query.WriteString(table.physicalName)
		break
	}
	return query.String()
}

func (table *Table) GetDml(dmlType dmlstatement.DmlStatement, fields []*Field) string {
	var result strings.Builder
	switch dmlType {
	case dmlstatement.Insert:
		return table.sqlInsert
	case dmlstatement.Update:
		result.Grow(len(table.sqlInsert))
		result.WriteString(dmlType.String())
		result.WriteString(dmlSpace)
		result.WriteString(table.physicalName)
		result.WriteString(dmlUpdateSet)
		for i := 0; i < len(fields); i++ {
			result.WriteString(fields[i].GetPhysicalName(table.provider))
			result.WriteString(operatorEqual)
			result.WriteString(table.getVariableName(i))
			if i < len(fields)-1 {
				result.WriteString(fieldListSeparator)
			}
		}
		result.WriteString(dqlWhere)
		table.addPrimaryKeyFilter(&result, len(fields))
		break
	case dmlstatement.Delete:
		result.Grow(len(table.physicalName) + 30)
		result.WriteString(dmlType.String())
		result.WriteString(dmlSpace)
		result.WriteString(table.physicalName)
		result.WriteString(dqlWhere)
		table.addPrimaryKeyFilter(&result, 0)
		break
	}
	return result.String()
}

// SELECT
func (table *Table) GetDql(whereClause string, orderClause string) string {
	capacity := len(dqlSelect) + len(dqlFrom) + len(table.fieldList) + len(table.physicalName)

	if len(whereClause) > 0 {
		capacity += len(dqlWhere)
		capacity += len(whereClause)
	}

	if len(orderClause) > 0 {
		capacity += len(dqlOrderBy)
		capacity += len(orderClause)
	}

	var result strings.Builder
	result.Grow(capacity)
	result.WriteString(dqlSelect)
	result.WriteString(table.fieldList)
	result.WriteString(dqlFrom)
	result.WriteString(table.physicalName)

	if len(whereClause) > 0 {
		result.WriteString(dqlWhere)
		result.WriteString(whereClause)
	}

	if len(orderClause) > 0 {
		result.WriteString(dqlOrderBy)
		result.WriteString(orderClause)
	}
	// check capacity
	//fmt.Printf("GetDql() capacity/len(str)/sql.Cap() %d /%d /%d\n", capacity, len(sql.String()), sql.Cap())
	return result.String()
}

/*
	Deep copy of table object
*/
func (table *Table) Clone() *Table {
	newTable := new(Table)
	var fields []Field
	var relations []Relation
	var indexes []Index

	// don't clone ToTable for reflexive relationship (recursive call)
	for i := 0; i < len(table.fields); i++ {
		fields = append(fields, *table.fields[i].Clone())
	}
	for i := 0; i < len(table.relations); i++ {
		relations = append(relations, *table.relations[i].Clone())
	}
	for i := 0; i < len(table.indexes); i++ {
		indexes = append(indexes, *table.indexes[i].Clone())
	}
	newTable.Init(table.id, table.name, table.description, fields, relations, indexes,
		table.physicalType, table.schemaId, "", table.tableType, table.provider, table.subject,
		table.cached, table.readonly, table.baseline, table.active)

	newTable.fieldList = table.fieldList
	newTable.physicalName = table.physicalName

	return newTable
}

func (table *Table) GetQueryResult(columnPointer []interface{}) []string {
	var capacity = len(table.fields)
	var result = make([]string, capacity, capacity)
	var currentField *Field
	var value interface{}
	var strValue string
	var index int
	var array []byte

	// manage fields first
	for i := 0; i < capacity; i++ {
		value = *columnPointer[i].(*interface{})
		currentField = table.fields[table.mapper[i]]
		// use a mapper instead
		//index = table.GetFieldIndexByName(currentField.name)
		index = int(table.mapper[i])
		if value == nil {
			strValue = ""
		} else {
			switch value.(type) {
			case string:
				strValue = value.(string)
				// long text==> mysql
				break
			case float32:
				// conversion issues
				strValue = fmt.Sprintf("%g", value.(float32))
				break
			case float64:
				strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
				break
			case int:
				strValue = strconv.Itoa(value.(int))
				break
			case uint:
				strValue = strconv.FormatUint(uint64(value.(uint)), 10)
				break
			case bool:
				strValue = strconv.FormatBool(value.(bool))
				break
			case int8:
				strValue = strconv.Itoa(int(value.(int8)))
				break
			case int16:
				strValue = strconv.Itoa(int(value.(int16)))
				break
			case int32:
				strValue = strconv.Itoa(int(value.(int32)))
				break
			case int64:
				strValue = strconv.FormatInt(value.(int64), 10)
				break
			case uint8:
				strValue = strconv.FormatUint(uint64(value.(uint8)), 10)
				break
			case uint16:
				strValue = strconv.FormatUint(uint64(value.(uint16)), 10)
				break
			case uint32:
				strValue = strconv.FormatUint(uint64(value.(uint32)), 10)
				break
			case uint64:
				strValue = strconv.FormatUint(value.(uint64), 10)
				break
			case []uint8:
				array = value.([]byte) //my SQL longtext & varchar
				strValue = string(array)
				break
			case time.Time:
				strValue = currentField.GetDateTimeString(value.(time.Time))
			default:
				strValue = value.(string)
			}
		}
		result[index] = strValue
	}
	return result
}

func (table *Table) Vacuum(jobId int64, full bool) error {
	var sql string

	switch table.provider {
	case databaseprovider.PostgreSql:
		if full == true {
			sql = fmt.Sprintf(postGreVacuum, postGreFullVacuum)
			sql += table.physicalName
		} else {
			sql = fmt.Sprintf(postGreVacuum, table.physicalName)
		}
		break
	}
	if len(sql) > 0 {
		var metaQuery = metaQuery{}
		var eventId int32 = 20
		var schema = GetSchemaById(table.schemaId)

		metaQuery.query = sql
		metaQuery.Init(schema, table)
		return metaQuery.vacuum(eventId, jobId, table)
	}
	return nil
}

func (table *Table) Analyze(jobId int64) error {
	var sql string
	switch table.provider {
	case databaseprovider.PostgreSql:
		sql = fmt.Sprintf(postGreAnalyze, table.physicalName)
		break
	}
	if len(sql) > 0 {
		var metaQuery = metaQuery{}
		var eventId int32 = 21
		var schema = GetSchemaById(table.schemaId)

		metaQuery.query = sql
		metaQuery.Init(schema, table)
		return metaQuery.analyze(eventId, jobId, table)
	}
	return nil
}

func (table *Table) Truncate(jobId int64) error {
	var eventId int32 = 63

	// truncate not allowed on meta tables
	if table.schemaId == 0 {
		// ? create an error ?
		return nil
	}
	var metaQuery = metaQuery{}
	var schema = GetSchemaById(table.schemaId)

	metaQuery.Init(schema, table)
	metaQuery.query = table.GetDdl(ddlstatement.Truncate, nil, nil)
	return metaQuery.truncate(eventId, jobId, table)
}

func (table *Table) String() string {
	var b strings.Builder
	for i := 0; i < len(table.fields); i++ {
		b.WriteString(table.fields[i].String())
		b.WriteString(crlf)
	}
	return b.String()
}

func (table *Table) GetInsertParameters(record []string) []interface{} {
	count := len(table.fields) + int(table.searchableCount)
	result := make([]interface{}, count, count)

	index := 0
	searchableCount := 0

	var value string
	var field *Field

	for i := 0; i < len(table.fields); i++ {

		index = int(table.mapper[i])
		value = record[index]
		field = table.fields[index]

		if len(value) > 0 {
			result[i+searchableCount] = field.GetParameterValue(value)
		} else if field.notNull == true {
			result[i+searchableCount] = field.GetParameterValue(field.GetDefaultValue())
		}
		if field.caseSensitive == false && field.fieldType == fieldtype.String {
			searchableCount++
			result[i+searchableCount] = field.GetSearchableValue(value, searchabletype.None)
		}

	}

	return result
}

func (table *Table) GetDeleteParameters(id int64) []interface{} {
	result := make([]interface{}, 1, 1)
	result[0] = id
	return result
}

//******************************
// private methods
//******************************
func (table *Table) toMeta() *meta {
	var metaTable = new(meta)

	// key
	metaTable.id = table.id
	metaTable.refId = table.schemaId
	metaTable.objectType = int8(entitytype.Table)

	// others
	metaTable.dataType = 0
	metaTable.name = table.name // max length 30 !! must be validated before
	metaTable.description = table.description
	metaTable.value = table.subject

	// flags
	metaTable.flags = 0
	metaTable.setEntityBaseline(table.baseline)
	metaTable.setTableCached(table.cached)
	metaTable.setTableReadonly(table.readonly)
	metaTable.enabled = table.active

	return metaTable
}

func (table *Table) getDdlAlter(field *Field) string {
	if field != nil {
		var query strings.Builder
		query.WriteString(ddlstatement.Alter.String())
		query.WriteString(ddlSpace)
		query.WriteString(entitytype.Table.String())
		query.WriteString(ddlSpace)
		query.WriteString(table.physicalName)
		query.WriteString(ddlSpace)

		if table.GetFieldByName(field.GetName()) == nil && table.GetRelationByName(field.GetName()) == nil {
			query.WriteString(postGreAddColumn)
			query.WriteString(ddlSpace)
			query.WriteString(postGreColumn)
			query.WriteString(field.GetDdl(table.provider, table.tableType))
		} else {
			query.WriteString(postGreDropColumn)
			query.WriteString(ddlSpace)
			query.WriteString(postGreColumn)
			query.WriteString(field.GetPhysicalName(table.provider))
			query.WriteString(ddlSpace)
			query.WriteString(postGreCascade)
		}
		return query.String()
	}
	return ""
}

func (table *Table) addVariables(query *strings.Builder) {

	variableName, index := table.getVariableInfo()
	count := len(table.fields) + int(table.searchableCount)
	for i := 0; i < count; i++ {
		query.WriteString(variableName)
		if table.provider == databaseprovider.PostgreSql {
			query.WriteString(strconv.Itoa(index))
		}
		// check last element
		if i < count-1 {
			query.WriteString(",")
		}
		index++
	}
}
func (table *Table) getVariableInfo() (string, int) {
	switch table.provider {
	case databaseprovider.PostgreSql:
		return postGreParameterName, 1
	case databaseprovider.MySql:
		return mysqlParameterName, 0
	}
	return "", 0
}
func (table *Table) getPhysicalName(physicalSchemaName string) string {
	var physicalName strings.Builder

	physicalName.Grow(len(table.name) + len(physicalSchemaName) + len(schemaSeparator) + len(physicalTablePrefix) + 2)
	physicalName.WriteString(physicalSchemaName)

	if len(physicalSchemaName) > 0 {
		physicalName.WriteString(schemaSeparator)
	}
	if table.tableType == tabletype.Business {
		physicalName.WriteString(sqlfmt.FormatEntityName(table.provider, physicalTablePrefix+table.name))
	} else {
		physicalName.WriteString(sqlfmt.FormatEntityName(table.provider, table.name))
	}
	return physicalName.String()
}

func (table *Table) getFieldList() string {
	// reduce memory usage
	// capacity
	var b strings.Builder
	var field *Field
	for i := 0; i < len(table.fields); i++ {
		field = table.fields[table.mapper[i]]
		b.WriteString(field.GetPhysicalName(table.provider))
		if i < len(table.fields)-1 {
			b.WriteString(fieldListSeparator)
		}
	}
	// add relations
	if table.tableType == tabletype.Mtm {
		relationsById := make([]*Relation, len(table.relations), len(table.relations))
		for i := 0; i < len(table.relations); i++ {
			relationsById[i] = table.relations[i]
		}
		sort.Slice(relationsById, func(i, j int) bool {
			return relationsById[i].id < relationsById[j].id
		})
		for i := 0; i < len(table.relations); i++ {
			b.WriteString(table.relations[i].GetPhysicalName(table.provider))
		}

	}
	return b.String()
}

func (table *Table) getUniqueKey() *Index {
	switch table.tableType {
	case tabletype.Meta:
		return table.GetIndexByName(metaTableName)
	case tabletype.MetaId:
		return table.GetIndexByName(metaIdTableName)

	}
	return nil
}

func (table *Table) getUniqueFieldList() string {
	// reduce memory usage
	// detect ==> capacity
	var b strings.Builder
	switch table.tableType {
	case tabletype.Business, tabletype.Lexicon:
		b.WriteString(table.GetPrimaryKey().GetPhysicalName(table.provider))
		break
	case tabletype.Meta, tabletype.MetaId, tabletype.LexiconItem:
		var ukIndex = table.GetIndexByName(table.name)
		for i := 0; i < len(ukIndex.fields); i++ {
			b.WriteString(table.GetFieldByName(ukIndex.fields[i]).GetPhysicalName(table.provider))
			b.WriteString(fieldListSeparator)
		}
		break
	case tabletype.Mtm:
		for i := 0; i < len(table.relations); i++ {
			b.WriteString(table.relations[i].GetPhysicalName(table.provider))
			b.WriteString(fieldListSeparator)
		}
		break
	}
	return strings.Trim(b.String(), fieldListSeparator)
}

// return -1 if not found
func (table *Table) findPrimaryKey(fields []Field) (int, *Field) {
	var invalidFieldCount = 0
	for i := 0; i < len(fields); i++ {
		if fields[i].IsValid() == false {
			invalidFieldCount++
			continue
		}
		if strings.EqualFold(fields[i].name, primaryKeyFieldName) {
			return i - invalidFieldCount, &fields[i]
		}
	}
	return -1, nil
}

func (table *Table) copyFields(fields []Field) {
	// copy fields
	for i := 0; i < len(fields); i++ {
		// append only valid fields
		if fields[i].IsValid() == true {
			table.fields = append(table.fields, &fields[i]) // sorted by name
		}
	}
}
func (table *Table) copyRelations(relations []Relation) {
	for i := 0; i < len(relations); i++ {
		// append only valid fields
		// relation are always valid
		table.relations = append(table.relations, &relations[i]) // sorted by name
	}

}
func (table *Table) copyIndexes(indexes []Index) {
	var validIndex = false

	for i := 0; i < len(indexes); i++ {
		// indexes are always valid

		if indexes[i].fields == nil || len(indexes[i].fields) == 0 {
			//TODO add logger here for i := 0; i < len(indexes); i++ {
			validIndex = false
		} else {
			validIndex = true
			for _, field := range indexes[i].fields {

				if table.GetFieldByName(field) == nil && table.GetRelationByName(field) == nil {
					//TODO add logger here for i := 0; i < len(indexes); i++ {
					validIndex = false
					break
				}
			}
		}
		if validIndex == true {
			table.indexes = append(table.indexes, &indexes[i])
		}
	}
}

func (table *Table) sortFields() {
	// sort structures
	if len(table.fields) > 1 {
		// sort fields by name
		sort.Slice(table.fields, func(i, j int) bool {
			return table.fields[i].name < table.fields[j].name
		})
	}
}

func (table *Table) sortRelations() {
	if len(table.relations) > 1 {
		// sort fields by id
		sort.Slice(table.relations, func(i, j int) bool {
			return table.relations[i].GetName() < table.relations[j].GetName()
		})
	}
}

func (table *Table) sortIndexes() {
	if len(table.indexes) > 1 {
		// sort fields by id
		sort.Slice(table.indexes, func(i, j int) bool {
			return table.indexes[i].name < table.indexes[j].name
		})
	}
}

func (table *Table) loadFields(fields []Field, tableType tabletype.TableType) {
	// copy slice -- func make([]T, len, cap) []T
	if fields != nil {
		var field = new(Field)
		var capacity = len(fields)
		var primaryKey *Field = nil
		var primaryKeyIndex = -1

		// missing primaryKey ? for business tables
		if tableType == tabletype.Business {
			primaryKeyIndex, primaryKey = table.findPrimaryKey(fields)
			if primaryKey == nil {
				capacity++
			}
		}

		// allow structures
		table.fields = make([]*Field, 0, capacity)

		// add missing primary key
		if primaryKey == nil && tableType == tabletype.Business {
			field.setType(fieldtype.Long)
			var defaultPrimaryKey = field.getDefaultPrimaryKey()
			table.fields = append(table.fields, defaultPrimaryKey) // sorted by name
		}

		table.copyFields(fields)

		// replace primary key
		if tableType == tabletype.Business && primaryKey != nil {
			field.setType(primaryKey.GetType())
			var defaultPrimaryKey = field.getDefaultPrimaryKey()
			table.fields[primaryKeyIndex] = defaultPrimaryKey
		}

		table.sortFields()
	} else {
		table.fields = make([]*Field, 0, 1)
	}
}

func (table *Table) loadRelations(relations []Relation) {
	table.relations = make([]*Relation, 0, len(relations))
	table.copyRelations(relations)
	table.sortRelations()
}

// check if each fields exist
func (table *Table) loadIndexes(indexes []Index) {
	table.indexes = make([]*Index, 0, len(indexes))
	table.copyIndexes(indexes)
	// replace searchable field by eg. s_name

	table.sortIndexes()
}

func (table *Table) loadMapper() {
	table.mapper = make([]uint16, len(table.fields), len(table.fields))
	fieldsById := make([]*Field, len(table.fields), len(table.fields))

	for i := 0; i < len(table.fields); i++ {
		fieldsById[i] = table.fields[i]
	}
	sort.Slice(fieldsById, func(i, j int) bool {
		return fieldsById[i].id < fieldsById[j].id
	})
	for i := 0; i < len(table.fields); i++ {
		table.mapper[i] = uint16(table.GetFieldIndexByName(fieldsById[i].name))
	}
}

// used by metaquery/ catalogue to generate variable value
func (table *Table) getVariableName(index int) string {
	result := ""
	switch table.provider {
	case databaseprovider.PostgreSql:
		result = postGreParameterName
		result += strconv.Itoa(index + 1)
		break
	case databaseprovider.MySql:
		result = mysqlParameterName
		break
	}
	return result
}

func (table *Table) getOptions() string {
	switch table.provider {
	case databaseprovider.PostgreSql:
		return postGreCreateOptions
	}
	return ""
}

func (table *Table) addPrimaryKeyFilter(query *strings.Builder, index int) {
	var variableName string
	var addedIndex = index

	switch table.provider {
	case databaseprovider.PostgreSql:
		variableName = postGreParameterName
		addedIndex++
		break
	case databaseprovider.MySql:
		variableName = mysqlParameterName
		break
	}
	switch table.tableType {
	case tabletype.Meta, tabletype.MetaId:
		//case :
		var ukIndex = table.getUniqueKey()
		for i := 0; i < len(ukIndex.fields); i++ {
			query.WriteString(table.GetFieldByName(ukIndex.fields[i]).GetPhysicalName(table.provider))
			query.WriteString(operatorEqual)
			query.WriteString(variableName)
			if table.provider != databaseprovider.MySql {
				query.WriteString(strconv.Itoa(addedIndex))
			}
			if i < len(ukIndex.fields)-1 {
				query.WriteString(filterSeparator)
				addedIndex++
			}
		}
		break
	case tabletype.Business, tabletype.Lexicon:
		query.WriteString(defaultPrimaryKeyInt64.GetName())
		query.WriteString(operatorEqual)
		query.WriteString(variableName)
		if table.provider != databaseprovider.MySql {
			query.WriteString(strconv.Itoa(addedIndex))
		}
		break
	}
}

// create physical table
func (table *Table) create(jobId int64) error {
	var metaQuery = metaQuery{}
	var schema = table.getSchema()
	var eventId int32 = 17
	var err error

	metaQuery.Init(schema, table)
	metaQuery.query = table.GetDdl(ddlstatement.Create, schema.findTablespace(table, nil, nil), nil)

	// create table
	err = metaQuery.create(eventId, jobId, table)
	if err != nil {
		return err
	}

	// create indexes except for @meta & meta_id tables
	table.createIndexes(jobId, schema)

	// we cannot create foreign key constraints: target table is may be not yet created
	// add primary key:
	if table.tableType != tabletype.Mtm {
		var primaryKey = new(constraint)
		var logger = schema.getLogger()
		primaryKey.Init(constrainttype.PrimaryKey, table, nil, nil)
		err := primaryKey.create(jobId, schema)
		if err != nil {
			logger.Error(-1, 0, err)
		}
	}
	return err
}

func (table *Table) drop(jobId int64) error {
	var metaquery = metaQuery{}
	var eventId int32 = 39
	var schema = table.getSchema()

	metaquery.Init(schema, table)
	metaquery.query = table.GetDdl(ddlstatement.Drop, nil, nil)
	//TODO drop contraints ??
	return metaquery.drop(eventId, jobId, table)
}

func (newTable *Table) alter(jobId int64, currentTable *Table) error {

	if newTable.fieldsEqual(currentTable) == false {
		newTable.alterFields(jobId, currentTable)
	}
	//TODO alter relation type!!!
	if newTable.relationsEqual(currentTable) == false {
		newTable.alterRelations(jobId, currentTable)
	}
	//!!! indexes after fields and relationships
	if newTable.indexesEqual(currentTable) == false {
		newTable.alterIndexes(jobId, currentTable)
	}
	return nil
}

func (newTable *Table) alterFields(jobId int64, currentTable *Table) error {
	//TODO rollback @meta if table alteration not working
	//TODO manage & test logical alteration only

	// generate dico map [physical_name] *field
	newDico := make(map[string]*Field, len(newTable.fields))
	curDico := make(map[string]*Field, len(currentTable.fields))
	newProvider := newTable.GetDatabaseProvider()
	currentProvider := currentTable.GetDatabaseProvider()

	// build dico
	for i := 0; i < len(newTable.fields); i++ {
		newDico[strings.ToUpper(newTable.fields[i].GetPhysicalName(newProvider))] =
			newTable.fields[i]
	}
	for i := 0; i < len(currentTable.fields); i++ {
		curDico[strings.ToUpper(currentTable.fields[i].GetPhysicalName(currentProvider))] =
			currentTable.fields[i]
	}
	// add new fields
	for i := 0; i < len(newTable.fields); i++ {
		if _, ok := curDico[strings.ToUpper(newTable.fields[i].GetPhysicalName(newProvider))]; !ok {
			newTable.addField(jobId, currentTable, newTable.fields[i])
		}
	}
	// drop obsolete fields
	for i := 0; i < len(currentTable.fields); i++ {
		if _, ok := newDico[strings.ToUpper(currentTable.fields[i].GetPhysicalName(currentProvider))]; !ok {
			// be sure field is available into currentTable
			newTable.dropField(jobId, currentTable, currentTable.fields[i])
		}
	}
	//TODO alter fields
	return nil
}

func (newTable *Table) alterIndexes(jobId int64, currentTable *Table) error {
	// generate dico map [index key] bool
	newDico := make(map[string]bool, len(newTable.fields))
	curDico := make(map[string]bool, len(currentTable.fields))
	schema := newTable.getSchema()

	var err error

	// generate dico of keys
	for i := 0; i < len(newTable.indexes); i++ {
		newDico[strings.ToUpper(newTable.indexes[i].getFieldList(newTable))] = true
	}
	for i := 0; i < len(currentTable.indexes); i++ {
		curDico[strings.ToUpper(currentTable.indexes[i].getFieldList(currentTable))] = true
	}
	// drop indexes!!! drop first
	for i := 0; i < len(currentTable.indexes); i++ {
		if _, ok := newDico[strings.ToUpper(currentTable.indexes[i].getFieldList(currentTable))]; !ok {
			err = currentTable.indexes[i].drop(jobId, schema, currentTable)
		}
	}
	// add new indexes
	for i := 0; i < len(newTable.indexes); i++ {
		if _, ok := curDico[strings.ToUpper(newTable.indexes[i].getFieldList(newTable))]; !ok {
			err = newTable.indexes[i].create(jobId, schema, newTable)
		}
	}

	//TODO modify index
	return err
}

func (newTable *Table) alterRelations(jobId int64, currentTable *Table) error {

	newDico := make(map[string]*Relation, len(newTable.relations))
	curDico := make(map[string]*Relation, len(currentTable.relations))
	newProvider := newTable.GetDatabaseProvider()
	currentProvider := currentTable.GetDatabaseProvider()

	var err error

	// generate dico of keys
	for i := 0; i < len(newTable.relations); i++ {
		newDico[strings.ToUpper(newTable.relations[i].GetPhysicalName(newProvider))] =
			newTable.relations[i]
	}
	for i := 0; i < len(currentTable.relations); i++ {
		curDico[strings.ToUpper(currentTable.relations[i].GetPhysicalName(currentProvider))] =
			currentTable.relations[i]
	}
	// add new relations
	for i := 0; i < len(newTable.relations); i++ {
		if _, ok := curDico[strings.ToUpper(newTable.relations[i].GetPhysicalName(newProvider))]; !ok {
			newTable.addRelation(jobId, currentTable, newTable.relations[i])
		}
	}
	// drop obsolete relations
	for i := 0; i < len(currentTable.relations); i++ {
		if _, ok := newDico[strings.ToUpper(currentTable.relations[i].GetPhysicalName(currentProvider))]; !ok {
			// be sure field is available into currentTable
			newTable.dropRelation(jobId, currentTable, currentTable.relations[i])
		}
	}
	//TODO alter relations
	return err
}

func (newTable *Table) dropField(jobId int64, currentTable *Table, field *Field) error {
	var queryRmv = new(metaQuery)
	var schema = currentTable.getSchema()
	var err error
	var eventId int32 = 43

	queryRmv.Init(schema, currentTable)
	queryRmv.query = currentTable.GetDdl(ddlstatement.Alter, nil, field)

	// alter table add field
	err = queryRmv.alter(eventId, jobId, currentTable, field)

	// searchable field?
	if err == nil && field.GetType() == fieldtype.String && field.IsCaseSensitive() == false {
		// generate ddl query 4 searchable field
		tableCopy := currentTable.Clone()
		searchableField := tableCopy.GetFieldByName(field.GetName())
		searchableField.setName(field.getSearchableFieldName())
		tableCopy.sortFields() // sort fields
		queryRmv = new(metaQuery)
		queryRmv.Init(schema, tableCopy)
		queryRmv.query = tableCopy.GetDdl(ddlstatement.Alter, nil, searchableField)
		err = queryRmv.alter(eventId, jobId, currentTable, searchableField)
	}
	return err
}

func (newTable *Table) addField(jobId int64, currentTable *Table, field *Field) error {

	var metaquery = new(metaQuery)
	var schema = currentTable.getSchema()
	var err error
	var eventId int32 = 41

	metaquery.Init(schema, newTable)
	metaquery.query = currentTable.GetDdl(ddlstatement.Alter, nil, field)

	// alter table add field
	err = metaquery.alter(eventId, jobId, newTable, field)

	// searchable field?
	if err == nil && field.GetType() == fieldtype.String && field.IsCaseSensitive() == false {
		tableCopy := newTable.Clone()
		searchableField := tableCopy.GetFieldByName(field.GetName())
		searchableField.setName(field.getSearchableFieldName())
		tableCopy.sortFields() // sort fields
		metaquery = new(metaQuery)
		metaquery.Init(schema, tableCopy)
		metaquery.query = currentTable.GetDdl(ddlstatement.Alter, nil, searchableField)
		err = metaquery.alter(eventId, jobId, newTable, searchableField)
	}

	return err
}

func (newTable *Table) dropRelation(jobId int64, currentTable *Table, relation *Relation) error {
	var relationType = relation.GetType()
	if relationType == relationtype.Mto || relationType == relationtype.Otop {
		return newTable.dropField(jobId, currentTable, relation.toField())
	}
	if relationType == relationtype.Mtm {
		// avoid two drops on mtm table take min relation id
		if relation.GetId() < relation.GetInverseRelation().GetId() {
			return relation.GetMtmTable().drop(jobId)
		}
	}
	return nil
}

func (newTable *Table) addRelation(jobId int64, currentTable *Table, relation *Relation) error {
	var metaquery = metaQuery{}
	var schema = currentTable.getSchema()
	var logger = schema.getLogger()
	var foreignKeyConstraint = new(constraint)
	var notNullConstraint = new(constraint)
	var relationType = relation.GetType()
	var eventId int32 = 45
	var err error

	metaquery.Init(schema, newTable)

	if relationType == relationtype.Mto || relationType == relationtype.Otop {
		var field = relation.toField()
		metaquery.query = currentTable.GetDdl(ddlstatement.Alter, nil, field)
		// alter table add relation
		err = metaquery.alter(eventId, jobId, currentTable, field)
	} else if relationType == relationtype.Mtm && relation.GetMtmTable().exists() == false {
		relation.GetMtmTable().create(jobId)
		relation.GetMtmTable().createConstraints(jobId, schema)
		return nil
	} else {
		// do nothing
		return nil
	}

	foreignKeyConstraint.Init(constrainttype.ForeignKey, newTable, nil, relation)
	notNullConstraint.Init(constrainttype.NotNull, newTable, nil, relation)

	if relation.HasConstraint() == true {
		err := foreignKeyConstraint.create(jobId, schema)
		if err != nil {
			logger.Error(-1, jobId, err)
		}
	}
	if relation.IsNotNull() == true {
		err := notNullConstraint.create(jobId, schema)
		if err != nil {
			logger.Error(-1, jobId, err)
		}
	}
	return err
}

func (table *Table) createIndexes(jobId int64, schema *Schema) {
	var logger = schema.getLogger()

	for i := 0; i < len(table.indexes); i++ {
		index := table.indexes[i]
		if (table.tableType == tabletype.Meta || table.tableType == tabletype.MetaId ||
			table.tableType == tabletype.LexiconItem) && index.IsUnique() {
			continue
		}
		err := index.create(jobId, schema, table)
		if err != nil {
			logger.Error(-1, 0, err)
		}
	}

}

func (table *Table) getInsertSql(provider databaseprovider.DatabaseProvider) string {
	var result strings.Builder

	result.WriteString(dmlstatement.Insert.String())
	result.WriteString(dmlSpace)
	result.WriteString(table.physicalName)
	result.WriteString(dmlInsertStart)
	// compute field list
	for i := 0; i < len(table.fields); i++ {
		var field = table.fields[table.mapper[i]]
		result.WriteString(field.GetPhysicalName(provider))
		if field.fieldType == fieldtype.String && field.caseSensitive == false {
			result.WriteString(fieldListSeparator)
			result.WriteString(field.getSearchableFieldName())
		}
		if i < len(table.fields)-1 {
			result.WriteString(fieldListSeparator)
		}
	}
	result.WriteString(dmlInsertValues)
	table.addVariables(&result)
	result.WriteString(dmlInsertEnd)

	return result.String()
}

func (table *Table) createConstraints(jobId int64, schema *Schema) {
	table.createFieldConstraints(jobId, schema)
	table.createRelationConstraints(jobId, schema)
}

func (table *Table) createFieldConstraints(jobId int64, schema *Schema) {
	var checkConstraint = new(constraint)
	checkConstraint.Init(constrainttype.Check, table, nil, nil)
	var notNullConstraint = new(constraint)
	notNullConstraint.Init(constrainttype.NotNull, table, nil, nil)
	var logger = schema.getLogger()

	// not null constraints
	for i := 0; i < len(table.fields); i++ {
		field := table.fields[i]
		checkConstraint.setField(field)
		err := checkConstraint.create(jobId, schema)
		if err != nil {
			logger.Error(-1, 0, err)
		}
		if field.IsNotNull() {
			notNullConstraint.setField(field)
			err = notNullConstraint.create(jobId, schema)
			if err != nil {
				logger.Error(-1, 0, err)
			}
		}
	}
}

func (table *Table) createRelationConstraints(jobId int64, schema *Schema) {
	var foreignKeyConstraint = new(constraint)
	var notNullConstraint = new(constraint)
	var logger = schema.getLogger()

	foreignKeyConstraint.Init(constrainttype.ForeignKey, table, nil, nil)
	notNullConstraint.Init(constrainttype.NotNull, table, nil, nil)

	// foreign Key constraints
	for i := 0; i < len(table.relations); i++ {
		relation := table.relations[i]

		if relation.GetType() != relationtype.Mto && relation.GetType() != relationtype.Otof {
			continue
		}

		foreignKeyConstraint.setRelation(relation)
		notNullConstraint.setRelation(relation)

		if relation.HasConstraint() == true {
			err := foreignKeyConstraint.create(jobId, schema)
			if err != nil {
				logger.Error(-1, 0, err)
			}
		}
		if relation.IsNotNull() == true {
			err := notNullConstraint.create(jobId, schema)
			if err != nil {
				logger.Error(-1, 0, err)
			}
		}
	}
}

// is physical table exists
func (table *Table) exists() bool {
	var schema = table.getSchema()
	cata := new(catalogue)
	return cata.exists(schema, table)
}

// used during upgrade to detect table physical changes
func (tableA *Table) equal(tableB *Table) bool {
	return tableA.fieldsEqual(tableB) && tableA.indexesEqual(tableB) && tableA.relationsEqual(tableB)
}

func (tableA *Table) fieldsEqual(tableB *Table) bool {
	if len(tableA.fields) != len(tableB.fields) {
		return false
	}
	var fieldA *Field
	var fieldB *Field
	for i := 0; i < len(tableA.fields); i++ {
		// id of field cannot change (create after creation)
		fieldA = tableA.GetFieldIdByIndex(i)
		fieldB = tableB.GetFieldIdByIndex(i)
		if fieldA.equal(fieldB) == false {
			return false
		}
	}
	return true
}

// ==> O(n²) complexity
func (tableA *Table) relationsEqual(tableB *Table) bool {
	if len(tableA.relations) != len(tableB.relations) {
		return false
	}
	var relationA *Relation
	var relationB *Relation
	for i := 0; i < len(tableA.relations); i++ {
		relationA = tableA.relations[i]
		relationB = tableB.GetRelationByNameI(relationA.name)
		if relationB == nil || relationA.equal(relationB) == false {
			return false
		}
	}
	return true
}

func (tableA *Table) indexesEqual(tableB *Table) bool {
	if len(tableA.indexes) != len(tableB.indexes) {
		return false
	}
	if len(tableA.indexes) > 0 {
		dico := make(map[int32]*Index, len(tableA.indexes))
		for i := 0; i < len(tableA.indexes); i++ {
			dico[tableA.indexes[i].id] = tableA.indexes[i]
		}
		for i := 0; i < len(tableB.indexes); i++ {
			if indexA, ok := dico[tableB.indexes[i].id]; ok {
				if !indexA.equal(tableB.indexes[i]) {
					fmt.Println(indexA.String())
					fmt.Println(tableB.indexes[i].String())
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

func (table *Table) getSchema() *Schema {
	var schema = GetSchemaById(table.schemaId)
	if schema == nil {
		schema = getUpgradingSchema()
	}
	return schema
}

func (table *Table) getMtmTable(schema *Schema, relation *Relation, name string) *Table {
	var fields = make([]Field, 0, 0)
	var relations = make([]Relation, 0, 2)
	var indexes = make([]Index, 1, 1)
	var result = new(Table)
	var inverseRelation = relation.GetInverseRelation()
	//var metaSchema = GetSchemaByName(metaSchemaName)
	//var indexIdSeq = metaSchema.GetSequenceByName(sequenceIndexId)
	var indexId int32 = 1 //FIX: int32(indexIdSeq.NextValue())
	// index
	var uk = Index{}
	var indexedFields = []string{relation.GetName(), inverseRelation.GetName()}

	uk.Init(indexId, relation.GetName(), "", indexedFields, false, true, true, true)
	indexes = append(indexes, uk)

	// relations
	var relationA = Relation{}
	var relationB = Relation{}

	relationA.Init(1, relation.GetName(), "", relation.GetToTable(), relationtype.Mto, relation.HasConstraint(), true, true, true)
	relationB.Init(2, inverseRelation.GetName(), "", inverseRelation.GetToTable(), relationtype.Mto, inverseRelation.HasConstraint(),
		true, true, true)

	relations = append(relations, relationA)
	relations = append(relations, relationB)

	result.Init(relation.GetId(), name, "", fields, relations, indexes, physicaltype.Table,
		schema.GetId(), schema.GetPhysicalName(), tabletype.Mtm, schema.GetDatabaseProvider(), "",
		true, false, true, true)

	return result
}

func (table *Table) getMetaIdTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields = make([]Field, 0, 4)
	var relations = make([]Relation, 0, 0)
	var indexes = make([]Index, 0, 0)
	var result = new(Table)

	// physical_name is built later
	//  == metaId table
	var id = Field{}
	var schemaId = Field{}
	var objectType = Field{}
	var value = Field{}
	var uk = Index{}

	// !!!! id field must be greater than 0 !!!!
	id.Init(1103, metaFieldId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1117, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1151, metaObjectType, "", fieldtype.Byte, 0, "", true, true, true, false, true)
	value.Init(1181, metaValue, "", fieldtype.Long, 0, "", true, true, true, false, true)

	var indexedFields = []string{metaFieldId, metaSchemaId, metaObjectType}
	uk.Init(1, metaIdTableName, "", indexedFields, false, true, true, true)

	fields = append(fields, id)
	fields = append(fields, schemaId)
	fields = append(fields, objectType)
	fields = append(fields, value)

	indexes = append(indexes, uk)

	result.Init(int32(tabletype.MetaId), metaIdTableName, "", fields, relations, indexes,
		physicaltype.Table, 0, schemaPhysicalName, tabletype.MetaId, provider, "", true, false, true, true)

	return result
}

func (table *Table) getMetaTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var result = new(Table)
	var uk Index
	var nuk Index

	// physical_name is built later
	//  == metaId table
	var id = Field{}
	var schemaId = Field{}
	var objectType = Field{}
	var referenceId = Field{}
	var dataType = Field{}

	var flags = Field{}
	var value = Field{}
	var name = Field{}
	var physicalName = Field{}
	var description = Field{}
	var active = Field{}

	// !!!! id field must be greater than 0 !!!!
	id.Init(1009, metaFieldId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1013, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1019, metaObjectType, "", fieldtype.Byte, 0, "", true, true, true, false, true)
	referenceId.Init(1021, metaReferenceId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	dataType.Init(1031, metaDataType, "", fieldtype.Int, 0, "", true, true, true, false, true)

	flags.Init(1039, metaFlags, "", fieldtype.Long, 0, "", true, true, true, false, true)
	name.Init(1061, metaName, "", fieldtype.String, 30, "", true, true, true, false, true)
	// metaName size * 2 ~ schema.name(max 30) + "." + table_name (max 28)
	description.Init(1087, metaDescription, "", fieldtype.String, 0, "", true, false, true, false, true)
	value.Init(1093, metaValue, "", fieldtype.String, 0, "", true, false, true, false, true)
	active.Init(1103, metaActive, "", fieldtype.Boolean, 0, "", true, true, true, false, true)

	// unique key (1)      id; schema_id; reference_id; object_type
	var indexedUkFields = []string{id.GetName(), schemaId.GetName(), objectType.GetName(), referenceId.GetName()}
	var indexedNukFields = []string{schemaId.GetName(), objectType.GetName()}

	nuk.Init(1, metaTableName, "", indexedNukFields, false, false, true, true)
	uk.Init(2, metaTableName, "", indexedUkFields, false, true, true, true)

	indexes = append(indexes, uk)
	indexes = append(indexes, nuk)

	fields = append(fields, id)           //1
	fields = append(fields, schemaId)     //2
	fields = append(fields, objectType)   //3
	fields = append(fields, referenceId)  //4
	fields = append(fields, dataType)     //5
	fields = append(fields, flags)        //6
	fields = append(fields, name)         //7
	fields = append(fields, physicalName) //8
	fields = append(fields, description)  //9
	fields = append(fields, value)        //10
	fields = append(fields, active)       //11

	result.Init(int32(tabletype.Meta), metaTableName, "", fields, relations, indexes,
		physicaltype.Table, 0, schemaPhysicalName, tabletype.Meta, provider, "", true, false, true, true)

	return result
}

// test table for unitesting

func (table *Table) getLogTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var result = new(Table)
	var idxEntryTime = Index{}
	var idxJobId = Index{}

	// physical_name is built later
	//  == metaId table
	var id = Field{}
	var entryTime = Field{}
	var levelId = Field{}
	var schemaId = Field{}
	var threadId = Field{}
	var callSite = Field{}
	var jobId = Field{}
	var method = Field{}
	var lineNumber = Field{}
	var message = Field{}
	var description = Field{}

	// "id","entry_time","level_id","thread_id","call_site","message","description","machine_name"
	id.Init(2111, metaLogId, "", fieldtype.Long, 0, "", true, true, true, false, true)
	entryTime.Init(2129, metaLogEntryTime, "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	levelId.Init(2131, metaLogLevelId, "", fieldtype.Short, 0, "", true, true, true, false, true)
	schemaId.Init(2137, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	threadId.Init(2143, metaLogThreadId, "", fieldtype.Int, 0, "", true, false, true, false, true)
	callSite.Init(2161, metaLogCallSite, "", fieldtype.String, 255, "", true, false, true, false, true)
	jobId.Init(2203, metaLogJobId, "", fieldtype.Long, 0, "", true, false, true, false, true)
	method.Init(2213, metaLogMethod, "", fieldtype.String, 80, "", true, false, true, false, true)
	lineNumber.Init(2221, metaLoglineNumber, "", fieldtype.Int, 0, "", true, false, true, false, true)
	message.Init(2237, metaLogMessage, "", fieldtype.String, 255, "", true, false, true, false, true)
	description.Init(2243, metaDescription, "", fieldtype.String, 0, "", true, false, true, false, true)

	fields = append(fields, id)          //1
	fields = append(fields, entryTime)   //2
	fields = append(fields, levelId)     //3
	fields = append(fields, threadId)    //4
	fields = append(fields, callSite)    //5
	fields = append(fields, jobId)       //6
	fields = append(fields, method)      //7
	fields = append(fields, schemaId)    //8
	fields = append(fields, message)     //9
	fields = append(fields, description) //10
	fields = append(fields, lineNumber)  //11

	// indexes
	var indexedFields = []string{entryTime.GetName()}
	idxEntryTime.Init(1, entryTime.GetName(), "", indexedFields, false, false, true, true)
	indexedFields = []string{jobId.GetName()}
	idxJobId.Init(2, jobId.GetName(), "", indexedFields, false, false, true, true)
	indexes = append(indexes, idxEntryTime)
	indexes = append(indexes, idxJobId)

	result.Init(int32(tabletype.Log), metaLogTableName, "", fields, relations, indexes, physicaltype.Table, 0, schemaPhysicalName,
		tabletype.Log, provider, "", false, false, true, true)
	return result
}

func (table *Table) getLexiconTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var result = new(Table)
	var idxUuid = Index{}
	var idxUk = Index{}
	var id = Field{}
	var schemaId = Field{}
	var name = Field{}
	var description = Field{}
	var uuid = Field{}
	var tableId = Field{}
	var fromFieldId = Field{}
	var toFieldId = Field{}
	var relationId = Field{}
	var relationValue = Field{}
	var updateStmp = Field{}
	var cached = Field{}
	var active = Field{}

	id.Init(5441, metaFieldId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(5449, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	name.Init(5477, metaName, "", fieldtype.String, 80, "", true, true, false, false, true)
	description.Init(5483, metaDescription, "", fieldtype.String, 0, "", true, false, true, false, true)
	uuid.Init(5503, metaUuid, "", fieldtype.String, 36, "", true, true, false, false, true)
	tableId.Init(5519, metaLexTableId, "", fieldtype.Int, 0, "", true, false, true, false, true)
	fromFieldId.Init(5527, metaLexFromFieldId, "", fieldtype.Int, 0, "", true, false, true, false, true)
	toFieldId.Init(5557, metaLexToFieldId, "", fieldtype.Int, 0, "", true, false, true, false, true)
	relationId.Init(5569, metaLexRelationId, "", fieldtype.Int, 0, "", true, false, true, false, true)
	relationValue.Init(5581, metaLexRelationValue, "", fieldtype.Long, 0, "", true, false, true, false, true)
	updateStmp.Init(5623, metaLexModifyStmp, "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	cached.Init(5641, metaCached, "", fieldtype.Boolean, 0, "", true, true, true, false, true)
	active.Init(5651, metaActive, "", fieldtype.Boolean, 0, "", true, true, true, false, true)

	fields = append(fields, id)            //1
	fields = append(fields, schemaId)      //2
	fields = append(fields, name)          //3
	fields = append(fields, description)   //4
	fields = append(fields, uuid)          //5
	fields = append(fields, tableId)       //6
	fields = append(fields, fromFieldId)   //7
	fields = append(fields, toFieldId)     //8
	fields = append(fields, relationId)    //9
	fields = append(fields, relationValue) //10
	fields = append(fields, updateStmp)    //11
	fields = append(fields, cached)        //12
	fields = append(fields, active)        //13

	/*
	   -----
	   unique key(1)      guid
	   unique key(2)      schema_id, table_id, from_field_id, to_field_id,relation_id, relation_value
	*/

	// indexes
	var indexedUuidFields = []string{uuid.GetName()}
	idxUuid.Init(1, uuid.GetName(), "", indexedUuidFields, false, true, true, true)

	var indexedUkFields = []string{schemaId.GetName(), tableId.GetName(), fromFieldId.GetName(),
		toFieldId.GetName(), relationId.GetName(), relationValue.GetName()}
	idxUk.Init(2, schemaId.GetName(), "", indexedUkFields, false, true, true, true)

	indexes = append(indexes, idxUuid)
	indexes = append(indexes, idxUk)

	result.Init(int32(tabletype.Lexicon), metaLexTableName, "", fields, relations, indexes, physicaltype.Table,
		0, schemaPhysicalName, tabletype.Lexicon, provider, "", false, false, true, true)

	return result
}

func (table *Table) getLexiconItemTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var result = new(Table)
	var idxUk = Index{}
	var id = Field{}
	var lexiconId = Field{}
	var langId = Field{}
	var value = Field{}

	id.Init(5801, metaFieldId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	lexiconId.Init(5813, metaLexItmLexId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	langId.Init(5827, metaLangId, "", fieldtype.Short, 0, "", true, true, true, false, true)
	value.Init(5843, metaValue, "", fieldtype.String, 0, "", true, true, false, false, true)

	/*
	   (0)  id,           // line of excel
	   (1)  lang_id,      // current language Id
	   (2)  lexicon_id,
	   (3)  value,
	*/

	fields = append(fields, id)        //1
	fields = append(fields, lexiconId) //2
	fields = append(fields, langId)    //3
	fields = append(fields, value)     //4

	/*
	   -----
	   unique key(1)    id, lexicon_id, lang_id
	*/

	// indexes
	var indexedUkFields = []string{id.GetName(), lexiconId.GetName(), langId.GetName()}
	idxUk.Init(1, metaLexItmTableName, "", indexedUkFields, false, true, true, true)
	indexes = append(indexes, idxUk)

	result.Init(int32(tabletype.LexiconItem), metaLexItmTableName, "", fields, relations, indexes, physicaltype.Table,
		0, schemaPhysicalName, tabletype.LexiconItem, provider, "", false, false, true, true)

	return result
}

func (table *Table) getSearchableCount() uint8 {
	var result uint8 = 0
	for i := 0; i < len(table.fields); i++ {
		var field = table.fields[i]
		if field.fieldType == fieldtype.String && field.caseSensitive == false {
			result++
		}
	}
	return result
}

// table used to get result from aggregate queries
func (table *Table) getLongTable() *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var result = new(Table)

	// physical_name is built later
	//  == metaId table
	var value = Field{}
	value.Init(4283, metaValue, "", fieldtype.Long, 0, "", true, true, true, false, true)
	fields = append(fields, value) //1

	result.Init(int32(tabletype.Logical), metaLongTableName, "", fields, relations, indexes, physicaltype.Logical, 0, "",
		tabletype.Logical, databaseprovider.Undefined, "", false, false, true, true)
	return result
}

func (table *Table) getTable(provider databaseprovider.DatabaseProvider, physicalSchemaName string,
	schemaId int32, metaList []*meta) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var metaTable *meta
	var result = new(Table)
	var fieldCount = 0
	var relationCount = 0
	var indexCount = 0

	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		var metaType = metaData.GetEntityType()
		switch metaType {
		case entitytype.Table:
			metaTable = metaData
			break
		case entitytype.Index:
			indexCount++
			break
		case entitytype.Field:
			fieldCount++
			break
		case entitytype.Relation:
			relationCount++
			break
		}
	}

	fields = make([]Field, 0, fieldCount)
	relations = make([]Relation, 0, relationCount)
	indexes = make([]Index, 0, indexCount)

	for i := 0; i < len(metaList); i++ {
		var metaData = metaList[i]
		var metaType = metaData.GetEntityType()
		switch metaType {
		case entitytype.Index:
			indexes = append(indexes, *metaData.toIndex())
			break
		case entitytype.Field:
			fields = append(fields, *metaData.toField())
			break
		case entitytype.Relation:
			// load table information later
			relations = append(relations, *metaData.toRelation(nil))
			break
		}
	}

	result.Init(metaTable.id, metaTable.name, metaTable.description, fields, relations, indexes, physicaltype.Table, schemaId,
		physicalSchemaName, tabletype.Business, provider, metaTable.value, metaTable.IsTableCached(),
		metaTable.IsTableReadonly(), metaTable.IsEntityBaseline(), metaTable.enabled)

	return result
}
