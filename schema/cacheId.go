package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"strings"
	"sync"
)

type cacheId struct {
	currentId     int64
	maxId         int64
	reservedRange uint32
	syncRoot      sync.Mutex
	query         *metaQuery
}

const (
	sqlPostgresReturning string = "RETURNING"
	maxReservedRange     uint32 = 1073741824           // 2^30
	initialMaxId         int64  = -9223372036854775808 // -1*((2^64)+1) - min int64 value
)

func (cacheid *cacheId) Init(metaSchema *Schema, targetSchemaId int32, targetEntity entity) {

	metaIdTable := metaSchema.GetTableByName(metaIdTableName)
	metaLongTable := metaSchema.GetTableByName(metaLongTableName)

	// max should be equal to int64 min value - forcing to initialize
	cacheid.maxId = initialMaxId

	// initialize query
	cacheid.query = new(metaQuery)
	cacheid.query.Init(metaSchema, metaLongTable)

	// added cacheIdSchema check to avoid unitesting crash!
	cacheid.query.query = cacheid.getDml(dmlstatement.Update, metaIdTable)
	cacheid.query.addParam(int32(1)) // default reserve range
	cacheid.query.addParam(targetEntity.GetId())
	cacheid.query.addParam(targetSchemaId)
	cacheid.query.addParam(int8(targetEntity.GetEntityType()))
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
func (cacheid *cacheId) setCurrentId(value int64) {
	cacheid.currentId = value
}
func (cacheid *cacheId) setReservedRange(value uint32) {
	cacheid.reservedRange = value
}

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
// used for Tables ==>
// Generate id with cache management for BulKsave (Hi/Lo algorithm)
// return last value generated
// id range can be lost multi slot of id should be managed
func (cacheid *cacheId) getNewRangeId(idRange uint32) int64 {
	// duplicate ==> GetNewId()
	/* ==> no check
	if reservedRange < 1 {
		// invalid range value
		// add logging here
		return 0
	}
	*/
	var result int64
	cacheid.syncRoot.Lock()

	//TODO manage cycles and overflows !!
	result = cacheid.currentId + int64(idRange)

	if result > cacheid.maxId {
		// take max reserve range
		if cacheid.reservedRange > idRange {
			cacheid.query.setParamValue(cacheid.reservedRange, 0)
			result = int64(idRange) - int64(cacheid.reservedRange)
		} else {
			// here we loosing id due to multiple interval ==>
			// 	]currentId,maxId]U[newId-reservedRange,newId]
			cacheid.query.setParamValue(int32(idRange), 0)
			result = 0
		}

		// never loaded
		cacheid.query.run(1)
		cacheid.maxId = cacheid.query.getInt64Value()
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

// used for Sequences ==>
// Generate id with cache management for BulKsave (Hi/Lo algorithm)
// !! here we cannot loose id on the same session
func (cacheid *cacheId) getNewId() int64 {
	var result int64

	cacheid.syncRoot.Lock()

	//TODO manage cycles and overflows !!
	result = cacheid.currentId + 1

	if result > cacheid.maxId {

		// compute reserve range
		if cacheid.reservedRange > 1 {
			cacheid.query.setParamValue(cacheid.reservedRange, 0)
			result = 1 - int64(cacheid.reservedRange)
		} else {
			// set default parameter for returning value
			cacheid.query.setParamValue(1, 0)
			result = 0
		}

		// never loaded
		cacheid.query.run(1)
		cacheid.maxId = cacheid.query.getInt64Value()
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

func (cacheid *cacheId) getDml(dmlType dmlstatement.DmlStatement, table *Table) string {
	var result strings.Builder
	var provider = table.GetDatabaseProvider()

	if dmlType == dmlstatement.Update && provider == databaseprovider.PostgreSql {
		//UPDATE rpg_sheet_test."@meta_id" SET "value"="value"+$1
		// WHERE id=$2 AND schema_id=$3 AND object_type=$4 RETURNING "value"
		var field = table.GetFieldByName(metaValue)
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
		result.WriteString(sqlPostgresReturning)
		result.WriteString(dmlSpace)
		result.WriteString(field.GetPhysicalName(provider))
	}

	return result.String()
}
