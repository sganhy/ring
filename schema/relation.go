package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/relationtype"
	"ring/schema/sqlfmt"
	"strconv"
	"strings"
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

const (
	relationToStringFormat string = "name=%s; description=%s; type=%s; to=%s; baseline=%t; not_null=%t; inverse_relation=%s"
	mtmTableNamePrefix     string = "@mtm"
	mtmSeperator           string = "_"
	mtmLeftPadding         string = "0"
)

func (relation *Relation) Init(id int32, name string, description string, inverseRelationName string,
	toTable *Table, relationType relationtype.RelationType, notNull bool, baseline bool, active bool) {
	relation.id = id
	relation.name = name
	relation.description = description
	relation.inverseRelationName = inverseRelationName
	relation.toTable = toTable
	relation.relationType = relationType
	relation.notNull = notNull
	relation.baseline = baseline
	relation.active = active
	// at the end ==>
	if toTable != nil {
		relation.physicalName = relation.getPhysicalName(toTable.GetDatabaseProvider(), name)
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

func (relation *Relation) GetEntityType() entitytype.EntityType {
	return entitytype.Relation
}

func (relation *Relation) setToTable(table *Table) {
	relation.toTable = table
}

func (relation *Relation) GetInverseRelation() *Relation {
	if relation.toTable != nil {
		return relation.toTable.GetRelationByName(relation.inverseRelationName)
	}
	return nil
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

func (relation *Relation) Clone() *Relation {
	newRelation := new(Relation)
	/*
		id int32, name string, description string, inverseRelationName string,
			mtmTable string, toTable *Table, relationType relationType.RelationType, notNull bool, baseline bool, active bool
	*/
	// don't clone ToTable for reflexive relationship (recursive call)
	newRelation.Init(relation.id, relation.name, relation.description,
		relation.inverseRelationName, relation.toTable, relation.relationType, relation.notNull, relation.baseline,
		relation.active)
	return newRelation
}

func (relation *Relation) String() string {
	//"name=%s; description=%s; type=%s; "
	return fmt.Sprintf(relationToStringFormat, relation.name, relation.description, relation.relationType.String(),
		relation.getToTableName(), relation.baseline, relation.notNull, relation.GetInverseRelationName())
}

//******************************
// private methods
//******************************
func (relation *Relation) getToTableName() string {
	if relation.toTable == nil {
		return ""
	}
	return relation.toTable.name
}

func (relation *Relation) getPhysicalName(provider databaseprovider.DatabaseProvider, name string) string {
	return sqlfmt.FormatEntityName(provider, name)
}

func (relation *Relation) toMeta(tableId int32) *meta {
	// we cannot have error here
	var result = new(meta)

	// key
	result.id = relation.id
	result.refId = tableId
	result.objectType = int8(entitytype.Relation)

	// others
	if relation.toTable != nil {
		result.dataType = relation.toTable.GetId() // to table id
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

func (relation *Relation) loadMtm(fromTableId int32) {
	if relation.relationType == relationtype.Mtm && relation.toTable != nil && relation.mtmTable == "" {
		relation.mtmTable = relation.getMtmName(fromTableId)
	}
}

func (relation *Relation) getMtmName(fromTableId int32) string {
	var b strings.Builder
	var toTableid = relation.toTable.id
	var strTo = strconv.FormatInt(int64(toTableid), 10)
	var strFrom = strconv.FormatInt(int64(fromTableId), 10)
	var strRelId string
	var inverseRel = relation.GetInverseRelation()

	b.Grow(30)
	b.WriteString(mtmTableNamePrefix)
	b.WriteString(mtmSeperator)

	if fromTableId < toTableid {
		strRelId = strconv.FormatInt(int64(relation.id), 10)
		b.WriteString(sqlfmt.PadLeft(strFrom, mtmLeftPadding, 5))
		b.WriteString(mtmSeperator)
		b.WriteString(sqlfmt.PadLeft(strTo, mtmLeftPadding, 5))
	} else {
		strRelId = strconv.FormatInt(int64(inverseRel.id), 10)
		b.WriteString(sqlfmt.PadLeft(strTo, mtmLeftPadding, 5))
		b.WriteString(mtmSeperator)
		b.WriteString(sqlfmt.PadLeft(strFrom, mtmLeftPadding, 5))
	}
	b.WriteString(mtmSeperator)

	if fromTableId != toTableid {
		b.WriteString(sqlfmt.PadLeft(strRelId, mtmLeftPadding, 3))
	} else {
		b.WriteString(sqlfmt.PadLeft(relation.getMinId(inverseRel.id, relation.id), mtmLeftPadding, 3))
	}

	return b.String()
}

func (relation *Relation) getMinId(valA int32, valB int32) string {
	if valA < valB {
		return strconv.FormatInt(int64(valA), 10)
	}
	return strconv.FormatInt(int64(valB), 10)
}
