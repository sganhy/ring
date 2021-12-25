package fieldtype

import (
	"strings"
)

type FieldType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	strLong           string = "Long"
	strInt            string = "Int"
	strShort          string = "Short"
	strByte           string = "Byte"
	strFloat          string = "Float"
	strDouble         string = "Double"
	strString         string = "String"
	strShortDateTime  string = "Short DateTime"
	strShortDateTime2 string = "ShortDateTime"
	strDateTime       string = "DateTime"
	strLongDateTime   string = "Long DateTime"
	strLongDateTime2  string = "LongDateTime"
	strArray          string = "Array"
	strBoolean        string = "Boolean"
	strUndefined      string = "Not Defined"
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
	Undefined     FieldType = 127
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
	case Undefined:
		return strUndefined
	}
	return ""
}

func GetFieldType(value string) FieldType {
	switch strings.ToLower(value) {
	case strings.ToLower(strLong):
		return Long
	case strings.ToLower(strInt):
		return Int
	case strings.ToLower(strShort):
		return Short
	case strings.ToLower(strByte):
		return Byte
	case strings.ToLower(strFloat):
		return Float
	case strings.ToLower(strDouble):
		return Double
	case strings.ToLower(strString):
		return String
	case strings.ToLower(strShortDateTime), strings.ToLower(strShortDateTime2):
		return ShortDateTime
	case strings.ToLower(strDateTime):
		return DateTime
	case strings.ToLower(strLongDateTime), strings.ToLower(strLongDateTime2):
		return LongDateTime
	case strings.ToLower(strBoolean):
		return Boolean
	}
	return Undefined
}

func GetFieldTypeById(entityId int) FieldType {
	if entityId <= 127 && entityId >= -128 {
		var newId = FieldType(entityId)
		if newId == Long || newId == Int || newId == Short || newId == Byte || newId == Float ||
			newId == Double || newId == String || newId == ShortDateTime || newId == DateTime || newId == LongDateTime ||
			newId == Array || newId == Boolean || newId == LongString {
			return newId
		}
	}
	return Undefined
}
