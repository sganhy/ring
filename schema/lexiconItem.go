package schema

type lexiconItem struct {
	id        int32
	lexiconId int32
	refId     int32
	value     string
}

//******************************
// getters and setters
//******************************
func (lexiconItm *lexiconItem) GetId() int32 {
	return lexiconItm.id
}

func (lexiconItm *lexiconItem) GetlexiconId() int32 {
	return lexiconItm.lexiconId
}

func (lexiconItm *lexiconItem) GetRefId() int32 {
	return lexiconItm.refId
}

func (lexiconItm *lexiconItem) GetValue() string {
	return lexiconItm.value
}

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
