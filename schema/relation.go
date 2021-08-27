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
	id              int32
	name            string
	description     string
	mtmTable        *Table
	inverseRelation *Relation
	toTable         *Table
	relationType    relationtype.RelationType
	constraint      bool // add foreign key
	notNull         bool
	baseline        bool
	active          bool
}

const (
	relationToStringFormat string = "name=%s; description=%s; type=%s; to=%s; baseline=%t; not_null=%t; inverse_relation=%s"
	mtmTableNamePrefix     string = "@mtm"
	mtmSeperator           string = "_"
	mtmLeftPadding         string = "0"
)

func (relation *Relation) Init(id int32, name string, description string, toTable *Table, relationType relationtype.RelationType,
	constraint bool, notNull bool, baseline bool, active bool) {
	relation.id = id
	relation.name = name
	relation.description = description
	relation.toTable = toTable
	relation.relationType = relationType
	relation.notNull = notNull
	relation.baseline = baseline
	relation.constraint = constraint
	relation.active = active
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

func (relation *Relation) GetMtmTable() *Table {
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

func (relation *Relation) HasConstraint() bool {
	return relation.constraint
}

func (relation *Relation) GetEntityType() entitytype.EntityType {
	return entitytype.Relation
}

func (relation *Relation) GetInverseRelation() *Relation {
	return relation.inverseRelation
}

func (relation *Relation) setToTable(table *Table) {
	relation.toTable = table
}

func (relation *Relation) setInverseRelation(inverseRelation *Relation) {
	relation.inverseRelation = inverseRelation
}

func (relation *Relation) setMtmTable(table *Table) {
	relation.mtmTable = table
}

func (relation *Relation) setId(id int32) {
	relation.id = id
}

//******************************
// public methods
//******************************
func (relation *Relation) GetPhysicalName(provider databaseprovider.DatabaseProvider) string {
	return relation.getPhysicalName(provider, relation.name)
}

func (relation *Relation) GetDdl(provider databaseprovider.DatabaseProvider) string {
	if relation.toTable != nil {
		targetPrimaryKey := relation.toTable.GetPrimaryKey()
		if targetPrimaryKey != nil {
			datatype := targetPrimaryKey.getSqlDataType(provider)
			if datatype != unknownFieldDataType {
				return relation.GetPhysicalName(provider) + ddlSpace + datatype
			}
		}
	}
	return ""
}

func (relation *Relation) Clone() *Relation {
	newRelation := new(Relation)

	// don't clone ToTable for reflexive relationship (recursive call)
	newRelation.Init(relation.id, relation.name, relation.description, relation.toTable, relation.relationType,
		relation.constraint, relation.notNull, relation.baseline, relation.active)

	return newRelation
}

func (relation *Relation) String() string {
	//"name=%s; description=%s; type=%s; "
	return fmt.Sprintf(relationToStringFormat, relation.name, relation.description, relation.relationType.String(),
		relation.getToTableName(), relation.baseline, relation.notNull, relation.GetInverseRelation().GetName())
}

func (relation *Relation) GetTableId() int32 {
	if relation.inverseRelation != nil && relation.inverseRelation.toTable != nil {
		return relation.inverseRelation.toTable.GetId()
	}
	return -1
}

//******************************
// private methods
//******************************
func (relation *Relation) toField() *Field {
	if relation.toTable != nil {
		targetPrimaryKey := relation.toTable.GetPrimaryKey()
		if targetPrimaryKey != nil {
			var result = targetPrimaryKey.Clone()
			result.setName(relation.name)
			return result
		}
	}
	return nil
}

func (relation *Relation) getToTableName() string {
	if relation.toTable == nil {
		return ""
	}
	return relation.toTable.GetName()
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
	}

	result.flags = 0
	result.setEntityBaseline(relation.baseline)
	// add flags for relation type
	result.setRelationNotNull(relation.notNull)
	result.setRelationType(relation.relationType)

	if relation.inverseRelation != nil {
		result.value = relation.inverseRelation.name
	}

	result.name = relation.name // max length 30 !! must be validated before
	result.description = relation.description
	result.enabled = relation.active

	return result
}

func (relation *Relation) getMtmName(fromTableId int32) string {
	var b strings.Builder
	var toTableid = relation.toTable.GetId()
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
