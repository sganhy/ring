package main

import (
	"fmt"
	"ring/data"
	"ring/data/operationtype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/sourcetype"
	"ring/schema/tabletype"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const minint32 string = "-2147483648"

//**************************
// configuration methods
//**************************
func main() {
	rcd := new(data.Record)
	schema.Init(databaseprovider.PostgreSql, "host=localhost port=5432 user=postgres password=sa dbname=postgres sslmode=disable", 10, 20)

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

	var importFile = schema.Import{}
	importFile.Init(sourcetype.XmlDocument, "C:\\Temp\\Coding\\rpg_schema.xml")
	//importFile.Init(sourcetype.XmlDocument, "C:\\Temp\\schema.xml")
	//importFile.Load()

	// Create an empty user and make the sql query (using $1 for the parameter)
	var br = new(data.BulkRetrieve)
	br.SimpleQuery(0, "@meta")
	br.AppendFilter(0, "schema_id", operationtype.Equal, 0)
	br.SimpleQuery(1, "@log")
	br.AppendFilter(1, "schema_id", operationtype.NotEqual, 2)
	br.AppendFilter(1, "schema_id", operationtype.NotEqual, 3)
	br.AppendFilter(1, "schema_id", operationtype.NotEqual, 4)
	br.RetrieveRecords()

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
	elemf0 := schema.Field{}
	elemf0.Init(11, "aField Test", "AField Test", fieldtype.Long, 5, "test default", true, false, false, true, true)

	meta := elemf0.ToMeta(12)
	elemf00 := meta.ToField()
	fmt.Println(elemf0.IsNotNull())
	fmt.Println(meta.IsFieldNotNull())
	fmt.Println(elemf00.GetSize())

	elemf := schema.Field{}
	elemf.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)

	var lang = schema.Language{}
	lang.Init("FR")
	fmt.Println(meta.IsEntityBaseline())
	fmt.Println(elemf00.GetDefaultValue())
	fmt.Println(elemf.IsNumeric())
	fmt.Println(elemf.GetSearchableValue("žůžo", lang))
	fmt.Println(elemf.GetSearchableValue("zuzo", lang))
	fmt.Println(elemf.GetSearchableValue("Français", lang))
	fmt.Println(elemf.GetSearchableValue("", lang))
	fmt.Println(elemf.GetSearchableValue("a", lang))

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

	elemr.Init(21, "rel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)
	elemr0 := schema.Relation{}
	elemr0.Init(-23, "arel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)

	var aarr = []string{"Gga", "Zorba"}
	elemi := schema.Index{}
	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool, actif bool
	elemi.Init(21, "rel test", "hellkzae", aarr, false, true, true, true)
	fmt.Println(elemr.GetName())
	var fields = []schema.Field{}
	//var prim = schema.GetDefaultPrimaryKey()
	fields = append(fields, elemf)
	fields = append(fields, elemf2)
	fields = append(fields, elemf3)
	fields = append(fields, elemf4)
	fields = append(fields, elemf5)
	fields = append(fields, elemf0)
	//aarr2 = append(aarr2, *prim)
	var relations = []schema.Relation{}
	relations = append(relations, elemr)
	relations = append(relations, elemr0)

	var indexes = []schema.Index{}
	indexes = append(indexes, elemi)

	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool,  bool
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes,
		"schema.t_site", physicaltype.Table, 64, tabletype.Business, "subject test", true, false, true, false)
	sql2, _ := elemt.GetDdlSql(databaseprovider.PostgreSql, nil)
	fmt.Println(sql2)
	fmt.Println(elemt.GetFieldByName("Field Test").GetName())
	//fmt.Println(elemt.GetFieldById(4).GetName())
	//fmt.Println(elemt.GetPrimaryKey().GetName())
	/*
		var elemt2 = *schema.GetMetaTable()
		fmt.Println(elemt.GetPhysicalName())
		fmt.Println(elemt2.GetPhysicalName())
	*/

	reg := []string{"a", "b", "c"}
	fmt.Println(strings.Join(reg, ","))
	//time.Sleep(20 * time.Second)
	fmt.Println("Finished!")
}
