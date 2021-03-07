package operationtype

import (
	"ring/schema/databaseprovider"
)

type OperationType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

var strEqual = "="
var strNotEqual = "<>"
var isNotNull = " IS NOT NULL "
var isNull = " IS NULL "
var strGreater = ">"
var strGreaterOrEqual = ">="
var strLess = "<"
var strLessOrEqual = "<="
var strLike = " LIKE "
var strNotLike = " NOT LIKE "
var strIn = " IN "

const (
	Equal          OperationType = 1
	NotEqual       OperationType = 2
	Greater        OperationType = 3
	GreaterOrEqual OperationType = 10
	Less           OperationType = 11
	LessOrEqual    OperationType = 12
	Like           OperationType = 13
	NotLike        OperationType = 14
	//SoundsLike = 15, -- one day ...
	In OperationType = 17

	// WARNING !!!sort operation
	// sync with sortordertype.Descending & Ascending value
	//TODO unitesting !!
	ascendingSort  OperationType = 101
	descendingSort OperationType = 102
)

func (operation OperationType) ToSql(provider databaseprovider.DatabaseProvider, value string) string {
	switch operation {
	case Equal:
		if value == "" {
			return isNull
		}
		return strEqual
	case NotEqual:
		if value == "" {
			return isNotNull
		}
		return strNotEqual
	case Greater:
		return strGreater
	case GreaterOrEqual:
		return strGreaterOrEqual
	case Less:
		return strLess
	case LessOrEqual:
		return strLessOrEqual
	case Like:
		return strLike
	case NotLike:
		return strNotLike
	case In:
		return strIn
	}
	return ""

}
