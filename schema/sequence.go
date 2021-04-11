package schema

type Sequence struct {
	id          int32
	name        string
	description string
	schemaId    int32
	maxValue    int64
	value       *CacheId
	baseline    bool
	active      bool
}

const sequenceJobId = "@job_id"
const sequenceLexId = "@lexicon_id"

func (sequence *Sequence) Init(id int32, name string, description string, schemaId int32,
	maxValue int64) {
	sequence.id = id
	sequence.name = name
	sequence.description = description
	sequence.schemaId = schemaId
	sequence.maxValue = maxValue
	sequence.value = new(CacheId)
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

func (sequence *Sequence) GetValue() *CacheId {
	return sequence.value
}

//******************************
// public methods
//******************************

//******************************
// private methods
//******************************
func (sequence *Sequence) getJobId() *Sequence {
	result := new(Sequence)
	result.id = 0
	result.name = sequenceJobId
	result.description = "Unique job number assigned based on auto-numbering definition"
	result.baseline = true
	result.active = true
	result.value = new(CacheId)
	result.value.CurrentId = 101007
	return result
}

func (sequence *Sequence) exist() bool {
	return true
}
