package schema

import (
	"testing"
)

//test: GetTableBySchemaName, GetSchemaByName, GetSchemaById
func Test__log__setters(t *testing.T) {
	logger := new(log)
	logger.setMethod("555")
	if logger.getMethod() != "555" {
		t.Errorf("log.setMethod() ==> method <> '555'")
	}

	// test right truncate max 80 char
	logger.setMethod("The hour of the reference time is 15, or 3PM. The layout can express it either way, ...")
	expectedText := "The hour of the reference time is 15, or 3PM. The layout can express it either w"
	if logger.getMethod() != expectedText {
		t.Errorf("log.setMethod() ==> method is not equal to '%s'", expectedText)
	}

	logger.setMessage("555")
	if logger.getMessage() != "555" {
		t.Errorf("log.setMessage() ==> message <> '555'")
	}
	longtext := "I have been fascinated with how the mind structures information for as long as I can remember. As a kid, my all-time favorite activity in middle school was diagramming sentences with their parts of speech. Perhaps it’s not surprising, then, that I ended up at MIT earning my doctorate on formal models of language and cognition."
	// set right truncate max 255
	logger.setMessage(longtext)
	if len(logger.getMessage()) > 255 {
		t.Errorf("log.setMessage() ==> len(logger.message) is greater than 255")
	}
	logger.setCallSite("555")
	if logger.getClassSite() != "555" {
		t.Errorf("log.setCallSite() ==> callSite <> '555'")
	}
	logger.setCallSite(longtext)
	if len(logger.getClassSite()) > 255 {
		t.Errorf("log.setCallSite() ==> len(logger.callSite) is greater than 255")
	}

}
