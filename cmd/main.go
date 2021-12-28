package main

import (
	"fmt"
	"ring/data"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/sourcetype"
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

	rcd.SetRecordType("RpgSheet.skill")
	rcd.SetField("id", 16)
	rcd.SetField("name", "tété")
	var bs = new(data.BulkSave)
	bs.InsertRecord(rcd)
	//	bs.InsertRecord(rcd)
	bs.Save()

	var importFile = schema.Import{}
	importFile.Init(sourcetype.XmlDocument, "C:\\Temp\\Coding\\rpg_schema.xml")
	importFile.Load()
	importFile.Upgrade()

	var metaSchema = schema.GetSchemaByName("@meta")
	fmt.Println(metaSchema.GetName())
	var meto = metaSchema.ToMeta()
	var lexicon = metaSchema.GetTableByName("@lexicon")
	fmt.Println(lexicon.GetName())
	for i := 0; i < len(meto); i++ {
		if meto[i].GetEntityType() == 0 {
			fmt.Println(meto[i].String())
		}
	}
	/*
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
	var br = new(data.BulkRetrieve)
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

}
