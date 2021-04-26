package physicaltype

type PhysicalType int8

const (
	Table       PhysicalType = 1
	View        PhysicalType = 3
	Measurement PhysicalType = 5
	Logical     PhysicalType = 11
)
