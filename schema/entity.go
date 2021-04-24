package schema

import (
	"ring/schema/entitytype"
)

type entity interface {
	GetId() int32
	GetName() string
	GetPhysicalName() string
	GetEntityType() entitytype.EntityType
}
