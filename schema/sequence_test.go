package schema

import (
	"ring/schema/entitytype"
	"testing"
)

// INIT
func Test__Sequence__Init(t *testing.T) {
	var sequence = Sequence{}

	//Init(id int32, name string, description string, schemaId int32, maxValue int64, baseline bool, active bool)
	sequence.Init(212, "test name", "test description", 3, 100000, true, true)

	if sequence.GetName() != "test name" {
		t.Errorf("Sequence.Init() ==> name <> GetName()")
	}
	if sequence.GetId() != 212 {
		t.Errorf("Sequence.Init() ==> id <> GetId()")
	}
	if sequence.GetDescription() != "test description" {
		t.Errorf("Sequence.Init() ==> description <> GetDescription()")
	}
	if sequence.GetSchemaId() != 3 {
		t.Errorf("Sequence.Init() ==> schemaId <> GetSchemaId()")
	}
	if sequence.GetMaxValue() != 100000 {
		t.Errorf("Sequence.Init() ==> maxValue <> GetMaxValue()")
	}
	if sequence.GetEntityType() != entitytype.Sequence {
		t.Errorf("Sequence.Init() ==> entitytype.Sequence <> GetEntityType()")
	}
}

func Test__Sequence__toMeta(t *testing.T) {
	var sequence01 = Sequence{}
	sequence01.Init(212, "test name", "test description", 3, 100000, true, true)

	metaData := sequence01.toMeta()
	sequence02 := metaData.toSequence(3)
	if sequence01.GetName() != sequence02.GetName() {
		t.Errorf("Sequence.Init() ==> s1.GetName() must be equal to s2.GetName()")
	}
	if sequence01.GetId() != sequence02.GetId() {
		t.Errorf("Sequence.toMeta() ==> s1.GetName() must be equal to s2.GetName()")
	}
	if sequence01.GetDescription() != sequence02.GetDescription() {
		t.Errorf("Sequence.toMeta() ==> s1.GetDescription() must be equal to s2.GetDescription()")
	}
	if sequence01.GetSchemaId() != sequence02.GetSchemaId() {
		t.Errorf("Sequence.toMeta() ==> s1.GetSchemaId() must be equal to s2.GetSchemaId()")
	}
	if sequence01.GetMaxValue() != sequence02.GetMaxValue() {
		t.Errorf("Sequence.toMeta() ==> s1.GetMaxValue() must be equal to s2.GetMaxValue()")
	}
	if sequence01.GetEntityType() != sequence02.GetEntityType() {
		t.Errorf("Sequence.toMeta() ==> s1.GetEntityType() must be equal to s2.GetEntityType()")
	}

}

func Test__Sequence__Clone(t *testing.T) {
	var sequence01 = Sequence{}
	sequence01.Init(212, "test name", "test description", 3, 100000, true, true)
	sequence02 := sequence01.Clone()

	if sequence01.GetName() != sequence02.GetName() {
		t.Errorf("Sequence.Clone() ==> s1.GetName() must be equal to s2.GetName()")
	}
	if sequence01.GetId() != sequence02.GetId() {
		t.Errorf("Sequence.Clone() ==> s1.GetName() must be equal to s2.GetName()")
	}
	if sequence01.GetDescription() != sequence02.GetDescription() {
		t.Errorf("Sequence.Clone() ==> s1.GetDescription() must be equal to s2.GetDescription()")
	}
	if sequence01.GetSchemaId() != sequence02.GetSchemaId() {
		t.Errorf("Sequence.Clone() ==> s1.GetSchemaId() must be equal to s2.GetSchemaId()")
	}
	if sequence01.GetMaxValue() != sequence02.GetMaxValue() {
		t.Errorf("Sequence.Clone() ==> s1.GetMaxValue() must be equal to s2.GetMaxValue()")
	}
	if sequence01.GetEntityType() != sequence02.GetEntityType() {
		t.Errorf("Sequence.Clone() ==> s1.GetEntityType() must be equal to s2.GetEntityType()")
	}
}
