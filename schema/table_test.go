package schema

import (
	"math/rand"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"strings"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Test__Table__Init(t *testing.T) {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)
	var uk Index = Index{}

	// physical_name is builded later
	//  == metaId table
	var id Field = Field{}
	var schemaId Field = Field{}
	var objectType Field = Field{}
	var referenceId Field = Field{}
	var dataType Field = Field{}

	var flags Field = Field{}
	var value Field = Field{}
	var name Field = Field{}
	var description Field = Field{}
	var active Field = Field{}

	// elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true, true)
	// !!!! id field must be greater than 0 !!!!
	id.Init(1009, "id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	schemaId.Init(1013, "schema_id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	objectType.Init(1019, "", "", fieldtype.Byte, 0, "", true, true, true, false, true)
	referenceId.Init(1021, "reference_id", "", fieldtype.Int, 0, "", true, true, true, false, true)
	dataType.Init(1031, "data_type", "", fieldtype.Int, 0, "", true, false, true, false, true)

	flags.Init(1039, "flags", "", fieldtype.Long, 0, "", true, false, true, false, true)
	name.Init(1061, "name", "", fieldtype.String, 30, "", true, false, true, false, true)
	description.Init(1069, "description", "", fieldtype.String, 0, "", true, false, true, false, true)
	value.Init(1087, metaValue, "", fieldtype.String, 0, "", true, false, true, false, true)
	active.Init(1093, "active", "", fieldtype.Boolean, 0, "", true, false, true, false, true)

	// elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
	// unique key (1)      id; schema_id; reference_id; object_type
	var indexedFields = []string{id.name, schemaId.name, objectType.name, referenceId.name}
	uk.Init(1, "pk_@meta", "ATable Test", indexedFields, int32(tabletype.MetaId), false, true, true, true)

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
	table.Init(int32(tabletype.MetaId), "@meta", "ATable Test", fields, relations, indexes, "schema.@meta",
		physicaltype.Table, -111, tabletype.MetaId, "[subject]", true, false, true, false)

	if table.GetName() != "@meta" {
		t.Errorf("Table.Init() ==> name <> GetName()")
	}
	if table.GetId() != int32(tabletype.MetaId) {
		t.Errorf("Table.Init() ==> id <> GetId()")
	}
	if table.GetDescription() != "ATable Test" {
		t.Errorf("Table.Init() ==> description <> GetDescription()")
	}
	if table.GetPhysicalName() != "schema.@meta" {
		t.Errorf("Table.Init() ==> physical <> GetPhysicalName()")
	}
	if table.GetType() != tabletype.MetaId {
		t.Errorf("Table.Init() ==> description <> GetDescription()")
	}
	if table.GetPhysicalType() != physicaltype.Table {
		t.Errorf("Table.Init() ==> physicaltype <> GetPhysicalType()")
	}
	if table.GetSubject() != "[subject]" {
		t.Errorf("Table.Init() ==> subject <> GetSubject()")
	}
	if table.GetSchemaId() != -111 {
		t.Errorf("Table.Init() ==> GetSchemaId() <> -111")
	}
	if table.IsCached() != true {
		t.Errorf("Table.Init() ==> IsCached() <> true")
	}
	if table.IsReadonly() != false {
		t.Errorf("Table.Init() ==> IsReadonly() <> false")
	}
	if table.IsBaseline() != true {
		t.Errorf("Table.Init() ==> IsBaseline() <> true")
	}
	if table.IsActive() != false {
		t.Errorf("Table.Init() ==> IsActive() <> false")
	}

}

func Test__Table__GetFieldByName(t *testing.T) {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)
	const FIELD_COUNT = 20000

	// added invalid fields (101)
	for i := -100; i <= FIELD_COUNT; i++ {
		field := new(Field)
		nameLenght := (abs(i) % 30) + 2
		// fixture
		fieldName := randStringBytes(nameLenght)
		// to force adding new primary key
		if strings.ToUpper(fieldName) != "ID" {
			field.Init(int32(i), fieldName, "", fieldtype.Int, 0, "", true, true, true, false, true)
			fields = append(fields, *field)
		}
	}

	//t.Errorf("fields.Count ==> %d ", len(fields))
	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes, "schema.@meta",
		physicaltype.Table, -111, tabletype.Business, "", true, false, true, true)

	//t.Errorf("Table.fields.Count ==> %d /%d", len(table.fields), cap(table.fields))
	//t.Errorf("Table.fieldsById.Count ==> %d /%d", len(table.fieldsById), cap(table.fieldsById))

	for i := 0; i < len(fields); i++ {
		// test valid field only
		if fields[i].IsValid() == true {
			fieldName := fields[i].name
			field := table.GetFieldByName(fieldName)
			if field == nil {
				t.Errorf("Table.GetFieldByName() ==> fields[i].name; i=%d, name=%s, id=%d", i, fieldName, fields[i].id)
				break
			} else {
				if fieldName != field.GetName() {
					t.Errorf("Table.GetFieldByName() ==> fields[i].name; i=%d, name=%s, found=%s", i, fieldName, table.GetFieldByName(fields[i].name).GetName())
					break
				}
			}
		}
	}

	//find primary key
	field := table.GetFieldByName("id")
	if field == nil {
		t.Errorf("Table.GetFieldByName() ==> Cannot find primary key")
	} else if field.IsPrimaryKey() == false {
		t.Errorf("Table.GetFieldByName() ==> Cannot find primary key reference")
	}
	field = table.GetPrimaryKey()
	if field == nil {
		t.Errorf("Table.GetFieldByName() ==> Cannot find primary key")
	} else if field.IsPrimaryKey() == false {
		t.Errorf("Table.GetFieldByName() ==> Cannot find primary key reference")
	}

	field = table.GetFieldByName("111")
	if field != nil {
		t.Errorf("Table.GetFieldByName() ==> field '111' cannot be found")
	}
	//Testing GetFieldByNameI
	field = table.GetFieldByNameI(strings.ToLower(table.fieldsById[3].name))
	if field == nil {
		t.Errorf("Table.GetFieldByNameI() ==> ToLower(fields[2].name); i=%d, name=%s", 2, strings.ToLower(fields[2].name))
	}
	field = table.GetFieldByNameI(strings.ToUpper(table.fieldsById[3].name))
	if field == nil {
		t.Errorf("Table.GetFieldByNameI() ==> ToUpper(fields[2].name); i=%d, name=%s", 2, strings.ToLower(fields[2].name))
	}
	field = table.GetFieldByNameI("111")
	if field != nil {
		t.Errorf("Table.GetFieldByNameI() ==> field '111' cannot be found")
	}
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
