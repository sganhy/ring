package schema

import (
	"ring/schema/ddlstatement"
	"ring/schema/entitytype"
)

type entity interface {
	GetId() int32
	GetName() string
	GetPhysicalName() string
	GetEntityType() entitytype.EntityType
	logStatment(statment ddlstatement.DdlStatement) bool
}
