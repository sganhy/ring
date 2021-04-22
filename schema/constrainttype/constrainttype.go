package constrainttype

type ConstraintType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	PrimaryKey ConstraintType = 1
	UniqueKey  ConstraintType = 2
	Check      ConstraintType = 3
	NotDefined ConstraintType = 125
)

func (constraintType ConstraintType) String() string {
	switch constraintType {
	case PrimaryKey:
		return "PRIMARY KEY"
	case Check:
		return "CHECK"
	}
	return ""
}
