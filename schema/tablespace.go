package schema

type Tablespace struct {
	id          int32
	name        string
	description string
	filName     string
	tableName   string
	table       bool
	index       bool
}

func (tablespace *Tablespace) Init(id int32, name string, description string, tableName string, table bool, index bool) {
	tablespace.id = id
	tablespace.name = name
	tablespace.description = description
	tablespace.tableName = tableName
	tablespace.table = table
	tablespace.index = index
}

//******************************
// getters
//******************************

func (tablespace *Tablespace) GetId() int32 {
	return tablespace.id
}
func (tablespace *Tablespace) GetName() string {
	return tablespace.name
}
func (tablespace *Tablespace) GetDescription() string {
	return tablespace.description
}

//******************************
// public methods
//******************************
func (tablespace *Tablespace) Clone() *Tablespace {
	newTablespace := new(Tablespace)
	newTablespace.Init(tablespace.id, tablespace.name, tablespace.description,
		tablespace.tableName, tablespace.table, tablespace.index)
	return newTablespace
}
