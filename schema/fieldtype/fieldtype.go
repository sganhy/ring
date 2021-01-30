package fieldtype

type FieldType int8

const (
	Long          FieldType = 0
	Int           FieldType = 1
	Short         FieldType = 2
	Byte          FieldType = 3
	Float         FieldType = 4
	Double        FieldType = 5
	String        FieldType = 6
	ShortDateTime FieldType = 7
	DateTime      FieldType = 8
	LongDateTime  FieldType = 9
	Array         FieldType = 11
	Boolean       FieldType = 13
	NotDefined    FieldType = 127
)
