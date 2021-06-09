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

const (
	filterSeparator string = " AND "
	operatorEqual   string = "="
	operatorPlus    string = "+"
)

var (
	metaQueryLogger = new(log)
)

func (query *metaQuery) Init(schema *Schema, table *Table) {
	query.table = table
	query.schema = schema
	if query.resultCount == nil {
		query.resultCount = new(int)
	}
}

//******************************
// getters and setters
//******************************
func (query *metaQuery) setSchema(schemaName string) {
	query.schema = GetSchemaByName(schemaName)
	if query.resultCount == nil {
		query.resultCount = new(int)
	}
}

func (query *metaQuery) setTable(tableName string) {
	if query.schema == nil {
		// get meta schema by default
		query.schema = GetSchemaByName(metaSchemaName)
	}
	query.table = query.schema.GetTableByName(tableName)
	if query.resultCount == nil {
		query.resultCount = new(int)
	}
	if query.table == nil {
		fmt.Errorf("Unknown table %s for schema %s", tableName, query.schema.GetName())
	}
}

func (query *metaQuery) setParamValue(param interface{}, index int) {
	(*query.params)[index] = param
}

func (query *metaQuery) getTable() *Table {
	return query.table
}

//******************************
// public methods
//******************************
func (query metaQuery) Execute(dbConnection *sql.DB) error {
	var sqlQuery = query.query
	var columns []interface{}
	var columnPointers []interface{}
	var rows *sql.Rows
	var err error

	if query.ddl == true {
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err = dbConnection.ExecContext(ctx, sqlQuery)
		fmt.Println(sqlQuery)
		return err
	} else if query.dml == true {
		fmt.Println(sqlQuery)
		rows, err = query.executeQuery(dbConnection, sqlQuery)

		// avoid==> panic: pq: sorry, too many clients already
		rows.Close() //WARN: don't forget rows.Close()
		return err
	} else {
		fmt.Println(sqlQuery)
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

//******************************
// private methods
//******************************
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

func (query *metaQuery) addParam(param interface{}) {
	if query.params == nil {
		params := make([]interface{}, 0, 2)
		query.params = &params
	}
	*query.params = append(*query.params, param)
}

func (query *metaQuery) addFilter(fieldName string, operator string, operand interface{}) {
	var field = query.getField(fieldName)
	if field != nil {
		if query.filters == nil {
			query.filters = make([]string, 0, 2)
		}
		if query.params == nil {
			params := make([]interface{}, 0, 2)
			query.params = &params
		}
		var variable = query.table.getVariableName(len(*query.params))
		var queryFilter strings.Builder
		var fieldPhysicalName = field.GetPhysicalName(query.table.provider)
		// add single cote for varchar values
		queryFilter.Grow(len(fieldPhysicalName) + len(operator) + 8 + len(variable))
		queryFilter.WriteString(fieldPhysicalName)
		queryFilter.WriteString(operator)
		queryFilter.WriteString(variable)
		*query.params = append(*query.params, operand)
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

func (query *metaQuery) getMetaList() []meta {
	if query.table != nil && query.table.GetName() == metaTableName {

		var resultCount = *query.resultCount
		var fieldCount = len(query.table.fields)
		var result = make([]meta, resultCount, resultCount)
		var tempMeta int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < resultCount; i++ {
			record = (*query.result)[i].([]string)
			for j = 0; j < fieldCount; j++ {
				field = query.table.fields[j]
				switch field.GetName() {
				case metaFieldId:
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
				case metaPhysicalName:
					result[i].physicalName = record[j]
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

func (query *metaQuery) getInt64Value() int64 {
	if *query.resultCount > 0 {
		val, _ := strconv.ParseInt(((*query.result)[0].([]string))[0], 10, 64)
		return val
	} else {
		return 0
	}
}

func (query *metaQuery) getMetaIdList() []metaId {
	if query.table != nil && query.table.GetName() == metaIdTableName {

		var resultCount = *query.resultCount
		var fieldCount = len(query.table.fields)
		var result = make([]metaId, resultCount, resultCount)
		var tempMetaId int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < resultCount; i++ {
			record = (*query.result)[i].([]string)
			for j = 0; j < fieldCount; j++ {
				field = query.table.fields[j]
				switch field.GetName() {
				case metaFieldId:
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
//    resultCount==0 for unknown expected result count
func (query *metaQuery) run(resultCount int) error {
	if query.result == nil {
		if resultCount == 0 {
			result := make([]interface{}, 0, 4)
			query.result = &result
		} else {
			result := make([]interface{}, 0, resultCount)
			query.result = &result
		}
		if query.query == "" {
			query.query = query.getQuery()
		}
	} else {
		*query.result = (*query.result)[:0]
	}
	query.ddl = false
	query.dml = false
	query.returnResultList = true
	return query.schema.execute(query)
}

// is table exist
func (query *metaQuery) exists() (bool, error) {
	// sql empty ?
	if query.query == "" {
		query.query = query.getQuery()
	}
	query.returnResultList = false
	query.ddl = false
	query.dml = false
	if query.resultCount == nil {
		query.resultCount = new(int)
	}
	query.schema.execute(query)
	return *query.resultCount > 0, nil
}

// create ddl
func (query *metaQuery) create() error {
	query.ddl = true
	query.dml = false
	return query.schema.execute(query)
}

func (query *metaQuery) truncate() error {
	return query.create()
}

func (query *metaQuery) vacuum() error {
	return query.create()
}

func (query *metaQuery) analyze() error {
	return query.create()
}

// insert log
func (query *metaQuery) insert(params []interface{}) error {
	query.returnResultList = false
	query.query = query.table.GetDml(dmlstatement.Insert, nil)
	query.params = &params
	query.dml = true
	query.ddl = false
	return query.schema.execute(query)
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
	params[0] = newLog.getId()
	params[1] = newLog.getEntryTime()
	params[2] = newLog.getLevel()
	params[3] = newLog.getSchemaId()
	params[4] = newLog.getThreadId()
	params[5] = newLog.getClassSite()
	params[6] = newLog.getJobId()
	params[7] = newLog.getMethod()
	params[8] = newLog.getLineNumber()
	params[9] = newLog.getMessage()
	params[10] = newLog.getDescription()

	return query.insert(params)
}

func (query *metaQuery) insertMeta(metaData *meta, schemaId int32) error {
	var params []interface{}
	params = make([]interface{}, len(query.table.fields), len(query.table.fields))
	params[0] = metaData.id
	params[1] = schemaId
	params[2] = metaData.objectType
	params[3] = metaData.refId
	params[4] = metaData.dataType
	params[5] = metaData.flags
	params[6] = metaData.name
	params[7] = metaData.physicalName
	params[8] = metaData.description
	params[9] = metaData.value
	params[10] = metaData.enabled
	return query.insert(params)
}

func (query *metaQuery) insertMetaId(metaid *metaId) error {
	var params []interface{}
	params = make([]interface{}, len(query.table.fields), len(query.table.fields))
	params[0] = metaid.id
	params[1] = metaid.schemaId
	params[2] = metaid.objectType
	params[3] = metaid.value
	return query.insert(params)
}
