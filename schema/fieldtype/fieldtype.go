package fieldtype

type FieldType int8

//!!! reserved value for unitesting {4, 5, 6} !!!

const (
	Long          FieldType = 0
	Int           FieldType = 1
	Short         FieldType = 2
	Byte          FieldType = 3
	Float         FieldType = 14
	Double        FieldType = 15
	String        FieldType = 16
	ShortDateTime FieldType = 17
	DateTime      FieldType = 18
	LongDateTime  FieldType = 19
	Array         FieldType = 21
	Boolean       FieldType = 23
	NotDefined    FieldType = 125
)
