package schema

import (
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"strings"
	"sync"
)

type CacheId struct {
	CurrentId     int64
	MaxId         int64
	ReservedRange int32
	syncRoot      sync.Mutex
	metaquery     *metaQuery
}

const (
	sqlPosgreSqlReturning string = "RETURNING"
)

var (
	cacheIdSchema *Schema
	cacheIdTable  *Table
	cacheIdQuery  string
)

func InitCacheId(schema *Schema, table *Table, resultTable *Table) {
	cacheIdSchema = schema
	cacheIdTable = resultTable
	cacheid := new(CacheId)
	cacheIdQuery = cacheid.GetDml(dmlstatement.Update, table)
}

func (cacheId *CacheId) Init(objid int32, schemaId int32, entityType entitytype.EntityType) {
	cacheId.metaquery = new(metaQuery)
	cacheId.metaquery.Init(cacheIdSchema, cacheIdTable)

	// added cacheIdSchema check to avoid unitesting crash!
	if cacheIdSchema != nil {
		cacheId.metaquery.query = cacheIdQuery
		// SET value=value+$1 WHERE id=$2 AND schema_id=$3 AND object_type=$4 RETURNING $5
		cacheId.metaquery.addParam(int64(1))
		cacheId.metaquery.addParam(objid)
		cacheId.metaquery.addParam(schemaId)
		cacheId.metaquery.addParam(int8(entityType))
	}
}

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (cacheId *CacheId) GetDml(dmlType dmlstatement.DmlStatement, table *Table) string {
	var result strings.Builder

	if dmlType == dmlstatement.Update {
		//TODO manage query for Mysql
		var field = table.GetFieldByName(metaValue)
		result.Grow(int(table.sqlCapacity))
		result.WriteString(dmlType.String())
		result.WriteString(dmlSpace)
		result.WriteString(table.physicalName)
		result.WriteString(dmlUpdateSet)
		result.WriteString(field.GetPhysicalName(table.provider))
		result.WriteString(operatorEqual)
		result.WriteString(field.GetPhysicalName(table.provider))
		result.WriteString(operatorPlus)
		result.WriteString(table.getVariableName(0))
		result.WriteString(dqlWhere)
		table.addPrimaryKeyFilter(&result, 1)
		// returning for postgresql ==> RETURNING value
		result.WriteString(dmlSpace)
		result.WriteString(sqlPosgreSqlReturning)
		result.WriteString(dmlSpace)
		result.WriteString(field.GetPhysicalName(table.provider))
	}

	return result.String()
}

//******************************
// private methods
//******************************
func (cacheId *CacheId) toMetaId(objectType entitytype.EntityType, objectId int32, schemaId int32) *MetaId {
	metaId := new(MetaId)
	metaId.id = objectId
	metaId.schemaId = schemaId
	metaId.objectType = int8(objectType)
	if cacheId.CurrentId == 0 {
		metaId.value = 1
	} else {
		metaId.value = cacheId.CurrentId
	}
	return metaId
}

func (cacheId *CacheId) create(objectType entitytype.EntityType, objectId int32, schemaId int32) error {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaIdTableName)
	return query.insertMetaId(cacheId.toMetaId(objectType, objectId, schemaId))
}

func (cacheId *CacheId) exists(objectType entitytype.EntityType, objectId int32, schemaId int32) bool {
	query := new(metaQuery)

	query.setSchema(metaSchemaName)
	query.setTable(metaIdTableName)

	query.addFilter(metaId, operatorEqual, objectId)
	query.addFilter(metaSchemaId, operatorEqual, schemaId)
	query.addFilter(metaObjectType, operatorEqual, int8(objectType))

	result, _ := query.exists()
	return result
}

func (cacheId *CacheId) GetNewId() bool {
	cacheId.metaquery.run()
	cacheId.CurrentId = cacheId.metaquery.getInt64Value()
	return true
}
