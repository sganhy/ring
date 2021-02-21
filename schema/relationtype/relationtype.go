package relationtype

type RelationType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Otop       RelationType = 1
	Otm        RelationType = 2
	Mtm        RelationType = 3
	Mto        RelationType = 11
	Otof       RelationType = 12
	NotDefined RelationType = 123
)
