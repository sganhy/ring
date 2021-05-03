package sqlfmt

import (
	"testing"
)

func Test__Sqlfmt__ToCamelCase(t *testing.T) {
	//======================
	//==== testing camelcase convention
	//======================
	expectedResult := "tableTesta"
	if ToCamelCase("TABLE_TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be equal to %s", expectedResult)
	}
	if ToCamelCase("TABLE-TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be equal to %s", expectedResult)

	}
	if ToCamelCase("TABLE TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be equal to %s", expectedResult)

	}
	if ToCamelCase("TABLE   -  TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be equal to %s", expectedResult)

	}
	if ToCamelCase("__TABLE-TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be equal to %s", expectedResult)
	}
	if ToCamelCase("__TABLE TeSTA") != expectedResult {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be equal to %s", expectedResult)
	}
	if ToCamelCase("") != "" {
		t.Errorf("Sqlfmt.ToCamelCase() ==> result must be empty")
	}
	//======================
	//==== testing Pascal convention
	//======================
	expectedResult = "TableTesta"
	if ToPascalCase("TABLE_TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be equal to %s", expectedResult)
	}
	if ToPascalCase("TABLE-TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be equal to %s", expectedResult)

	}
	if ToPascalCase("TABLE TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be equal to %s", expectedResult)

	}
	if ToPascalCase("TABLE   -  TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be equal to %s", expectedResult)

	}
	if ToPascalCase("__TABLE-TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be equal to %s", expectedResult)
	}
	if ToPascalCase("__TABLE TeSTA") != expectedResult {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be equal to %s", expectedResult)
	}
	if ToPascalCase("_") != "" {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be empty")
	}
	if ToPascalCase("") != "" {
		t.Errorf("Sqlfmt.ToPascalCase() ==> result must be empty")
	}

}