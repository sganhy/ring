package schema

import (
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"testing"
)

//test mappers Meta to parameter to Meta
func Test__Parameter__toMeta(t *testing.T) {

	elemp0 := parameter{}

	// id int32, name string, description string, parameterType entitytype.EntityType, fieldType fieldtype.FieldType, value string
	elemp0.Init(111778, "test parameter", "test description", 7777, entitytype.Constraint, fieldtype.Boolean, "value test!!")

	metaData := elemp0.toMeta(7777)
	elemp1 := metaData.toParameter(7777)

	if elemp0.GetId() != elemp1.GetId() {
		t.Errorf("Parameter.toMeta() ==> p0.GetId() must be equal to p1.GetId()")
	}
	if elemp0.GetName() != elemp1.GetName() {
		t.Errorf("Parameter.toMeta() ==> p0.GetName() must be equal to p1.GetName()")
	}
	if elemp0.GetDescription() != elemp1.GetDescription() {
		t.Errorf("Parameter.toMeta() ==> p0.GetDescription() must be equal to p1.GetDescription()")
	}
	if elemp0.GetDescription() != elemp1.GetDescription() {
		t.Errorf("Parameter.toMeta() ==> p0.GetDescription() must be equal to p1.GetDescription()")
	}
	if elemp0.GetEntityType() != elemp1.GetEntityType() {
		t.Errorf("Parameter.toMeta() ==> p0.GetEntityType() must be equal to p1.GetEntityType()")
	}
	if elemp0.GetDataType() != elemp1.GetDataType() {
		t.Errorf("Parameter.toMeta() ==> p0.GetDataType() must be equal to p1.GetDataType()")
	}
	if elemp0.GetValue() != elemp1.GetValue() {
		t.Errorf("Parameter.toMeta() ==> p0.GetValue() must be equal to p1.GetValue()")
	}

	elemp0.Init(111778, "test parameter", "test description", 5555, entitytype.Constraint, fieldtype.Undefined, "value test!!")
	if elemp0.GetDataType() != fieldtype.String {
		t.Errorf("Parameter.toMeta() ==> p0.GetDataType() must be equal to fieldtype.String")
	}

	elemp0.setValue("5564654")
	elemp1.setValue("5564654")
	if elemp0.GetValue() != elemp1.GetValue() {
		t.Errorf("Parameter.toMeta() ==> p0.GetValue() must be equal to p1.GetValue()")
	}
}

func Test__Parameter__Clone(t *testing.T) {
	elemp0 := parameter{}

	// id int32, name string, description string, parameterType entitytype.EntityType, fieldType fieldtype.FieldType, value string
	elemp0.Init(1111, "test parameter", "test description", 7777, entitytype.Constraint, fieldtype.Boolean, "value test!!")
	elemp1 := elemp0.Clone()

	if elemp0.GetId() != elemp1.GetId() {
		t.Errorf("Parameter.Clone() ==> p0.GetId() must be equal to p1.GetId()")
	}
	if elemp0.GetName() != elemp1.GetName() {
		t.Errorf("Parameter.Clone() ==> p0.GetName() must be equal to p1.GetName()")
	}
	if elemp0.GetDescription() != elemp1.GetDescription() {
		t.Errorf("Parameter.Clone() ==> p0.GetDescription() must be equal to p1.GetDescription()")
	}
	if elemp0.GetDescription() != elemp1.GetDescription() {
		t.Errorf("Parameter.Clone() ==> p0.GetDescription() must be equal to p1.GetDescription()")
	}
	if elemp0.GetEntityType() != elemp1.GetEntityType() {
		t.Errorf("Parameter.Clone() ==> p0.GetEntityType() must be equal to p1.GetEntityType()")
	}
	if elemp0.GetDataType() != elemp1.GetDataType() {
		t.Errorf("Parameter.Clone() ==> p0.GetDataType() must be equal to p1.GetDataType()")
	}
	if elemp0.GetValue() != elemp1.GetValue() {
		t.Errorf("Parameter.Clone() ==> p0.GetValue() must be equal to p1.GetValue()")
	}
}
