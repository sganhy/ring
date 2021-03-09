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

//return operator string, has bind variable yes/no
func (operation OperationType) ToSql(provider databaseprovider.DatabaseProvider, value string) (string, bool) {
	switch operation {
	case Equal:
		if value == "" {
			return isNull, false
		}
		return strEqual, true
	case NotEqual:
		if value == "" {
			return isNotNull, false
		}
		return strNotEqual, true
	case Greater:
		return strGreater, true
	case GreaterOrEqual:
		return strGreaterOrEqual, true
	case Less:
		return strLess, true
	case LessOrEqual:
		return strLessOrEqual, true
	case Like:
		if value == "" {
			return isNull, false
		}
		return strLike, true
	case NotLike:
		if value == "" {
			return isNull, false
		}
		return strNotLike, true
	case In:
		return strIn, true
	}
	return "", false

}
