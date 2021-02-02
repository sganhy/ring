package schema

import (
	"fmt"
	"os"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"sort"
	"strings"
)

// call exemple elemr.Init(21, "rel test", "hellkzae", "hell1", "52", nil, relationtype.Mto, false, true, false)
type Table struct {
	id           int32
	name         string
	description  string
	fields       []*Field    // sorted by name
	fieldsById   []*Field    // sorted by id
	relations    []*Relation // sorted by name
	indexes      []*Index    // sorted by name
	physicalName string
	physicalType physicaltype.PhysicalType
	schemaId     int32
	tableType    tabletype.TableType
	subject      string
	cacheId      CacheId
	cached       bool
	readonly     bool
	baseline     bool
	active       bool
	//internal readonly LexiconIndex[] LexiconIndexes;
}

const postgreSqlSchema string = "information_schema"
const fieldNotFound int = -1
const metaSchemaId string = "schema_id"
const metaId string = "id"
const metaObjectType string = "object_type"
const metaValue string = "value"
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

// // call exemple elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes, "schema.t_site", physicaltype.Table, 64, tabletype.Business, "subject test", true, false, true, false)
func (table *Table) Init(id int32, name string, description string, fields []Field, relations []Relation, indexes []Index, physicalName string,
	physicalType physicaltype.PhysicalType, schemaId int32, tableType tabletype.TableType, subject string,
	cached bool, readonly bool, baseline bool, active bool) {
	table.id = id
	table.name = name
	table.description = description
	table.loadFields(fields, tableType)
	table.loadRelations(relations)
	table.loadIndexes(indexes) // run at the end only
	// initialize cacheId
	table.cacheId.CurrentId = 0
	table.cacheId.MaxId = 0
	table.cacheId.ReservedRange = 0
	table.physicalName = physicalName
	table.tableType = tableType
	table.schemaId = schemaId
	table.tableType = tableType
	table.subject = subject
	table.physicalType = physicalType
	table.cached = cached
	table.readonly = readonly
	table.baseline = baseline
	table.active = active
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
	return &table.cacheId
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

//******************************
// public methods
//******************************
func (table *Table) GetFieldByName(name string) *Field {
	var indexerLeft, indexerRigth, indexerMiddle, indexerCompare int = 0, len(table.fields) - 1, 0, 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.fields[indexerMiddle].name)
		if indexerCompare == 0 {
			return table.fields[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
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
	var indexerLeft, indexerRigth, indexerMiddle int = 0, len(table.fieldsById) - 1, 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		if table.fieldsById[indexerMiddle].id == id {
			return table.fieldsById[indexerMiddle]
		} else if table.fieldsById[indexerMiddle].id > id {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	return nil
}

//return -1 if not found
func (table *Table) GetFieldIndexByName(name string) int {
	var indexerLeft, indexerRigth, indexerMiddle, indexerCompare int = 0, len(table.fields) - 1, 0, 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.fields[indexerMiddle].name)
		if indexerCompare == 0 {
			return indexerMiddle
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	return fieldNotFound
}

func (table *Table) GetRelationByName(name string) *Relation {
	var indexerLeft, indexerRigth, indexerMiddle, indexerCompare int = 0, len(table.relations) - 1, 0, 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.relations[indexerMiddle].name)
		if indexerCompare == 0 {
			return table.relations[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	return nil
}

func (table *Table) GetIndexByName(name string) *Index {
	var indexerLeft, indexerRigth, indexerMiddle, indexerCompare int = 0, len(table.indexes) - 1, 0, 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.indexes[indexerMiddle].name)
		if indexerCompare == 0 {
			return table.indexes[indexerMiddle]
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	return nil
}

//return -1 if not found
func (table *Table) GetRelationIndexByName(name string) int {
	var indexerLeft, indexerRigth, indexerMiddle, indexerCompare int = 0, len(table.relations) - 1, 0, 0
	for indexerLeft <= indexerRigth {
		indexerMiddle = indexerLeft + indexerRigth
		indexerMiddle >>= 1 // indexerMiddle <-- indexerMiddle /2
		indexerCompare = strings.Compare(name, table.relations[indexerMiddle].name)
		if indexerCompare == 0 {
			return indexerMiddle
		} else if indexerCompare > 0 {
			indexerLeft = indexerMiddle + 1
		} else {
			indexerRigth = indexerMiddle - 1
		}
	}
	return fieldNotFound
}

func (table *Table) GetPrimaryKey() *Field {
	if len(table.fieldsById) > 0 && table.tableType == tabletype.Business {
		return table.fieldsById[0]
	}
	return nil
}

func (table *Table) GetSql(provider databaseprovider.DatabaseProvider, tablespace *Tablespace) string {
	return ""
}

func (table *Table) ToFile(filename string) {
	f, _ := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	for i := 0; i < len(table.fieldsById); i++ {
		field := table.fieldsById[i]
		message := fmt.Sprintf("field (%d) name=%s \n", field.GetId(), field.GetName())
		f.WriteString(message)
	}
}

//******************************
// private methods
//******************************

// return -1 if not found
func findPrimaryKey(fields []Field) (int, *Field) {
	var invalidFieldCount int = 0
	for i := 0; i < len(fields); i++ {
		if fields[i].IsValid() == false {
			invalidFieldCount++
			continue
		}
		if strings.EqualFold(fields[i].name, primaryKeyFielName) {
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
		var capacity int = len(fields)
		var primaryKey *Field = nil
		var primaryKeyIndex int = -1

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
	// allo structure
	table.relations = make([]*Relation, 0, len(relations))
	table.copyRelations(relations)
	table.sortRelations()
}

//!!! Must be ran aftet loadFields() method
// check if each fields exist
func (table *Table) loadIndexes(indexes []Index) {
	table.indexes = make([]*Index, 0, len(indexes))
	table.copyIndexes(indexes)
	table.sortIndexes()
}

func GetMetaIdTable(provider databaseprovider.DatabaseProvider) *Table {
	var fields = make([]Field, 0, 16)
	var relations = make([]Relation, 0, 16)
	var indexes = make([]Index, 0, 16)
	var table = new(Table)

	// physical_name is builded later
	//  == metaId table
	var id Field = Field{}
	var schemaId Field = Field{}
	var objectType Field = Field{}
	var value Field = Field{}
	var uk Index = Index{}

	// elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true, true)
	// !!!! id field must be greater than 0 !!!!
	id.Init(1103, metaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1117, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1151, metaObjectType, "", fieldtype.Byte, 0, "", true, true, true, false, true)
	value.Init(1181, metaValue, "", fieldtype.Long, 0, "", true, true, true, false, true)

	// elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
	var indexedFields = []string{metaId, metaSchemaId, metaObjectType}
	uk.Init(1, "pk_@meta_id", "", indexedFields, int32(tabletype.MetaId), false, true, true, true)

	fields = append(fields, id)
	fields = append(fields, schemaId)
	fields = append(fields, objectType)
	fields = append(fields, value)

	indexes = append(indexes, uk)

	table.Init(int32(tabletype.MetaId), metaIdTableName, "", fields, relations, indexes, getPhysicalName(provider, metaIdTableName),
		physicaltype.Table, 0, tabletype.MetaId, "", true, false, true, true)

	return table
}

func GetMetaTable(provider databaseprovider.DatabaseProvider) *Table {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)
	var uk Index = Index{}

	// physical_name is builded later
	//  == metaId table
	var id Field = Field{}
	var schemaId Field = Field{}
	var objectType Field = Field{}
	var referenceId Field = Field{}
	var dataType Field = Field{}

	var flags Field = Field{}
	var value Field = Field{}
	var name Field = Field{}
	var description Field = Field{}
	var active Field = Field{}

	// elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true, true)
	// !!!! id field must be greater than 0 !!!!
	id.Init(1009, metaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1013, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1019, metaObjectType, "", fieldtype.Byte, 0, "", true, true, true, false, true)
	referenceId.Init(1021, "reference_id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	dataType.Init(1031, "data_type", "", fieldtype.Int, 0, "", true, false, true, false, true)

	flags.Init(1039, "flags", "", fieldtype.Long, 0, "", true, false, true, false, true)
	name.Init(1061, "name", "", fieldtype.String, 30, "", true, false, true, false, true)
	description.Init(1069, "description", "", fieldtype.String, 0, "", true, false, true, false, true)
	value.Init(1087, metaValue, "", fieldtype.String, 0, "", true, false, true, false, true)
	active.Init(1093, "active", "", fieldtype.Boolean, 0, "", true, false, true, false, true)

	// elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
	// unique key (1)      id; schema_id; reference_id; object_type
	var indexedFields = []string{id.name, schemaId.name, objectType.name, referenceId.name}
	uk.Init(1, "pk_@meta", "", indexedFields, int32(tabletype.MetaId), false, true, true, true)

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

	table.Init(int32(tabletype.MetaId), metaTableName, "", fields, relations, indexes, getPhysicalName(provider, metaTableName),
		physicaltype.Table, 0, tabletype.MetaId, "", true, false, true, true)

	return table
}

func GetLogTable() *Table {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)

	// physical_name is builded later
	//  == metaId table
	var id Field = Field{}
	var entryTime Field = Field{}
	var levelId Field = Field{}
	var schemaId Field = Field{}
	var threadId Field = Field{}
	var callSite Field = Field{}
	var jobId Field = Field{}
	var method Field = Field{}
	var message Field = Field{}

	// elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true, true)
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

	table.Init(int32(tabletype.MetaId), metaTableName, "", fields, relations, indexes, "", physicaltype.Table, 0, tabletype.MetaId, "",
		false, false, true, true)
	return table
}

func getPhysicalName(provider databaseprovider.DatabaseProvider, name string) string {
	var physicalName = ""
	//TODO implement other provider
	switch provider {
	case databaseprovider.PostgreSql:
		physicalName = postgreSqlSchema + ".\"" + name + "\""
	case databaseprovider.MySql:
		physicalName = "MySql"
	case databaseprovider.Oracle:
		physicalName = "Oracle"
	case databaseprovider.Sqlite3:
		physicalName = "Sqlite"
	}
	return physicalName
}
