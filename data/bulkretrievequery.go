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
const maxFilterCount = 255
const filterSeparator = " AND "
const errorMessageMaxVariable = "Too many filters (%d max)"

//type queryFunc func(dbConnection *sql.DB, sql string, params []interface{}) (*sql.Rows, error)

// all this structure is readonly
type bulkRetrieveQuery struct {
	targetObject *schema.Table
	queryType    bulkretrievetype.BulkRetrieveType
	filterCount  *int
	items        *[]*bulkRetrieveQueryItem
	result       *List // pointer is mandatory
}

var (
	dqlQuery *bulkQuery
)

func init() {
	dqlQuery = new(bulkQuery)
}

//******************************
// public methods (Interface schema.Query implementations)
//******************************
func (query bulkRetrieveQuery) Execute(dbConnection *sql.DB) error {
	var provider = query.targetObject.GetDatabaseProvider()
	var whereClause, parameters = query.getWhereClause(provider)
	var orderClause = query.getOrderClause(provider)
	var sqlQuery = query.targetObject.GetDql(whereClause, orderClause)

	rows, err := dqlQuery.Execute(dbConnection, sqlQuery, parameters)

	if err != nil {
		//fmt.Println(err)
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

/* C# source code to generate ==> executeQuery() body
   const int maxCount = 257;
        static void Main(string[] args)
        {
            for (var i = 27; i < maxCount; ++i)
            {
                Console.WriteLine(i.ToString() + " : executeQuery{0}".Replace("{0}", i.ToString().PadLeft(4, '0'))+",");
            }
            for (var i = 255; i < maxCount; ++i)
             {
                string method = string.Empty;
                method+="func executeQuery{0}(con *sql.DB, sql string, params []interface{}) (*sql.Rows, error) {".Replace("{0}", i.ToString().PadLeft(4, '0'))+"\n";
                method+= "		return con.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7]," + "\n";
                method+= "			params[8], params[9], params[10], params[11], params[12], params[13], params[14], params[15], params[16], params[17]," + "\n";

                var j = 18;
                    string str=string.Empty;
                    for (; j < i - 9; j+=10)
                    {
                        var template = "			params[{0}], params[{1}], params[{2}], params[{3}], params[{4}], params[{5}], params[{6}], params[{7}], params[{8}], params[{9}],\n";
                        str += string.Format(template, j , j + 1, j+2 , j + 3, j+4, j + 5,  j + 6,  j + 7,  j + 8,  j + 9);
                    }
                    if (j >= i)
                    {
                        str = str.Substring(0, str.Length - 2);
                        method += str;
                    }
                    else
                    {

                        method += str;
                        method +=  "			";
                        str = string.Empty;
                        while (j < i)
                        {
                            str+="params[" + j + "], ";
                            ++j;
                        }
                        str = str.Substring(0, str.Length - 2);
                    method += str;
                    }
                    method += ")\n";
                    method += "}\n";
                    Console.Write(method.Replace("params", "p" + i.ToString()));
                }

            }
*/
