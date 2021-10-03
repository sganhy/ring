package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
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
	var meta = getFieldMeta(importFile)
	importFile.metaList = append(importFile.metaList, meta)
	valid.fieldNameUnique(importFile)
	if importFile.errorCount != 1 {
		t.Errorf("validator.entityNameUnique() ==> importFile.errorCount <> 0 ")
	}
	// 2> duplicate relation
	meta = getFieldMeta(importFile)
	importFile.errorCount = 0 // reset counter
	meta.objectType = int8(entitytype.Relation)
	importFile.metaList = append(importFile.metaList, meta)
	valid.fieldNameUnique(importFile)
	if importFile.errorCount != 2 {
		// should be equal to 2 ==>  relation name + duplicate field name
		t.Errorf("validator.entityNameUnique() ==> importFile.errorCount <> 2 ")
	}

}

func getMetaImportFile() *Import {
	var result = new(Import)
	schema := new(Schema)
	schema = schema.getMetaSchema(databaseprovider.PostgreSql, "", 0, 0, true)
	result.provider = databaseprovider.NotDefined
	logger := new(log)
	logger.Init(0, 0, true)
	result.logger = logger
	result.metaList = schema.toMeta()
	return result
}

func getTableMeta(importFile *Import) *meta {
	for i := 0; i < len(importFile.metaList); i++ {
		if importFile.metaList[i].GetEntityType() == entitytype.Table {
			return importFile.metaList[i].Clone()
		}
	}
	return nil
}

func getFieldMeta(importFile *Import) *meta {
	for i := 0; i < len(importFile.metaList); i++ {
		if importFile.metaList[i].GetEntityType() == entitytype.Field {
			return importFile.metaList[i].Clone()
		}
	}
	return nil
}
