package schema

type Tablespace struct {
	schemaId  int32
	name      string
	filName   string
	tableName string
	table     bool
	index     bool
}

//******************************
// getters
//******************************

func (tablespace *Tablespace) GetName() string {
	return tablespace.name
}
