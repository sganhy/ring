package relationtype

type RelationType int8

const (
	Otop       RelationType = 1
	Otm        RelationType = 2
	Mtm        RelationType = 3
	Mto        RelationType = 4
	Otof       RelationType = 5
	NotDefined RelationType = 127
)
