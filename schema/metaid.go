package schema

import (
	"fmt"
)

type metaId struct {
	id         int32
	schemaId   int32
	objectType int8
	value      int64
}

//******************************
// getters and setters
//******************************
func (metaid *metaId) GetId() int32 {
	return metaid.id
}
func (metaid *metaId) GetSchemaId() int32 {
	return metaid.schemaId
}
func (metaid *metaId) GetObjectType() int8 {
	return metaid.objectType
}
func (metaid *metaId) GetValue() int64 {
	return metaid.value
}

//******************************
// public methods
//******************************
func (metaid *metaId) String() string {
	return fmt.Sprintf("id: %d; schema_id: %d; object_type: %d; value: %d",
		metaid.id, metaid.schemaId, metaid.objectType, metaid.value)
}

//******************************
// private methods
//******************************
