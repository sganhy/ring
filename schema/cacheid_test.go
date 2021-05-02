package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/dmlstatement"
	"ring/schema/entitytype"
	"testing"
)

func Test__CacheId__GetDdl(t *testing.T) {
	cacheId := new(CacheId)
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	table := schema.GetTableByName(metaIdTableName)

	InitCacheId(schema, schema.GetTableByName("@meta_id"), schema.GetTableByName("@long"))
	cacheId.Init(1, 0, entitytype.Table)
	//======================
	//==== testing PostgreSql
	//======================
	expectedSql := "UPDATE rpg_sheet_test.\"@meta_id\" SET \"value\"=\"value\"+$1 WHERE id=$2 AND schema_id=$3 AND object_type=$4 RETURNING \"value\""
	if cacheId.GetDml(dmlstatement.Update, table) != expectedSql {
		t.Errorf("cacheId.GetDml() ==> query must be equal to " + expectedSql)
	}
}

//test setters/ getters
func Test__CacheId__setters(t *testing.T) {
	cacheId := new(CacheId)
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)

	InitCacheId(schema, schema.GetTableByName("@meta_id"), schema.GetTableByName("@long"))
	cacheId.Init(1, 0, entitytype.Table)

	cacheId.SetCache(true)
	if cacheId.reservedRange != 1 {
		t.Errorf("cacheId.SetCache() ==> reservedRange must be equal to 1")
	}
	cacheId.SetCache(false)
	if cacheId.reservedRange != 0 {
		t.Errorf("cacheId.SetCache() ==> reservedRange must be equal to 0")
	}
	if cacheId.IsInitialized() != false {
		t.Errorf("cacheId.IsInitialized() ==> isInitialized must be equal to FALSE")
	}
	if cacheId.GetCurrentId() != 0 {
		t.Errorf("cacheId.GetCurrentId() ==> GetCurrentId must be equal to 0")
	}
}

func Test__CacheId__toMetaId(t *testing.T) {
	cacheId := new(CacheId)
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	InitCacheId(schema, schema.GetTableByName("@meta_id"), schema.GetTableByName("@long"))
	cacheId.Init(1, 0, entitytype.Table)
	metaId := cacheId.toMetaId(entitytype.Sequence, 11, 111)

	if metaId.schemaId != 111 {
		t.Errorf("cacheId.toMetaId() ==> schemaId must be equal to 111")
	}
	if metaId.id != 11 {
		t.Errorf("cacheId.toMetaId() ==> schemaId must be equal to 11")
	}
	if metaId.objectType != int8(entitytype.Sequence) {
		t.Errorf("cacheId.toMetaId() ==> schemaId must be equal to int8(entitytype.Sequence)")
	}
}
