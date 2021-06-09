package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
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
	elemr0.Init(-23, "arel test", "hellkzae", "hell1", elemt, relationtype.Mto, false, true, false)

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
	if elemr0.GetPhysicalName() != "arel test" {
		t.Errorf("Relationeld.Init() ==> GetPhysicalName() <> 'arel test'")
	}
	if elemr0.GetEntityType() != entitytype.Relation {
		t.Errorf("Relationeld.Init() ==> GetEntityType() <> entitytype.Relation")
	}
}

//test mappers Meta to Relation, and Relation to Meta
func Test__Relation__toMeta(t *testing.T) {
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
	elemr0.Init(23, "rel test", "hellkzae", "hell1", &elemt, relationtype.Otop, false, true, false)

	metaData := elemr0.toMeta(777)
	elemr1 := metaData.toRelation(&elemt)

	if elemr0.GetId() != elemr1.GetId() {
		t.Errorf("Relation.toMeta() ==> r0.GetId() must be equal to r1.GetId()")
	}
	if elemr0.GetName() != elemr1.GetName() {
		t.Errorf("Relation.toMeta() ==> r0.GetName() must be equal to r1.GetName()")
	}
	if elemr0.GetDescription() != elemr1.GetDescription() {
		t.Errorf("Relation.toMeta() ==> r0.GetDescription() must be equal to r1.GetDescription()")
	}
	if elemr0.GetInverseRelationName() != elemr1.GetInverseRelationName() {
		t.Errorf("Relation.toMeta() ==> r0.GetInverseRelationName() must be equal to r1.GetInverseRelationName()")
	}
	if elemr0.GetType() != elemr1.GetType() {
		t.Errorf("Relation.toMeta() ==> r0.GetType() must be equal to r1.GetType()")
	}
	// check reference of table must be the same
	if elemr0.GetToTable() != elemr1.GetToTable() {
		t.Errorf("Relation.toMeta() ==> r0.GetToTable() reference must be equal to r1.GetToTable()")
	}
	if elemr0.IsBaseline() != elemr1.IsBaseline() {
		t.Errorf("Relation.toMeta() ==> r0.IsBaseline() must be equal to r1.IsBaseline()")
	}
	if elemr0.IsNotNull() != elemr1.IsNotNull() {
		t.Errorf("Relation.toMeta() ==> r0.IsNotNull() must be equal to r1.IsNotNull()")
	}
	if elemr0.IsActive() != elemr1.IsActive() {
		t.Errorf("Relation.toMeta() ==> r0.IsActive() must be equal to r1.IsActive()")
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
	elemr0.Init(23, "rel test", "hellkzae", "hell1", &elemt, relationtype.Otop, false, true, false)

	var sql = elemr0.GetDdl(databaseprovider.PostgreSql)
	if strings.ToUpper(sql) != "REL TEST INT8" {
		t.Errorf("Field.GetSql() ==> (1) sql should be equal to REL TEST INT8")
	}

	elemr0.setToTable(nil)
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
	elemr0.Init(23, "rel test", "hellkzae", "hell1", &elemt, relationtype.Otop, false, true, false)
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
	if elemr0.IsNotNull() != elemr1.IsNotNull() {
		t.Errorf("Relation.Clone() ==> r0.IsNotNull() must be equal to r1.notNull")
	}
	if elemr0.IsBaseline() != elemr1.IsBaseline() {
		t.Errorf("Relation.Clone() ==> r0.IsBaseline() must be equal to r1.IsBaseline()")
	}
}

func Test__Relation__getMtmName(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}

	elemr0 := Relation{}
	elemr0.Init(23, "test", "hellkzae", "test inv", nil, relationtype.Mtm, false, true, false)

	//*****
	//***** REFLEXIVE MTM RELATIONS (fromTableId == toTableId)
	//*****
	elemr1 := Relation{}
	elemr1.Init(24, "test inv", "hellkzae", "test", nil, relationtype.Mtm, false, true, false)

	relations = append(relations, elemr0)
	relations = append(relations, elemr1)

	elemt01 := Table{}
	elemt01.Init(22, "rel test", "hellkzae", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.NotDefined, "subject test",
		true, false, true, false)

	elemr0.setToTable(&elemt01)
	elemr1.setToTable(&elemt01)

	if elemr0.getMtmName(22) != elemr1.getMtmName(22) {
		t.Errorf("Relation.getMtmName() ==>	r1 should be equal to r0")
	}

	if elemr0.getMtmName(22) != "@mtm_00022_00022_023" {
		t.Errorf("Relation.getMtmName() ==>	r0 should be equal to '@mtm_00022_00022_023'")
	}

	//*****
	//***** MTM RELATIONS (fromTableId > toTableId || fromTableId < toTableId)
	//*****
	// fromTableId > toTableId
	//TABLE 1
	relations = make([]Relation, 1, 1)
	elemr3 := Relation{}
	elemr3.Init(25, "test2", "[description]", "test2 inv", nil, relationtype.Mtm, false, true, false)
	relations[0] = elemr3
	elemt01.Init(22, "rel test", "[description]", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.NotDefined, "subject test",
		true, false, true, false)

	//TABLE 2
	relations = make([]Relation, 1, 1)
	elemr4 := Relation{}
	elemr4.Init(24, "test2 inv", "[description]", "test2", nil, relationtype.Mtm, false, true, false)
	relations[0] = elemr4
	elemt02 := Table{}
	elemt02.Init(23, "rel test33", "[description]", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.NotDefined, "subject test",
		true, false, true, false)

	elemt01.relations[0].setToTable(&elemt02)
	elemt02.relations[0].setToTable(&elemt01)

	elemt01.relations[0].GetInverseRelation()
	if elemt01.relations[0].getMtmName(22) != "@mtm_00022_00023_025" {
		t.Errorf("Relation.getMtmName() ==>	r3 should be equal to '@mtm_00022_00023_025'")
	}
	if elemt01.relations[0].getMtmName(22) != elemt02.relations[0].getMtmName(23) {
		t.Errorf("Relation.getMtmName() ==>	r3 should be equal to r4")
	}
	// fromTableId < toTableId
	elemt02.id = 21
	if elemt01.relations[0].getMtmName(22) != "@mtm_00021_00022_024" {
		t.Errorf("Relation.getMtmName() ==>	r3 should be equal to '@mtm_00021_00022_024'")
	}
	if elemt01.relations[0].getMtmName(22) != elemt02.relations[0].getMtmName(21) {
		t.Errorf("Relation.getMtmName() ==>	r3 should be equal to r4")
	}
}
