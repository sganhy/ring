package data

import (
	"fmt"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"strconv"
	"time"
)

type bulkRetrieveQueryItem struct {
	field               *schema.Field
	operation           operationtype.OperationType
	operand             string
	operands            *[]string
	caseSensitiveSearch bool
}

//******************************
// private methods
//******************************

func newQueryFilter(field *schema.Field, operation operationtype.OperationType, operand interface{}) (*bulkRetrieveQueryItem, error) {
	var filter = new(bulkRetrieveQueryItem)

	filter.field = field
	filter.operation = operation
	switch operand.(type) {
	case string:
		filter.operand = operand.(string)
		break
	case float32:
		// conversion issues
		filter.operand = fmt.Sprintf("%g", operand.(float32))
		break
	case float64:
		filter.operand = strconv.FormatFloat(operand.(float64), 'f', -1, 64)
		break
	case int:
		filter.operand = strconv.Itoa(operand.(int))
		break
	case uint:
		filter.operand = strconv.FormatUint(uint64(operand.(int)), 10)
		break
	case bool:
		filter.operand = strconv.FormatBool(operand.(bool))
		break
	case int8:
		filter.operand = strconv.Itoa(int(operand.(int8)))
		break
	case int16:
		filter.operand = strconv.Itoa(int(operand.(int16)))
		break
	case int32:
		filter.operand = strconv.Itoa(int(operand.(int32)))
		break
	case int64:
		filter.operand = strconv.FormatInt(operand.(int64), 10)
		break
	case uint8:
		filter.operand = strconv.FormatUint(uint64(operand.(uint8)), 10)
		break
	case uint16:
		filter.operand = strconv.FormatUint(uint64(operand.(uint16)), 10)
		break
	case uint32:
		filter.operand = strconv.FormatUint(uint64(operand.(uint32)), 10)
		break
	case uint64:
		filter.operand = strconv.FormatUint(operand.(uint64), 10)
		break
	case time.Time:
		filter.operand = field.GetDateTimeString(operand.(time.Time))
		break
	}
	return filter, nil
}

func newQuerySort(field *schema.Field, sortType sortordertype.SortOrderType) *bulkRetrieveQueryItem {
	var sort = new(bulkRetrieveQueryItem)
	sort.field = field
	sort.operation = operationtype.OperationType(sortType)
	return sort
}

func (bulkItem *bulkRetrieveQueryItem) Clone() *bulkRetrieveQueryItem {
	result := new(bulkRetrieveQueryItem)
	result.field = bulkItem.field
	result.operation = bulkItem.operation
	result.operand = bulkItem.operand
	result.operands = bulkItem.operands
	result.caseSensitiveSearch = bulkItem.caseSensitiveSearch
	return result
}
