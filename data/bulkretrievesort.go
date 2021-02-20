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
