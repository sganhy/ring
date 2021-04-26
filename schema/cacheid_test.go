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
