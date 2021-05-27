package schema

import (
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"strings"
	"sync"
)

type CacheId struct {
	currentId     int64
	maxId         int64
	reservedRange int32
	syncRoot      sync.Mutex
	metaquery     *metaQuery
}

const (
	sqlPosgreSqlReturning string = "RETURNING"
	maxReservedRange      int32  = 1073741824           // 2^30
	maxExtReservedRange   uint32 = 2147483647           // 2^32 - max external defined Reserve value
	initialMaxId          int64  = -9223372036854775808 // -1*((2^64)+1) - min int64 value
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
		cacheId.metaquery.addParam(int32(1))
		cacheId.metaquery.addParam(objid)
		cacheId.metaquery.addParam(schemaId)
		cacheId.metaquery.addParam(int8(entityType))
	}

	// max should be equal to int64 min value - forcing to initialize
	cacheId.maxId = initialMaxId
	cacheId.currentId = 0
	cacheId.reservedRange = 0
}

//******************************
// getters and setters
//******************************
func (cacheId *CacheId) GetCurrentId() int64 {
	return cacheId.currentId
}
func (cacheId *CacheId) IsInitialized() bool {
	return cacheId.maxId != initialMaxId
}
func (cacheId *CacheId) SetCurrentId(value int64) {
	cacheId.currentId = value
}

//dynamic range (Cache I)
func (cacheId *CacheId) SetCache(value bool) {
	if value == true {
		cacheId.reservedRange = 1
	} else {
		cacheId.reservedRange = 0
	}
}

//******************************
// public methods
//******************************
func (cacheId *CacheId) GetDml(dmlType dmlstatement.DmlStatement, table *Table) string {
	var result strings.Builder
	var provider = table.GetDatabaseProvider()

	if dmlType == dmlstatement.Update {
		//TODO manage query for Mysql
		var field = table.GetFieldByName(metaValue)
		result.Grow(int(table.getSqlCapacity()))
		result.WriteString(dmlType.String())
		result.WriteString(dmlSpace)
		result.WriteString(table.GetPhysicalName())
		result.WriteString(dmlUpdateSet)
		result.WriteString(field.GetPhysicalName(provider))
		result.WriteString(operatorEqual)
		result.WriteString(field.GetPhysicalName(provider))
		result.WriteString(operatorPlus)
		result.WriteString(table.getVariableName(0))
		result.WriteString(dqlWhere)
		table.addPrimaryKeyFilter(&result, 1)
		// returning for postgresql ==> RETURNING value
		result.WriteString(dmlSpace)
		result.WriteString(sqlPosgreSqlReturning)
		result.WriteString(dmlSpace)
		result.WriteString(field.GetPhysicalName(provider))
	}

	return result.String()
}

// Generate id with cache management for BulKsave (Hi/Lo algorithm)
func (cacheId *CacheId) GetNewId() int64 {
	var result int64
	cacheId.syncRoot.Lock()

	//TODO manage cycles and overflows !!
	result = cacheId.currentId + 1
	if result > cacheId.maxId {
		// compute reserve range
		if cacheId.reservedRange > 1 {
			cacheId.metaquery.setParamValue(cacheId.reservedRange, 0)
			result = 1 - int64(cacheId.reservedRange)
		} else {
			// set default parameter for returning value
			cacheId.metaquery.setParamValue(int32(1), 0)
			result = 0
		}

		// never loaded
		cacheId.metaquery.run(1)
		cacheId.maxId = cacheId.metaquery.getInt64Value()
		result += cacheId.maxId

		// multiply by 2 next reserved range if cacheId.reservedRange is less than 2^30
		if cacheId.reservedRange < maxReservedRange {
			cacheId.reservedRange <<= 1
		}
	}
	cacheId.currentId = result
	cacheId.syncRoot.Unlock()
	return result
}

// Generate id with cache management for BulKsave (Hi/Lo algorithm)
func (cacheId *CacheId) GetNewRangeId(reservedRange uint32) int64 {
	// duplicate ==> GetNewId()
	if reservedRange > maxExtReservedRange || reservedRange < 1 {
		// invalid range value
		// add logging here
		return 0
	}

	var resultRange int64
	cacheId.syncRoot.Lock()

	//TODO manage cycles and overflows !!
	resultRange = cacheId.currentId + int64(reservedRange)
	if resultRange > cacheId.maxId {
		// compute reserve range
		if cacheId.reservedRange > int32(reservedRange) {
			cacheId.metaquery.setParamValue(cacheId.reservedRange, 0)
			resultRange = 1 - int64(cacheId.reservedRange)
		} else {
			cacheId.metaquery.setParamValue(int32(reservedRange), 0)
			resultRange = 1 - int64(reservedRange)
		}

		// never loaded
		cacheId.metaquery.run(1)
		cacheId.maxId = cacheId.metaquery.getInt64Value()
		resultRange += cacheId.maxId

		// multiply by 2 next reserved range if cacheId.reservedRange is less than 2^30
		if cacheId.reservedRange < maxReservedRange {
			cacheId.reservedRange <<= 1
		}
	}
	cacheId.currentId = resultRange
	cacheId.syncRoot.Unlock()
	return resultRange
}

//******************************
// private methods
//******************************
func (cacheId *CacheId) toMetaId(objectType entitytype.EntityType, objectId int32, schemaId int32) *MetaId {
	metaId := new(MetaId)
	metaId.id = objectId
	metaId.schemaId = schemaId
	metaId.objectType = int8(objectType)
	metaId.value = cacheId.currentId
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
