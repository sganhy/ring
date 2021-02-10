package data

import (
	"ring/schema/fieldtype"
	"testing"
)

func Test__Record__isValidInteger(t *testing.T) {
	var rcd = new(Record)
	rcd.SetRecordType("@meta")
	if isValidInteger("55451", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> 55451 is not a valid integer (32 bits)")
	}
	if isValidInteger("-55451", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> -55451 is not a valid integer (32 bits)")
	}
	if isValidInteger("0", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> 0 is not a valid integer (32 bits)")
	}
	if isValidInteger("2147483647", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> 2147483647 is not a valid integer (32 bits)")
	}
	if isValidInteger("2147483646", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> 2147483646 is not a valid integer (32 bits)")
	}
	if isValidInteger("-7483646", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> -7483646 is not a valid integer (32 bits)")
	}
	if isValidInteger("-2147483648", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> -2147483648 is not a valid integer (32 bits)")
	}
	if isValidInteger("-214748364", fieldtype.Int) == false {
		t.Errorf("Record.isValidInteger() ==> -214748364 is not a valid integer (32 bits)")
	}
	if isValidInteger("-2147483649", fieldtype.Int) != false {
		t.Errorf("Record.isValidInteger() ==> -2147483649 is not a valid integer (32 bits)")
	}
	if isValidInteger("", fieldtype.Int) != false {
		t.Errorf("Record.isValidInteger() ==> empty string is not a valid integer (32 bits)")
	}

}
