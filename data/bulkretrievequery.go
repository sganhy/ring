package data

import (
	"database/sql"
	"fmt"
	"ring/data/bulkretrievetype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"strconv"
	"strings"
)

const defaultPostGreParameterName = "$"
const filterSeparator = " AND "

//type queryFunc func(dbConnection *sql.DB, sql string, params []interface{}) (*sql.Rows, error)

// all this structure is readonly
type bulkRetrieveQuery struct {
	targetObject *schema.Table
	queryType    bulkretrievetype.BulkRetrieveType
	filterCount  *int
	items        *[]*bulkRetrieveQueryItem
	result       *List // pointer is mandatory
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************

func (query bulkRetrieveQuery) Execute(dbConnection *sql.DB, transaction *sql.Tx) error {
	var provider = query.targetObject.GetDatabaseProvider()
	var whereClause, parameters = query.getWhereClause(provider)
	var orderClause = query.getOrderClause(provider)
	var sqlQuery = query.targetObject.GetDql(whereClause, orderClause)

	rows, err := dbConnection.Query(sqlQuery, parameters...)

	if rows != nil {
		rows.Close()
	}
	if err != nil {

		return err
	}
	var rowIndex = 0
	count := query.targetObject.GetFieldCount()
	query.result.data = make([]*Record, 0, 10)

	columns := make([]interface{}, count)
	columnPointers := make([]interface{}, count)
	for rows.Next() {
		var record = new(Record)
		record.recordType = query.targetObject
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			rows.Close()
			return err
		}
		record.data = query.targetObject.GetQueryResult(columnPointers)
		query.result.appendItem(record)
		rowIndex++
	}
	rows.Close()
	return nil
}

//******************************
// private methods
//******************************
func (query *bulkRetrieveQuery) getWhereClause(provider databaseprovider.DatabaseProvider) (string, []interface{}) {
	var result strings.Builder
	var operator string
	var hasVariable bool
	var item *bulkRetrieveQueryItem
	var variableId = 0
	var parameterId = 0
	var parameters []interface{}

	result.Grow((*query.filterCount) * 30)
	//TODO may be two pass to reduce allocations
	for i := 0; i < len(*query.items); i++ {
		item = (*query.items)[i]
		operator, hasVariable = item.operation.ToSql(provider, item.operand)
		if operator != "" {
			result.WriteString(item.field.GetPhysicalName(provider))
			result.WriteString(operator)

			if hasVariable == true {
				// get parameter
				query.getParameterName(provider, &result, variableId)
				if item.operand != "" {
					parameters = append(parameters, item.field.GetParameterValue(item.operand))
				}
				variableId++
			}
			if parameterId+1 < *query.filterCount {
				result.WriteString(filterSeparator)
			}
			parameterId++
		}
	}
	//fmt.Printf("Parameters(%d)\n", len(parameters))
	return result.String(), parameters
}

func (query *bulkRetrieveQuery) getOrderClause(provider databaseprovider.DatabaseProvider) string {
	var result strings.Builder
	var capacity = len(*query.items) - *query.filterCount
	var descId = int8(sortordertype.Descending)
	var ascId = int8(sortordertype.Ascending)
	var item *bulkRetrieveQueryItem
	var parameterId = 0

	if capacity > 0 {
		result.Grow(capacity * 30)
		for i := 0; i < len(*query.items); i++ {
			item = (*query.items)[i]
			if int8(item.operation) == descId || int8(item.operation) == ascId {
				result.WriteString(item.field.GetPhysicalName(provider))
				if int8(item.operation) == descId {
					result.WriteString(" DESC")
				}
				operator, _ := item.operation.ToSql(provider, item.operand)
				result.WriteString(operator)
				if parameterId < capacity-1 {
					result.WriteString(",")
				}
				parameterId++
			}
		}
	}
	return result.String()
}

func (query *bulkRetrieveQuery) getParameterName(provider databaseprovider.DatabaseProvider, params *strings.Builder, index int) {
	switch provider {
	case databaseprovider.PostgreSql:
		params.WriteString(defaultPostGreParameterName)
		params.WriteString(strconv.Itoa(index + 1))
		break
	}
}

func (query *bulkRetrieveQuery) clearItems() {
	var items = make([]*bulkRetrieveQueryItem, 0, 2)
	query.items = &items
	*query.filterCount = 0
}

func newSimpleQuery(table *schema.Table) schema.Query {
	var query = new(bulkRetrieveQuery)
	var items = make([]*bulkRetrieveQueryItem, 0, 2)

	query.targetObject = table
	query.queryType = bulkretrievetype.SimpleQuery
	query.filterCount = new(int)
	*query.filterCount = 0
	query.items = &items
	query.result = new(List)

	return *query
}

func (query *bulkRetrieveQuery) addFilter(item *bulkRetrieveQueryItem) {
	*query.filterCount = *query.filterCount + 1
	*query.items = append(*query.items, item)
}

func (query *bulkRetrieveQuery) addSort(item *bulkRetrieveQueryItem) {
	*query.items = append(*query.items, item)
}
