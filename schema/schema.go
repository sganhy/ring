package schema

type Schema struct {
	id               int32
	name             string
	description      string
	connectionString string
	language         Language
	baseline         bool
	active          bool
}

func (schema *Schema) Init(id int32, name string, description string, connectionString string) {
	schema.id = id
	schema.name = name
	schema.description = description
	schema.connectionString = connectionString

}

//******************************
// getters
//******************************
func (schema *Schema) GetId() int32 {
	return schema.id
}

func (schema *Schema) GetName() string {
	return schema.name
}

func (schema *Schema) GetDescription() string {
	return schema.description
}

func (schema *Schema) GetConnectionString() string {
	return schema.connectionString
}

func (schema *Schema) IsBaseline() bool {
	return schema.baseline
}

func (schema *Schema) IsActive() bool {
	return schema.active
}
