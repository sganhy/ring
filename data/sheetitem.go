package data

type sheetItem struct {
	id   int
	name string
}

func (item *sheetItem) Init(id int, name string) {
	item.id = id
	item.name = name
}

//******************************
// getters / setters
//******************************

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
