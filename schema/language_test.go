package schema

import (
	"ring/schema/entitytype"
	"testing"
)

//test mappers Meta to Language to Meta
func Test__Language__toMeta(t *testing.T) {

	eleml0 := Language{}

	// id int32, name string, description string, parameterType entitytype.EntityType, fieldType fieldtype.FieldType, value string
	eleml0.Init(1, "es-MX")

	metaData := eleml0.toMeta()
	eleml1 := metaData.toLanguage()

	if eleml0.GetId() != eleml1.GetId() {
		t.Errorf("Language.toMeta() ==> l0.GetId() must be equal to l1.GetId()")
	}
	if eleml0.GetName() != eleml1.GetName() {
		t.Errorf("Language.toMeta() ==> l0.GetName() must be equal to l1.GetName()")
	}
	if eleml0.GetDescription() != eleml1.GetDescription() {
		t.Errorf("Language.toMeta() ==> l0.GetDescription() must be equal to l1.GetDescription()")
	}
	if eleml0.GetCode() != eleml1.GetCode() {
		t.Errorf("Language.toMeta() ==> l0.GetCode() must be equal to l1.GetCode()")
	}
	if eleml0.GetNativeName() != eleml1.GetNativeName() {
		t.Errorf("Language.toMeta() ==> l0.GetNativeName() must be equal to l1.GetNativeName()")
	}
	if eleml0.GetEntityType() != entitytype.Language {
		t.Errorf("Language.toMeta() ==> l0.GetEntityType() must be equal to entitytype.Language")
	}

	metaData.objectType = int8(entitytype.Constraint)
	eleml2 := metaData.toLanguage()
	if eleml2 != nil {
		t.Errorf("Language.toMeta() ==> l2 must be equal NULL")
	}

}

func Test__Language__IsCodeValid(t *testing.T) {
	eleml0 := Language{}

	//===== TEST 1
	valid, err := eleml0.IsCodeValid("e1-M1")
	if err == nil {
		t.Errorf("Language.IsCodeValid() ==> l0.err must not be equal to NULL")
	}
	if valid == true {
		t.Errorf("Language.IsCodeValid() ==> l0.valid must not be equal to false")
	}
	//===== TEST 2
	valid, err = eleml0.IsCodeValid("en-US")
	if err != nil {
		t.Errorf("Language.IsCodeValid() ==> l0.err must be equal to NULL")
	}
	if valid == false {
		t.Errorf("Language.IsCodeValid() ==> l0.valid must not be equal to true")
	}
	//===== TEST 3
	valid, err = eleml0.IsCodeValid("")
	if err == nil {
		t.Errorf("Language.IsCodeValid() ==> l0.err must not be equal to NULL")
	}
	if valid == true {
		t.Errorf("Language.IsCodeValid() ==> l0.valid must not be equal to false")
	}
	//===== TEST 4
	valid, err = eleml0.IsCodeValid("FR")
	if err != nil {
		t.Errorf("Language.IsCodeValid() ==> l0.err must be equal to NULL")
	}
	if valid == false {
		t.Errorf("Language.IsCodeValid() ==> l0.valid must not be equal to true")
	}
	//===== TEST 5
	valid, err = eleml0.IsCodeValid("F8")
	if err == nil {
		t.Errorf("Language.IsCodeValid() ==> l0.err must not be equal to NULL")
	}
	if valid == true {
		t.Errorf("Language.IsCodeValid() ==> l0.valid must not be equal to false")
	}
	//===== TEST 6
	valid, err = eleml0.IsCodeValid("uk-be")
	if err == nil {
		t.Errorf("Language.IsCodeValid() ==> l0.err must not be equal to NULL")
	}
	if valid == true {
		t.Errorf("Language.IsCodeValid() ==> l0.valid must not be equal to false")
	}
}
