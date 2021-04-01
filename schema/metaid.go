package schema

import (
	"fmt"
)

type MetaId struct {
	id         int32
	schemaId   int32
	objectType int8
	value      int64
}

//******************************
// getters
//******************************
func (metaid *MetaId) GetId() int32 {
	return metaid.id
}
func (metaid *MetaId) GetSchemaId() int32 {
	return metaid.schemaId
}
func (metaid *MetaId) GetObjectType() int8 {
	return metaid.objectType
}
func (metaid *MetaId) GetValue() int64 {
	return metaid.value
}

//******************************
// public methods
//******************************
func (metaid *MetaId) String() string {
	return fmt.Sprintf(" id: %d; schema_id: %d; object_type: %d; value: %d",
		metaid.id, metaid.schemaId, metaid.objectType, metaid.value)
}

//******************************
// private methods
//******************************
