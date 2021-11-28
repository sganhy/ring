package schema

import (
	"time"
)

type Lexicon struct {
	id            int32
	schemaId      int32
	name          string
	description   string
	uuid          string
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

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
