package schema

import (
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"strings"
	"sync"
)

type cacheId struct {
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

func initCacheId(schema *Schema, table *Table, resultTable *Table) {
	cacheIdSchema = schema
	cacheIdTable = resultTable
	cacheid := new(cacheId)
	cacheIdQuery = cacheid.GetDml(dmlstatement.Update, table)
}

func (cacheid *cacheId) Init(objid int32, schemaId int32, entityType entitytype.EntityType) {
	cacheid.metaquery = new(metaQuery)
	cacheid.metaquery.Init(cacheIdSchema, cacheIdTable)

	// added cacheIdSchema check to avoid unitesting crash!
	if cacheIdSchema != nil {
		cacheid.metaquery.query = cacheIdQuery
		cacheid.metaquery.addParam(int32(1))
		cacheid.metaquery.addParam(objid)
		cacheid.metaquery.addParam(schemaId)
		cacheid.metaquery.addParam(int8(entityType))
	}

	// max should be equal to int64 min value - forcing to initialize
	cacheid.maxId = initialMaxId
	cacheid.currentId = 0
	cacheid.reservedRange = 0
}

//******************************
// getters and setters
//******************************
func (cacheid *cacheId) GetCurrentId() int64 {
	return cacheid.currentId
}
func (cacheid *cacheId) IsInitialized() bool {
	return cacheid.maxId != initialMaxId
}
func (cacheid *cacheId) SetCurrentId(value int64) {
	cacheid.currentId = value
}

//dynamic range (Cache I)
func (cacheid *cacheId) SetCache(value bool) {
	if value == true {
		cacheid.reservedRange = 1
	} else {
		cacheid.reservedRange = 0
	}
}

//******************************
// public methods
//******************************
func (cacheid *cacheId) GetDml(dmlType dmlstatement.DmlStatement, table *Table) string {
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
func (cacheid *cacheId) GetNewId() int64 {
	var result int64
	cacheid.syncRoot.Lock()

	//TODO manage cycles and overflows !!
	result = cacheid.currentId + 1
	if result > cacheid.maxId {
		// compute reserve range
		if cacheid.reservedRange > 1 {
			cacheid.metaquery.setParamValue(cacheid.reservedRange, 0)
			result = 1 - int64(cacheid.reservedRange)
		} else {
			// set default parameter for returning value
			cacheid.metaquery.setParamValue(int32(1), 0)
			result = 0
		}

		// never loaded
		cacheid.metaquery.run(1)
		cacheid.maxId = cacheid.metaquery.getInt64Value()
		result += cacheid.maxId

		// multiply by 2 next reserved range if cacheId.reservedRange is less than 2^30
		if cacheid.reservedRange < maxReservedRange {
			cacheid.reservedRange <<= 1
		}
	}
	cacheid.currentId = result
	cacheid.syncRoot.Unlock()
	return result
}

// Generate id with cache management for BulKsave (Hi/Lo algorithm)
func (cacheid *cacheId) GetNewRangeId(reservedRange uint32) int64 {
	// duplicate ==> GetNewId()
	if reservedRange > maxExtReservedRange || reservedRange < 1 {
		// invalid range value
		// add logging here
		return 0
	}

	var resultRange int64
	cacheid.syncRoot.Lock()

	//TODO manage cycles and overflows !!
	resultRange = cacheid.currentId + int64(reservedRange)
	if resultRange > cacheid.maxId {
		// compute reserve range
		if cacheid.reservedRange > int32(reservedRange) {
			cacheid.metaquery.setParamValue(cacheid.reservedRange, 0)
			resultRange = 1 - int64(cacheid.reservedRange)
		} else {
			cacheid.metaquery.setParamValue(int32(reservedRange), 0)
			resultRange = 1 - int64(reservedRange)
		}

		// never loaded
		cacheid.metaquery.run(1)
		cacheid.maxId = cacheid.metaquery.getInt64Value()
		resultRange += cacheid.maxId

		// multiply by 2 next reserved range if cacheId.reservedRange is less than 2^30
		if cacheid.reservedRange < maxReservedRange {
			cacheid.reservedRange <<= 1
		}
	}
	cacheid.currentId = resultRange
	cacheid.syncRoot.Unlock()
	return resultRange
}

//******************************
// private methods
//******************************
func (cacheid *cacheId) toMetaId(objectType entitytype.EntityType, objectId int32, schemaId int32) *metaId {
	metaid := new(metaId)
	metaid.id = objectId
	metaid.schemaId = schemaId
	metaid.objectType = int8(objectType)
	metaid.value = cacheid.currentId
	return metaid
}

func (cacheid *cacheId) create(objectType entitytype.EntityType, objectId int32, schemaId int32) error {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaIdTableName)
	return query.insertMetaId(cacheid.toMetaId(objectType, objectId, schemaId))
}

func (cacheid *cacheId) exists(objectType entitytype.EntityType, objectId int32, schemaId int32) bool {
	query := new(metaQuery)

	query.setSchema(metaSchemaName)
	query.setTable(metaIdTableName)

	query.addFilter(metaFieldId, operatorEqual, objectId)
	query.addFilter(metaSchemaId, operatorEqual, schemaId)
	query.addFilter(metaObjectType, operatorEqual, int8(objectType))

	result, _ := query.exists()
	return result
}
