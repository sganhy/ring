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
	NotDefined EntityType = 127
)
