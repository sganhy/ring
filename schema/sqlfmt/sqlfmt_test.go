package sqlfmt

import (
	"ring/schema/databaseprovider"
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
	//======================
	//==== testing Snake convention
	//======================
	expectedResult = "table_testa"
	if ToSnakeCase("TableTesta") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("tableTesta") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("table Testa") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("_TableTesta") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("Table_Testa") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("Table____Testa") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("Table  Testa") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("Table      Testa") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("TABLE_TESTA") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("Table_testa") != expectedResult {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be equal to %s", expectedResult)
	}
	if ToSnakeCase("") != "" {
		t.Errorf("Sqlfmt.ToSnakeCase() ==> result must be empty")
	}

}

func Test__Sqlfmt__FormatEntityName(t *testing.T) {
	//======================
	//==== testing postgreSQL
	//======================
	provider := databaseprovider.PostgreSql
	expectedResult := "table_testa"
	if FormatEntityName(provider, "TableTesta") != expectedResult {
		t.Errorf("Sqlfmt.FormatEntityName() ==> result must be equal to %s", expectedResult)
	}

	expectedResult = "\"@meta\""
	t.Errorf(FormatEntityName(provider, "@meta"))
	if FormatEntityName(provider, "@meta") != expectedResult {
		t.Errorf("Sqlfmt.FormatEntityName() ==> result must be equal to %s", expectedResult)
	}

}
