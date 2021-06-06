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
	if metaData.GetEntityType() != entitytype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> objectType==4 must return entitytype.NotDefined")
	}
	metaData.objectType = 5
	if metaData.GetEntityType() != entitytype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> objectType==5 must return entitytype.NotDefined")
	}

	metaData.objectType = 6
	if metaData.GetEntityType() != entitytype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> objectType==6 must return entitytype.NotDefined")
	}
}

func Test__Meta__GetRelationType(t *testing.T) {
	var metaData = new(meta)

	// key
	metaData.id = 12123
	metaData.refId = 545434

	metaData.setRelationType(4)
	if metaData.GetRelationType() != relationtype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> relationType==4 must return entitytype.NotDefined")
	}
	metaData.setRelationType(5)
	if metaData.GetRelationType() != relationtype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> relationType==5 must return entitytype.NotDefined")
	}

	metaData.setRelationType(6)
	if metaData.GetRelationType() != relationtype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> relationType==6 must return entitytype.NotDefined")
	}
}

func Test__Meta__GetFieldType(t *testing.T) {
	var metaData = new(meta)

	// key
	metaData.id = 12123
	metaData.refId = 545434

	metaData.dataType = 4
	if metaData.GetFieldType() != fieldtype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> fieldType==4 must return entitytype.NotDefined")
	}
	metaData.dataType = 5
	if metaData.GetFieldType() != fieldtype.NotDefined {
		t.Errorf("meta.GetEntityType() ==> fieldType==5 must return entitytype.NotDefined")
	}

	metaData.dataType = 6
	if metaData.GetFieldType() != fieldtype.NotDefined {
		t.Errorf("meta.GetFieldType() ==> fieldType==6 must return entitytype.NotDefined")
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
	if metaData.toParameter() != nil {
		t.Errorf("meta.toParameter() ==> objectType==4 must return nil")
	}
}
