package schema

import (
	"database/sql"
	"strings"
)

const filterSeparator = " AND "

type metaQuery struct {
	table   *Table
	result  *[][]string
	filters []string
}

func (query metaQuery) Execute(dbConnection *sql.DB) error {
	return nil
}

func (query *metaQuery) setTable(tableName string) {
	var metaSchema = GetSchemaByName(metaSchemaName)
	query.table = metaSchema.GetTableByName(tableName)

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
		if query.filters == nil {
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
		return query.table.GetDql(query.getFilters(), "")
	}
	return ""
}

func (query *metaQuery) run() error {
	var metaSchema = GetSchemaByName(metaSchemaName) // get meta schema
	result := make([][]string, 0, 4)
	query.result = &result
	queries := make([]Query, 1, 1)
	queries[0] = *query
	return metaSchema.Execute(queries)
}
