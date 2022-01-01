package schema

import (
	"context"
	"database/sql"
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"ring/schema/sqlfmt"
	"strconv"
	"strings"
	"time"
)

type metaQuery struct {
	schema         *Schema
	table          *Table
	result         []interface{}
	params         []interface{}
	columns        []interface{}
	columnsPointer []interface{}
	filters        []string
	sort           string
	query          string
	resultCount    int
	ddl            bool
	dml            bool
}

const (
	filterSeparator       string = " AND "
	operatorEqual         string = "="
	operatorPlus          string = "+"
	logTimeDescription    string = "time=%dms"
	logDefaultDescription string = "%s (done) | %s"
)

var (
	metaQueryLogger    = new(log)
	metaQueryUpdateSet []*Field
)

func init() {

	table := new(Table)
	metaTable := table.getMetaTable(databaseprovider.PostgreSql, "")
	// @meta table cache
	metaQueryUpdateSet = make([]*Field, 0, 6)
	metaQueryUpdateSet = append(metaQueryUpdateSet, metaTable.GetFieldByName(metaName))        // name
	metaQueryUpdateSet = append(metaQueryUpdateSet, metaTable.GetFieldByName(metaDescription)) // description
	metaQueryUpdateSet = append(metaQueryUpdateSet, metaTable.GetFieldByName(metaValue))       // value
	metaQueryUpdateSet = append(metaQueryUpdateSet, metaTable.GetFieldByName(metaActive))      // active
	metaQueryUpdateSet = append(metaQueryUpdateSet, metaTable.GetFieldByName(metaDataType))    // data_type
	metaQueryUpdateSet = append(metaQueryUpdateSet, metaTable.GetFieldByName(metaFlags))       // flags

}

func (query *metaQuery) Init(schema *Schema, table *Table) {
	query.table = table
	query.schema = schema
	query.resultCount = -1
}

//******************************
// getters and setters
//******************************
func (query *metaQuery) setSchema(schemaName string) {
	query.schema = GetSchemaByName(schemaName)
	if query.schema == nil {
		query.schema = getUpgradingSchema()
	}
	query.resultCount = -1
}

func (query *metaQuery) setTable(tableName string) {
	if query.schema == nil {
		// get meta schema by default
		query.schema = GetSchemaByName(metaSchemaName)
	}
	query.table = query.schema.GetTableByName(tableName)
	if query.table == nil {
		panic(fmt.Errorf("Unknown table %s for schema %s", tableName, query.schema.GetName()))
	}
	query.resultCount = -1
}

func (query *metaQuery) setParamValue(param interface{}, index int) {
	query.params[index] = param
}

//******************************
// public methods
//******************************
func (query *metaQuery) Execute(dbConnection *sql.DB, transaction *sql.Tx) error {
	var rows *sql.Rows
	var err error

	if query.ddl == true {
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err = dbConnection.ExecContext(ctx, query.query)
		fmt.Println(query.query)
		return err
	} else if query.dml == true {
		fmt.Println(query.query)
		rows, err = query.executeQuery(dbConnection, query.query)
		// avoid==> panic: pq: sorry, too many clients already
		rows.Close() //WARN: don't forget rows.Close()
		return err
	}

	fmt.Println(query.query)
	rows, err = query.executeQuery(dbConnection, query.query)

	if err != nil {
		fmt.Println("ERROR ==>")
		fmt.Println(err.Error())
		return err
	}
	query.resultCount = 0
	//fmt.Println(sqlQuery)
	if query.columns != nil {
		for rows.Next() {
			for i := range query.columns {
				query.columnsPointer[i] = &query.columns[i]
			}
			if err := rows.Scan(query.columnsPointer...); err != nil {
				fmt.Println(err)
				rows.Close()
				return err
			}
			query.result = append(query.result, query.table.GetQueryResult(query.columnsPointer))
		}
		query.resultCount = len(query.result)
	} else {
		for rows.Next() {
			query.resultCount++
		}
	}
	rows.Close()
	return nil
}

//******************************
// private methods
//******************************
func (query *metaQuery) getResult(index int, fieldName string) string {
	if query.result != nil && index >= 0 && index < len(query.result) {
		var record = query.result[index].([]string)
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
		query.params = make([]interface{}, 0, 2)
	}
	query.params = append(query.params, param)
}

func (query *metaQuery) addFilter(fieldName string, operator string, operand interface{}) {
	var field = query.getField(fieldName)
	if field != nil {
		if query.filters == nil {
			query.filters = make([]string, 0, 2)
		}
		if query.params == nil {
			query.params = make([]interface{}, 0, 2)
		}
		var variable = query.table.getVariableName(len(query.params))
		var queryFilter strings.Builder
		var fieldPhysicalName = field.GetPhysicalName(query.table.GetDatabaseProvider())
		// add single cote for varchar values
		queryFilter.Grow(len(fieldPhysicalName) + len(operator) + 8 + len(variable))
		queryFilter.WriteString(fieldPhysicalName)
		queryFilter.WriteString(operator)
		queryFilter.WriteString(variable)
		query.params = append(query.params, operand)
		query.filters = append(query.filters, queryFilter.String())
	}
}

func (query *metaQuery) addSort(fieldName string, ascending bool) {
	var field = query.getField(fieldName)
	if field != nil {
		if len(query.sort) > 0 {
			query.sort = query.sort + "," + field.GetPhysicalName(query.table.GetDatabaseProvider())
		} else {
			query.sort = field.GetPhysicalName(query.table.GetDatabaseProvider())
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

		var fieldCount = len(query.table.fields)
		var result = make([]meta, query.resultCount, query.resultCount)
		var tempMeta int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < query.resultCount; i++ {
			record = query.result[i].([]string)
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
	if query.resultCount > 0 {
		val, _ := strconv.ParseInt((query.result[0].([]string))[0], 10, 64)
		return val
	} else {
		return 0
	}
}

func (query *metaQuery) getMetaIdList() []metaId {
	if query.table != nil && query.table.GetName() == metaIdTableName {

		var fieldCount = len(query.table.fields)
		var result = make([]metaId, query.resultCount, query.resultCount)
		var tempMetaId int64 = 0
		var field *Field
		var j = 0
		var record []string

		for i := 0; i < query.resultCount; i++ {
			record = query.result[i].([]string)
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
			query.result = make([]interface{}, 0, 4)
		} else {
			query.result = make([]interface{}, 0, resultCount)
		}
		if len(query.query) <= 0 {
			query.query = query.getQuery()
		}
	} else {
		query.result = query.result[:0]
	}
	query.ddl = false
	query.dml = false
	if query.columns == nil {
		count := len(query.table.fields)
		query.columns = make([]interface{}, count)
		query.columnsPointer = make([]interface{}, count)
	}
	return query.schema.execute(query)
}

// is table exist
func (query *metaQuery) exists() (bool, error) {
	// sql empty ?
	if len(query.query) <= 0 {
		query.query = query.getQuery()
	}

	query.columns = nil
	query.columnsPointer = nil
	query.ddl = false
	query.dml = false
	query.schema.execute(query)

	return query.resultCount > 0, nil
}

// get
func (query *metaQuery) logDdl(id int32, jobId int64, ent entity, creationTime time.Time,
	statement ddlstatement.DdlStatement, err error, operation string) {
	var entityName = ent.GetPhysicalName()
	var logger = query.schema.logger
	var message string

	if ent.GetEntityType() == entitytype.Index {
		// get physical name
		var index = query.table.GetIndexByName(ent.GetName())
		entityName = index.getPhysicalName(query.table)
	}

	//logs
	if statement == ddlstatement.Undefined {
		message = sqlfmt.ToPascalCase(operation) + " " + sqlfmt.ToCamelCase(ent.GetEntityType().String())
	} else {
		message = sqlfmt.ToPascalCase(statement.String()) + " " + sqlfmt.ToCamelCase(ent.GetEntityType().String())
	}
	description := query.getLogDescription(entityName, ent, creationTime, statement, operation)

	if err == nil {
		if len(ent.GetPhysicalName()) > 0 {
			var log = logger.getNewLog(id, levelInfo, logger.schemaId, jobId, message, description)
			logger.writeToDb(log)
		}
	} else {
		logger.writePartialLog(id, levelError, jobId, message, description)
		logger.writePartialLog(id, levelError, jobId, err)
	}
}

func (query *metaQuery) getLogDescription(entityName string, ent entity, creationTime time.Time, statement ddlstatement.DdlStatement,
	operation string) string {

	var description string
	switch statement {
	case ddlstatement.Create, ddlstatement.Drop, ddlstatement.Truncate, ddlstatement.Undefined:
		description = query.getLogCreateDescription(entityName, ent, creationTime)
		break
	case ddlstatement.Alter:
		description = query.getLogAlterDescription(ent, creationTime, operation)
		break
	}
	return description
}

func (query *metaQuery) getLogCreateDescription(entityName string, ent entity, creationTime time.Time) string {
	description := entityName
	if query.table != nil &&
		(ent.GetEntityType() == entitytype.Index ||
			ent.GetEntityType() == entitytype.Constraint) {
		description += " on " + query.table.GetPhysicalName()
	}
	return fmt.Sprintf(logDefaultDescription, description, query.getLogTime(creationTime))
}

func (query *metaQuery) getLogAlterDescription(ent entity, creationTime time.Time, operation string) string {
	var description strings.Builder

	description.WriteString(ent.GetPhysicalName())
	description.WriteString(dqlSpace)

	if strings.Contains(query.query, dqlSpace+postGreDropColumn+dqlSpace) {
		description.WriteString(strings.ToLower(postGreDropColumn))
	} else {
		description.WriteString(strings.ToLower(postGreAddColumn))
	}
	if query.table != nil {
		if query.table.GetFieldByName(operation) != nil {
			description.WriteString(dqlSpace)
			description.WriteString(strings.ToLower(entitytype.Field.String()))
		}
		if query.table.GetRelationByName(operation) != nil {
			description.WriteString(dqlSpace)
			description.WriteString(strings.ToLower(entitytype.Relation.String()))
		}
	}
	description.WriteString(dqlSpace)
	description.WriteString(operation)

	return fmt.Sprintf(logDefaultDescription, description.String(), query.getLogTime(creationTime))
}

func (query *metaQuery) getLogTime(creationTime time.Time) string {
	duration := time.Now().Sub(creationTime)
	// for unitesting ==> we make sure that the duration is equal to zero
	if duration.Seconds() < 0 {
		duration = 0
	}
	return fmt.Sprintf(logTimeDescription, int(duration.Seconds()*1000))
}

// create ddl
func (query *metaQuery) create(id int32, jobId int64, ent entity) error {
	var creationTime = time.Now()

	query.ddl = true
	query.dml = false
	err := query.schema.execute(query)
	if ent.logStatement(ddlstatement.Create) {
		query.logDdl(id, jobId, ent, creationTime, ddlstatement.Create, err, operatorEqual)
	}
	return err
}

func (query *metaQuery) drop(id int32, jobId int64, ent entity) error {
	var eventTime = time.Now()

	query.ddl = true
	query.dml = false
	err := query.schema.execute(query)
	if ent.logStatement(ddlstatement.Drop) {
		query.logDdl(id, jobId, ent, eventTime, ddlstatement.Drop, err, operatorEqual)
	}
	return err
}

func (query *metaQuery) alter(id int32, jobId int64, ent entity, field *Field) error {
	var createTime = time.Now()
	query.ddl = true
	query.dml = false
	err := query.schema.execute(query)

	if ent.logStatement(ddlstatement.Alter) {
		query.logDdl(id, jobId, ent, createTime, ddlstatement.Alter, err, field.GetName())
	}
	return err
}

func (query *metaQuery) truncate(id int32, jobId int64, ent entity) error {
	var createTime = time.Now()
	query.ddl = true
	query.dml = false
	err := query.schema.execute(query)

	if ent.logStatement(ddlstatement.Truncate) {
		query.logDdl(id, jobId, ent, createTime, ddlstatement.Truncate, err, operatorEqual)
	}
	return err
}

func (query *metaQuery) vacuum(id int32, jobId int64, ent entity) error {
	var createTime = time.Now()
	query.ddl = true
	query.dml = false
	err := query.schema.execute(query)

	if ent.logStatement(ddlstatement.Undefined) {
		query.logDdl(id, jobId, ent, createTime, ddlstatement.Undefined, err, "Vacuum")
	}

	return err
}

func (query *metaQuery) analyze(id int32, jobId int64, ent entity) error {
	var eventTime = time.Now()
	query.ddl = true
	query.dml = false
	err := query.schema.execute(query)

	if ent.logStatement(ddlstatement.Undefined) {
		query.logDdl(id, jobId, ent, eventTime, ddlstatement.Undefined, err, "Analyze")
	}

	return err
}

// insert log
func (query *metaQuery) insert(params []interface{}) error {
	query.columns = nil
	query.columnsPointer = nil
	query.query = query.table.GetDml(dmlstatement.Insert, nil)
	query.params = params
	query.dml = true
	query.ddl = false
	return query.schema.execute(query)
}

func (query *metaQuery) update(params []interface{}) error {
	query.columns = nil
	query.columnsPointer = nil
	query.query = query.table.GetDml(dmlstatement.Update, metaQueryUpdateSet)
	query.params = params
	query.dml = true
	query.ddl = false
	return query.schema.execute(query)
}

func (query *metaQuery) delete(params []interface{}) error {
	query.columns = nil
	query.columnsPointer = nil
	query.query = query.table.GetDml(dmlstatement.Delete, nil)
	query.params = params
	query.dml = true
	query.ddl = false
	return query.schema.execute(query)
}

func (query *metaQuery) executeQuery(dbConn *sql.DB, sql string) (*sql.Rows, error) {
	if query.params == nil {
		return dbConn.Query(sql)
	}
	var params = query.params
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
	params[7] = metaData.description
	params[8] = metaData.value
	params[9] = metaData.enabled
	return query.insert(params)
}

func (query *metaQuery) updateMeta(metaData *meta, schemaId int32) error {
	var params []interface{}
	params = make([]interface{}, len(query.table.fields), len(query.table.fields))

	//meta(1) ==> SET name=$1, description=$2, value= $3, active= $5
	params[0] = metaData.name
	params[1] = metaData.description
	params[2] = metaData.value
	params[3] = metaData.enabled
	//
	params[4] = metaData.dataType
	params[5] = metaData.flags

	//meta(2) ==> key (id=$12 AND schema_id=$13 AND object_type=$14 AND reference_id=$15)
	params[6] = metaData.id
	params[7] = schemaId
	params[8] = metaData.objectType
	params[9] = metaData.refId

	return query.update(params)
}

func (query *metaQuery) deleteMeta(metaData *meta, schemaId int32, forceDelete bool) error {
	if forceDelete == false {
		metaData.enabled = false
		return query.updateMeta(metaData, schemaId)
	}
	var params []interface{}
	params = make([]interface{}, 4, 4)
	params[0] = metaData.id
	params[1] = schemaId
	params[2] = metaData.objectType
	params[3] = metaData.refId
	return query.delete(params)
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
