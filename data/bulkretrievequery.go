package data

import (
	"database/sql"
	"fmt"
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
// public methods (Interface schema.Query implementations)
//******************************

func (query bulkRetrieveQuery) Execute(dbConnection *sql.DB) error {
	fmt.Println("Execute query")
	return nil
}

//******************************
// private methods
//******************************
func newSimpleQuery(table *schema.Table) schema.Query {
	var query = new(bulkRetrieveQuery)
	var filters = make([]*bulkRetrieveFilter, 0, 2)
	var sorts = make([]*bulkRetrieveSort, 0, 1)

	query.targetObject = table
	query.queryType = bulkquerytype.SimpleQuery
	query.filters = &filters
	query.sorts = &sorts

	return *query
}

func (query *bulkRetrieveQuery) addFilter(filter *bulkRetrieveFilter) {
	*query.filters = append(*query.filters, filter)
}

func (query *bulkRetrieveQuery) addSort(sort *bulkRetrieveSort) {
	*query.sorts = append(*query.sorts, sort)
}
