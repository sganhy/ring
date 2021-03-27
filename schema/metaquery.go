package schema

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

const filterSeparator = " AND "

type metaQuery struct {
	table   *Table
	result  *[]interface{}
	filters []string
	sort    string
}

func (query metaQuery) Execute(dbConnection *sql.DB) error {
	var sqlQuery = query.getQuery()

	rows, err := dbConnection.Query(sqlQuery)
	if err != nil {
		return err
	}

	fmt.Println(sqlQuery)
	count := len(query.table.fields)
	columns := make([]interface{}, count)
	columnPointers := make([]interface{}, count)
	for rows.Next() {
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			rows.Close()
			return err
		}
		*query.result = append(*query.result, query.table.GetQueryResult(columnPointers))
	}
	rows.Close()
	return nil
}

func (query *metaQuery) setTable(tableName string) {
	var metaSchema = GetSchemaByName(metaSchemaName)
	query.table = metaSchema.GetTableByName(tableName)
}

func (query *metaQuery) getResult(index int, fieldName string) string {
	if query.result != nil && index >= 0 && index < len(*query.result) {
		var record = (*query.result)[index].([]string)
		var index = query.table.GetFieldIndexByName(fieldName)
		return record[index]
	}
	return ""
}

func (query *metaQuery) resultCount() int {
	if query.result != nil {
		return len(*query.result)
	}
	return 0
}

func (query *metaQuery) getField(fieldName string) *Field {
	if query.table != nil {
		return query.table.GetFieldByName(fieldName)
	}
	return nil
}

func (query *metaQuery) addFilter(fieldName string, operator string, operand string) {
	var field = query.getField(fieldName)
	if field != nil {
		if query.filters == nil {
			query.filters = make([]string, 0, 2)
		}
		var queryFilter strings.Builder
		// add single cote for varchar values
		queryFilter.Grow(len(field.name) + len(operator) + 8 + len(operand))
		queryFilter.WriteString(field.GetPhysicalName(query.table.provider))
		queryFilter.WriteString(" ")
		queryFilter.WriteString(operator)
		queryFilter.WriteString(" ")
		queryFilter.WriteString(operand)
		query.filters = append(query.filters, queryFilter.String())
	}
}

func (query *metaQuery) addSort(fieldName string, ascending bool) {
	var field = query.getField(fieldName)
	if field != nil {
		if query.sort != "" {
			query.sort = query.sort + "," + field.GetPhysicalName(query.table.provider)
		} else {
			query.sort = field.GetPhysicalName(query.table.provider)
		}
		if ascending == false {
			query.sort = " DESC"
		}
	}
}

func (query *metaQuery) getFilters() string {
	if query.table != nil && query.filters != nil && len(query.filters) > 0 {
		var result strings.Builder
		for i := 0; i < len(query.filters); i++ {
			result.WriteString(query.filters[i])
			if i < len(query.filters)-1 {
				result.WriteString(filterSeparator)
			}
		}
		return result.String()
	}
	return ""
}

func (query *metaQuery) getQuery() string {
	if query.table != nil {
		return query.table.GetDql(query.getFilters(), query.sort)
	}
	return ""
}

func (query *metaQuery) getMetaList() []Meta {
	if query.table != nil && query.table.name == metaTableName {

		var resultCount = query.resultCount()
		var fieldCount = len(query.table.fields)
		var result = make([]Meta, resultCount, resultCount)
		var temp int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < resultCount; i++ {
			record = (*query.result)[i].([]string)
			for j = 0; j < fieldCount; j++ {
				field = query.table.fields[j]
				switch field.name {
				case metaId:
					temp, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].id = int32(temp)
					break
				case metaDataType:
					temp, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].dataType = int32(temp)
					break
				case metaName:
					result[i].name = record[j]
					break
				case metaDescription:
					result[i].description = record[j]
					break
				case metaFlags:
					temp, _ = strconv.ParseInt(record[j], 10, 64)
					result[i].flags = uint64(temp)
					break
				case metaObjectType:
					temp, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].objectType = int8(temp)
					break
				case metaReferenceId:
					temp, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].refId = int32(temp)
					break
				case metaValue:
					result[i].value = record[j]
					break
				}
			}
			result[i].lineNumber = 0
			result[i].enabled = true
		}
		return result
	}
	return nil
}

func (query *metaQuery) run() error {
	var metaSchema = GetSchemaByName(metaSchemaName) // get meta schema
	result := make([]interface{}, 0, 4)
	query.result = &result
	queries := make([]Query, 1, 1)
	queries[0] = *query
	return metaSchema.Execute(queries)
}
