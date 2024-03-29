package schema

import (
	"fmt"
	"ring/schema/entitytype"
)

type metaId struct {
	id         int32
	schemaId   int32
	objectType int8
	value      int64
}

const (
	metaIdToStringFormat string = "id: %d; schema_id: %d; object_type: %d; value: %d"
)

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
	return fmt.Sprintf(metaIdToStringFormat, metaid.id, metaid.schemaId, metaid.objectType, metaid.value)
}

//******************************
// private methods
//******************************
func (metaid *metaId) saveMetaIdList(schemaId int32, metaList []*meta) (int64, error) {
	queryExist := new(metaQuery)
	queryInsert := new(metaQuery)
	var count int64 = 0

	metaid.schemaId = schemaId
	metaid.objectType = int8(entitytype.Table)
	metaid.value = 0

	queryExist.setSchema(metaSchemaName)
	queryExist.setTable(metaIdTableName)
	queryInsert.setSchema(metaSchemaName)
	queryInsert.setTable(metaIdTableName)

	queryExist.addFilter(metaFieldId, operatorEqual, 0)
	queryExist.addFilter(metaSchemaId, operatorEqual, schemaId)
	queryExist.addFilter(metaObjectType, operatorEqual, int8(entitytype.Table))

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			queryExist.setParamValue(metaData.id, 0)
			exist, err := queryExist.exists()
			if err != nil {
				return count, err
			}
			if exist == false {
				metaid.id = metaData.id
				queryInsert.insertMetaId(metaid)
				count++
			}
		}
	}

	return count, nil
}
