package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"strings"
	"testing"
)

func Test__Validator__isValidName(t *testing.T) {
	valid := new(validator)

	if valid.isValidName("element") == false {
		t.Errorf("validator.isValidName() ==> 'element' must be valid")
	}
	if valid.isValidName("element1") == false {
		t.Errorf("validator.isValidName() ==> 'element1' must be valid")
	}
	if valid.isValidName("element_°") == true {
		t.Errorf("validator.isValidName() ==> 'element_°' must be invalid")
	}
	if valid.isValidName("element 11") == true {
		t.Errorf("validator.isValidName() ==> 'element 11' must be invalid")
	}
	if valid.isValidName("element-11") == true {
		t.Errorf("validator.isValidName() ==> 'element-11' must be invalid")
	}
}

func Test__Validator__fieldNameUnique(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	valid.fieldNameUnique(importFile)

	if len(importFile.metaList) <= 0 {
		t.Errorf("validator.entityNameUnique() ==> importFile.metaList is empty we are not able to unitest the method")
		// escape here !!!!
		return
	}

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.entityNameUnique() ==> importFile.errorCount <> 0 ")
	}

	// negative test
	// 1> duplicate field
	var meta = getFieldMeta(importFile, true)
	importFile.metaList = append(importFile.metaList, meta)
	valid.fieldNameUnique(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.entityNameUnique() ==> importFile.errorCount <> 0 ")
	}
	// 2> duplicate relation
	meta = getFieldMeta(importFile, true)
	importFile.errorCount = 0 // reset counter
	meta.objectType = int8(entitytype.Relation)
	importFile.metaList = append(importFile.metaList, meta)
	valid.fieldNameUnique(importFile)
	if importFile.errorCount != 2 {
		// should be equal to 2 ==>  relation name + duplicate field name
		t.Errorf("validator.entityNameUnique() ==> importFile.errorCount <> 2 ")
	}

}

func Test__Validator__duplicateMetaKey(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	valid.duplicateMetaKey(importFile)

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.duplicateMetaKey() ==> importFile.errorCount <> 0 ")
	}

	// negative test -- create duplicate key
	var meta = getTableMeta(importFile, true)
	importFile.metaList = append(importFile.metaList, meta)
	valid.duplicateMetaKey(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.duplicateMetaKey() ==> importFile.errorCount <> 0 ")
	}
}

// test checkEntityName and entityNameValid
func Test__Validator__checkEntityName(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	var meta = getTableMeta(importFile, true)
	meta.name = "test1"
	valid.checkEntityName(importFile, meta, entitytype.Table)

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.checkEntityName(Table) ==> importFile.errorCount <> 0")
	}
	// negative test -- tables
	meta.name = "@meta"
	valid.checkEntityName(importFile, meta, entitytype.Table)
	if importFile.errorCount != 1 {
		t.Errorf("validator.checkEntityName(Table) ==> importFile.errorCount <> 1")
	}
	// negative test -- schema
	meta.name = "@meta"
	meta.dataType = int32(entitytype.Schema)
	importFile.errorCount = 0
	valid.checkEntityName(importFile, meta, entitytype.Schema)
	if importFile.errorCount != 1 {
		t.Errorf("validator.checkEntityName(Schema) ==> importFile.errorCount <> 1")
	}
	// negative test -- schema: metaData.name len > 30
	meta.name = "test012345678901234567890123456789"
	importFile.errorCount = 0
	valid.checkEntityName(importFile, meta, entitytype.Schema)
	if importFile.errorCount != 1 {
		t.Errorf("validator.checkEntityName(Schema) ==> importFile.errorCount <> 1")
	}
	// negative test -- field: metaData.name len > 28
	meta = getFieldMeta(importFile, true)
	meta.name = "test0123456789012345678945899"
	importFile.errorCount = 0
	valid.checkEntityName(importFile, meta, entitytype.Field)
	if importFile.errorCount != 1 {
		t.Errorf("validator.checkEntityName(Field) ==> importFile.errorCount <> 1")
	}
	// positive test: entityNameValid
	importFile = getTestImportFile()
	valid.entityNameValid(importFile)
	if importFile.errorCount != 0 {
		t.Errorf("validator.entityNameValid() ==> importFile.errorCount <> 0")
	}
	// negative test: entityNameValid
	importFile = getTestImportFile()
	importFile.metaList[0].name = " "
	valid.entityNameValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.entityNameValid() ==> importFile.errorCount <> 1")
	}

}

// test indexValid and indexValueValid
func Test__Validator__indexValid(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	valid.indexValid(importFile)
	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.indexValid() ==> importFile.errorCount <> 0")
	}
	valid.indexValueValid(importFile)
	if importFile.errorCount != 0 {
		t.Errorf("validator.indexValueValid() ==> importFile.errorCount <> 0")
	}
	// negative test
	var meta = getIndexMeta(importFile)
	meta.value = "test012345678901234567890123456789"
	valid.indexValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.indexValid() ==> importFile.errorCount <> 1")
	}
	meta.value = "  "
	importFile.errorCount = 0
	valid.indexValueValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.indexValid() ==> importFile.errorCount <> 1")
	}
}

func Test__Validator__tableValueValid(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	valid.tableValueValid(importFile)

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.tableValueValid() ==> importFile.errorCount <> 0")
	}

	// negative test
	var meta = getTableMeta(importFile, true)
	meta.id = -789564224
	importFile.metaList = append(importFile.metaList, meta)
	valid.tableValueValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.tableValueValid() ==> importFile.errorCount <> 1")
	}
}

func Test__Validator__languageCodeValid(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	valid.languageCodeValid(importFile)

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.languageCodeValid() ==> importFile.errorCount <> 0")
	}

}

func Test__Validator__entityTypeValid(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)
	valid.entityTypeValid(importFile)

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.entityTypeValid() ==> importFile.errorCount <> 0")
	}

	// negative test: field
	field := getFieldMeta(importFile, false)
	field.setFieldType(4)
	valid.entityTypeValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.entityTypeValid() ==> importFile.errorCount <> 1")
	}

	// negative test: relation
	importFile.errorCount = 0
	relation := getFieldMeta(importFile, false)
	relation.objectType = int8(entitytype.Relation)
	valid.entityTypeValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.entityTypeValid() ==> importFile.errorCount <> 1")
	}

}

func Test__Validator__tableSpaceValueValid(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)

	tablespace := getTableMeta(importFile, false)
	tablespace.objectType = int8(entitytype.Tablespace)
	valid.tableSpaceValueValid(importFile)

	// negative test
	if importFile.errorCount != 1 {
		t.Errorf("validator.tableSpaceValueValid() ==> importFile.errorCount <> 1")
	}

	// negative test nul
	tablespace.value = "/NUL/"
	importFile.errorCount = 0
	valid.tableSpaceValueValid(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.tableSpaceValueValid() ==> importFile.errorCount <> 1")
	}
}

func Test__Validator__tableNameUnique(t *testing.T) {
	importFile := getMetaImportFile()
	valid := new(validator)

	tablespace := getTableMeta(importFile, false)
	tablespace.objectType = int8(entitytype.Tablespace)
	valid.tableNameUnique(importFile)

	// positive test
	if importFile.errorCount != 0 {
		t.Errorf("validator.tableNameUnique() ==> importFile.errorCount <> 0")
	}

	// negative test
	meta := getTableMeta(importFile, true)
	importFile.metaList = append(importFile.metaList, meta)
	valid.tableNameUnique(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.tableNameUnique() ==> importFile.errorCount <> 1")
	}

}

func getMetaImportFile() *Import {
	var result = new(Import)
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	result.provider = databaseprovider.Undefined
	logger := new(log)
	logger.Init(0, 0, true)
	result.logger = logger
	result.metaList = schema.toMeta()
	return result
}
func getTableMeta(importFile *Import, clone bool) *meta {
	for i := 0; i < len(importFile.metaList); i++ {
		if importFile.metaList[i].GetEntityType() == entitytype.Table {
			if clone {
				return importFile.metaList[i].Clone()
			} else {
				return importFile.metaList[i]
			}
		}
	}
	return nil
}
func getFieldMeta(importFile *Import, clone bool) *meta {
	for i := 0; i < len(importFile.metaList); i++ {
		if importFile.metaList[i].GetEntityType() == entitytype.Field {
			if clone {
				return importFile.metaList[i].Clone()
			} else {
				return importFile.metaList[i]
			}
		}
	}
	return nil
}

//important don't clone indexes
func getIndexMeta(importFile *Import) *meta {
	for i := 0; i < len(importFile.metaList); i++ {
		if importFile.metaList[i].GetEntityType() == entitytype.Index {
			return importFile.metaList[i]
		}
	}
	return nil
}
func getTestImportFile() *Import {
	var result = new(Import)
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	result.provider = databaseprovider.Undefined
	logger := new(log)
	logger.Init(0, 0, true)
	result.logger = logger
	result.metaList = schema.toMeta()
	for i := 0; i < len(result.metaList); i++ {
		result.metaList[i].name = strings.ReplaceAll(result.metaList[i].name, "@", "")
	}
	return result
}
