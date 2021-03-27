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

const metaMaxInt32 int64 = 2147483647
const metaMaxInt8 int64 = 127
const metaIndexSeparator string = ";"

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
	bitPositionFieldCaseSensitive   uint8 = 2
	bitPositionFieldNotNull         uint8 = 3
	bitPositionFieldMultilingual    uint8 = 4
	bitPositionIndexBitmap          uint8 = 9
	bitPositionIndexUnique          uint8 = 10
	bitPositionEntityBaseline       uint8 = 14
	bitPositionFirstPositionSize    uint8 = 17 // max value bit pos for field=16 !!!
	bitPositionFirstPositionRelType uint8 = 18 // max value bit pos for field=17 !!!
	bitPositionRelationNotNull      uint8 = 4  // max value bit pos for relation =17 !!!
	bitPositionTableCached          uint8 = 9
	bitPositionTableReadonly        uint8 = 10
)

//******************************
// getters
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
func (meta *Meta) GetFieldType() fieldtype.FieldType {
	var result = fieldtype.FieldType(meta.dataType & 127)
	if result != fieldtype.Array &&
		result != fieldtype.Boolean &&
		result != fieldtype.Byte &&
		result != fieldtype.DateTime &&
		result != fieldtype.Double &&
		result != fieldtype.Float &&
		result != fieldtype.Int &&
		result != fieldtype.Long &&
		result != fieldtype.LongDateTime &&
		result != fieldtype.Short &&
		result != fieldtype.ShortDateTime &&
		result != fieldtype.String {
		result = fieldtype.NotDefined
	}
	return result
}
func (meta *Meta) GetEntityType() entitytype.EntityType {
	var result = entitytype.EntityType(meta.objectType & 127)
	if result != entitytype.Field &&
		result != entitytype.Relation &&
		result != entitytype.Index &&
		result != entitytype.Schema &&
		result != entitytype.Sequence &&
		result != entitytype.Table {
		result = entitytype.NotDefined
	}
	return result
}

func (meta *Meta) GetRelationType() relationtype.RelationType {
	var result = relationtype.RelationType((meta.flags >> bitPositionFirstPositionRelType) & 127)
	if result != relationtype.Mtm &&
		result != relationtype.Mto &&
		result != relationtype.Otm &&
		result != relationtype.Otof &&
		result != relationtype.Otop {
		result = relationtype.NotDefined
	}
	return result
}

// mappers
func (meta *Meta) ToField() *Field {
	if entitytype.EntityType(meta.objectType) == entitytype.Field {
		var field = new(Field)
		field.Init(meta.id, meta.name, meta.description,
			meta.GetFieldType(), meta.GetFieldSize(), meta.value,
			meta.IsEntityBaseline(), meta.IsFieldNotNull(), meta.IsFieldCaseSensitive(),
			meta.IsFieldMultilingual(), meta.enabled)
		return field
	}
	return nil
}

func (meta *Meta) ToRelation(table *Table) *Relation {
	if entitytype.EntityType(meta.objectType) == entitytype.Relation {
		var relation = new(Relation)
		relation.Init(meta.id, meta.name, meta.description,
			meta.value, meta.value, table, meta.GetRelationType(),
			meta.IsRelationNotNull(), meta.IsEntityBaseline(), meta.enabled)
		return relation
	}
	return nil
}

func (meta *Meta) ToIndex() *Index {
	if entitytype.EntityType(meta.objectType) == entitytype.Index {
		var index = new(Index)
		var arr = strings.Split(meta.value, metaIndexSeparator)
		index.Init(meta.id, meta.name, meta.description, arr, meta.IsIndexBitmap(), meta.IsIndexUnique(),
			meta.IsEntityBaseline(), meta.enabled)
		return index
	}
	return nil
}

func (meta *Meta) ToTable(fields []Field, relations []Relation, indexes []Index) *Table {
	if entitytype.EntityType(meta.objectType) == entitytype.Table {
		var table = new(Table)
		/*
			t1.Init(1154, "@meta", "ATable Test", fields, relations, indexes,
			physicaltype.Table, -111, metaSchemaName, tabletype.Fake, databaseprovider.NotDefined, "", true, false, true, true)
		*/
		table.Init(meta.id, meta.name, meta.name, fields, relations, indexes,
			physicaltype.Table, 0, metaSchemaName, tabletype.Business,
			databaseprovider.NotDefined, "", true, false, true, true)
		return table
	}
	return nil
}

// Build partial schema object
func (meta *Meta) ToSchema() *Schema {
	if entitytype.EntityType(meta.objectType) == entitytype.Schema {
		var schema = new(Schema)
		schema.id = meta.id
		schema.name = meta.name
		schema.description = meta.description
		return schema
	}
	return nil
}

func (meta *Meta) String() string {
	// used for debug only
	/*
		id          int32
		dataType    int32
		description string
		flags       uint64
		lineNumber  int64
		name        string
		objectType  int8
		refId       int32 // ref id to Id
		value       string
		enabled     bool
	*/
	return fmt.Sprintf(" id: %d; name: %s; object_type: %d; reference_id: %d; dataType: %d; flags: %d; value: %s; description: %s",
		meta.id, meta.name, meta.objectType, meta.refId, meta.dataType, meta.flags, meta.value, meta.description)
}

//******************************
// private methods
//******************************

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
			fieldsMap[metaList[i].refId] = append(fieldsMap[metaList[i].refId], *metaList[i].ToField())
			break
		case entitytype.Relation:
			relationsMap[metaList[i].refId] = append(relationsMap[metaList[i].refId], *metaList[i].ToRelation(nil))
			break
		case entitytype.Index:
			indexesMap[metaList[i].refId] = append(indexesMap[metaList[i].refId], *metaList[i].ToIndex())
			break
		}
	}
	return fieldsMap, relationsMap, indexesMap
}

func getTables(schemaId int32, metaList []Meta) []Table {
	var objectType entitytype.EntityType
	fieldsMap, relationsMap, IndexesMap := initTableMappers(metaList)
	result := []Table{}

	for i := 0; i < len(metaList); i++ {
		objectType = entitytype.EntityType(metaList[i].objectType)
		if objectType == entitytype.Table {
			var tableId = metaList[i].id
			var table = *metaList[i].ToTable(fieldsMap[tableId], relationsMap[tableId], IndexesMap[tableId])
			// define
			table.schemaId = schemaId
			//table.provider = schemaId
			result = append(result, table)
		}
	}
	return result
}
