package schema

import (
	"ring/schema/ddlstatement"
	"time"
)

type Lexicon struct {
	id            int32
	name          string
	description   string
	uuid          string
	schemaId      int32
	tableId       int32
	fromFieldId   int32
	toFieldId     int32
	relationId    int32
	relationValue int64
	lastUpdate    *time.Time
	enabled       bool
}

/*
func (lexicon *Lexicon) Init(id int32, code string) {

}
*/

//******************************
// getters and setters
//******************************
func (lexicon *Lexicon) GetId() int32 {
	return lexicon.id
}

func (lexicon *Lexicon) GetName() string {
	return lexicon.name
}

func (lexicon *Lexicon) GetDescription() string {
	return lexicon.description
}

func (lexicon *Lexicon) GetPhysicalName() string {
	return lexicon.GetName()
}

func (lexicon *Lexicon) GetUuid() string {
	return lexicon.uuid
}

func (lexicon *Lexicon) GetSchemaId() int32 {
	return lexicon.schemaId
}

func (lexicon *Lexicon) GetTableId() int32 {
	return lexicon.tableId
}

func (lexicon *Lexicon) logStatement(statment ddlstatement.DdlStatement) bool {
	return true
}

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
