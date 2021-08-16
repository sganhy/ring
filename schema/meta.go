package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/tabletype"
	"strconv"
	"strings"
)

type meta struct {
	id          int32
	dataType    int32
	name        string
	description string
	flags       uint64
	lineNumber  int64
	objectType  int8
	refId       int32 // ref id to Id
	value       string
	enabled     bool
}

const (
	bitPositionFieldCaseSensitive     uint8  = 2
	bitPositionFieldNotNull           uint8  = 3
	bitPositionFieldMultilingual      uint8  = 4
	bitPositionIndexBitmap            uint8  = 9
	bitPositionIndexUnique            uint8  = 10
	bitPositionEntityBaseline         uint8  = 14
	bitPositionFirstPositionSize      uint8  = 17 // max value bit pos for field=16 !!!
	bitPositionFirstPositionRelType   uint8  = 18 // max value bit pos for field=17 !!!
	bitPositionRelationNotNull        uint8  = 4  // max value bit pos for relation =17 !!!
	bitPositionRelationConstraint     uint8  = 5
	bitPositionFirstPositionParamType uint8  = 14 // max value bit pos for field=16 !!!
	bitPositionTableCached            uint8  = 9
	bitPositionTableReadonly          uint8  = 10
	bitPositionTablespaceIndex        uint8  = 11
	bitPositionTablespaceTable        uint8  = 12
	metaMaxInt32                      int64  = 2147483647
	metaMaxInt8                       int64  = 127
	metaIndexSeparator                string = ";"
)

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (metaData *meta) GetFieldSize() uint32 {
	return uint32((metaData.flags >> bitPositionFirstPositionSize) & uint64(metaMaxInt32))
}
func (metaData *meta) IsFieldMultilingual() bool {
	return metaData.readFlag(bitPositionFieldMultilingual)
}
func (metaData *meta) IsFieldCaseSensitive() bool {
	return metaData.readFlag(bitPositionFieldCaseSensitive)
}
func (metaData *meta) IsFieldNotNull() bool {
	return metaData.readFlag(bitPositionFieldNotNull)
}
func (metaData *meta) IsEntityBaseline() bool {
	return metaData.readFlag(bitPositionEntityBaseline)
}
func (metaData *meta) IsRelationNotNull() bool {
	return metaData.readFlag(bitPositionRelationNotNull)
}
func (metaData *meta) IsRelationConstraint() bool {
	return metaData.readFlag(bitPositionRelationConstraint)
}
func (metaData *meta) IsIndexUnique() bool {
	return metaData.readFlag(bitPositionIndexUnique)
}
func (metaData *meta) IsIndexBitmap() bool {
	return metaData.readFlag(bitPositionIndexBitmap)
}
func (metaData *meta) IsTableCached() bool {
	return metaData.readFlag(bitPositionTableCached)
}

func (metaData *meta) IsTableReadonly() bool {
	return metaData.readFlag(bitPositionTableReadonly)
}

func (metaData *meta) IsTablespaceIndex() bool {
	return metaData.readFlag(bitPositionTablespaceIndex)
}

func (metaData *meta) IsTablespaceTable() bool {
	return metaData.readFlag(bitPositionTablespaceTable)
}

func (metaData *meta) GetFieldType() fieldtype.FieldType {
	return fieldtype.GetFieldTypeById(int(metaData.dataType & 127))
}
func (metaData *meta) GetEntityType() entitytype.EntityType {
	return entitytype.GetEntityTypeById(int(metaData.objectType & 127))
}

func (metaData *meta) GetRelationType() relationtype.RelationType {
	return relationtype.GetRelationTypeById(int((metaData.flags >> bitPositionFirstPositionRelType) & 127))
}

func (metaData *meta) String() string {
	// used for debug only
	return fmt.Sprintf("id: %d; name: %s; object_type: %d; reference_id: %d; dataType: %d; flags: %d; value: %s; line_number: %d; description: %s",
		metaData.id, metaData.name, metaData.objectType, metaData.refId, metaData.dataType, metaData.flags,
		metaData.value, metaData.lineNumber, metaData.description)
}

func (metaData *meta) GetParameterType() entitytype.EntityType {
	return entitytype.GetEntityTypeById(int((metaData.flags >> bitPositionFirstPositionParamType) & 127))
}

//******************************
// private methods
//******************************

// mappers
func (metaData *meta) toField() *Field {
	if metaData.GetEntityType() == entitytype.Field {
		var field = new(Field)
		field.Init(metaData.id, metaData.name, metaData.description,
			metaData.GetFieldType(), metaData.GetFieldSize(), metaData.value,
			metaData.IsEntityBaseline(), metaData.IsFieldNotNull(), metaData.IsFieldCaseSensitive(),
			metaData.IsFieldMultilingual(), metaData.enabled)
		return field
	}
	return nil
}

func (metaData *meta) toRelation(table *Table) *Relation {
	if metaData.GetEntityType() == entitytype.Relation {
		var relation = new(Relation)
		relation.Init(metaData.id, metaData.name, metaData.description, table, metaData.GetRelationType(),
			metaData.IsRelationConstraint(), metaData.IsRelationNotNull(), metaData.IsEntityBaseline(), metaData.enabled)
		return relation
	}
	return nil
}

func (metaData *meta) toIndex() *Index {
	if metaData.GetEntityType() == entitytype.Index {
		var index = new(Index)
		var arr = strings.Split(metaData.value, metaIndexSeparator)
		index.Init(metaData.id, metaData.name, metaData.description, arr, metaData.IsIndexBitmap(), metaData.IsIndexUnique(),
			metaData.IsEntityBaseline(), metaData.enabled)
		return index
	}
	return nil
}

func (metaData *meta) toTablespace() *tablespace {
	if metaData.GetEntityType() == entitytype.Tablespace {
		var tableSpace = new(tablespace)
		// id int32, name string, description string, fileName string, table bool, index bool
		tableSpace.Init(metaData.id, metaData.name, metaData.description, metaData.value, metaData.IsTablespaceTable(), metaData.IsTablespaceIndex())
		return tableSpace
	}
	return nil
}

func (metaData *meta) toSequence(schemaId int32) *Sequence {
	if metaData.GetEntityType() == entitytype.Sequence {
		var sequence = new(Sequence)
		// id int32, name string, description string, schemaId int32, maxValue int64, baseline bool, active bool
		value, _ := strconv.ParseInt(metaData.value, 10, 64)
		sequence.Init(metaData.id, metaData.name, metaData.description, schemaId, value,
			metaData.IsEntityBaseline(), metaData.enabled)
		return sequence
	}
	return nil
}

// loosing schemaId and databaseprovider
func (metaData *meta) toTable(fields []Field, relations []Relation, indexes []Index) *Table {
	if metaData.GetEntityType() == entitytype.Table {
		var table = new(Table)
		/*
			t1.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
			physicaltype.Table, -111, metaSchemaName, tabletype.Fake, databaseprovider.NotDefined, "", true, false, true, true)
		*/
		table.Init(metaData.id, metaData.name, metaData.description, fields, relations, indexes,
			physicaltype.Table, 0, metaSchemaName, tabletype.Business,
			databaseprovider.NotDefined, metaData.value, metaData.IsTableCached(), metaData.IsTableReadonly(), metaData.IsEntityBaseline(),
			metaData.enabled)
		return table
	}
	return nil
}

func (metaData *meta) toParameter(schemaId int32) *parameter {
	if metaData.GetEntityType() == entitytype.Parameter {
		var param = new(parameter)
		var entityType = metaData.GetParameterType()
		var fieldType = fieldtype.FieldType(int(metaData.dataType))
		param.Init(metaData.id, metaData.name, metaData.description, schemaId, entityType, fieldType, metaData.value)
		return param
	}
	return nil
}

func (metaData *meta) toLanguage() *Language {
	if metaData.GetEntityType() == entitytype.Language {
		var lang = new(Language)
		lang.Init(metaData.id, metaData.value)
		return lang
	}
	return nil
}

func (metaA *meta) equal(metaB *meta) bool {
	if metaA != nil && metaB != nil {
		/* ==> don't compare key
		dataType    int32
		name        string
		description string
		flags       uint64
		value       string
		*/
		if metaA.dataType == metaB.dataType && metaA.name == metaB.name &&
			metaA.description == metaB.description && metaA.flags == metaB.flags &&
			metaA.value == metaB.value {
			return true
		}
		return false
	}
	return metaA == metaB
}

// flags
func (metaData *meta) setFieldMultilingual(value bool) {
	metaData.writeFlag(bitPositionFieldMultilingual, value)
}
func (metaData *meta) setFieldNotNull(value bool) {
	metaData.writeFlag(bitPositionFieldNotNull, value)
}
func (metaData *meta) setFieldCaseSensitive(value bool) {
	metaData.writeFlag(bitPositionFieldCaseSensitive, value)
}
func (metaData *meta) setFieldSize(size uint32) {
	var temp = uint64(size & uint32(metaMaxInt32))
	// maxInt32 & size << ()
	// reset flags 16 first bits using  65.535
	metaData.flags &= 65535
	temp <<= bitPositionFirstPositionSize
	//temp = metaData.flags & temp // reset size to 0;
	metaData.flags += temp
}

func (metaData *meta) setParameterType(parameterType entitytype.EntityType) {
	var temp = uint64(parameterType)
	temp <<= bitPositionFirstPositionParamType
	metaData.flags &= 65535
	metaData.flags += temp
}

func (metaData *meta) setEntityBaseline(value bool) {
	metaData.writeFlag(bitPositionEntityBaseline, value)
}
func (metaData *meta) setRelationNotNull(value bool) {
	metaData.writeFlag(bitPositionRelationNotNull, value)
}
func (metaData *meta) setRelationConstraint(value bool) {
	metaData.writeFlag(bitPositionRelationConstraint, value)
}
func (metaData *meta) setIndexBitmap(value bool) {
	metaData.writeFlag(bitPositionIndexBitmap, value)
}
func (metaData *meta) setIndexUnique(value bool) {
	metaData.writeFlag(bitPositionIndexUnique, value)
}
func (metaData *meta) setTableCached(value bool) {
	metaData.writeFlag(bitPositionTableCached, value)
}
func (metaData *meta) setTableReadonly(value bool) {
	metaData.writeFlag(bitPositionTableReadonly, value)
}
func (metaData *meta) setTablespaceIndex(value bool) {
	metaData.writeFlag(bitPositionTablespaceIndex, value)
}
func (metaData *meta) setTablespaceTable(value bool) {
	metaData.writeFlag(bitPositionTablespaceTable, value)
}

func (metaData *meta) setRelationType(relationType relationtype.RelationType) {
	var temp = uint64(uint32(relationType) & uint32(metaMaxInt8))
	// maxInt32 & size << ()
	// reset flags 16 first bits using  65.535
	metaData.flags &= 65535
	temp <<= bitPositionFirstPositionRelType
	//temp = metaData.flags & temp // reset size to 0;
	metaData.flags += temp
}

// bit position: ]1,64[
func (metaData *meta) writeFlag(bitPosition uint8, value bool) {
	var mask uint64 = 0
	if bitPosition < 64 {
		mask = 1
		mask <<= bitPosition - 1
		if value == true {
			metaData.flags |= mask
		} else {
			metaData.flags &= ^mask
		}
	}
}

func (metaData *meta) readFlag(bitPosition uint8) bool {
	return ((metaData.flags >> (bitPosition - 1)) & 1) > 0
}

func (metaData *meta) saveMetaList(schemaId int32, metaList []*meta) error {
	prevMetaList := getMetaList(schemaId)
	var err error

	err = nil
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	// create new schema - separated case to reduce allocation
	if len(prevMetaList) == 0 {
		for i := 0; i < len(metaList) && err == nil; i++ {
			err = query.insertMeta(metaList[i], schemaId)
		}
	} else {
		// create dictionary : object_type_id, [reference_id], [upper(name)]
		dico := metaData.getMetaDictionary(prevMetaList)

		//disallow slice
		prevMetaList = []meta{}

		// metaList should have valid objectTypeId >=0 && objectTypeId <= entitytype.MaxEntityTypeId
		err = metaData.updateMetaList(schemaId, metaList, dico)

		// detect deleted meta item
		if err == nil {
			err = metaData.deleteMetaList(schemaId, dico, true)
		}
	}

	return err
}

func (metaData *meta) updateMetaList(schemaId int32, metaList []*meta, dico []map[int32]map[string]*meta) error {
	var err error

	err = nil
	queryInsert := new(metaQuery)
	queryInsert.setSchema(metaSchemaName)
	queryInsert.setTable(metaTableName)

	queryUpdate := new(metaQuery)
	queryUpdate.setSchema(metaSchemaName)
	queryUpdate.setTable(metaTableName)

	for i := 0; i < len(metaList); i++ {
		var currMeta = metaList[i]
		var prevMeta = dico[currMeta.objectType][currMeta.refId][strings.ToUpper(currMeta.name)]

		if prevMeta == nil {
			// insert
			/*
				fmt.Println("== INSERT META ==")
				fmt.Println(currMeta.String())
			*/
			metaData.setMetaDataId(currMeta, dico)
			err = queryInsert.insertMeta(currMeta, schemaId)
		} else {
			//update
			prevMeta.enabled = true // flag this meta is in the prev schema
			if prevMeta.equal(currMeta) == false {
				/*
					fmt.Println("== UPDATE META ==")
					fmt.Println(currMeta.String())
				*/
				queryUpdate.updateMeta(currMeta, schemaId)
			}
		}
		if err != nil {
			break
		}
	}

	return err
}

func (metaData *meta) deleteMetaList(schemaId int32, dico []map[int32]map[string]*meta, forceDelete bool) error {
	//current schema:  dico
	//warning: one type of dml by metaQuery object !!!!
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	for i := 0; i < len(dico); i++ {
		currMap := dico[i]
		for _, mapId := range currMap {
			for _, mapMeta := range mapId {
				if mapMeta.enabled == false {
					/*
						fmt.Println("== DELETE META ==")
						fmt.Println(mapMeta.String())
					*/
					if forceDelete == false {
						query.updateMeta(mapMeta, schemaId)
					} else {
						query.deleteMeta(mapMeta, schemaId)
					}
				}
			}
		}
	}

	return nil
}

func (metaData *meta) setMetaDataId(newMeta *meta, dico []map[int32]map[string]*meta) {
	// dictionary : object_type_id, [reference_id], [upper(name)]

	// cannot change id of tables
	if newMeta.GetEntityType() == entitytype.Table {
		return
	}

	// object type exists always
	if dico[newMeta.objectType][newMeta.refId] == nil {
		dico[newMeta.objectType][newMeta.refId] = make(map[string]*meta)
	}
	var subMap = dico[newMeta.objectType][newMeta.refId]
	var maxId int32 = 0

	for _, item := range subMap {
		if item.id > maxId {
			maxId = item.id
		}
	}

	maxId++
	newMeta.id = maxId

	dico[newMeta.objectType][newMeta.refId][strings.ToUpper(newMeta.name)] = newMeta
}

func (metaData *meta) getMetaDictionary(currentMetaList []meta) []map[int32]map[string]*meta {
	result := make([]map[int32]map[string]*meta, entitytype.MaxEntityTypeId+1, entitytype.MaxEntityTypeId+1)

	for i := 0; i <= entitytype.MaxEntityTypeId; i++ {
		result[i] = make(map[int32]map[string]*meta)
	}

	for i := 0; i < len(currentMetaList); i++ {
		var currMeta = &currentMetaList[i]
		if currMeta.objectType >= 0 && currMeta.objectType <= entitytype.MaxEntityTypeId {
			if result[currMeta.objectType][currMeta.refId] == nil {
				result[currMeta.objectType][currMeta.refId] = make(map[string]*meta)
			}
			currMeta.enabled = false // set false to detect later if item deleted
			result[currMeta.objectType][currMeta.refId][strings.ToUpper(currMeta.name)] = currMeta
		}
	}

	return result
}
