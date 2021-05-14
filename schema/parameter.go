package schema

import (
	"fmt"
	"ring/schema/entitytype"
)

type parameter struct {
	id            int32
	name          string
	description   string
	parameterType entitytype.EntityType
	value         string
}

const (
	parameterToStringFormat string = "name=%s; description=%s"
	parameterVersion        string = "@version"
)

func (param *parameter) Init(id int32, name string, description string, value string) {
	param.id = id
	param.name = name
	param.description = description
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

func (param *parameter) String() string {
	// "name=%s; description=%s"
	return fmt.Sprintf(parameterToStringFormat, param.name, param.description)
}

//******************************
// public methods
//******************************
