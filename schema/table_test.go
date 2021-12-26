package schema

import (
	"math/rand"
	"ring/schema/databaseprovider"
	"ring/schema/ddlstatement"
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/tabletype"
	"strconv"
	"strings"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

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
	var indexedFields = []string{id.GetName(), schemaId.GetName(), objectType.GetName(), referenceId.GetName()}
	uk.Init(1, "pk_@meta", "ATable Test", indexedFields, false, true, true, true)

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

	if table.GetName() != "@meta" {
		t.Errorf("Table.Init() ==> name <> GetName()")
	}
	if table.GetId() != int32(tabletype.MetaId) {
		t.Errorf("Table.Init() ==> id <> GetId()")
	}
	if table.GetDescription() != "ATable Test" {
		t.Errorf("Table.Init() ==> description <> GetDescription()")
	}
	if table.GetPhysicalName() != "information_schema.\"@meta\"" {
		t.Errorf("Table.Init() ==> physical name <> information_schema.\"@meta\"")
	}
	if table.GetType() != tabletype.MetaId {
		t.Errorf("Table.Init() ==> GetType <> tabletype.MetaId")
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
	if table.GetDatabaseProvider() != databaseprovider.PostgreSql {
		t.Errorf("Table.Init() ==> GetDatabaseProvider() <> databaseprovider.Oracle")
	}
	if table.getCacheId() != nil {
		t.Errorf("Table.Init() ==> GetCacheId() cannot be null")
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
	if table.GetEntityType() != entitytype.Table {
		t.Errorf("Table.Init() ==>  GetEntityType() <> entitytype.Table")
	}
	// should be equal to true always
	if table.logStatement(ddlstatement.Create) != true {
		t.Errorf("Table.Init() ==>  logStatment(Create) <> true")
	}
	if table.logStatement(ddlstatement.Drop) != true {
		t.Errorf("Table.Init() ==>  logStatment(Drop) <> true")
	}
	if table.logStatement(ddlstatement.Alter) != true {
		t.Errorf("Table.Init() ==>  logStatment(Alter) <> true")
	}
	if table.logStatement(ddlstatement.Truncate) != true {
		t.Errorf("Table.Init() ==>  logStatment(Alter) <> true")
	}
	if table.GetFieldIdByIndex(1).GetName() != "schema_id" {
		t.Errorf("Table.Init() ==>  GetFieldIdByIndex(1) <> 'schema_id'")
	}
	if table.GetFieldByIndex(2).GetName() != "data_type" {
		t.Errorf("Table.Init() ==>  GetFieldByIndex(2) <> 'data_type'")
	}

}

//test: GetFieldByName, GetFieldByNameI, GetFieldById, and GetFieldIndexByName
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
	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, "@meta", tabletype.Business, databaseprovider.Undefined, "", true, false, true, true)

	//t.Errorf("Table.fields.Count ==> %d /%d", len(table.fields), cap(table.fields))
	//t.Errorf("Table.fieldsById.Count ==> %d /%d", len(table.fieldsById), cap(table.fieldsById))
	for i := 0; i < len(fields); i++ {
		// test valid field only
		if fields[i].IsValid() == true {
			fieldName := fields[i].GetName()
			field := table.GetFieldByName(fieldName)
			if field == nil {
				t.Errorf("Table.GetFieldByName() ==> fields[i].name; i=%d, name=%s, id=%d", i, fieldName, fields[i].GetId())
				break
			} else {
				if fieldName != field.GetName() {
					t.Errorf("Table.GetFieldByName() ==> fields[i].name; i=%d, name=%s, found=%s", i, fieldName, table.GetFieldByName(fields[i].GetName()).GetName())
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
	field = table.GetFieldByNameI(strings.ToLower(table.fields[table.mapper[3]].GetName()))
	if field == nil {
		t.Errorf("Table.GetFieldByNameI() ==> ToLower(fields[2].name); i=%d, name=%s", 2, strings.ToLower(fields[2].GetName()))
	}
	field = table.GetFieldByNameI(strings.ToUpper(table.fields[table.mapper[3]].GetName()))
	if field == nil {
		t.Errorf("Table.GetFieldByNameI() ==> ToUpper(fields[2].name); i=%d, name=%s", 2, strings.ToLower(fields[2].GetName()))
	}
	field = table.GetFieldByNameI("111")
	if field != nil {
		t.Errorf("Table.GetFieldByNameI() ==> field '111' cannot be found!!")
	}

	//Testing GetFieldIndexByName
	position := table.GetFieldIndexByName("111")
	if position != fieldNotFound {
		t.Errorf("Table.GetFieldIndexByName() ==> field '111' index cannot be found!!")
	}
	table = table.getMetaTable(databaseprovider.MySql, metaSchemaName)
	position = table.GetFieldIndexByName("description")
	if position != 2 {
		t.Errorf("Table.GetFieldIndexByName() ==> field '%s' index should be equal to 2", field.GetName())
	}

}

//test: GetRelationByName, GetRelationIndexByName, and GetPrimaryKey
func Test__Table__GetRelationByName(t *testing.T) {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)
	const RELATION_COUNT = 20000

	for i := -100; i <= RELATION_COUNT; i++ {
		relation := new(Relation)
		nameLenght := (abs(i) % 30) + 2
		// fixture
		relationName := randStringBytes(nameLenght)
		relation.Init(int32(i), relationName, "11", nil, relationtype.Mtm, false, false, true, false)
		relations = append(relations, *relation)
	}

	//t.Errorf("fields.Count ==> %d ", len(fields))
	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, metaSchemaName, tabletype.Mtm, databaseprovider.Undefined, "", true, false, true, true)

	for i := 0; i < len(relations); i++ {
		// test valid field only
		relationName := relations[i].GetName()
		relation := table.GetRelationByName(relationName)
		if relation == nil {
			t.Errorf("Table.GetRelationByName() ==> relations[i].name; i=%d, name=%s, id=%d", i, relationName, relations[i].GetId())
			break
		} else {
			if relationName != relation.GetName() {
				t.Errorf("Table.GetRelationByName() ==> relations[i].name; i=%d, name=%s, found=%s", i, relationName, table.GetRelationByName(fields[i].GetName()).GetName())
				break
			}
		}
	}
	// test nil
	relation := table.GetRelationByName("22222")
	if relation != nil {
		t.Errorf("Table.GetRelationByName() ==> relation '22222' cannot be found!!")
	}

	position := table.GetRelationIndexByName("22222")
	if position != relationNotFound {
		t.Errorf("Table.GetRelationIndexByName() ==> relation '22222' index cannot be found!!")
	}
	position = table.GetRelationIndexByName(relations[RELATION_COUNT>>2].GetName())
	if position == relationNotFound {
		t.Errorf("Table.GetRelationIndexByName() ==> relation index must be found!!")
	}

	field := table.GetPrimaryKey()
	if field != nil {
		t.Errorf("Table.GetPrimaryKey() ==> field cannot be found")
	}
}

func Test__Table__GetIndexByName(t *testing.T) {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var table = new(Table)
	const INDEX_COUNT = 20000
	aarr := []string{"Gga", "Zorba", "testllk", "testllk22"}

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)
	field1 := Field{}
	field1.Init(2, "Zorba", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field1)
	field2 := Field{}
	field2.Init(3, "testllk", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field2)
	field3 := Field{}
	field3.Init(4, "testllk22", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field3)

	for i := -100; i <= INDEX_COUNT; i++ {
		index := new(Index)
		nameLenght := (abs(i) % 30) + 2
		// fixture
		indexName := randStringBytes(nameLenght)
		index.Init(21, indexName, "hellkzae", aarr, false, true, false, true)
		indexes = append(indexes, *index)
	}

	table.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, metaSchemaName, tabletype.Fake, databaseprovider.Undefined, "", true, false, true, true)

	//t.Errorf("indexes.Count ==> %d ", len(table.indexes))

	for i := 0; i < len(indexes); i++ {
		// test valid field only
		indexName := indexes[i].GetName()
		index := table.GetIndexByName(indexName)
		if index == nil {
			t.Errorf("Table.GetIndexByName() ==> indexes[i].name; i=%d, name=%s, id=%d", i, indexName, indexes[i].GetId())
			break
		} else {
			if indexName != index.GetName() {
				t.Errorf("Table.GetIndexByName() ==> indexes[i].name; i=%d, name=%s, found=%s", i, indexName, table.GetIndexByName(indexes[i].GetName()).GetName())
				break
			}
		}
	}

	// test nil
	index := table.GetIndexByName("22222")
	if index != nil {
		t.Errorf("Table.GetIndexByName() ==> index '22222' cannot be found!!")
	}

}

func Test__Table__Clone(t *testing.T) {
	var fields = []Field{}
	var relations = []Relation{}
	var indexes = []Index{}
	var uk Index = Index{}

	var t1 = new(Table)

	// creating fields
	field0 := Field{}
	field0.Init(1, "Gga", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field0)
	field1 := Field{}
	field1.Init(2, "Zorba", "", fieldtype.Int, 0, "", true, true, true, false, true)
	fields = append(fields, field1)

	var indexedFields = []string{"Zorba"}
	uk.Init(1, "uk_test", "ATable Test", indexedFields, false, false, true, true)
	indexes = append(indexes, uk)

	t1.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
		physicaltype.Table, -111, metaSchemaName, tabletype.Fake, databaseprovider.Undefined, "", true, false, true, true)

	t2 := t1.Clone()

	//TODO add more test
	//TODO check if we have a deep copy!!!!!!!!!!!!!!!!!!!!
	if t1.GetName() != t2.GetName() {
		t.Errorf("Table.Clone() ==> t1.GetName() <> t2.GetName()")
	}
	if t1.GetId() != t2.GetId() {
		t.Errorf("Table.Clone() ==> t1.GetId() <> t2.GetId()")
	}
	if t1.GetDescription() != t2.GetDescription() {
		t.Errorf("Table.Clone() ==> t1.GetDescription() <> t2.GetDescription()")
	}
	if t1.GetFieldCount() != t2.GetFieldCount() {
		t.Errorf("Table.Clone() ==> t1.GetFieldCount() <> t2.GetFieldCount()")
	}
	if t1.GetPhysicalName() != t2.GetPhysicalName() {
		t.Errorf("Table.Clone() ==> t1.GetPhysicalName() <> t2.GetPhysicalName()")
	}
	// no reference copy
	if t1.GetPrimaryKey() != t2.GetPrimaryKey() {
		t.Errorf("Table.Clone() ==> t1.GetPrimaryKey() reference <> t2.GetPrimaryKey() reference")
	}
	if t1.GetDatabaseProvider() != t2.GetDatabaseProvider() {
		t.Errorf("Table.Clone() ==> t1.GetDatabaseProvider() <> t2.GetDatabaseProvider()")
	}
	if t1.GetDdl(ddlstatement.Create, nil, nil) != t2.GetDdl(ddlstatement.Create, nil, nil) {
		t.Errorf("Table.Clone() ==> t1.GetDdlSql()<> t2.GetDdlSql()")
	}
}

// test mapper
func Test__Table__toMeta(t *testing.T) {
	var relations = []Relation{}
	var indexes = []Index{}
	var fields = []Field{}

	elemf00 := Field{}
	elemf00.Init(21, "Field Test", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)
	fields = append(fields, elemf00)

	elemf01 := Field{}
	elemf01.Init(22, "Field Test2", "Field Test", fieldtype.Double, 5, "", true, false, false, true, true)
	fields = append(fields, elemf01)

	elemt := Table{}
	elemt.Init(2222, "rel test", "hellkzae", fields, relations, indexes, physicaltype.Table, 64, "", tabletype.Business, databaseprovider.Influx, "subject test",
		true, true, true, true)

	elemr0 := Relation{}
	//provider databaseprovider.DatabaseProvider, tableType tabletype.TableType
	elemr0.Init(23, "rel test", "hellkzae", &elemt, relationtype.Otop, true, false, true, false)

	metaData := elemt.toMeta()
	newTable := metaData.toTable(fields, relations, indexes)

	if elemt.GetId() != newTable.GetId() {
		t.Errorf("Table.toMeta() ==> t0.GetId() must be equal to t1.GetId()")
	}
	if elemt.GetName() != newTable.GetName() {
		t.Errorf("Table.toMeta() ==> t0.GetName() must be equal to t1.GetName()")
	}
	if elemt.GetDescription() != newTable.GetDescription() {
		t.Errorf("Table.toMeta() ==> t0.GetDescription() must be equal to t1.GetDescription()")
	}
	if elemt.GetPhysicalType() != newTable.GetPhysicalType() {
		t.Errorf("Table.toMeta() ==> t0.GetPhysicalType() must be equal to t1.GetPhysicalType()")
	}
	if elemt.GetSubject() != newTable.GetSubject() {
		t.Errorf("Table.toMeta() ==> t0.GetSubject() must be equal to t1.GetSubject()")
	}
	if elemt.IsCached() != newTable.IsCached() {
		t.Errorf("Table.toMeta() ==> t0.IsCached() must be equal to t1.IsCached()")
	}
	if elemt.IsBaseline() != newTable.IsBaseline() {
		t.Errorf("Table.toMeta() ==> t0.IsBaseline() must be equal to t1.IsBaseline()")
	}
	if elemt.IsReadonly() != newTable.IsReadonly() {
		t.Errorf("Table.toMeta() ==> t0.IsReadonly() must be equal to t1.IsReadonly()")
	}
	if elemt.IsActive() != newTable.IsActive() {
		t.Errorf("Table.toMeta() ==> t0.IsActive() must be equal to t1.IsActive()")
	}
}

func Test__Table__GetDml(t *testing.T) {
	tbl := new(Table)
	//======================
	//==== testing PostgreSql
	//======================
	// table @log
	table := tbl.getLogTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl := "INSERT INTO information_schema.\"@log\" (id,entry_time,level_id,schema_id,thread_id,call_site,job_id,\"method\",line_number,message,description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	if table.GetDml(dmlstatement.Insert, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	// table @meta
	table = tbl.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl = "INSERT INTO information_schema.\"@meta\" (id,schema_id,object_type,reference_id,data_type,flags,\"name\",description,\"value\",active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	if table.GetDml(dmlstatement.Insert, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	expectedSQl = "DELETE FROM information_schema.\"@meta\" WHERE id=$1 AND schema_id=$2 AND object_type=$3 AND reference_id=$4"
	if table.GetDml(dmlstatement.Delete, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	expectedSQl = "UPDATE information_schema.\"@meta\" SET \"value\"=$1,id=$2,schema_id=$3,object_type=$4,reference_id=$5,data_type=$6,flags=$7,\"name\"=$8,description=$9,\"value\"=$10,active=$11 WHERE id=$12 AND schema_id=$13 AND object_type=$14 AND reference_id=$15"
	field := table.GetFieldByName("value")
	field1 := table.GetFieldByName("id")
	field2 := table.GetFieldByName("schema_id")
	field3 := table.GetFieldByName("object_type")
	field4 := table.GetFieldByName("reference_id")
	field5 := table.GetFieldByName("data_type")
	field6 := table.GetFieldByName("flags")
	field7 := table.GetFieldByName("name")
	field8 := table.GetFieldByName("description")
	field9 := table.GetFieldByName("value")
	field10 := table.GetFieldByName("active")
	fields := []*Field{field, field1, field2, field3, field4, field5, field6, field7, field8, field9, field10}
	if table.GetDml(dmlstatement.Update, fields) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	// table @meta_id
	table = tbl.getMetaIdTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl = "DELETE FROM information_schema.\"@meta_id\" WHERE id=$1 AND schema_id=$2 AND object_type=$3"
	if table.GetDml(dmlstatement.Delete, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	expectedSQl = "UPDATE information_schema.\"@meta_id\" SET \"value\"=$1 WHERE id=$2 AND schema_id=$3 AND object_type=$4"
	field = table.GetFieldByName("value")
	fields = []*Field{field}
	if table.GetDml(dmlstatement.Update, fields) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	// table @lexicon
	table = tbl.getLexiconTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl = "INSERT INTO information_schema.\"@lexicon\" (id,schema_id,\"name\",s_name,description,uuid,s_uuid,table_id,source_field_id,target_field_id,relation_id,relation_value,modify_stmp,cached,active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)"
	if table.GetDml(dmlstatement.Insert, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}

	//======================
	//==== testing Mysql
	//======================
	// table @log
	table = tbl.getLogTable(databaseprovider.MySql, "information_schema")
	expectedSQl = "INSERT INTO information_schema.`@log` (id,entry_time,level_id,schema_id,thread_id,call_site,job_id,`method`,line_number,message,description) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	if table.GetDml(dmlstatement.Insert, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	// table @meta_id
	table = tbl.getMetaIdTable(databaseprovider.MySql, "information_schema")
	expectedSQl = "UPDATE information_schema.`@meta_id` SET `value`=? WHERE id=? AND schema_id=? AND object_type=?"
	field = table.GetFieldByName("value")
	fields = []*Field{field}
	if table.GetDml(dmlstatement.Update, fields) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
	// table @meta
	table = tbl.getMetaTable(databaseprovider.MySql, "information_schema")
	expectedSQl = "DELETE FROM information_schema.`@meta` WHERE id=? AND schema_id=? AND object_type=? AND reference_id=?"
	if table.GetDml(dmlstatement.Delete, nil) != expectedSQl {
		t.Errorf("Table.GetDml() ==> query must be equal to " + expectedSQl)
	}
}

func Test__Table__GetDql(t *testing.T) {
	tbl := new(Table)
	//======================
	//==== testing PostgreSql
	//======================
	// table @log
	table := tbl.getLogTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl := "SELECT id,entry_time,level_id,schema_id,thread_id,call_site,job_id,\"method\",line_number,message,description FROM information_schema.\"@log\""
	if table.GetDql("", "") != expectedSQl {
		t.Errorf("Table.GetDql() ==> query must be equal to " + expectedSQl)
	}
	expectedSQl = "SELECT id,entry_time,level_id,schema_id,thread_id,call_site,job_id,\"method\",line_number,message,description FROM information_schema.\"@log\" WHERE level_id>=? AND thread_id=? ORDER BY entry_time DESC"
	if table.GetDql("level_id>=? AND thread_id=?", "entry_time DESC") != expectedSQl {
		t.Errorf("Table.GetDql() ==> query must be equal to " + expectedSQl)
	}
	if table.GetDql("level_id>=? AND thread_id=?", "entry_time DESC") != expectedSQl {
		t.Errorf("Table.GetDql() ==> query must be equal to " + expectedSQl)
	}
}

func Test__Table__GetDdl(t *testing.T) {
	tbl := new(Table)
	tblspc := new(tablespace)
	tblspc.Init(1, "Test", "", "", true, true)

	//======================
	//==== testing PostgreSql
	//======================
	// table @log: Truncate
	table := tbl.getLogTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl := "TRUNCATE TABLE information_schema.\"@log\""
	if table.GetDdl(ddlstatement.Truncate, nil, nil) != expectedSQl {
		t.Errorf("Table.GetDdl(Truncate) ==> query must be equal to " + expectedSQl)
	}
	// table @test: Create
	table = getTestTable(databaseprovider.PostgreSql, "information_schema")
	createScript := table.GetDdl(ddlstatement.Create, tblspc, nil)
	createScript = strings.ReplaceAll(createScript, "\n", " ")
	expectedSQl = "CREATE TABLE information_schema.\"t_@test\" ( id int8, entry_time timestamp without time zone, level_id int2, schema_id int4, thread_id int2, call_site varchar(255) COLLATE \"C\", s_call_site varchar(255) COLLATE \"C\", test2test int8 ) "
	expectedSQl += "WITH (autovacuum_enabled=false)  TABLESPACE Test"

	if createScript != expectedSQl {
		t.Errorf("Table.GetDdl(Create) ==> query must be equal to " + expectedSQl)
	}
	// table @test: Drop
	table = getTestTable(databaseprovider.PostgreSql, "information_schema")
	expectedSQl = "DROP TABLE information_schema.\"t_@test\" CASCADE"
	if table.GetDdl(ddlstatement.Drop, nil, nil) != expectedSQl {
		t.Errorf("Table.GetDdl(Drop) ==> query must be equal to " + expectedSQl)
	}
	// table @test: Alter add
	table = getTestTable(databaseprovider.PostgreSql, "information_schema")
	field := new(Field)
	field.Init(2, "test11", "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	expectedSQl = "ALTER TABLE information_schema.\"t_@test\" ADD COLUMN test11 timestamp without time zone"
	if table.GetDdl(ddlstatement.Alter, nil, field) != expectedSQl {
		t.Errorf("Table.GetDdl(Alter add) ==> query must be equal to " + expectedSQl)
	}
	// table @test: Alter drop
	field = new(Field)
	field.Init(2, metaSchemaId, "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	expectedSQl = "ALTER TABLE information_schema.\"t_@test\" DROP COLUMN schema_id CASCADE"
	if table.GetDdl(ddlstatement.Alter, nil, field) != expectedSQl {
		t.Errorf("Table.GetDdl(Alter drop) ==> query must be equal to " + expectedSQl)
	}
	expectedSQl = ""
	if table.GetDdl(ddlstatement.Alter, nil, nil) != expectedSQl {
		t.Errorf("Table.GetDdl(Alter drop with null field) ==> query must be empty")
	}
}

func Test__Table__getUniqueFieldList(t *testing.T) {
	tbl := new(Table)
	table := tbl.getLogTable(databaseprovider.PostgreSql, "information_schema")
	// business table
	table.setTableType(tabletype.Business)
	if table.getUniqueFieldList() != "id" {
		t.Errorf("Table.getUniqueFieldList() ==> Must be equal to 'id'")
	}
	// @meta table
	table = tbl.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	if table.getUniqueFieldList() != "id,schema_id,object_type,reference_id" {
		t.Errorf("Table.getUniqueFieldList() ==> Must be equal to 'id,schema_id,object_type,reference_id'")
	}
	// @meta_id table
	table = tbl.getMetaIdTable(databaseprovider.PostgreSql, "information_schema")
	if table.getUniqueFieldList() != "id,schema_id,object_type" {
		t.Errorf("Table.getUniqueFieldList() ==> Must be equal to 'id,schema_id,object_type'")
	}
	// @mtm table
	table = getMtmTable()
	if table.getUniqueFieldList() != "relationA,relationB" {
		t.Errorf("Table.getUniqueFieldList() ==> Must be equal to 'relationA,relationB'")
	}
}

func Test__Table__GetQueryResult(t *testing.T) {
	tbl := new(Table)
	table := tbl.getMetaTable(databaseprovider.PostgreSql, "information_schema")
	var result []interface{}
	var resultPtr []interface{}
	var id int32 = 5
	//var schemaId int32 = 1
	var objectType int8 = 15
	var referenceId int = 2
	var dataType int16 = 3
	var flags int32 = 554545
	var value int64 = 2222222
	var name string = "testName"
	var nameDesc string = "testName DESC"
	var active bool = true

	//======================
	//==== TEST 1
	//======================
	result = append(result, id)          //1
	result = append(result, nil)         //2 - schema_id
	result = append(result, objectType)  //3
	result = append(result, referenceId) //4
	result = append(result, dataType)    //5
	result = append(result, flags)       //6
	result = append(result, name)        //7
	result = append(result, nameDesc)    //8
	result = append(result, value)       //9
	result = append(result, active)      //10

	for i := 0; i < len(result); i++ {
		resultPtr = append(resultPtr, &result[i])
	}
	arr := table.GetQueryResult(resultPtr)

	if arr[table.GetFieldIndexByName("id")] != strconv.Itoa(int(id)) {
		t.Errorf("Table.GetQueryResult() ==> id must be equal to '%s'", strconv.Itoa(int(id)))
	}
	if arr[table.GetFieldIndexByName("object_type")] != strconv.Itoa(int(objectType)) {
		t.Errorf("Table.GetQueryResult() ==> object_type must be equal to '%s'", strconv.Itoa(int(objectType)))
	}
	if arr[table.GetFieldIndexByName("schema_id")] != "" {
		t.Errorf("Table.GetQueryResult() ==> schema_id must be equal to null")
	}
	if arr[table.GetFieldIndexByName("value")] != strconv.Itoa(int(value)) {
		t.Errorf("Table.GetQueryResult() ==> value must be equal to '%s'", strconv.Itoa(int(value)))
	}
	if arr[table.GetFieldIndexByName("flags")] != strconv.Itoa(int(flags)) {
		t.Errorf("Table.GetQueryResult() ==> flags must be equal to '%s'", strconv.Itoa(int(flags)))
	}
	if arr[table.GetFieldIndexByName("name")] != name {
		t.Errorf("Table.GetQueryResult() ==> name must be equal to '%s'", name)
	}

	//======================
	//==== TEST 2
	//======================
	var idU uint = 4294967295
	var objectTypeU uint8 = 150
	var referenceIdU uint16 = 20
	var dataTypeU uint32 = 30
	var flagsU uint64 = 18446744073709551615 // test max value
	result[0] = idU
	result[2] = objectTypeU
	result[3] = referenceIdU
	result[4] = dataTypeU
	result[5] = flagsU
	arr = table.GetQueryResult(resultPtr)

	if arr[table.GetFieldIndexByName("id")] != strconv.FormatUint(uint64(idU), 10) {
		t.Errorf("Table.GetQueryResult() ==> id must be equal to '%s'", strconv.Itoa(int(objectTypeU)))
	}
	if arr[table.GetFieldIndexByName("object_type")] != strconv.Itoa(int(objectTypeU)) {
		t.Errorf("Table.GetQueryResult() ==> object_type must be equal to '%s'", strconv.Itoa(int(objectTypeU)))
	}
	if arr[table.GetFieldIndexByName("flags")] != strconv.FormatUint(flagsU, 10) {
		t.Errorf("Table.GetQueryResult() ==> flags must be equal to '%s'", strconv.FormatUint(flagsU, 10))
	}
	if arr[table.GetFieldIndexByName("data_type")] != strconv.Itoa(int(dataTypeU)) {
		t.Errorf("Table.GetQueryResult() ==> data_type must be equal to '%s'", strconv.Itoa(int(dataTypeU)))
	}

	var objectTypeF float32 = 150
	var referenceIdF float64 = 20
	result[2] = objectTypeF
	result[3] = referenceIdF
	arr = table.GetQueryResult(resultPtr)

	if arr[table.GetFieldIndexByName("object_type")] != strconv.Itoa(int(objectTypeF)) {
		t.Errorf("Table.GetQueryResult() ==> object_type must be equal to '%s'", strconv.Itoa(int(objectTypeF)))
	}
}

func Test__Table__getVariableInfo(t *testing.T) {
	schema := new(Schema)
	table := new(Table)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)

	//======================
	//==== testing PostgreSql
	//======================
	table = table.getMetaTable(databaseprovider.PostgreSql, "zorba")
	name, initialId := table.getVariableInfo()
	if initialId != 1 {
		t.Errorf("Table.getVariableInfo() ==> initial id must be equal to 1")
	}
	if name != "$" {
		t.Errorf("Table.getVariableInfo() ==> initial name must be equal to '$'")
	}
	//======================
	//==== testing Mysql
	//======================
	table = table.getMetaTable(databaseprovider.MySql, "zorba")
	name, initialId = table.getVariableInfo()
	if initialId != 0 {
		t.Errorf("Table.getVariableInfo() ==> initial id must be equal to 1")
	}
	if name != "?" {
		t.Errorf("Table.getVariableInfo() ==> initial name must be equal to '$'")
	}
	//======================
	//==== testing Others
	//======================
	table = table.getMetaTable(databaseprovider.Undefined, "zorba")
	name, initialId = table.getVariableInfo()
	if initialId != 0 {
		t.Errorf("Table.getVariableInfo() ==> initial id must be equal to 0")
	}
	if name != "" {
		t.Errorf("Table.getVariableInfo() ==> initial name must be equal to '$'")
	}

}

func Test__Table__getTable(t *testing.T) {
	schema := new(Schema)
	table := new(Table)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	table = table.getMetaTable(databaseprovider.PostgreSql, "")
	tableId := table.GetId()
	metaList := getTableMetas(tableId)
	//t.Errorf("%d", len(metaList))
	table = table.getTable(databaseprovider.PostgreSql, "", 0, metaList)
	if table.name != metaTableName {
		t.Errorf("Table.getTable() ==> name is diferent than '@meta'")
	}
	if table.id != tableId {
		t.Errorf("Table.getTable() ==> id is diferent than 3")
	}
	if table.IsBaseline() != true {
		t.Errorf("Table.getTable() ==> IsBaseline() is diferent than true")
	}
	if table.IsActive() != true {
		t.Errorf("Table.getTable() ==> IsActive() is diferent than true")
	}
	if table.GetType() != tabletype.Business {
		t.Errorf("Table.getTable() ==> GetType() is diferent than tabletype.Meta")
	}
	if table.GetDatabaseProvider() != databaseprovider.PostgreSql {
		t.Errorf("Table.getTable() ==> GetDatabaseProvider() is diferent than databaseprovider.PostgreSql")
	}
	if table.GetEntityType() != entitytype.Table {
		t.Errorf("Table.getTable() ==> GetEntityType() is diferent than entitytype.Table")
	}

}

func Test__Table__getPhysicalName(t *testing.T) {
	table := new(Table)

	table.setName("hello")
	table.setTableType(tabletype.Business)
	table.setDatabaseProvider(databaseprovider.PostgreSql)

	if table.getPhysicalName("test") != "test.t_hello" {
		t.Errorf("Table.getPhysicalName() ==> getPhysicalName() is diferent than 'test.t_hello'")
	}

	table.setName("@meta_id")
	table.setTableType(tabletype.MetaId)

	if table.getPhysicalName("test") != "test.\"@meta_id\"" {
		t.Errorf("Table.getPhysicalName() ==> getPhysicalName() is diferent than 'test.\"meta_id\"'")
	}
}

func Test__Table__equal(t *testing.T) {
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	table1 := schema.GetTableByName(metaTableName)
	table2 := table1.Clone()

	// positive test
	if table1.equal(table2) == false {
		t.Errorf("Table.equal() ==> {1} table1 <> table2")
	}

	// positive - modify case of a field name
	table2 = table1.Clone()
	field2 := table2.GetFieldByIndex(2)
	field2.setName(strings.ToUpper(field2.name))
	temp := make([]Field, len(table2.fields), len(table2.fields))
	for i := 0; i < len(table2.fields); i++ {
		temp[i] = *table2.fields[i]
	}
	table2.loadFields(temp, tabletype.Meta)
	table2.loadMapper() // !!! load after loadFields()
	if table1.equal(table2) == false {
		t.Errorf("Table.equal() ==> {2} table1 <> table2")
	}

	// negative - change a field
	table2 = table1.Clone()
	field2 = table2.GetFieldByIndex(0)
	field2.setSize(555)
	if table1.equal(table2) == true {
		t.Errorf("Table.equal() ==> {3} table1 == table2")
	}

	// negative - add a field
	table2 = table1.Clone()
	field2 = table2.GetFieldByIndex(0)
	table2.fields = append(table2.fields, field2.Clone())
	if table1.equal(table2) == true {
		t.Errorf("Table.equal() ==> {4} table1 == table2")
	}
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}

func getTestTable(provider databaseprovider.DatabaseProvider, schemaPhysicalName string) *Table {
	var fields []Field
	var relations []Relation
	var indexes []Index
	var result = new(Table)

	// physical_name is built later
	var id = Field{}
	var entryTime = Field{}
	var levelId = Field{}
	var schemaId = Field{}
	var threadId = Field{}
	var callSite = Field{}

	// "id","entry_time","level_id","thread_id","call_site","message","description","machine_name"
	id.Init(1, metaLogId, "", fieldtype.Long, 0, "", true, true, true, false, true)
	entryTime.Init(2, metaLogEntryTime, "", fieldtype.DateTime, 0, "", true, true, true, false, true)
	levelId.Init(3, metaLogLevelId, "", fieldtype.Short, 0, "", true, true, true, false, true)
	schemaId.Init(4, metaSchemaId, "", fieldtype.Int, 0, "", true, true, true, false, true)
	threadId.Init(5, metaLogThreadId, "", fieldtype.Byte, 0, "", true, false, true, false, true)
	callSite.Init(6, metaLogCallSite, "", fieldtype.String, 255, "", true, false, false, false, true)

	fields = append(fields, id)        //1
	fields = append(fields, entryTime) //2
	fields = append(fields, levelId)   //3
	fields = append(fields, schemaId)  //4
	fields = append(fields, threadId)  //5
	fields = append(fields, callSite)  //6

	result.Init(-2, "@test", "", fields, relations, indexes, physicaltype.Table, 0, schemaPhysicalName,
		tabletype.Business, provider, "", false, false, true, true)

	elemr0 := Relation{}
	elemr0.Init(1, "test2test", "", result, relationtype.Mto, false, false, true, false)
	relations = append(relations, elemr0)
	result = new(Table)
	result.Init(-3, "@test", "", fields, relations, indexes, physicaltype.Table, 0, schemaPhysicalName,
		tabletype.Business, provider, "", false, false, true, true)

	return result
}

func getMtmTable() *Table {

	var fields = make([]Field, 0, 0)
	var relations = make([]Relation, 0, 2)
	var indexes = make([]Index, 1, 1)
	var result = new(Table)

	// relations
	var relationA = Relation{}
	var relationB = Relation{}

	relationA.Init(1, "relationA", "", nil, relationtype.Mto, false, true, true, true)
	relationB.Init(2, "relationB", "", nil, relationtype.Mto, false, true, true, true)

	relations = append(relations, relationA)
	relations = append(relations, relationB)

	result.Init(77, "@test", "", fields, relations, indexes, physicaltype.Table,
		0, "test", tabletype.Mtm, databaseprovider.PostgreSql, "",
		true, false, true, true)

	return result
}

func getTableMetas(referenceId int32) []*meta {
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	metaList := schema.toMeta()
	result := make([]*meta, 0)
	var fieldIndex = 0
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table && meta.id == referenceId {
			result = append(result, meta)
		}
		if meta.refId == referenceId {
			result = append(result, meta)
			if meta.GetEntityType() == entitytype.Field {
				if fieldIndex == 1 {
					meta.objectType = int8(entitytype.Relation)
				}
				fieldIndex++
			}
		}
	}
	return result
}
