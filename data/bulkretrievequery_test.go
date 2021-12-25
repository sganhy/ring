package data

import (
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

//test Execute()
func Test__bulkRetrieveQuery__Execute(t *testing.T) {
	schema.Init(databaseprovider.PostgreSql, "", 0, 0)
	sche := schema.GetDefaultSchema()
	table := sche.GetTableByName("@meta")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"id", "schema_id", "object_type", "reference_id", "data_type", "flags", "name", "description",
		"value", "active"}
	rs := mock.NewRows(columns)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "ability", "", "ability", true)
	baseSql := "SELECT id,schema_id,object_type,reference_id,data_type,flags,\"name\",description,\"value\",active FROM " +
		table.GetPhysicalName()

	//======================
	//==== 0 arguments + order by
	//======================
	bulkQuery := newSimpleQuery(table).(bulkRetrieveQuery)
	mock.ExpectQuery(baseSql + " ORDER BY object_type,reference_id DESC").
		WithArgs().
		WillReturnRows(rs)
	field := table.GetFieldByName("object_type")
	bulkQuery.addSort(newQuerySort(field, sortordertype.Ascending))
	field = table.GetFieldByName("reference_id")
	bulkQuery.addSort(newQuerySort(field, sortordertype.Descending))
	err = bulkQuery.Execute(db, nil)
	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}
	if err != nil {
		t.Errorf("ERROR: '%s'", err)
	}

	//======================
	//==== 1 argument
	//======================
	bulkQuery.result = new(List)
	db, mock, err = sqlmock.New()
	rs = mock.NewRows(columns)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "ability", "", "ability", true)

	bulkQuery = newSimpleQuery(table).(bulkRetrieveQuery) // reset query
	mock.ExpectQuery(baseSql + " WHERE schema_id=\\$1").
		WithArgs(1).
		WillReturnRows(rs)
	queryItem2 := new(bulkRetrieveQueryItem)
	queryItem2.field = table.GetFieldByName("schema_id")
	queryItem2.operation = operationtype.Equal
	queryItem2.operand = "1"
	bulkQuery.addFilter(queryItem2)

	bulkQuery.Execute(db, nil)
	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}

	//======================
	//==== 2 arguments
	//======================
	bulkQuery.result = new(List)
	db, mock, err = sqlmock.New()
	rs = mock.NewRows(columns)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "ability", "", "ability", true)

	mock.ExpectQuery(baseSql+" WHERE schema_id=\\$1 AND \"name\" LIKE \\$2").
		WithArgs(1, "abi%").
		WillReturnRows(rs)
	queryItem3 := new(bulkRetrieveQueryItem)
	queryItem3.field = table.GetFieldByName("name")
	queryItem3.operation = operationtype.Like
	queryItem3.operand = "abi%"
	bulkQuery.addFilter(queryItem3)

	bulkQuery.Execute(db, nil)
	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}

	//======================
	//==== 17 arguments
	//======================
	bulkQuery.result = new(List)
	db, mock, err = sqlmock.New()
	rs = mock.NewRows(columns)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "ability", "", "ability", true)
	bulkQuery.clearItems()
	mock.ExpectQuery(baseSql+" WHERE schema_id=\\$1 AND schema_id=\\$2 AND schema_id=\\$3 AND schema_id=\\$4"+
		" AND schema_id=\\$5 AND schema_id=\\$6 AND schema_id=\\$7 AND schema_id=\\$8 AND schema_id=\\$9 AND schema_id=\\$10"+
		" AND schema_id=\\$11 AND schema_id=\\$12 AND schema_id=\\$13 AND schema_id=\\$14 AND schema_id=\\$15"+
		" AND schema_id=\\$16 AND schema_id=\\$17 ORDER BY reference_id DESC").
		WithArgs(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17).
		WillReturnRows(rs)
	for i := 0; i < 17; i++ {
		var newfilter = queryItem2.Clone()
		newfilter.operand = strconv.Itoa(i + 1)
		bulkQuery.addFilter(newfilter)
	}
	field = table.GetFieldByName("reference_id")
	bulkQuery.addSort(newQuerySort(field, sortordertype.Descending))
	bulkQuery.Execute(db, nil)
	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}
	//======================
	//==== 100 arguments
	//======================
	bulkQuery.result = new(List)
	db, mock, err = sqlmock.New()
	rs = mock.NewRows(columns)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "skill", "", "skill", true)
	bulkQuery.clearItems()
	mock.ExpectQuery(baseSql+" WHERE schema_id=\\$1 AND schema_id=\\$2 AND schema_id=\\$3 AND schema_id=\\$4"+
		" AND schema_id=\\$5 AND schema_id=\\$6").
		WithArgs(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1).
		WillReturnRows(rs)
	for i := 0; i < 100; i++ {
		bulkQuery.addFilter(queryItem2)
	}
	bulkQuery.Execute(db, nil)
	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}

	//======================
	//==== 255 arguments
	//======================
	bulkQuery.result = new(List)
	db, mock, err = sqlmock.New()
	rs = mock.NewRows(columns)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "skill", "", "skill", true)
	bulkQuery.clearItems()
	mock.ExpectQuery(baseSql+" WHERE schema_id=\\$1 AND schema_id=\\$2 AND schema_id=\\$3 AND schema_id=\\$4"+
		" AND schema_id=\\$5").
		WithArgs(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1).
		WillReturnRows(rs)
	for i := 0; i < 255; i++ {
		bulkQuery.addFilter(queryItem2)
	}
	bulkQuery.Execute(db, nil)
	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}

	//======================
	//==== loop 1 TO 254 arguments
	//======================
	bulkQuery.clearItems()
	queryItem2.operand = "444"
	var expectedQuery = "SELECT id,schema_id,object_type,reference_id,data_type,flags,name,description,value,active FROM " +
		"information_schema.\"@meta\" WHERE schema_id=\\$1"
	for i := 0; i < 254; i++ {
		rs = mock.NewRows(columns)
		rs.AddRow(1, 1, 0, 0, 0, 159744, "skill", "", "skill", true)
		mock.ExpectQuery(expectedQuery).
			WithArgs().
			WillReturnRows(rs)
		bulkQuery.addFilter(queryItem2)
		bulkQuery.Execute(db, nil)
		if bulkQuery.result.Count() <= 0 {
			t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
		}
		expectedQuery = expectedQuery + " AND schema_id=\\$" + strconv.Itoa(i+2)
	}

	//======================
	//==== out of range
	//======================
	bulkQuery.addFilter(queryItem2)
	bulkQuery.addFilter(queryItem2)
	rs.AddRow(1, 1, 0, 0, 0, 159744, "skill", "", "skill", true)
	mock.ExpectQuery("SELECT id,schema_id,object_type,reference_id,data_type,flags,name,description,value,active FROM").
		WithArgs().
		WillReturnRows(rs)
	err = bulkQuery.Execute(db, nil)
	if err == nil {
		t.Errorf("bulkRetrieveQuery.Execute() ==> must return an error")
	}

}
