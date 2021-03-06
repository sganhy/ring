package data

import (
	"database/sql"
	"fmt"
	"ring/data/bulkquerytype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"strconv"
	"strings"
)

const defaultParameterName = "B"
const defaultParameterPrefix = ":"

type bulkRetrieveQuery struct {
	targetObject *schema.Table
	queryType    bulkquerytype.BulkQueryType
	items        *[]*bulkRetrieveQueryItem
	result       *List // pointer is mandatory
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************

func (query bulkRetrieveQuery) Execute(provider databaseprovider.DatabaseProvider, dbConnection *sql.DB) error {
	var whereClause, parameters = query.getWhereClause(provider)
	var sql = query.targetObject.GetDql(provider, whereClause)

	rows, err := dbConnection.Query(sql, parameters)
	fmt.Println(sql)
	cols, _ := rows.Columns()

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
		record.setRecordType(query.targetObject)

		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			rows.Close()
			return err
		}

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			record.SetField(colName, *val)
		}
		query.result.appendItem(record)
		rowIndex++
	}
	rows.Close()
	return nil
}

func (query *bulkRetrieveQuery) getWhereClause(provider databaseprovider.DatabaseProvider) (string, []interface{}) {
	var sql strings.Builder
	var operator string
	var item *bulkRetrieveQueryItem
	var parameterId = 0
	var parameters []interface{}

	// may be two pass to reduce allocations
	for i := 0; i < len(*query.items); i++ {
		item = (*query.items)[i]
		operator = item.operation.ToSql(provider, item.operand)
		if operator != "" {
			sql.WriteString(item.field.GetPhysicalName(provider))
			sql.WriteString(operator)
			sql.WriteString(defaultParameterPrefix)
			sql.WriteString(defaultParameterName)
			sql.WriteString(strconv.Itoa(parameterId))
			parameterId++
		}
	}
	return sql.String(), parameters
}

//******************************
// private methods
//******************************
func newSimpleQuery(table *schema.Table) schema.Query {
	var query = new(bulkRetrieveQuery)
	var items = make([]*bulkRetrieveQueryItem, 0, 2)

	query.targetObject = table
	query.queryType = bulkquerytype.SimpleQuery
	query.items = &items
	query.result = new(List)

	return *query
}

func (query *bulkRetrieveQuery) addItem(item *bulkRetrieveQueryItem) {
	*query.items = append(*query.items, item)
}
