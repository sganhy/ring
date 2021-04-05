package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/tabletype"
	"strings"
	"testing"
)

func Test__Index__Init(t *testing.T) {
	var aarr = []string{"Gga", "Zorba"}
	elemi := Index{}
	// id int32, name string, description string, fields []string, tableId int32, bitmap bool, unique bool, baseline bool, actif bool
	elemi.Init(21, "rel test", "hellkzae", aarr, -17, false, true, true, true)

	if elemi.GetName() != "rel test" {
		t.Errorf("Index.Init() ==> name <> GetName()")
	}
	if elemi.GetId() != 21 {
		t.Errorf("Index.Init() ==> id <> GetId()")
	}
	if elemi.GetDescription() != "hellkzae" {
		t.Errorf("Index.Init() ==> description <> GetDescription()")
	}
	if elemi.GetFields() != nil {
		if strings.Join(elemi.GetFields(), ";") != strings.Join(aarr, ";") {
			t.Errorf("Index.Init() ==> fields <> %s", strings.Join(aarr, ";"))
		}
	} else {
		t.Errorf("Index.Init() ==> fields cannot be null")
	}
	if elemi.GetTableId() != -17 {
		t.Errorf("Index.Init() ==> tableId <> GetTableId()")
	}
	if elemi.IsBaseline() != true {
		t.Errorf("Index.Init() ==> IsBaseline() <> true")
	}
	if elemi.IsActive() != true {
		t.Errorf("Index.Init() ==> IsActive() <> true")
	}
	if elemi.IsUnique() != true {
		t.Errorf("Index.Init() ==> IsUnique() <> true")
	}
	if elemi.IsBitmap() != false {
		t.Errorf("Index.Init() ==> IsBitmap() <> false")
	}

	elemi.Init(21, "rel test", "hellkzae", nil, 0, false, true, true, true)
	if elemi.GetFields() == nil {
		t.Errorf("Index.Init() ==> fields cannot be null")
	}
}

// test mappers Meta to Index, and Index to Meta
func Test__Index__ToMeta(t *testing.T) {

	elemi0 := Index{}
	aarr := []string{"Gga", "Zorba", "testllk", "testllk22"}

	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemi0.Init(21, "rel test", "hellkzae", aarr, -11, false, true, false, true)

	meta := elemi0.ToMeta()
	elemi1 := meta.ToIndex()

	if elemi0.GetId() != elemi1.GetId() {
		t.Errorf("Index.ToMeta() ==> i0.GetId() must be equal to i1.GetId()")
	}
	if elemi0.GetName() != elemi1.GetName() {
		t.Errorf("Index.ToMeta() ==> i0.GetName() must be equal to i1.GetName()")
	}
	if elemi0.GetDescription() != elemi1.GetDescription() {
		t.Errorf("Index.ToMeta() ==> i0.GetDescription() must be equal to i1.GetDescription()")
	}
	if strings.Join(elemi0.GetFields(), "#") != strings.Join(elemi1.GetFields(), "#") {
		t.Errorf("Index.ToMeta() ==> i0.GetFields() must be equal to i1.GetFields()")
	}
	if elemi0.IsBitmap() != elemi1.IsBitmap() {
		t.Errorf("Index.ToMeta() ==> i0.IsBitmap() must be equal to i1.IsBitmap()")
	}
	if elemi0.IsUnique() != elemi1.IsUnique() {
		t.Errorf("Index.ToMeta() ==> i0.IsUnique() must be equal to i1.IsUnique()")
	}
	if elemi0.IsBaseline() != elemi1.IsBaseline() {
		t.Errorf("Index.ToMeta() ==> i0.IsBaseline() must be equal to i1.IsBaseline()")
	}
	if elemi0.IsActive() != elemi1.IsActive() {
		t.Errorf("Index.ToMeta() ==> i0.IsActive() must be equal to i1.IsActive()")
	}
	if elemi0.GetTableId() != elemi1.GetTableId() {
		t.Errorf("Index.ToMeta() ==> i0.GetTableId() must be equal to i1.GetTableId()")
	}
	// test fields
	if elemi1.fields == nil {
		t.Errorf("Index.ToMeta() ==> i1.fields cannot be nil")
	} else {
		// keep ";" hardcoded to detectec metaIndexSeparator constant change
		arr0str := strings.Join(elemi0.fields, ";")
		arr1str := strings.Join(elemi1.fields, ";")
		if arr0str != arr1str {
			t.Errorf("Index.ToMeta() ==> elemi0.fields is not equal to elemi1.fields")
		}
	}
}

func Test__Index__Clone(t *testing.T) {
	elemi0 := Index{}
	aarr := []string{"Gga", "Zorba", "testllk", "testllk22", "xxxxxx", "x44"}

	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemi0.Init(21, "rel test", "hellkzae", aarr, 77, false, true, false, true)
	elemi1 := elemi0.Clone()

	if elemi0.GetId() != elemi1.GetId() {
		t.Errorf("Index.Clone() ==> i0.GetId() must be equal to i1.GetId()")
	}
	if elemi0.GetName() != elemi1.GetName() {
		t.Errorf("Index.Clone() ==> i0.GetName() must be equal to i1.GetName()")
	}
	if elemi0.GetDescription() != elemi1.GetDescription() {
		t.Errorf("Index.Clone() ==> i0.GetDescription() must be equal to i1.GetDescription()")
	}
	if strings.Join(elemi1.GetFields(), "|") != strings.Join(aarr, "|") {
		t.Errorf("Index.Clone() ==> i0.GetFields() must be equal to i1.GetFields()")
	}
	if elemi0.IsBitmap() != elemi1.IsBitmap() {
		t.Errorf("Index.Clone() ==> i0.IsBitmap() must be equal to i1.IsBitmap()")
	}
	if elemi0.IsUnique() != elemi1.IsUnique() {
		t.Errorf("Index.Clone() ==> i0.IsUnique() must be equal to i1.IsUnique()")
	}
	if elemi0.IsBaseline() != elemi1.IsBaseline() {
		t.Errorf("Index.Clone() ==> i0.IsBaseline() must be equal to i1.IsBaseline()")
	}
	if elemi0.IsActive() != elemi1.IsActive() {
		t.Errorf("Index.Clone() ==> i0.IsActive() must be equal to i1.IsActive()")
	}
	if elemi0.IsActive() != elemi1.IsActive() {
		t.Errorf("Index.Clone() ==> i0.IsActive() must be equal to i1.IsActive()")
	}
}

func Test__Index__GetDdl(t *testing.T) {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)
	var uk Index = Index{}
	var tablespace = Tablespace{}

	tablespace.name = "rpg_index"

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
	objectType.Init(1019, "object_type", "", fieldtype.Byte, 0, "", true, true, true, false, true)
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
	uk.Init(1, "pk_@meta", "ATable Test", indexedFields, int32(tabletype.MetaId), false, false, true, true)

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
	table.Init(int32(tabletype.MetaId), "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "information_schema", tabletype.MetaId, databaseprovider.PostgreSql,
		"[subject]", true, false, true, false)

	var query = "CREATE INDEX idx_meta_001 ON information_schema.\"@meta\" USING btree (id,schema_id,object_type,reference_id)"
	// without tablespace
	if uk.GetDdl(ddlstatement.Create, table, nil) != query {
		t.Errorf("Index.GetDdl() ==> should be equal to %s", query)
	}

	query = "CREATE UNIQUE INDEX idx_meta_001 ON information_schema.\"@meta\" USING btree (id,schema_id,object_type,reference_id)"
	//if CREATE UNIQUE INDEX "pk_@meta" ON information_schema."@meta" (id,schema_id,,reference_id)
	uk.Init(1, "pk_@meta", "ATable Test", indexedFields, 55, false, true, true, true)
	if uk.GetDdl(ddlstatement.Create, table, nil) != query {
		t.Errorf("Index.GetDdl() ==> should be equal to %s", query)
	}

	// test with tablespace
	query = "CREATE UNIQUE INDEX idx_meta_001 ON information_schema.\"@meta\" USING btree (id,schema_id,object_type,reference_id) TABLESPACE rpg_index"
	if uk.GetDdl(ddlstatement.Create, table, &tablespace) != query {
		t.Errorf("Index.GetDdl() ==> should be equal to %s", query)
	}

	table.Init(1077, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "information_schema", tabletype.Business, databaseprovider.PostgreSql,
		"[subject]", true, false, true, false)

	if uk.GetPhysicalName(table) != "idx_1077_0001" {
		t.Errorf("Index.GetPhysicalName() ==> should be equal to idx_1077_0001")
	}

	table.Init(1077, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "information_schema", tabletype.Mtm, databaseprovider.PostgreSql,
		"[subject]", true, false, true, false)

	if uk.GetPhysicalName(table) != "idx_mtm_0001" {
		t.Errorf("Index.GetPhysicalName() ==> should be equal to idx_mtm_0001")
	}

	// test DROP
	query = "DROP INDEX information_schema.idx_mtm_0001"
	if uk.GetDdl(ddlstatement.Drop, table, nil) != query {
		t.Errorf("Index.GetDdl() ==> should be equal to %s", query)
	}

	// no alter on indexes !!
	if uk.GetDdl(ddlstatement.Alter, table, nil) != "" {
		t.Errorf("Index.GetDdl() ==> should be equal to null")
	}

	// unknown db provider
	table.Init(1077, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "information_schema", tabletype.Business, databaseprovider.NotDefined,
		"[subject]", true, false, true, false)

	if uk.GetDdl(ddlstatement.Create, table, nil) != "" {
		t.Errorf("Index.GetDdl() ==> should be equal to null")
	}

	if uk.GetDdl(ddlstatement.Drop, table, nil) != "" {
		t.Errorf("Index.GetDdl() ==> should be equal to null")
	}

}
