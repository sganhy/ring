package data

import (
	"ring/data/sortordertype"
	"ring/schema"
)

type bulkRetrieveSort struct {
	field     *schema.Field
	orderType sortordertype.SortOrderType
}

//******************************
// private methods
//******************************

func newQuerySort(field *schema.Field, sortType sortordertype.SortOrderType) *bulkRetrieveSort {
	var sort = new(bulkRetrieveSort)
	sort.field = field
	sort.orderType = sortType
	return sort
}
