package schema

type validator struct {
	importFile *Import
}

func (valid *validator) Init(importFile *Import) {
	valid.importFile = importFile
}

//******************************
// getters
//******************************

//******************************
// public methods
//******************************
