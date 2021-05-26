package schema

import (
	"fmt"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"strings"
	"time"
)

type parameter struct {
	id            int32
	name          string
	description   string
	dataType      fieldtype.FieldType
	parameterType entitytype.EntityType
	value         string
}

const (
	parameterToStringFormat string = "name=%s; description=%s"
	parameterVersion        string = "@version"
	parameterCreationTime   string = "@creation_time"
	parameterLastUpgrade    string = "@last_upgrade"
)

func (param *parameter) Init(id int32, name string, description string, parameterType entitytype.EntityType, fieldType fieldtype.FieldType, value string) {
	param.id = id
	param.name = name
	param.description = description
	param.parameterType = parameterType
	if fieldType == fieldtype.NotDefined {
		param.dataType = fieldtype.String
	} else {
		param.dataType = fieldType
	}

	param.value = value
}

//******************************
// getters
//******************************
func (param *parameter) GetId() int32 {
	return param.id
}

func (param *parameter) GetName() string {
	return param.name
}

func (param *parameter) GetDescription() string {
	return param.description
}

func (param *parameter) GetEntityType() entitytype.EntityType {
	return entitytype.Parameter
}

func (param *parameter) GetDataType() fieldtype.FieldType {
	return param.dataType
}
func (param *parameter) GetValue() string {
	return param.value
}

//******************************
// public methods
//******************************
func (param *parameter) Clone() *parameter {
	var newParam = new(parameter)
	//id int32, name string, description string, parameterType entitytype.EntityType, value string
	newParam.Init(param.id, param.name, param.description, param.parameterType, param.dataType, param.value)
	return newParam
}

func (param *parameter) String() string {
	// "name=%s; description=%s"
	return fmt.Sprintf(parameterToStringFormat, param.name, param.description)
}

//******************************
// private methods
//******************************
func (param *parameter) exists(schema *Schema) bool {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	query.addFilter(metaId, operatorEqual, param.id)
	query.addFilter(metaSchemaId, operatorEqual, schema.id)
	query.addFilter(metaObjectType, operatorEqual, int8(entitytype.Parameter))
	query.addFilter(metaReferenceId, operatorEqual, schema.id)

	result, _ := query.exists()
	return result
}

func (param *parameter) create(schema *Schema) error {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)
	return query.insertMeta(param.toMeta(schema.id), schema.id)
}

func (param *parameter) toMeta(objectId int32) *Meta {
	var metaParam = new(Meta)

	// key
	metaParam.id = param.id
	metaParam.name = param.name
	metaParam.description = param.description
	metaParam.refId = objectId
	metaParam.objectType = int8(entitytype.Parameter)
	metaParam.dataType = int32(param.parameterType)
	metaParam.setEntityBaseline(true)
	metaParam.setEntityBaseline(true)
	metaParam.setParameterDataType(param.dataType)
	metaParam.value = param.value
	metaParam.enabled = true

	return metaParam
}

func (param *parameter) getVersionParameter(refId int32, parameterType entitytype.EntityType, value string) *parameter {
	result := new(parameter)
	switch parameterType {
	case entitytype.Schema:
		if refId > 0 {
			result.Init(refId+41, parameterVersion, "Schema version", parameterType, fieldtype.String, value)
		} else {
			var ver = new(version)
			result.Init(3, parameterVersion, "Ring version", parameterType, fieldtype.String, ver.GetCurrentVersion())
		}
		break
	}
	return result
}

func (param *parameter) getCreationTimeParameter(refId int32, parameterType entitytype.EntityType) *parameter {
	result := new(parameter)
	value := time.Now().UTC().Format(time.RFC3339)
	message := strings.Title(strings.ToLower(parameterType.String())) + " creation time"
	result.Init(1, parameterCreationTime, message, parameterType, fieldtype.DateTime, value)
	return result
}

func (param *parameter) getLastUpgradeParameter(refId int32, parameterType entitytype.EntityType) *parameter {
	result := new(parameter)
	value := time.Now().UTC().Format(time.RFC3339)
	message := "Last " + strings.Title(strings.ToLower(parameterType.String())) + " upgrade"
	result.Init(2, parameterLastUpgrade, message, parameterType, fieldtype.DateTime, value)
	return result
}
