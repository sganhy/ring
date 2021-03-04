package data

import (
	"database/sql"
	"fmt"
	"ring/data/bulkquerytype"
	"ring/schema"
	"ring/schema/databaseprovider"
)

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
	fmt.Println("Execute query")
	fmt.Println(query.targetObject.GetDql(provider))
	rows, err := dbConnection.Query(query.targetObject.GetDql(provider))
	cols, _ := rows.Columns()
	if err != nil {
		return err
	}
	var rowIndex = 0

	count := query.targetObject.GetFieldCount()
	query.result.data = make([]*Record, 0, 10)
	fmt.Println(count)
	fmt.Println(query.targetObject.GetName())

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
		fmt.Println(query.result.Count())
		rowIndex++
	}
	fmt.Println("query.result.Count()==")
	fmt.Println(query.result.Count())
	rows.Close()
	return nil
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
