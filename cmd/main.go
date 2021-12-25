package main

import (
	"fmt"
	"ring/data"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"ring/schema/searchabletype"
	"ring/schema/sourcetype"
	"runtime"
	"strings"
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

	schema.Init(databaseprovider.PostgreSql, "host=localhost port=5432 user=postgres password=sa dbname=postgres sslmode=disable", 10, 20)
	//schema.Init(databaseprovider.SqlServer, "server=localhost;User id=NA_USER;Password=NA_USER_PWD;database=CQL_CIV;port=1434", 10, 20)
	//server=SAKHALOO-PC;user id=sakhaloo;password=hoollehayerazi;database=webApp
	//schema.Init(databaseprovider.MySql, "root:root@tcp(127.0.0.1:3306)/mysql", 10, 20)

	var lang = new(schema.Language)
	var langList = lang.GetList()
	fmt.Println(len(langList))
	for i := 0; i < len(langList); i++ {
		var lang2 = langList[i]
		fmt.Println(lang2.String())
	}
	lang.Init("es-MX")
	lang.Init("es")
	fmt.Println(lang.String())

	var importFile = schema.Import{}
	importFile.Init(sourcetype.XmlDocument, "C:\\Temp\\Coding\\rpg_schema.xml")
	importFile.Load()
	importFile.Upgrade()

	_, _ = lang.IsCodeValid("FR")

	var metaSchema = schema.GetSchemaByName("@meta")
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

	rcd := new(data.Record)
	location, _ := time.LoadLocation("MST")
	ttt := time.Now().In(location)
	zone, offset := ttt.Zone()
	fmt.Println(offset)
	fmt.Println(zone)

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

	time.Sleep(time.Second * 10)
	runtime.GC()

	//time.Sleep(time.Second * 10)

	var lstTest = br.GetRecordList(0)
	fmt.Println("RECORD =======>")
	fmt.Println(lstTest.Count())
	fmt.Println(lstTest.ItemByIndex(0).String())
	//fmt.Println(lstTest.ItemByIndex(0).ToString())

	recordType := ".@meta2"
	var index = strings.Index(recordType, ".")
	fmt.Println(recordType[:index])
	fmt.Println(recordType[index+1:])
	rcd.SetRecordType("@meta")
	rcd.SetField("description", "758645454")
	rcd.SetField("value", 40.4)
	fmt.Println(rcd.GetField("description"))
	fmt.Println(rcd.GetField("reference_id"))

	fmt.Println("Successfully connected!")

	// FIELDS **********
	elemf := schema.Field{}
	elemf.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)

	fmt.Println(elemf.IsNumeric())
	var sss = elemf.GetSearchableValue("a", searchabletype.None)
	var sss2 = "A"
	fmt.Println(sss)
	fmt.Println(sss2)
	fmt.Println(elemf.GetSearchableValue("Français", searchabletype.None))
	fmt.Println(elemf.GetSearchableValue("", searchabletype.None))
	fmt.Println(elemf.GetSearchableValue("a", searchabletype.None))

	elemf2 := schema.Field{}
	elemf2.Init(4, "Zorba", "Zorba", fieldtype.Double, 5, "", true, false, false, true, true)

	elemf3 := schema.Field{}
	elemf3.Init(7, "Gga", "Gga", fieldtype.Double, 5, "", true, false, false, true, true)

	elemf4 := schema.Field{}
	elemf4.Init(7, "id", "id", fieldtype.DateTime, 5, "", true, false, false, true, true)

	elemf5 := schema.Field{}
	elemf5.Init(88, "id", "id", fieldtype.String, 5, "", true, false, false, true, true)

	// RELATIONS **********
	elemt := new(schema.Table)
	elemr := schema.Relation{}

	elemr.Init(21, "rel test", "hellkzae", elemt, relationtype.Mto, true, false, true, false)
	elemr0 := schema.Relation{}
	elemr0.Init(-23, "arel test", "hellkzae", elemt, relationtype.Mto, true, false, true, false)

}
