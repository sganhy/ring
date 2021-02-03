package schema

import (
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"testing"
)

func Test__Meta__GetEntityType(t *testing.T) {
	var meta = new(Meta)

	// key
	meta.id = 12123
	meta.refId = 545434

	meta.objectType = 4
	if meta.GetEntityType() != entitytype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> objectType==4 must return entitytype.NotDefined")
	}
	meta.objectType = 5
	if meta.GetEntityType() != entitytype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> objectType==5 must return entitytype.NotDefined")
	}

	meta.objectType = 6
	if meta.GetEntityType() != entitytype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> objectType==6 must return entitytype.NotDefined")
	}
}

func Test__Meta__GetRelationType(t *testing.T) {
	var meta = new(Meta)

	// key
	meta.id = 12123
	meta.refId = 545434

	meta.setRelationType(4)
	if meta.GetRelationType() != relationtype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> relationType==4 must return entitytype.NotDefined")
	}
	meta.setRelationType(5)
	if meta.GetRelationType() != relationtype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> relationType==5 must return entitytype.NotDefined")
	}

	meta.setRelationType(6)
	if meta.GetRelationType() != relationtype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> relationType==6 must return entitytype.NotDefined")
	}
}

func Test__Meta__GetFieldType(t *testing.T) {
	var meta = new(Meta)

	// key
	meta.id = 12123
	meta.refId = 545434

	meta.dataType = 4
	if meta.GetFieldType() != fieldtype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> fieldType==4 must return entitytype.NotDefined")
	}
	meta.dataType = 5
	if meta.GetFieldType() != fieldtype.NotDefined {
		t.Errorf("Meta.GetEntityType() ==> fieldType==5 must return entitytype.NotDefined")
	}

	meta.dataType = 6
	if meta.GetFieldType() != fieldtype.NotDefined {
		t.Errorf("Meta.GetFieldType() ==> fieldType==6 must return entitytype.NotDefined")
	}
}

//test ToField, ToRelation, and ToIndex
func Test__Meta__ToField(t *testing.T) {
	var meta = new(Meta)

	// testing nil return
	meta.objectType = 4
	if meta.ToField() != nil {
		t.Errorf("Meta.ToField() ==> objectType==4 must return nil")
	}
	if meta.ToRelation() != nil {
		t.Errorf("Meta.ToField() ==> objectType==4 must return nil")
	}
	if meta.ToIndex() != nil {
		t.Errorf("Meta.ToField() ==> objectType==4 must return nil")
	}
}
