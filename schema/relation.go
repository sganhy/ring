package schema

import (
	"ring/schema/entitytype"
	"ring/schema/relationtype"
)

type Relation struct {
	id                  int32
	name                string
	description         string
	inverseRelationName string
	mtmTable            string
	toTable             *Table
	relationType        relationtype.RelationType
	notNull             bool
	baseline            bool
	active              bool
}

// call exemple elemr.Init(21, "rel test", "hellkzae", "hell1", "52", nil, relationtype.Mto, false, true, false)
func (relation *Relation) Init(id int32, name string, description string, inverseRelationName string,
	mtmTable string, toTable *Table, relationType relationtype.RelationType, notNull bool, baseline bool, active bool) {
	relation.id = id
	relation.name = name
	relation.description = description
	relation.inverseRelationName = inverseRelationName
	relation.mtmTable = mtmTable
	relation.toTable = toTable
	relation.relationType = relationType
	relation.notNull = notNull
	relation.baseline = baseline
	relation.active = active
}

//******************************
// getters
//******************************
func (relation *Relation) GetId() int32 {
	return relation.id
}

func (relation *Relation) GetName() string {
	return relation.name
}

func (relation *Relation) GetDescription() string {
	return relation.description
}

func (relation *Relation) GetInverseRelationName() string {
	return relation.inverseRelationName
}

func (relation *Relation) GetMtmTable() string {
	return relation.inverseRelationName
}

func (relation *Relation) GetToTable() *Table {
	return relation.toTable
}

func (relation *Relation) GetType() relationtype.RelationType {
	return relation.relationType
}

func (relation *Relation) IsNotNull() bool {
	return relation.notNull
}

func (relation *Relation) IsBaseline() bool {
	return relation.baseline
}

func (relation *Relation) IsActive() bool {
	return relation.active
}

//******************************
// public methods
//******************************
func (relation *Relation) GetInverseRelation() *Relation {
	return relation.toTable.GetRelationByName(relation.inverseRelationName)
}

func (relation *Relation) ToMeta(tableId int32) *Meta {
	// we cannot have error here
	var result = new(Meta)

	// key
	result.id = relation.id
	result.refId = tableId
	result.objectType = int8(entitytype.Relation)

	// others
	if relation.toTable != nil {
		result.dataType = relation.toTable.id // to table id
	} else {
		//TODO LOG
	}
	result.flags = 0
	result.setEntityBaseline(relation.baseline)
	// add flags for relation tyoe
	result.setEntityEnabled(relation.active)
	result.setRelationNotNull(relation.notNull)
	result.setRelationType(relation.relationType)

	result.value = relation.inverseRelationName
	result.name = relation.name // max lenght 30 !! must be valided before
	result.description = relation.description

	return result
}
