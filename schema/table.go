package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Table struct {
	id           int32
	name         string
	description  string
	fields       []*Field    // sorted by name
	fieldsById   []*Field    // sorted by id
	relations    []*Relation // sorted by name
	indexes      []*Index    // sorted by name
	mapper       []uint16    // mapping from .fieldsById to .fields ==> max number of column 65535!!
	physicalName string
	physicalType physicaltype.PhysicalType
	schemaId     int32
	tableType    tabletype.TableType
	fieldList    string
	subject      string
	provider     databaseprovider.DatabaseProvider
	cacheId      *CacheId
	sqlCapacity  uint16
	cached       bool
	readonly     bool
	baseline     bool
	active       bool
	//internal readonly LexiconIndex[] LexiconIndexes;
}

const createTableSql string = "CREATE TABLE %s (\n%s\n)"
const fieldNotFound int = -1
const relationNotFound int = -1
const metaSchemaId string = "schema_id"
const metaId string = "id"
const metaObjectType string = "object_type"
const metaDataType string = "data_type"
const metaDescription string = "description"
const metaValue string = "value"
const metaReferenceId string = "reference_id"
const metaName string = "name"
const metaFlags string = "flags"
const metaIdTableName string = "@meta_id"
const metaTableName string = "@meta"
const metaLogId string = "id"
const metaLogEntryTime string = "entry_time"
const metaLogLevelId string = "level_id"
const metaLogThreadId string = "thread_id"
const metaLogCallSite string = "call_site"
const metaLogJobId string = "job_id"
const metaLogMethod string = "method"
const metaLogMessage string = "message"
const fieldListSeparator string = ","
const dmlInsert = "INSERT INTO "
const dmlInsertValues = ") VALUES ("
const dmlInsertStart = " ("
const dmlInsertEnd = ")"
const dqlSelect = "SELECT "
const dqlFrom = " FROM "
const dqlWhere = " WHERE "
const dqlOrderBy = " ORDER BY "
const maxNumberOfColumn = 65535
const postGreParameterName = "$"

func (table *Table) Init(id int32, name string, description string, fields []Field, relations []Relation, indexes []Index,
	physicalType physicaltype.PhysicalType, schemaId int32, schemaPhysicalName string, tableType tabletype.TableType, provider databaseprovider.DatabaseProvider,
	subject string, cached bool, readonly bool, baseline bool, active bool) {
	table.id = id
	table.name = name
	table.sqlCapacity = 16 // min value
	table.description = description
	table.loadFields(fields, tableType)
	table.loadRelations(relations)
	table.loadMapper()         // !!!load after loadFields
	table.loadIndexes(indexes) //!!!! run at the end only
	// initialize cacheId
	table.cacheId = new(CacheId)
	table.cacheId.CurrentId = 0
	table.cacheId.MaxId = 0
	table.cacheId.ReservedRange = 0
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
	if provider != databaseprovider.NotDefined {
		table.physicalName = table.getPhysicalName(provider, schemaPhysicalName)
		table.fieldList = table.getFieldList(provider)
		table.loadSqlCapacity(provider) // !!!load after loadFields
	}
}

//******************************
// getters
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

func (table *Table) GetCacheId() *CacheId {
	return table.cacheId
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
	for i := len(table.fieldsById) - 1; i >= 0; i-- {
		if strings.EqualFold(name, table.fields[i].name) {
			return table.fields[i]
		}
	}
	return nil
}

func (table *Table) GetFieldById(id int32) *Field {
	var indexerLeft, indexerRight, indexerMiddle = 0, len(table.fieldsById) - 1, 0
	for indexerLeft <= indexerRight {
		indexerMiddle = indexerLeft + indexerRight
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2

		if id == table.fieldsById[indexerMiddle].id {
			return table.fieldsById[indexerMiddle]
		} else if id > table.fieldsById[indexerMiddle].id {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRight = indexerMiddle - 1
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
func (table *Table) GetFieldByIndex(index int) *Field {
	return table.fields[index]
}

//
func (table *Table) GetFieldIdByIndex(index int) *Field {
	return table.fieldsById[index]
}

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
	if len(table.fieldsById) > 0 && table.tableType == tabletype.Business {
		return table.fieldsById[0]
	}
	return nil
}

func (table *Table) GetDdl(tablespace *Tablespace) string {
	var fields []string
	for i := 0; i < len(table.fields); i++ {
		fieldSql := table.fields[i].GetDdlSql(table.provider, table.tableType)
		fields = append(fields, fieldSql)
	}
	for i := 0; i < len(table.relations); i++ {
		relationSql := table.relations[i].GetDdlSql(table.provider)
		fields = append(fields, relationSql)
	}
	sql := fmt.Sprintf(createTableSql, table.physicalName, strings.Join(fields, ",\n"))
	if tablespace != nil {
		sql = sql + " " + getDdlTableSpace(table.provider, tablespace)
	}
	return sql
}

func (table *Table) GetDml(dmlType dmlstatement.DmlStatement) string {
	var result strings.Builder
	switch dmlType {
	case dmlstatement.Insert:
		result.Grow(int(table.sqlCapacity))
		result.WriteString(dmlInsert)
		result.WriteString(table.physicalName)
		result.WriteString(dmlInsertStart)
		result.WriteString(table.fieldList)
		result.WriteString(dmlInsertValues)
		table.addVariables(&result)
		result.WriteString(dmlInsertEnd)
		//fmt.Printf("GetDml() len(str)/sql.Cap() %d /%d\n", len(result.String()), result.Cap())
		break
	case dmlstatement.Update:
		break
	case dmlstatement.Delete:
		break
	}
	return result.String()
}

// SELECT
func (table *Table) GetDql(whereClause string, orderClause string) string {
	capacity := len(dqlSelect) + len(dqlFrom) + len(table.fieldList) + len(table.physicalName)

	if whereClause != "" {
		capacity += len(dqlWhere)
		capacity += len(whereClause)
	}

	if orderClause != "" {
		capacity += len(dqlOrderBy)
		capacity += len(orderClause)
	}

	var result strings.Builder
	result.Grow(capacity)
	result.WriteString(dqlSelect)
	result.WriteString(table.fieldList)
	result.WriteString(dqlFrom)
	result.WriteString(table.physicalName)

	if whereClause != "" {
		result.WriteString(dqlWhere)
		result.WriteString(whereClause)
	}

	if orderClause != "" {
		result.WriteString(dqlOrderBy)
		result.WriteString(orderClause)
	}
	// check capacity
	//fmt.Printf("GetDql() capacity/len(str)/sql.Cap() %d /%d /%d\n", capacity, len(sql.String()), sql.Cap())
	return result.String()
}

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

func (table *Table) ToMeta() []*Meta {
	var capacity = 1 + len(table.fields) + len(table.relations) + len(table.indexes)
	var result = make([]*Meta, 0, capacity)
	var metaTable = new(Meta)

	// key
	metaTable.id = table.id
	metaTable.refId = 0
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

	result = append(result, metaTable)

	for i := 0; i < len(table.fields); i++ {
		result = append(result, table.fields[i].ToMeta(table.id))
	}
	for i := 0; i < len(table.relations); i++ {
		result = append(result, table.relations[i].ToMeta(table.id))
	}
	for i := 0; i < len(table.indexes); i++ {
		result = append(result, table.indexes[i].ToMeta(table.id))
	}
	return result
}

func (table *Table) GetQueryResult(columnPointer []interface{}) []string {
	var capacity = len(table.fields)
	var result = make([]string, capacity, capacity)
	var currentField *Field
	var value interface{}
	var strValue string
	var index int

	// manage fields first
	for i := 0; i < capacity; i++ {
		value = *columnPointer[i].(*interface{})
		currentField = table.fieldsById[i]
		// use a mapper instead
		//index = table.GetFieldIndexByName(currentField.name)
		index = int(table.mapper[i])
		if value == nil {
			strValue = ""
		} else {
			switch value.(type) {
			case string:
				strValue = value.(string)
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

//4debug
func GetTable() *Table {
	return getLogTable(databaseprovider.PostgreSql, "information_schema")
}

//******************************
// private methods
//******************************
func (table *Table) addVariables(query *strings.Builder) {
	var variableName string
	var index = 0
	switch table.provider {
	case databaseprovider.PostgreSql:
		variableName = postGreParameterName
		index = 1
		break
	}
	//TODO improve performance ==> switch outside of loop
	for i := 0; i < len(table.fields); i++ {
		query.WriteString(variableName)
		query.WriteString(strconv.Itoa(index))
		// check last element
		if i < len(table.fields)-1 {
			query.WriteString(",")
		}
		index++
	}
}

func (table *Table) getPhysicalName(provider databaseprovider.DatabaseProvider, physicalSchemaName string) string {
	var physicalName = ""
	//TODO implement other provider
	switch provider {
	case databaseprovider.PostgreSql:
		if table.tableType != tabletype.Business {
			physicalName = physicalSchemaName + ".\"" + table.name + "\""
		} else {
			physicalName = physicalSchemaName + "." + table.name
		}
		break
	case databaseprovider.MySql:
		physicalName = "MySql"
		break
	case databaseprovider.Oracle:
		physicalName = "Oracle"
		break
	}
	return physicalName
}

func (table *Table) getFieldList(provider databaseprovider.DatabaseProvider) string {
	// reduce memory usage
	// capacity
	var b strings.Builder
	for i := 0; i < len(table.fieldsById); i++ {
		b.WriteString(table.fieldsById[i].GetPhysicalName(provider))
		if i < len(table.fieldsById)-1 {
			b.WriteString(fieldListSeparator)
		}
	}
	return b.String()
}

// return -1 if not found
func findPrimaryKey(fields []Field) (int, *Field) {
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
			table.fields = append(table.fields, &fields[i])         // sorted by name
			table.fieldsById = append(table.fieldsById, &fields[i]) // sorted by id
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
				if table.GetFieldIndexByName(field) == fieldNotFound {
					//TODO add logger here for i := 0; i < len(indexes); i++ {
					validIndex = false
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
		// sort fields by id
		sort.Slice(table.fieldsById, func(i, j int) bool {
			return table.fieldsById[i].id < table.fieldsById[j].id
		})
	}
}

func (table *Table) sortRelations() {
	if len(table.relations) > 1 {
		// sort fields by id
		sort.Slice(table.relations, func(i, j int) bool {
			return table.relations[i].name < table.relations[j].name
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
		var capacity = len(fields)
		var primaryKey *Field = nil
		var primaryKeyIndex = -1

		// missing primaryKey ? for business tables
		if tableType == tabletype.Business {
			primaryKeyIndex, primaryKey = findPrimaryKey(fields)
			if primaryKey == nil {
				capacity++
			}
		}

		// allow structures
		table.fields = make([]*Field, 0, capacity)
		table.fieldsById = make([]*Field, 0, capacity)

		// add missing primary key
		if primaryKey == nil && tableType == tabletype.Business {
			var defaultPrimaryKey = getDefaultPrimaryKey(fieldtype.Long)
			table.fields = append(table.fields, defaultPrimaryKey)         // sorted by name
			table.fieldsById = append(table.fieldsById, defaultPrimaryKey) // sorted by id
		}

		table.copyFields(fields)

		// replace primary key
		if tableType == tabletype.Business && primaryKey != nil {
			var defaultPrimaryKey = getDefaultPrimaryKey(primaryKey.fieldType)
			table.fields[primaryKeyIndex] = defaultPrimaryKey
			table.fieldsById[primaryKeyIndex] = defaultPrimaryKey
		}

		table.sortFields()
	} else {
		//TODO throw an error + logging
		table.fields = make([]*Field, 0, 1)
		table.fieldsById = make([]*Field, 0, 1)
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
	table.sortIndexes()
}

func (table *Table) loadMapper() {
	table.mapper = make([]uint16, len(table.fieldsById), len(table.fieldsById))
	for i := 0; i < len(table.fieldsById); i++ {
		table.mapper[i] = uint16(table.GetFieldIndexByName(table.fieldsById[i].name))
	}
}

func (table *Table) loadSqlCapacity(provider databaseprovider.DatabaseProvider) {
	table.sqlCapacity = 0
	table.sqlCapacity = uint16(len(dmlInsert) + len(table.physicalName) + len(dmlInsertStart) + len(table.fieldList))
	table.sqlCapacity += uint16(len(dmlInsertValues) + len(dmlInsertEnd))
	var capacity uint16 = 0
	for i := 0; i < len(table.fields); i++ {
		switch provider {
		case databaseprovider.PostgreSql:
			capacity += uint16(len(postGreParameterName))
			capacity += uint16(len(strconv.Itoa(i + 1)))
			break
		}
		// check last element
		if i < len(table.fields)-1 {
			capacity++
		}
	}
	table.sqlCapacity += capacity
}

func getMetaIdTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields = make([]Field, 0, 16)
	var relations = make([]Relation, 0, 16)
	var indexes = make([]Index, 0, 16)
	var table = new(Table)

	// physical_name is built later
	//  == metaId table
	var id = Field{}
	var schemaId = Field{}
	var objectType = Field{}
	var value = Field{}
	var uk = Index{}

	// !!!! id field must be greater than 0 !!!!
	id.Init(1103, metaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1117, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1151, metaObjectType, "", fieldtype.Byte, 0, "", true, true, true, false, true)
	value.Init(1181, metaValue, "", fieldtype.Long, 0, "", true, true, true, false, true)

	var indexedFields = []string{metaId, metaSchemaId, metaObjectType}
	uk.Init(1, "pk_@meta_id", "", indexedFields, false, true, true, true)

	fields = append(fields, id)
	fields = append(fields, schemaId)
	fields = append(fields, objectType)
	fields = append(fields, value)

	indexes = append(indexes, uk)

	table.Init(int32(tabletype.MetaId), metaIdTableName, "", fields, relations, indexes,
		physicaltype.Table, 0, schemaPhysicalName, tabletype.MetaId, provider, "", true, false, true, true)

	return table
}

func getMetaTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var table = new(Table)
	var uk Index

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
	var description = Field{}
	var active = Field{}

	// !!!! id field must be greater than 0 !!!!
	id.Init(1009, metaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1013, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1019, metaObjectType, "", fieldtype.Byte, 0, "", true, true, true, false, true)
	referenceId.Init(1021, metaReferenceId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	dataType.Init(1031, metaDataType, "", fieldtype.Int, 0, "", true, false, true, false, true)

	flags.Init(1039, metaFlags, "", fieldtype.Long, 0, "", true, true, true, false, true)
	name.Init(1061, metaName, "", fieldtype.String, 30, "", true, true, true, false, true)
	description.Init(1069, metaDescription, "", fieldtype.String, 0, "", true, false, true, false, true)
	value.Init(1087, metaValue, "", fieldtype.String, 0, "", true, false, true, false, true)
	active.Init(1093, "active", "", fieldtype.Boolean, 0, "", true, true, true, false, true)

	// unique key (1)      id; schema_id; reference_id; object_type
	var indexedFields = []string{id.name, schemaId.name, objectType.name, referenceId.name}
	uk.Init(1, "pk_@meta", "", indexedFields, false, true, true, true)

	fields = append(fields, id)          //1
	fields = append(fields, schemaId)    //2
	fields = append(fields, objectType)  //3
	fields = append(fields, referenceId) //4
	fields = append(fields, dataType)    //5
	fields = append(fields, flags)       //6
	fields = append(fields, name)        //7
	fields = append(fields, description) //8
	fields = append(fields, value)       //9
	fields = append(fields, active)      //10

	indexes = append(indexes, uk)

	table.Init(int32(tabletype.Meta), metaTableName, "", fields, relations, indexes,
		physicaltype.Table, 0, schemaPhysicalName, tabletype.MetaId, provider, "", true, false, true, true)

	return table
}

func getLogTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var table = new(Table)

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
	var message = Field{}

	// "id","entry_time","level_id","thread_id","call_site","message","description","machine_name"
	id.Init(2111, metaLogId, "", fieldtype.Long, 0, "", true, true, true, false, true)
	entryTime.Init(2129, metaLogEntryTime, "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	levelId.Init(2131, metaLogLevelId, "", fieldtype.Short, 0, "", true, false, true, false, true)
	schemaId.Init(2137, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	threadId.Init(2143, metaLogThreadId, "", fieldtype.Long, 0, "", true, false, true, false, true)
	callSite.Init(2161, metaLogCallSite, "", fieldtype.String, 255, "", true, false, true, false, true)
	jobId.Init(2203, metaLogJobId, "", fieldtype.Long, 0, "", true, false, true, false, true)
	method.Init(2213, metaLogMethod, "", fieldtype.String, 80, "", true, false, true, false, true)
	message.Init(2237, metaLogMessage, "", fieldtype.String, 255, "", true, false, true, false, true)

	fields = append(fields, id)        //1
	fields = append(fields, entryTime) //2
	fields = append(fields, levelId)   //3
	fields = append(fields, threadId)  //4
	fields = append(fields, callSite)  //5
	fields = append(fields, jobId)     //6
	fields = append(fields, method)    //7
	fields = append(fields, schemaId)  //8
	fields = append(fields, message)   //9

	table.Init(int32(tabletype.Log), "@log", "", fields, relations, indexes, physicaltype.Table, 0, schemaPhysicalName,
		tabletype.MetaId, provider, "", false, false, true, true)
	return table
}

func getDdlTableSpace(provider databaseprovider.DatabaseProvider, tablespace *Tablespace) string {
	var sql string
	switch provider {
	case databaseprovider.PostgreSql:
		sql = "TABLESPACE " + tablespace.name
		break
	case databaseprovider.MySql:
		sql = ""
		break
	case databaseprovider.Oracle:
		sql = ""
		break
	}
	return sql
}
