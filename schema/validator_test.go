package schema

import (
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
