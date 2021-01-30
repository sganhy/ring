package schema

type Sequence struct {
	id          int32
	name        string
	description string
	schemaId    int32
	maxValue    int64
	value       CacheId
}

func (sequence *Sequence) Init(id int32, name string, description string, schemaId int32,
	maxValue int64, value CacheId) {
	sequence.id = id
	sequence.name = name
	sequence.description = description
	sequence.schemaId = schemaId
	sequence.maxValue = maxValue
	sequence.value = value
}

//******************************
// getters
//******************************
func (sequence *Sequence) GetId() int32 {
	return sequence.id
}

func (sequence *Sequence) GetName() string {
	return sequence.name
}

func (sequence *Sequence) GetDescription() string {
	return sequence.description
}

func (sequence *Sequence) GetSchemaId() int32 {
	return sequence.schemaId
}

func (sequence *Sequence) GetMaxValue() int64 {
	return sequence.maxValue
}

func (sequence *Sequence) GetValue() CacheId {
	return sequence.value
}

//******************************
// public methods
//******************************
