package schema

import (
	"testing"
)

//test: GetTableBySchemaName, GetSchemaByName, GetSchemaById
func Test__log__setters(t *testing.T) {
	logger := new(log)
	logger.setMethod("555")
	if logger.method != "555" {
		t.Errorf("log.setMethod() ==> method <> '555'")
	}

	// test left truncate max 80 char
	logger.setMethod("The hour of the reference time is 15, or 3PM. The layout can express it either way, ...")
	if logger.method != "The hour of the reference time is 15, or 3PM. The layout can express it either w" {
		t.Errorf("log.setMethod() ==> method is not equal to 'The hour of the reference time is 15, or 3PM. The layout can express it either w'")
	}
}
