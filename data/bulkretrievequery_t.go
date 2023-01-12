package data

import (
	//	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

//test Execute()
func Test__bulkRetrieveQuery__Execute(t *testing.T) {

	table := getTestTable(databaseprovider.PostgreSql, "Test")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{"id", "test1", "test2", "test3", "test4"}
	rows := mock.NewRows(columns)
	rows.AddRow(1, "ability", "", "ability", 1)
	baseSql := "SELECT id,schema_id,object_type,reference_id,data_type,flags,\"name\",description,\"value\",active FROM " +
		table.GetPhysicalName()

	//======================
	//==== 0 arguments + order by
	//======================
	bulkQuery := newSimpleQuery(table).(bulkRetrieveQuery)
	mock.ExpectQuery(baseSql + " ORDER BY object_type,reference_id DESC").
		WillReturnRows(rows)
	field := table.GetFieldByName("tes1")
	bulkQuery.addSort(newQuerySort(field, sortordertype.Ascending))
	field = table.GetFieldByName("test2")
	bulkQuery.addSort(newQuerySort(field, sortordertype.Descending))
	err = bulkQuery.Execute(db, nil)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if bulkQuery.result.Count() <= 0 {
		t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() must be greater than 0")
	}
	if err != nil {
		t.Errorf("ERROR: '%s'", err)
	}
	defer db.Close()

	//======================
	//==== 1 argument
	//======================
	/*
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
			t.Errorf("bulkRetrieveQuery.Execute() ==> result.Count() should be greater than 0")
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
		/*
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
	*/
	//======================
	//==== out of range
	//======================
	/*
		bulkQuery.result = new(List)
		db, mock, err = sqlmock.New()
		rs = mock.NewRows(columns)
		rs.AddRow(1, 1, 0, 0, 0, 159744, "skill", "", "skill", true)
		bulkQuery.clearItems()
		mock.ExpectQuery("SELECT id,schema_id,object_type,reference_id,data_type,flags,name,description,value,active FROM").
			WithArgs().
			WillReturnRows(rs)
		err = bulkQuery.Execute(db, nil)
		if err == nil {
			t.Errorf("bulkRetrieveQuery.Execute() ==> must return an error")
		}
	*/
}

func getTestTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *schema.Table {
	var fields = make([]schema.Field, 0, 4)
	var relations = make([]schema.Relation, 0, 0)
	var indexes = make([]schema.Index, 0, 0)
	var result = new(schema.Table)

	// physical_name is built later
	//  == metaId table
	var id = schema.Field{}
	var test1 = schema.Field{}
	var test2 = schema.Field{}
	var test3 = schema.Field{}
	var test4 = schema.Field{}

	// !!!! id field must be greater than 0 !!!!
	id.Init(1103, "id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	test1.Init(1117, "test1", "", fieldtype.String, 50, "", true, true, true, false, true)
	test2.Init(1151, "test2", "", fieldtype.String, 50, "", true, true, true, false, true)
	test3.Init(1181, "test3", "", fieldtype.String, 50, "", true, true, true, false, true)
	test4.Init(1182, "test4", "", fieldtype.Int, 0, "", true, true, true, false, true)

	fields = append(fields, id)
	fields = append(fields, test1)
	fields = append(fields, test2)
	fields = append(fields, test3)
	fields = append(fields, test4)

	result.Init(22, "test", "", fields, relations, indexes, physicaltype.Table, 0, schemaPhysicalName, tabletype.Business, provider, "", true, false, true, true)

	return result
}
