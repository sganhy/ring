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

const (
	minJobIdValue     int64  = 101007
	maxJobIdValue     int64  = 9223372036854775807
	sequenceJobIdName string = "@job_id"
	sequenceLexId     string = "@lexicon_id"
)

func (sequence *Sequence) Init(id int32, name string, description string, schemaId int32, maxValue int64) {
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
func (sequence *Sequence) exists(schema *Schema) bool {
	return true
}

func (sequence *Sequence) create(schema *Schema) error {
	return nil
}

func (sequence *Sequence) getJobId(schemaId int32) *Sequence {
	// id int32, name string, description string, schemaId int32, maxValue int64
	result := new(Sequence)
	result.Init(0, sequenceJobIdName, "Unique job number assigned based on auto-numbering definition", schemaId, maxJobIdValue)
	result.active = true
	result.baseline = true
	result.value.CurrentId = minJobIdValue
	return result
}
