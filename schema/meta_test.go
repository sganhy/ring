package schema

import (
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"testing"
)

func Test__Meta__GetEntityType(t *testing.T) {
	var metaData = new(meta)

	// key
	metaData.id = 12123
	metaData.refId = 545434

	metaData.objectType = 4
	if metaData.GetEntityType() != entitytype.Undefined {
		t.Errorf("meta.GetEntityType() ==> objectType==4 must return entitytype.Undefined")
	}
	metaData.objectType = 5
	if metaData.GetEntityType() != entitytype.Undefined {
		t.Errorf("meta.GetEntityType() ==> objectType==5 must return entitytype.Undefined")
	}

	metaData.objectType = 6
	if metaData.GetEntityType() != entitytype.Undefined {
		t.Errorf("meta.GetEntityType() ==> objectType==6 must return entitytype.Undefined")
	}
}

func Test__Meta__GetRelationType(t *testing.T) {
	var metaData = new(meta)

	// key
	metaData.id = 12123
	metaData.refId = 545434

	metaData.setRelationType(4)
	if metaData.GetRelationType() != relationtype.Undefined {
		t.Errorf("meta.GetEntityType() ==> relationType==4 must return entitytype.Undefined")
	}
	metaData.setRelationType(5)
	if metaData.GetRelationType() != relationtype.Undefined {
		t.Errorf("meta.GetEntityType() ==> relationType==5 must return entitytype.Undefined")
	}

	metaData.setRelationType(6)
	if metaData.GetRelationType() != relationtype.Undefined {
		t.Errorf("meta.GetEntityType() ==> relationType==6 must return entitytype.Undefined")
	}
}

func Test__Meta__GetFieldType(t *testing.T) {
	var metaData = new(meta)

	// key
	metaData.id = 12123
	metaData.refId = 545434

	metaData.dataType = 4
	if metaData.GetFieldType() != fieldtype.Undefined {
		t.Errorf("meta.GetEntityType() ==> fieldType==4 must return entitytype.Undefined")
	}
	metaData.dataType = 5
	if metaData.GetFieldType() != fieldtype.Undefined {
		t.Errorf("meta.GetEntityType() ==> fieldType==5 must return entitytype.Undefined")
	}

	metaData.dataType = 6
	if metaData.GetFieldType() != fieldtype.Undefined {
		t.Errorf("meta.GetFieldType() ==> fieldType==6 must return entitytype.Undefined")
	}
}

//test ToField, ToRelation, ToIndex, ToTable, and ToSchema
func Test__Meta__toField(t *testing.T) {
	var metaData = new(meta)

	// testing nil return
	metaData.objectType = 4
	if metaData.toField() != nil {
		t.Errorf("meta.toField() ==> objectType==4 must return nil")
	}
	if metaData.toRelation(nil) != nil {
		t.Errorf("meta.toRelation() ==> objectType==4 must return nil")
	}
	if metaData.toIndex() != nil {
		t.Errorf("meta.toIndex() ==> objectType==4 must return nil")
	}
	if metaData.toTable(nil, nil, nil) != nil {
		t.Errorf("meta.toTable() ==> objectType==4 must return nil")
	}
	if metaData.toParameter(54) != nil {
		t.Errorf("meta.toParameter() ==> objectType==4 must return nil")
	}
	if metaData.toTablespace() != nil {
		t.Errorf("meta.toTablespace() ==> objectType==4 must return nil")
	}
	if metaData.toSequence(54) != nil {
		t.Errorf("meta.toSequence() ==> objectType==4 must return nil")
	}

}

//test equal
func Test__Meta__equal(t *testing.T) {
	var metaA = new(meta)
	var metaB = new(meta)

	metaA.dataType = 1
	metaA.name = "Test__Meta__equal"
	metaA.description = "Test__Meta__equal desc"
	metaA.flags = 7777
	metaA.value = "value test"
	metaA.enabled = true
	metaA.id = 1
	metaA.refId = 111

	metaB.dataType = 1
	metaB.name = "Test__Meta__equal"
	metaB.description = "Test__Meta__equal desc"
	metaB.flags = 7777
	metaB.value = "value test"
	metaB.enabled = true
	metaB.id = 2
	metaB.refId = 222

	if metaB.equal(metaA) == false {
		t.Errorf("meta.equal() ==> metaB should be equal to metaA")
	}
	if metaA.equal(metaB) == false {
		t.Errorf("meta.equal() ==> metaA should be equal to metaB")
	}
	metaA.enabled = false
	if metaB.equal(metaA) == true {
		t.Errorf("meta.equal() ==> metaA shouldn't be equal to metaB")
	}
	if metaB.equal(nil) == true {
		t.Errorf("meta.equal() ==> metaA shouldn't be equal to null")
	}
}

//test all setters
func Test__Meta__setters(t *testing.T) {
	var metaA = new(meta)

	//======================
	//==== testing flags
	//======================
	metaA.flags = 99777799

	metaA.setEntityBaseline(false)
	if metaA.IsEntityBaseline() == true {
		t.Errorf("meta.setEntityBaseline() ==> IsEntityBaseline shouldn't be equal to false")
	}
	metaA.setEntityBaseline(true)
	if metaA.IsEntityBaseline() == false {
		t.Errorf("meta.setEntityBaseline() ==> IsEntityBaseline shouldn't be equal to true")
	}
	metaA.setTablespaceIndex(true)
	if metaA.IsTablespaceIndex() == false {
		t.Errorf("meta.setTablespaceIndex() ==> IsTablespaceIndex shouldn't be equal to true")
	}
	metaA.setTablespaceIndex(false)
	if metaA.IsTablespaceIndex() == true {
		t.Errorf("meta.setTablespaceIndex() ==> IsTablespaceIndex shouldn't be equal to false")
	}
	metaA.setTablespaceTable(false)
	if metaA.IsTablespaceTable() == true {
		t.Errorf("meta.setTablespaceTable() ==> IsTablespaceTable shouldn't be equal to false")
	}
	metaA.setTablespaceTable(true)
	if metaA.IsTablespaceTable() == false {
		t.Errorf("meta.setTablespaceTable() ==> IsTablespaceTable shouldn't be equal to true")
	}
	metaA.setRelationConstraint(false)
	if metaA.IsRelationConstraint() == true {
		t.Errorf("meta.setRelationConstraint() ==> IsRelationConstraint shouldn't be equal to false")
	}
	metaA.setRelationConstraint(true)
	if metaA.IsRelationConstraint() == false {
		t.Errorf("meta.setRelationConstraint() ==> IsRelationConstraint shouldn't be equal to true")
	}
	metaA.setIndexUnique(false)
	if metaA.IsIndexUnique() == true {
		t.Errorf("meta.setIndexUnique() ==> IsIndexUnique shouldn't be equal to false")
	}
	metaA.setIndexUnique(true)
	if metaA.IsIndexUnique() == false {
		t.Errorf("meta.setIndexUnique() ==> IsIndexUnique shouldn't be equal to true")
	}
	metaA.setFieldNotNull(false)
	if metaA.IsFieldNotNull() == true {
		t.Errorf("meta.setFieldNotNull() ==> IsFieldNotNull shouldn't be equal to false")
	}
	metaA.setFieldNotNull(true)
	if metaA.IsFieldNotNull() == false {
		t.Errorf("meta.setFieldNotNull() ==> IsFieldNotNull shouldn't be equal to true")
	}
	metaA.setFieldCaseSensitive(false)
	if metaA.IsFieldCaseSensitive() == true {
		t.Errorf("meta.setFieldCaseSensitive() ==> IsFieldCaseSensitive shouldn't be equal to false")
	}
	metaA.setFieldCaseSensitive(true)
	if metaA.IsFieldCaseSensitive() == false {
		t.Errorf("meta.setFieldCaseSensitive() ==> IsFieldCaseSensitive shouldn't be equal to true")
	}

}
