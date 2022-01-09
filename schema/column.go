package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
)

type Column interface {
	GetId() int32
	GetName() string
	GetColumnType() fieldtype.FieldType
	GetPhysicalName(provider databaseprovider.DatabaseProvider) string
	IsCaseSensitive() bool
}
