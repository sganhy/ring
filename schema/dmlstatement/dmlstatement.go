package dmlstatement

type DmlStatement int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Insert          DmlStatement = 1
	Update          DmlStatement = 2
	Delete          DmlStatement = 3
	UpdateReturning DmlStatement = 9
)

func (statement DmlStatement) String() string {
	switch statement {
	case Insert:
		return "INSERT INTO"
	case Update:
		return "UPDATE"
	case Delete:
		return "DELETE FROM"
	}
	return ""
}
