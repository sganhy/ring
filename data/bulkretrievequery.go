package data

import (
	"ring/data/bulkquerytype"
	"ring/schema"
)

type bulkRetrieveQuery struct {
	targetObject *schema.Table
	queryType    bulkquerytype.BulkQueryType
	filters      *[]*bulkRetrieveFilter
	sorts        *[]*bulkRetrieveSort
}

//******************************
// private methods
//******************************

func newSimpleQuery(table *schema.Table) *bulkRetrieveQuery {
	var query = new(bulkRetrieveQuery)
	query.targetObject = table
	query.queryType = bulkquerytype.SimpleQuery
	query.filters = nil
	query.sorts = nil
	return query
}

func (query *bulkRetrieveQuery) addFilter(filter *bulkRetrieveFilter) {
	if query.filters == nil {
		filters := make([]*bulkRetrieveFilter, 0, 2)
		query.filters = &filters
	}
	*query.filters = append(*query.filters, filter)
}

func (query *bulkRetrieveQuery) addSort(sort *bulkRetrieveSort) {
	if query.sorts == nil {
		sorts := make([]*bulkRetrieveSort, 0, 1)
		query.sorts = &sorts
	}
	*query.sorts = append(*query.sorts, sort)
}
