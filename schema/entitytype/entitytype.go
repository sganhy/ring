package entitytype

type EntityType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Table      EntityType = 0
	Field      EntityType = 1
	Relation   EntityType = 2
	Index      EntityType = 3
	Schema     EntityType = 7
	Sequence   EntityType = 15
	Language   EntityType = 17
	Tablespace EntityType = 18
	Parameter  EntityType = 23
	Constraint EntityType = 101 // not stored in @meta table
	NotDefined EntityType = 127
)

func (entityType EntityType) String() string {
	switch entityType {
	case Table:
		return "TABLE"
	case Field:
		return "FIELD"
	case Relation:
		return "RELATION"
	case Index:
		return "INDEX"
	case Schema:
		return "SCHEMA"
	case Sequence:
		return "SEQUENCE"
	case Tablespace:
		return "TABLESPACE"
	case Parameter:
		return "PARAMETER"
	case Language:
		return "LANGUAGE"
	case Constraint:
		return "CONSTRAINT"
	}
	return ""
}

func GetEntityTypeById(entityId int) EntityType {
	if entityId <= 127 && entityId >= -128 {
		var newId = EntityType(entityId)
		if newId == Table || newId == Field || newId == Relation || newId == Index || newId == Schema || newId == Sequence ||
			newId == Language || newId == Tablespace || newId == Parameter || newId == Constraint {
			return newId
		}
	}
	return NotDefined
}
