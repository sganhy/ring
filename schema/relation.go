package schema

import (
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/relationtype"
	"ring/schema/sqlfmt"
)

type Relation struct {
	id                  int32
	name                string
	description         string
	physicalName        string
	inverseRelationName string
	mtmTable            string
	toTable             *Table
	relationType        relationtype.RelationType
	notNull             bool
	baseline            bool
	active              bool
}

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
	// at the end ==>
	if toTable != nil {
		relation.physicalName = relation.getPhysicalName(toTable.provider)
	} else {
		relation.physicalName = name
	}
}

//******************************
// getters and setters
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

func (relation *Relation) GetMtmTableName() string {
	return relation.mtmTable
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

func (relation *Relation) GetPhysicalName() string {
	return relation.physicalName
}

//******************************
// public methods
//******************************
func (relation *Relation) GetDdl(provider databaseprovider.DatabaseProvider) string {
	if relation.toTable != nil {
		targetPrimaryKey := relation.toTable.GetPrimaryKey()
		if targetPrimaryKey != nil {
			datatype := targetPrimaryKey.getSqlDataType(provider)
			if datatype != unknownFieldDataType {
				return relation.name + ddlSpace + datatype
			}
		}
	}
	return ""
}

func (relation *Relation) GetInverseRelation() *Relation {
	if relation.toTable != nil {
		return relation.toTable.GetRelationByName(relation.inverseRelationName)
	}
	return nil
}

func (relation *Relation) Clone() *Relation {
	newRelation := new(Relation)
	/*
		id int32, name string, description string, inverseRelationName string,
			mtmTable string, toTable *Table, relationType relationType.RelationType, notNull bool, baseline bool, active bool
	*/
	// don't clone ToTable for reflexive relationship (recursive call)
	newRelation.Init(relation.id, relation.name, relation.description,
		relation.inverseRelationName, relation.mtmTable, relation.toTable, relation.relationType, relation.notNull, relation.baseline,
		relation.active)
	return newRelation
}

//******************************
// private methods
//******************************
func (relation *Relation) getPhysicalName(provider databaseprovider.DatabaseProvider) string {
	return sqlfmt.FormatEntityName(provider, relation.name)
}

func (relation *Relation) toMeta(tableId int32) *Meta {
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
	// add flags for relation type
	result.setRelationNotNull(relation.notNull)
	result.setRelationType(relation.relationType)

	result.value = relation.inverseRelationName
	result.name = relation.name // max length 30 !! must be validated before
	result.description = relation.description
	result.enabled = relation.active
	return result
}
