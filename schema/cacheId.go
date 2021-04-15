package schema

import (
	"ring/schema/entitytype"
	"sync"
)

type CacheId struct {
	CurrentId     int64
	MaxId         int64
	ReservedRange int32
	syncRoot      sync.Mutex
}

//******************************
// getters
//******************************

//******************************
// public methods
//******************************
func (cacheId *CacheId) ToMetaId(objectType entitytype.EntityType, objectId int32, schemaId int32) *MetaId {
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

//******************************
// private methods
//******************************
func (cacheId *CacheId) create(objectType entitytype.EntityType, objectId int32, schemaId int32) error {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaIdTableName)
	return query.insertMetaId(cacheId.ToMetaId(objectType, objectId, schemaId))
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
