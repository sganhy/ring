package physicaltype

type PhysicalType int8

const (
	Table       PhysicalType = 1
	View        PhysicalType = 3
	Measurement PhysicalType = 5
	Logical     PhysicalType = 11
)

func (physicalType *PhysicalType) String() string {
	physType := *physicalType
	switch physType {
	case Table:
		return "table"
	case View:
		return "view"
	case Measurement:
		return "measurement"
	case Logical:
		return "logical"
	}

	return ""
}
