package main

import (
	"fmt"
	"ring/data"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/sourcetype"
	"ring/schema/tabletype"
	"runtime"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const minint32 string = "-2147483648"

//**************************
// configuration methods
//**************************
func main() {

	data.Test__bulkRetrieveQuery__Execute(nil)
	elemf0 := new(schema.Field)
	elemf0.Init(11, "aName Test", "AField Test", fieldtype.DateTime, 10, "test default", true, false, false, true, true)
	value, _ := elemf0.GetValue("2020-04-30T00:00:00.000Z")
	fmt.Println(value)
	fmt.Println("2020-04-29Tr22:00:00.000")
	rcd := new(data.Record)

	schema.Init(databaseprovider.PostgreSql, "host=localhost port=5432 user=postgres password=sa dbname=postgres sslmode=disable", 10, 20)
	//schema.Init(databaseprovider.SqlServer, "server=localhost;User id=NA_USER;Password=NA_USER_PWD;database=CQL_CIV;port=1434", 10, 20)
	//server=SAKHALOO-PC;user id=sakhaloo;password=hoollehayerazi;database=webApp
	//schema.Init(databaseprovider.MySql, "root:root@tcp(127.0.0.1:3306)/mysql", 10, 20)

	var br = new(data.BulkRetrieve)
	br.SimpleQuery(0, "@meta.@lexicon")

	rcd.SetRecordType("RpgSheet.skill")
	rcd.SetField("id", 4365)
	rcd.SetField("name", "Zorba le grec")
	rcd.SetField("category", "CAT")
	var impData = new(data.Import)
	impData.SetSchema("RpgSheet")
	impData.ParseFile("C:/Temp/Coding/rpg_sheet_dd3.xlsx")

	var bs = new(data.BulkSave)
	bs.UpdateRecord(rcd)
	rcd.SetField("category", "CAT6")
	bs.ChangeRecord(rcd)
	//	bs.InsertRecord(rcd)
	bs.Save()
	/*
		bs.Clear()
		bs.Clear()
		bs.InsertRecord(rcd)
		//	bs.InsertRecord(rcd)
		bs.Save()
		bs.Clear()
		bs.InsertRecord(rcd)
		bs.InsertRecord(rcd)
		bs.InsertRecord(rcd)
		//	bs.InsertRecord(rcd)
		bs.Save()
		bs.Clear()
		bs.InsertRecord(rcd)
		bs.Save()
	*/
	var importFile = schema.Import{}
	importFile.Init(sourcetype.XmlDocument, "C:\\Temp\\Coding\\quotify_schema.xml")
	importFile.Load()
	importFile.Upgrade()
	var metaSchema = schema.GetSchemaByName("@meta")
	var seq1 = metaSchema.GetSequenceByName("@job_id")
	fmt.Println(seq1.NextValue())
	fmt.Println(seq1.NextValue())
	fmt.Println(seq1.NextValue())
	fmt.Println(seq1.NextValue())
	fmt.Println(seq1.NextValue())

	/*

		fmt.Println(seq1.NextValue())
		fmt.Println(seq1.NextValue())
		fmt.Println(metaSchema.GetName())
		var meto = metaSchema.ToMeta()
		var lexicon = metaSchema.GetTableByName("@lexicon")
		fmt.Println(lexicon.GetName())
		for i := 0; i < len(meto); i++ {
			if meto[i].GetEntityType() == 0 {
				fmt.Println(meto[i].String())
			}
		}

		if metaSchema != nil {
			fmt.Println(metaSchema.GetId())
			var tableBook = metaSchema.GetTableByName("feat14")
			if tableBook != nil {
				fmt.Println(metaSchema.GetId())
				fmt.Println(tableBook.GetId())
			}
		}
	*/

	location, _ := time.LoadLocation("MST")
	ttt := time.Now().In(location)
	zone, offset := ttt.Zone()
	fmt.Println(offset)
	fmt.Println(zone)
	rcd = new(data.Record)
	rcd.SetRecordType("@log")
	rcd.SetField("entry_time", "2014-04-15T21:14:55")
	fmt.Println(rcd.GetField("entry_time"))
	rcd.SetField("entry_time", time.Now())
	fmt.Println(rcd.GetField("entry_time"))

	// Create an empty user and make the sql query (using $1 for the parameter)
	br = new(data.BulkRetrieve)
	br.SimpleQuery(0, "@meta")
	br.AppendFilter(0, "schema_id", operationtype.Equal, 0)
	br.AppendFilter(0, "object_type", operationtype.Equal, 15)
	br.AppendSort(0, "id", sortordertype.Descending)
	br.AppendSort(0, "schema_id", sortordertype.Ascending)
	br.AppendSort(0, "name", sortordertype.Descending)
	br.AppendFilter(0, "name", operationtype.NotEqual, nil)
	br.SimpleQuery(1, "@meta_id")
	br.AppendFilter(1, "schema_id", operationtype.Equal, 1)
	/*
		br.AppendFilter(1, "schema_id", operationtype.NotEqual, 2)
		br.AppendFilter(1, "schema_id", operationtype.NotEqual, 3)
		br.AppendFilter(1, "schema_id", operationtype.NotEqual, 4)
	*/
	br.RetrieveRecords()

	runtime.GC()
	time.Sleep(3 * time.Second)

}

func getTable() *schema.Table {
	var fields = []schema.Field{}
	var relations = []schema.Relation{}
	var indexes = []schema.Index{}
	var table = new(schema.Table)
	var uk schema.Index = schema.Index{}

	// physical_name is builded later
	//  == metaId table
	var id schema.Field = schema.Field{}
	var schemaId schema.Field = schema.Field{}
	var objectType schema.Field = schema.Field{}
	var referenceId schema.Field = schema.Field{}
	var dataType schema.Field = schema.Field{}

	var flags schema.Field = schema.Field{}
	var value schema.Field = schema.Field{}
	var name schema.Field = schema.Field{}
	var description schema.Field = schema.Field{}
	var active schema.Field = schema.Field{}

	// elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true, true)
	// !!!! id field must be greater than 0 !!!!
	id.Init(1009, "id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1013, "schema_id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1019, "object_type", "", fieldtype.Byte, 0, "", true, true, true, false, true)
	referenceId.Init(1021, "reference_id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	dataType.Init(1031, "data_type", "", fieldtype.Int, 0, "", true, false, true, false, true)

	flags.Init(1039, "flags", "", fieldtype.Long, 0, "", true, false, true, false, true)
	name.Init(1061, "name", "", fieldtype.String, 30, "", true, false, false, false, true)
	description.Init(1069, "description", "", fieldtype.String, 0, "", true, false, true, false, true)
	value.Init(1087, "value", "", fieldtype.String, 0, "", true, false, true, false, true)
	active.Init(1093, "active", "", fieldtype.Boolean, 0, "", true, false, true, false, true)

	// elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
	// unique key (1)      id; schema_id; reference_id; object_type
	var indexedFields = []string{id.GetName(), name.GetName(), objectType.GetName(), referenceId.GetName()}
	uk.Init(1, "pk_@meta", "ATable Test", indexedFields, false, false, true, true)

	fields = append(fields, id)          //1
	fields = append(fields, schemaId)    //2
	fields = append(fields, objectType)  //3
	fields = append(fields, referenceId) //4
	fields = append(fields, dataType)    //5
	fields = append(fields, flags)       //6
	fields = append(fields, name)        //7
	fields = append(fields, description) //8
	fields = append(fields, value)       //9
	fields = append(fields, active)      //10

	indexes = append(indexes, uk)

	//id int32, name string, description string, fields []Field, relations []Relation, indexes []Index, physicalName string,
	//  physicalType physicaltype.PhysicalType, schemaId int32, tableType tabletype.TableType, subject string,
	//  cached bool, readonly bool, baseline bool, active bool
	table.Init(21, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "information_schema", tabletype.MetaId, databaseprovider.PostgreSql,
		"[subject]", true, false, true, false)
	return table
}
