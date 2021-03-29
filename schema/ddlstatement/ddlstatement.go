package ddlstatement

type DdlStatement int8

const (
	Create     DdlStatement = 1
	Drop       DdlStatement = 2
	Alter      DdlStatement = 3
	NotDefined DdlStatement = 127
)

func (statement DdlStatement) String() string {
	switch statement {
	case Create:
		return "CREATE"
	case Drop:
		return "DROP"
	case Alter:
		return "ALTER"
	}
	return ""
}
