package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/tabletype"
	"strings"
	"testing"
)

// INIT
func Test__Relation__Init(t *testing.T) {
	table := new(Table)
	elemr0 := Relation{}
	elemt := table.getMetaTable(databaseprovider.PostgreSql, metaSchemaName)
	elemr0.Init(-23, "arel test", "hellkzae", "hell1", "52", elemt, relationtype.Mto, false, true, false)

	if elemr0.GetName() != "arel test" {
		t.Errorf("Relation.Init() ==> name <> GetName()")
	}
	if elemr0.GetId() != -23 {
		t.Errorf("Relation.Init() ==> id <> GetId()")
	}
	if elemr0.GetDescription() != "hellkzae" {
		t.Errorf("Relation.Init() ==> description <> GetDescription()")
	}
	if elemr0.GetInverseRelationName() != "hell1" {
		t.Errorf("Relationeld.Init() ==> inverseRelationNam <> GetInverseRelationName()")
	}
	if elemr0.GetMtmTableName() != "52" {
		t.Errorf("Relationeld.Init() ==> GetMtmTable() <> mtm table")
	}
	if elemr0.GetType() != relationtype.Mto {
		t.Errorf("Relationeld.Init() ==> type <> GetType()")
	}
	if elemr0.IsNotNull() != false {
		t.Errorf("Relationeld.Init() ==> IsNotNull() <> false")
	}
	if elemr0.IsBaseline() != true {
		t.Errorf("Relationeld.Init() ==> IsBaseline() <> true")
	}
	if elemr0.IsActive() != false {
		t.Errorf("Relationeld.Init() ==> IsActive() <> false")
	}
	if elemr0.GetToTable() != elemt {
		t.Errorf("Relationeld.Init() ==> GetToTable() <> table pointer")
	}
}

//test mappers Meta to Relation, and Relation to Meta
func Test__Relation__ToMeta(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}

	elemf := Field{}
	elemf.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)

	//var prim = schema.GetDefaultPrimaryKey()
	fields = append(fields, elemf)
	elemt := Table{}
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.NotDefined, "subject test",
		true, false, true, false)

	elemr0 := Relation{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemr0.Init(23, "rel test", "hellkzae", "hell1", "52", &elemt, relationtype.Otop, false, true, false)

	meta := elemr0.ToMeta(777)
	elemr1 := meta.ToRelation(&elemt)

	if elemr0.GetId() != elemr1.GetId() {
		t.Errorf("Relation.ToMeta() ==> r0.GetId() must be equal to r1.GetId()")
	}
	if elemr0.GetName() != elemr1.GetName() {
		t.Errorf("Relation.ToMeta() ==> r0.GetName() must be equal to r1.GetName()")
	}
	if elemr0.GetDescription() != elemr1.GetDescription() {
		t.Errorf("Relation.ToMeta() ==> r0.GetDescription() must be equal to r1.GetDescription()")
	}
	if elemr0.GetInverseRelationName() != elemr1.GetInverseRelationName() {
		t.Errorf("Relation.ToMeta() ==> r0.GetInverseRelationName() must be equal to r1.GetInverseRelationName()")
	}
	if elemr0.GetType() != elemr1.GetType() {
		t.Errorf("Relation.ToMeta() ==> r0.GetType() must be equal to r1.GetType()")
	}
	// check reference of table must be the same
	if elemr0.GetToTable() != elemr1.GetToTable() {
		t.Errorf("Relation.ToMeta() ==> r0.GetToTable() reference must be equal to r1.GetToTable()")
	}
	if elemr0.IsBaseline() != elemr1.IsBaseline() {
		t.Errorf("Relation.ToMeta() ==> r0.IsBaseline() must be equal to r1.IsBaseline()")
	}
	if elemr0.IsNotNull() != elemr1.IsNotNull() {
		t.Errorf("Relation.ToMeta() ==> r0.IsNotNull() must be equal to r1.IsNotNull()")
	}
	if elemr0.IsActive() != elemr1.IsActive() {
		t.Errorf("Relation.ToMeta() ==> r0.IsActive() must be equal to r1.IsActive()")
	}
	// test GetDdlSql

}

func Test__Relation__GetDdl(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}

	elemf := Field{}
	elemf.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)

	//var prim = schema.GetDefaultPrimaryKey()
	fields = append(fields, elemf)
	elemt := Table{}
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.NotDefined, "subject test",
		true, false, true, false)

	elemr0 := Relation{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemr0.Init(23, "rel test", "hellkzae", "hell1", "52", &elemt, relationtype.Otop, false, true, false)

	var sql = elemr0.GetDdl(databaseprovider.PostgreSql)
	if strings.ToUpper(sql) != "REL TEST INT8" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to REL TEST INT8")
	}

	elemr0.toTable = nil
	sql = elemr0.GetDdl(databaseprovider.PostgreSql)
	if sql != "" {
		t.Errorf("Field.GetSql() ==> (1) sql should be null")
	}
}

func Test__Relation__Clone(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}

	elemf := Field{}
	elemf.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)

	//var prim = schema.GetDefaultPrimaryKey()
	fields = append(fields, elemf)
	elemt := Table{}
	elemt.Init(22, "rel test", "hellkzae", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.NotDefined, "subject test",
		true, false, true, false)

	elemr0 := Relation{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemr0.Init(23, "rel test", "hellkzae", "hell1", "52", &elemt, relationtype.Otop, false, true, false)
	elemr1 := elemr0.Clone()

	if elemr0.GetId() != elemr1.GetId() {
		t.Errorf("Relation.Clone() ==> r0.GetId() must be equal to r1.GetId()")
	}
	if elemr0.GetName() != elemr1.GetName() {
		t.Errorf("Relation.Clone() ==> r0.GetName() must be equal to r1.GetName()")
	}
	if elemr0.GetDescription() != elemr1.GetDescription() {
		t.Errorf("Relation.Clone() ==> r0.GetDescription() must be equal to r1.GetDescription()")
	}
	if elemr0.GetType() != elemr1.GetType() {
		t.Errorf("Relation.Clone() ==> r0.GetType() must be equal to r1.GetType()")
	}
	if elemr0.notNull != elemr1.notNull {
		t.Errorf("Relation.Clone() ==> r0.notNull must be equal to r1.notNull")
	}
	if elemr0.IsBaseline() != elemr1.IsBaseline() {
		t.Errorf("Relation.Clone() ==> r0.IsBaseline() must be equal to r1.IsBaseline()")
	}
}
