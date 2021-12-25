package constrainttype

type ConstraintType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	PrimaryKey ConstraintType = 1
	UniqueKey  ConstraintType = 2
	Check      ConstraintType = 3
	NotNull    ConstraintType = 8
	ForeignKey ConstraintType = 9
	Undefined  ConstraintType = 127
)

func (constraintType ConstraintType) String() string {
	switch constraintType {
	case PrimaryKey:
		return "PRIMARY KEY"
	case Check:
		return "CHECK"
	case ForeignKey:
		return "FOREIGN KEY"
	case NotNull:
		return "NOT NULL"
	}
	return ""
}
