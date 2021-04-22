package main

import (
	"fmt"
	"ring/data"
	"ring/data/operationtype"
	"ring/data/sortordertype"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/searchabletype"
	"ring/schema/sourcetype"
	"ring/schema/tabletype"
	"runtime"
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
	runtime.LockOSThread()

	rcd := new(data.Record)
	frame := test4()
	fmt.Println(frame.File)
	schema.Init(databaseprovider.PostgreSql, "host=localhost port=5432 user=postgres password=sa dbname=postgres sslmode=disable", 10, 20)
	//schema.Init(databaseprovider.MySql, "root:root@/rpg_sheet", 10, 20)

	//ss.LogWarn(1, 544, "hello", "World")

	location, _ := time.LoadLocation("MST")
	ttt := time.Now().In(location)
	zone, offset := ttt.Zone()
	fmt.Println(offset)
	fmt.Println(zone)

	fmt.Println(test3())

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

	var lang = schema.Language{}
	lang.Init("FR")
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

	elemr.Init(21, "rel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)
	elemr0 := schema.Relation{}
	elemr0.Init(-23, "arel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)

	var aarr = []string{"Gga", "Zorba"}
	elemi := schema.Index{}
	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool, actif bool
	elemi.Init(21, "rel test", "hellkzae", aarr, 55, false, true, true, true)
	fmt.Println(elemr.GetName())
	var fields = []schema.Field{}
	//var prim = schema.GetDefaultPrimaryKey()
	fields = append(fields, elemf)
	fields = append(fields, elemf2)
	fields = append(fields, elemf3)
	fields = append(fields, elemf4)
	fields = append(fields, elemf5)
	//aarr2 = append(aarr2, *prim)
	var relations = []schema.Relation{}
	relations = append(relations, elemr)
	relations = append(relations, elemr0)

	var indexes = []schema.Index{}
	indexes = append(indexes, elemi)

	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool,  bool
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes,
		physicaltype.Table, 64, "@meta", tabletype.Business, databaseprovider.PostgreSql, "subject test", true, false, true, false)
	fmt.Println(elemt.GetFieldByName("Field Test").GetName())
	fmt.Println(elemt.GetDdl(ddlstatement.Create, nil))
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

func test2() string {
	return getFrame(1).Function
}

func test3() int {
	return getFrame(1).Line
}

func test4() runtime.Frame {
	return getFrame(1)
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}
	return frame
}
