package schema

import (
	"fmt"
	"ring/schema/entitytype"
	"strconv"
)

type Sequence struct {
	id          int32
	name        string
	description string
	schemaId    int32
	maxValue    int64
	value       *cacheId
	baseline    bool
	active      bool
}

const (
	maxJobIdValue        int64  = 9223372036854775807
	maxLexIdValue        int64  = 2147483647
	maxLangIdValue       int64  = 32767
	maxUserIdValue       int64  = 2147483647
	initialJobId         int64  = 101007
	sequenceJobIdName    string = "@job_id"
	sequenceLexId        string = "@lexicon_id"
	sequenceUserId       string = "@user_id"
	sequenceIndexId      string = "@index_id"
	sequenceEventId      string = "@event_id"
	sequenceStringFormat string = "id=%d; name=%s; description=%s; schemaId=%d; maxValue=%d"
)

func (sequence *Sequence) Init(id int32, name string, description string, schemaId int32, maxValue int64, baseline bool, active bool) {
	sequence.id = id
	sequence.name = name
	sequence.description = description
	sequence.schemaId = schemaId
	sequence.maxValue = maxValue
	sequence.value = new(cacheId)
	sequence.value.Init(id, schemaId, sequence.GetEntityType())
	sequence.baseline = baseline
	sequence.active = active
}

//******************************
// getters and setters
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

func (sequence *Sequence) GetEntityType() entitytype.EntityType {
	return entitytype.Sequence
}
func (sequence *Sequence) getCacheId() *cacheId {
	return sequence.value
}

//******************************
// public methods
//******************************
func (sequence *Sequence) NextValue() int64 {
	return sequence.value.GetNewId()
}

func (sequence *Sequence) Clone() *Sequence {
	// id int32, name string, description string, schemaId int32, maxValue int64, baseline bool, active bool
	newSequence := new(Sequence)
	newSequence.Init(sequence.id, sequence.name, sequence.description, sequence.schemaId, 0, sequence.baseline, sequence.active)
	// reference copy
	newSequence.value = sequence.value
	return sequence
}

func (sequence *Sequence) String() string {
	//	id=%d; name=%s; description=%s; schemaId=%d; maxValue=%d
	return fmt.Sprintf(sequenceStringFormat, sequence.id, sequence.name, sequence.description, sequence.schemaId,
		sequence.schemaId)
}

//******************************
// private methods
//******************************
func (sequence *Sequence) exists() bool {
	query := new(metaQuery)

	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)
	query.addFilter(metaFieldId, operatorEqual, sequence.id)
	query.addFilter(metaSchemaId, operatorEqual, sequence.schemaId)
	query.addFilter(metaObjectType, operatorEqual, int8(entitytype.Sequence))
	query.addFilter(metaReferenceId, operatorEqual, sequence.schemaId)
	result, _ := query.exists()

	return result
}

func (sequence *Sequence) create() error {

	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	return query.insertMeta(sequence.toMeta(), sequence.schemaId)
}

func (sequence *Sequence) getJobId(schemaId int32) *Sequence {
	result := new(Sequence)
	result.Init(0, sequenceJobIdName, "Unique job number assigned based on auto-numbering definition", schemaId, maxJobIdValue,
		true, true)
	result.value.Init(0, schemaId, entitytype.Sequence)
	result.value.SetCurrentId(initialJobId) // assign min value
	return result
}

func (sequence *Sequence) getLexiconId(schemaId int32) *Sequence {
	result := new(Sequence)
	result.Init(1, sequenceLexId, "Unique lexicon number assigned based on auto-numbering definition", schemaId, maxLexIdValue, true, true)
	result.value.SetCurrentId(103) // assign min value
	return result
}

func (sequence *Sequence) getUserId(schemaId int32) *Sequence {
	result := new(Sequence)
	result.Init(3, sequenceUserId, "Unique user number assigned based on auto-numbering definition", schemaId, maxUserIdValue, true, true)
	result.value.SetCurrentId(1003) // assign min value
	return result
}

func (sequence *Sequence) getIndexId(schemaId int32) *Sequence {
	result := new(Sequence)
	result.Init(4, sequenceIndexId, "Unique index number assigned based on auto-numbering definition", schemaId, maxUserIdValue, true, true)
	result.value.SetCurrentId(1105)
	return result
}

func (sequence *Sequence) getEventId(schemaId int32) *Sequence {
	result := new(Sequence)
	result.Init(5, sequenceEventId, "Unique event number assigned based on auto-numbering definition", schemaId, maxUserIdValue, true, true)
	result.value.SetCurrentId(198)
	return result
}

func (sequence *Sequence) toMeta() *meta {
	var metaTable = new(meta)

	// key
	metaTable.id = sequence.id
	metaTable.refId = sequence.schemaId
	metaTable.objectType = int8(entitytype.Sequence)

	// others
	metaTable.dataType = 0
	metaTable.name = sequence.name // max length 30 !! must be validated before
	metaTable.description = sequence.description
	metaTable.value = strconv.FormatInt(sequence.maxValue, 10)

	// flags
	metaTable.flags = 0
	metaTable.setEntityBaseline(sequence.baseline)
	metaTable.enabled = sequence.active

	return metaTable
}
