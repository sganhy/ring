package schema

import (
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/physicaltype"
	"ring/schema/relationtype"
	"ring/schema/tabletype"
	"strings"
)

type Meta struct {
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
	bitPositionFieldCaseSensitive    uint8  = 2
	bitPositionFieldNotNull          uint8  = 3
	bitPositionFieldMultilingual     uint8  = 4
	bitPositionIndexBitmap           uint8  = 9
	bitPositionIndexUnique           uint8  = 10
	bitPositionEntityBaseline        uint8  = 14
	bitPositionFirstPositionSize     uint8  = 17 // max value bit pos for field=16 !!!
	bitPositionFirstPositionRelType  uint8  = 18 // max value bit pos for field=17 !!!
	bitPositionRelationNotNull       uint8  = 4  // max value bit pos for relation =17 !!!
	bitPositionFirstPositionDataType uint8  = 14 // max value bit pos for field=16 !!!
	bitPositionTableCached           uint8  = 9
	bitPositionTableReadonly         uint8  = 10
	bitPositionTablespaceIndex       uint8  = 11
	bitPositionTablespaceTable       uint8  = 12
	metaMaxInt32                     int64  = 2147483647
	metaMaxInt8                      int64  = 127
	metaIndexSeparator               string = ";"
)

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (meta *Meta) GetFieldSize() uint32 {
	return uint32((meta.flags >> bitPositionFirstPositionSize) & uint64(metaMaxInt32))
}
func (meta *Meta) IsFieldMultilingual() bool {
	return meta.readFlag(bitPositionFieldMultilingual)
}
func (meta *Meta) IsFieldCaseSensitive() bool {
	return meta.readFlag(bitPositionFieldCaseSensitive)
}
func (meta *Meta) IsFieldNotNull() bool {
	return meta.readFlag(bitPositionFieldNotNull)
}
func (meta *Meta) IsEntityBaseline() bool {
	return meta.readFlag(bitPositionEntityBaseline)
}
func (meta *Meta) IsRelationNotNull() bool {
	return meta.readFlag(bitPositionRelationNotNull)
}
func (meta *Meta) IsIndexUnique() bool {
	return meta.readFlag(bitPositionIndexUnique)
}
func (meta *Meta) IsIndexBitmap() bool {
	return meta.readFlag(bitPositionIndexBitmap)
}
func (meta *Meta) IsTableCached() bool {
	return meta.readFlag(bitPositionTableCached)
}

func (meta *Meta) IsTableReadonly() bool {
	return meta.readFlag(bitPositionTableReadonly)
}

func (meta *Meta) IsTablespaceIndex() bool {
	return meta.readFlag(bitPositionTablespaceIndex)
}

func (meta *Meta) IsTablespaceTable() bool {
	return meta.readFlag(bitPositionTablespaceTable)
}

func (meta *Meta) GetFieldType() fieldtype.FieldType {
	return fieldtype.GetFieldTypeById(int(meta.dataType & 127))
}
func (meta *Meta) GetEntityType() entitytype.EntityType {
	return entitytype.GetEntityTypeById(int(meta.objectType & 127))
}

func (meta *Meta) GetRelationType() relationtype.RelationType {
	return relationtype.GetRelationTypeById(int((meta.flags >> bitPositionFirstPositionRelType) & 127))
}

func (meta *Meta) String() string {
	// used for debug only
	return fmt.Sprintf("id: %d; name: %s; object_type: %d; reference_id: %d; dataType: %d; flags: %d; value: %s; line_number: %d; description: %s",
		meta.id, meta.name, meta.objectType, meta.refId, meta.dataType, meta.flags, meta.value, meta.lineNumber, meta.description)
}

func (meta *Meta) GetParameterType() fieldtype.FieldType {
	return fieldtype.GetFieldTypeById(int((meta.flags >> bitPositionFirstPositionDataType) & 127))
}

//******************************
// private methods
//******************************

// mappers
func (meta *Meta) toField() *Field {
	if meta.GetEntityType() == entitytype.Field {
		var field = new(Field)
		field.Init(meta.id, meta.name, meta.description,
			meta.GetFieldType(), meta.GetFieldSize(), meta.value,
			meta.IsEntityBaseline(), meta.IsFieldNotNull(), meta.IsFieldCaseSensitive(),
			meta.IsFieldMultilingual(), meta.enabled)
		return field
	}
	return nil
}

func (meta *Meta) toRelation(table *Table) *Relation {
	if meta.GetEntityType() == entitytype.Relation {
		var relation = new(Relation)
		relation.Init(meta.id, meta.name, meta.description,
			meta.value, meta.value, table, meta.GetRelationType(),
			meta.IsRelationNotNull(), meta.IsEntityBaseline(), meta.enabled)
		return relation
	}
	return nil
}

func (meta *Meta) toIndex() *Index {
	if meta.GetEntityType() == entitytype.Index {
		var index = new(Index)
		var arr = strings.Split(meta.value, metaIndexSeparator)
		index.Init(meta.id, meta.name, meta.description, arr, meta.refId, meta.IsIndexBitmap(), meta.IsIndexUnique(),
			meta.IsEntityBaseline(), meta.enabled)
		return index
	}
	return nil
}

func (meta *Meta) toTablespace() *Tablespace {
	if meta.GetEntityType() == entitytype.Tablespace {
		var tablespace = new(Tablespace)
		// Init(id int32, name string, description string, fileName string, table bool, index bool) {
		tablespace.Init(meta.id, meta.name, meta.description, meta.value, meta.IsTablespaceTable(), meta.IsTablespaceIndex())
		return tablespace
	}
	return nil
}

// loosing schemaId and databaseprovider
func (meta *Meta) toTable(fields []Field, relations []Relation, indexes []Index) *Table {
	if meta.GetEntityType() == entitytype.Table {
		var table = new(Table)
		/*
			t1.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
			physicaltype.Table, -111, metaSchemaName, tabletype.Fake, databaseprovider.NotDefined, "", true, false, true, true)
		*/
		table.Init(meta.id, meta.name, meta.description, fields, relations, indexes,
			physicaltype.Table, 0, metaSchemaName, tabletype.Business,
			databaseprovider.NotDefined, meta.value, meta.IsTableCached(), meta.IsTableReadonly(), meta.IsEntityBaseline(),
			meta.enabled)
		return table
	}
	return nil
}

// Build partial schema object
func (meta *Meta) toSchema() *Schema {
	if meta.GetEntityType() == entitytype.Schema {
		var schema = new(Schema)
		schema.id = meta.id
		schema.name = meta.name
		schema.description = meta.description
		return schema
	}
	return nil
}

func (meta *Meta) toParameter() *parameter {
	if meta.GetEntityType() == entitytype.Parameter {
		var param = new(parameter)
		var entityType = entitytype.GetEntityTypeById(int(meta.dataType))
		var fieldType = meta.GetParameterType()
		param.Init(meta.id, meta.name, meta.description, entityType, fieldType, meta.value)
		return param
	}
	return nil
}

func (meta *Meta) toLanguage() *Language {
	if meta.GetEntityType() == entitytype.Language {
		var lang = new(Language)
		lang.Init(meta.id, meta.value)
		return lang
	}
	return nil
}

// flags
func (meta *Meta) setFieldMultilingual(value bool) {
	meta.writeFlag(bitPositionFieldMultilingual, value)
}
func (meta *Meta) setFieldNotNull(value bool) {
	meta.writeFlag(bitPositionFieldNotNull, value)
}
func (meta *Meta) setFieldCaseSensitive(value bool) {
	meta.writeFlag(bitPositionFieldCaseSensitive, value)
}
func (meta *Meta) setFieldSize(size uint32) {
	var temp = uint64(size & uint32(metaMaxInt32))
	// maxInt32 & size << ()
	// reset flags 16 first bits using  65.535
	meta.flags &= 65535
	temp <<= bitPositionFirstPositionSize
	//temp = meta.flags & temp // reset size to 0;
	meta.flags += temp
}

func (meta *Meta) setParameterDataType(fieldType fieldtype.FieldType) {
	var temp = uint64(fieldType)
	temp <<= bitPositionFirstPositionDataType
	meta.flags &= 65535
	meta.flags += temp
}

func (meta *Meta) setEntityBaseline(value bool) {
	meta.writeFlag(bitPositionEntityBaseline, value)
}
func (meta *Meta) setRelationNotNull(value bool) {
	meta.writeFlag(bitPositionRelationNotNull, value)
}

func (meta *Meta) setIndexBitmap(value bool) {
	meta.writeFlag(bitPositionIndexBitmap, value)
}
func (meta *Meta) setIndexUnique(value bool) {
	meta.writeFlag(bitPositionIndexUnique, value)
}
func (meta *Meta) setTableCached(value bool) {
	meta.writeFlag(bitPositionTableCached, value)
}
func (meta *Meta) setTableReadonly(value bool) {
	meta.writeFlag(bitPositionTableReadonly, value)
}
func (meta *Meta) setTablespaceIndex(value bool) {
	meta.writeFlag(bitPositionTablespaceIndex, value)
}
func (meta *Meta) setTablespaceTable(value bool) {
	meta.writeFlag(bitPositionTablespaceTable, value)
}

func (meta *Meta) setRelationType(relationType relationtype.RelationType) {
	var temp = uint64(uint32(relationType) & uint32(metaMaxInt8))
	// maxInt32 & size << ()
	// reset flags 16 first bits using  65.535
	meta.flags &= 65535
	temp <<= bitPositionFirstPositionRelType
	//temp = meta.flags & temp // reset size to 0;
	meta.flags += temp
}

// bit position: ]1,64[
func (meta *Meta) writeFlag(bitPosition uint8, value bool) {
	var mask uint64 = 0
	if bitPosition < 64 {
		mask = 1
		mask <<= bitPosition - 1
		if value == true {
			meta.flags |= mask
		} else {
			meta.flags &= ^mask
		}
	}
}

func (meta *Meta) readFlag(bitPosition uint8) bool {
	return ((meta.flags >> (bitPosition - 1)) & 1) > 0
}

func initTableMappers(metaList []Meta) (map[int32][]Field, map[int32][]Relation, map[int32][]Index) {
	var objectType entitytype.EntityType
	fieldsMap := make(map[int32][]Field)
	relationsMap := make(map[int32][]Relation)
	indexesMap := make(map[int32][]Index)

	// pass 1
	for i := 0; i < len(metaList); i++ {
		objectType = entitytype.EntityType(metaList[i].objectType)
		if objectType == entitytype.Table {
			tableId := metaList[i].id
			fieldsMap[tableId] = make([]Field, 0, 4)
			relationsMap[tableId] = make([]Relation, 0, 4)
			indexesMap[tableId] = make([]Index, 0, 2)
		}
	}

	// pass 2
	for i := 0; i < len(metaList); i++ {
		objectType = entitytype.EntityType(metaList[i].objectType)
		switch objectType {
		case entitytype.Field:
			fieldsMap[metaList[i].refId] = append(fieldsMap[metaList[i].refId], *metaList[i].toField())
			break
		case entitytype.Relation:
			relationsMap[metaList[i].refId] = append(relationsMap[metaList[i].refId], *metaList[i].toRelation(nil))
			break
		case entitytype.Index:
			indexesMap[metaList[i].refId] = append(indexesMap[metaList[i].refId], *metaList[i].toIndex())
			break
		}
	}
	return fieldsMap, relationsMap, indexesMap
}

func getTables(schema Schema, metaList []Meta) []Table {
	var objectType entitytype.EntityType
	fieldsMap, relationsMap, IndexesMap := initTableMappers(metaList)
	result := []Table{}

	for i := 0; i < len(metaList); i++ {
		objectType = entitytype.EntityType(metaList[i].objectType)
		if objectType == entitytype.Table {
			var tableId = metaList[i].id
			var table = *metaList[i].toTable(fieldsMap[tableId], relationsMap[tableId], IndexesMap[tableId])
			//TO define
			table.schemaId = schema.id
			//table.provider = schema.
			result = append(result, table)
		}
	}
	return result
}
