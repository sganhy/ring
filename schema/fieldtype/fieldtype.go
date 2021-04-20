package fieldtype

type FieldType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	strLong          string = "Long"
	strInt           string = "Int"
	strShort         string = "Short"
	strByte          string = "Byte"
	strFloat         string = "Float"
	strDouble        string = "Double"
	strString        string = "String"
	strShortDateTime string = "Short DateTime"
	strDateTime      string = "DateTime"
	strLongDateTime  string = "Long DateTime"
	strArray         string = "Array"
	strBoolean       string = "Boolean"
	strNotDefined    string = "Not Defined"
)

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
	LongString    FieldType = 27
	NotDefined    FieldType = 125
)

func (fieldType FieldType) String() string {
	switch fieldType {
	case Long:
		return strLong
	case Int:
		return strInt
	case Short:
		return strShort
	case Byte:
		return strByte
	case Float:
		return strFloat
	case Double:
		return strDouble
	case String, LongString:
		return strString
	case ShortDateTime:
		return strShortDateTime
	case DateTime:
		return strDateTime
	case LongDateTime:
		return strLongDateTime
	case Array:
		return strArray
	case Boolean:
		return strBoolean
	case NotDefined:
		return strNotDefined
	}
	return ""
}
