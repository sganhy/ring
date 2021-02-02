package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"ring/schema"
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/tabletype"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//**************************
// configuration methods
//**************************
func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "sa", "postgres")

	Test__Table__GetFieldByName()

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connection to database")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// FIELDS **********
	elemf0 := schema.Field{}
	elemf0.Init(11, "aField Test", "AField Test", fieldtype.DateTime, 5, "test default", true, false, false, true, true)

	meta := elemf0.ToMeta(12)
	elemf00 := meta.ToField()
	fmt.Println(meta.GetEntityType())
	fmt.Println(meta.GetFlags())
	fmt.Println(elemf0.IsNotNull())
	fmt.Println(meta.IsFieldNotNull())
	fmt.Println(elemf00.GetSize())

	elemf := schema.Field{}
	elemf.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)
	var lang = schema.NewLanguage("Fr")
	fmt.Println(meta.IsEntityBaseline())
	fmt.Println(elemf00.GetDefaultValue())
	fmt.Println(elemf.IsNumeric())
	fmt.Println(elemf.GetSearchableValue("žůžo", *lang))
	fmt.Println(elemf.GetSearchableValue("zuzo", *lang))
	fmt.Println(elemf.GetSearchableValue("Français", *lang))
	fmt.Println(elemf.GetSearchableValue("", *lang))
	fmt.Println(elemf.GetSearchableValue("a", *lang))

	elemf2 := schema.Field{}
	elemf2.Init(4, "Zorba", "Zorba", fieldtype.Double, 5, "", true, false, false, true, true)

	elemf3 := schema.Field{}
	elemf3.Init(7, "Gga", "Gga", fieldtype.Double, 5, "", true, false, false, true, true)

	elemf4 := schema.Field{}
	elemf4.Init(7, "id", "id", fieldtype.DateTime, 5, "", true, false, false, true, true)

	elemf5 := schema.Field{}
	elemf5.Init(88, "id", "id", fieldtype.DateTime, 5, "", true, false, false, true, true)

	// RELATIONS **********
	elemt := new(schema.Table)
	elemr := schema.Relation{}
	elemr.Init(21, "rel test", "hellkzae", "hell1", "52", nil, relationtype.Mto, false, true, false)
	elemr0 := schema.Relation{}
	elemr0.Init(-23, "arel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)

	var aarr = []string{"Gga", "Zorba"}
	elemi := schema.Index{}
	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool, actif bool
	elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
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
	fmt.Println(elemt.GetFieldByName("Field Test").GetName())
	//fmt.Println(elemt.GetFieldById(4).GetName())
	//fmt.Println(elemt.GetPrimaryKey().GetName())
	/*
		var elemt2 = *schema.GetMetaTable()
		fmt.Println(elemt.GetPhysicalName())
		fmt.Println(elemt2.GetPhysicalName())
	*/

	var meta1 = schema.GetMetaTable(databaseprovider.PostgreSql)
	fmt.Println(lang.GetNativeName())
	fmt.Println(meta1.GetName())
	fmt.Println(meta1.GetPhysicalName())
}

func Test__Table__GetFieldByName() {
	var fields = []schema.Field{}
	var relations = []schema.Relation{}
	var indexes = []schema.Index{}
	var table = new(schema.Table)
	const FIELD_COUNT = 10000

	// added invalid fields
	for i := -100; i <= FIELD_COUNT; i++ {
		field := new(schema.Field)
		nameLenght := (abs(i) % 30) + 2
		field.Init(int32(i), randStringBytes(nameLenght), "", fieldtype.Int, 0, "", true, true, true, false, true)
		fields = append(fields, *field)
	}

	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes, "schema.@meta",
		physicaltype.Table, -111, tabletype.Business, "", true, false, true, true)

	for i := len(fields) - 1; i > 0; i-- {
		fieldName := fields[i].GetName()
		field := table.GetFieldByName(fieldName)
		// test valid field only
		if fields[i].IsValid() == true {
			if field == nil {
				fmt.Println("Table.GetFieldByName() ==> fields[i].name; i=%d, name=%s", i, fieldName)
				break
			}
		}
	}

	table.ToFile("c:\\temp\\data.txt")

}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
