package ddlstatement

type DdlStatement int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Create    DdlStatement = 1
	Drop      DdlStatement = 2
	Alter     DdlStatement = 3
	Truncate  DdlStatement = 9
	Undefined DdlStatement = 127
)

func (statement DdlStatement) String() string {
	switch statement {
	case Create:
		return "CREATE"
	case Drop:
		return "DROP"
	case Alter:
		return "ALTER"
	case Truncate:
		return "TRUNCATE"
	}
	return ""
}
