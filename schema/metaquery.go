package schema

import (
	"context"
	"database/sql"
	"fmt"
	"ring/schema/dmlstatement"
	"strconv"
	"strings"
	"time"
)

const filterSeparator = " AND "

type metaQuery struct {
	schema           *Schema
	table            *Table
	result           *[]interface{}
	params           *[]interface{}
	filters          []string
	sort             string
	query            string
	resultCount      *int
	returnResultList bool
	ddl              bool
	dml              bool
}

func (query metaQuery) Execute(dbConnection *sql.DB) error {
	var sqlQuery = query.query
	var columns []interface{}
	var columnPointers []interface{}
	var rows *sql.Rows
	var err error

	if query.ddl == true {
		var sqlResult sql.Result
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		sqlResult, err = dbConnection.ExecContext(ctx, sqlQuery)
		fmt.Println(sqlResult)
		return err
	} else if query.dml == true {
		rows, err = query.executeQuery(dbConnection, sqlQuery)
		return err
	} else {
		rows, err = query.executeQuery(dbConnection, sqlQuery)
	}

	if err != nil {
		fmt.Println("ERROR ==>")
		fmt.Println(err.Error())
		return err
	}
	*query.resultCount = 0
	//fmt.Println(sqlQuery)

	count := len(query.table.fields)
	if query.returnResultList == true {
		columns = make([]interface{}, count)
		columnPointers = make([]interface{}, count)
	}
	for rows.Next() {
		if query.returnResultList == true {
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
		*query.resultCount++
	}
	rows.Close()
	return nil
}

func (query *metaQuery) setSchema(schemaName string) {
	query.schema = GetSchemaByName(schemaName)
	query.resultCount = new(int)
}

func (query *metaQuery) setTable(tableName string) {
	if query.schema == nil {
		// get meta schema by default
		query.schema = GetSchemaByName(metaSchemaName)
	}
	query.table = query.schema.GetTableByName(tableName)
	query.resultCount = new(int)
	if query.table == nil {
		fmt.Errorf("Unknown table %s for schema %s", tableName, query.schema.name)
	}
}

func (query *metaQuery) getResult(index int, fieldName string) string {
	if query.result != nil && index >= 0 && index < len(*query.result) {
		var record = (*query.result)[index].([]string)
		var index = query.table.GetFieldIndexByName(fieldName)
		return record[index]
	}
	return ""
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

		var resultCount = *query.resultCount
		var fieldCount = len(query.table.fields)
		var result = make([]Meta, resultCount, resultCount)
		var tempMeta int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < resultCount; i++ {
			record = (*query.result)[i].([]string)
			for j = 0; j < fieldCount; j++ {
				field = query.table.fields[j]
				switch field.name {
				case metaId:
					tempMeta, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].id = int32(tempMeta)
					break
				case metaDataType:
					tempMeta, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].dataType = int32(tempMeta)
					break
				case metaName:
					result[i].name = record[j]
					break
				case metaDescription:
					result[i].description = record[j]
					break
				case metaFlags:
					tempMeta, _ = strconv.ParseInt(record[j], 10, 64)
					result[i].flags = uint64(tempMeta)
					break
				case metaObjectType:
					tempMeta, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].objectType = int8(tempMeta)
					break
				case metaReferenceId:
					tempMeta, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].refId = int32(tempMeta)
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

func (query *metaQuery) getMetaIdList() []MetaId {
	if query.table != nil && query.table.name == metaIdTableName {

		var resultCount = *query.resultCount
		var fieldCount = len(query.table.fields)
		var result = make([]MetaId, resultCount, resultCount)
		var tempMetaId int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < resultCount; i++ {
			record = (*query.result)[i].([]string)
			for j = 0; j < fieldCount; j++ {
				field = query.table.fields[j]
				switch field.name {
				case metaId:
					tempMetaId, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].id = int32(tempMetaId)
					break
				case metaSchemaId:
					tempMetaId, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].schemaId = int32(tempMetaId)
					break
				case metaObjectType:
					tempMetaId, _ = strconv.ParseInt(record[j], 10, 32)
					result[i].objectType = int8(tempMetaId)
					break
				case metaValue:
					tempMetaId, _ = strconv.ParseInt(record[j], 10, 64)
					result[i].value = tempMetaId
					break
				}
			}
		}
		return result
	}
	return nil
}

// launch select query with result list
func (query *metaQuery) run() error {
	result := make([]interface{}, 0, 4)
	query.ddl = false
	query.dml = false
	query.query = query.getQuery()
	query.returnResultList = true
	query.result = &result
	queries := make([]Query, 1, 1)
	queries[0] = *query
	return query.schema.Execute(queries)
}

// is table exist
func (query *metaQuery) exists() (bool, error) {
	query.returnResultList = false
	queries := make([]Query, 1, 1)
	query.ddl = false
	query.dml = false
	if query.resultCount == nil {
		query.resultCount = new(int)
	}
	//fmt.Println(query.query)
	queries[0] = *query
	query.schema.Execute(queries)
	//fmt.Println("query count() ==> ")
	//fmt.Println(*query.resultCount)
	return *query.resultCount > 0, nil
}

// create ddl
func (query *metaQuery) create() error {
	query.ddl = true
	query.dml = false
	queries := make([]Query, 1, 1)
	queries[0] = *query
	return query.schema.Execute(queries)
}

// insert log
func (query *metaQuery) insert(params []interface{}) error {
	query.returnResultList = false
	query.query = query.table.GetDml(dmlstatement.Insert)
	query.params = &params
	query.dml = true
	query.ddl = false
	queries := make([]Query, 1, 1)
	queries[0] = *query
	return query.schema.Execute(queries)
}

func (query *metaQuery) executeQuery(dbConn *sql.DB, sql string) (*sql.Rows, error) {
	if query.params == nil {
		return dbConn.Query(sql)
	}
	var params = *query.params
	switch len(params) {
	case 0:
		return dbConn.Query(sql)
	case 1:
		return dbConn.Query(sql, params[0])
	case 2:
		return dbConn.Query(sql, params[0], params[1])
	case 3:
		return dbConn.Query(sql, params[0], params[1], params[2])
	case 4:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3])
	case 5:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4])
	case 6:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5])
	case 7:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5], params[6])
	case 8:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7])
	case 9:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7],
			params[8])
	case 10:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7],
			params[8], params[9])
	case 11:
		return dbConn.Query(sql, params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7],
			params[8], params[9], params[10])
	}
	return nil, nil
}

func (query *metaQuery) insertLog(newLog *log) error {
	var params []interface{}
	params = make([]interface{}, len(query.table.fields), len(query.table.fields))
	params[0] = newLog.id
	params[1] = newLog.entryTime
	params[2] = newLog.level
	params[3] = newLog.schemaId
	params[4] = newLog.threadId
	params[5] = newLog.callSite
	params[6] = newLog.jobId
	params[7] = newLog.method
	params[8] = newLog.lineNumber
	params[9] = newLog.message
	params[10] = newLog.description
	return query.insert(params)
}
