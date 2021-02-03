package schema

import (
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"strings"
)

const maxInt32 int64 = 2147483647
const maxInt8 int64 = 127
const metaIndexSeparator string = ";"

type Meta struct {
	id          int32
	dataType    int32
	description string
	flags       uint64
	lineNumber  int64
	name        string
	objectType  int8
	refId       int32 // ref id to Id
	value       string
}

func (meta *Meta) Init(flags uint64) {
	meta.flags = flags

}

const (
	bitPositionFieldCaseSensitive   uint8 = 2
	bitPositionFieldNotNull         uint8 = 3
	bitPositionFieldMultilingual    uint8 = 4
	bitPositionIndexBitmap          uint8 = 9
	bitPositionIndexUnique          uint8 = 10
	bitPositionEntityEnabled        uint8 = 13
	bitPositionEntityBaseline       uint8 = 14
	bitPositionFirstPositionSize    uint8 = 17 // max value bit pos for field=16 !!!
	bitPositionFirstPositionRelType uint8 = 18 // max value bit pos for field=17 !!!
	bitPositionRelationNotNull      uint8 = 4  // max value bit pos for relation =17 !!!
)

//******************************
// getters
//******************************
func (meta *Meta) GetFlags() uint64 {
	return meta.flags
}
func (meta *Meta) GetId() int32 {
	return meta.id
}

//******************************
// public methods
//******************************
func (meta *Meta) GetFieldSize() uint32 {
	return uint32((meta.flags >> bitPositionFirstPositionSize) & uint64(maxInt32))
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

func (meta *Meta) IsEntityBaseline() bool {
	return meta.readFlag(bitPositionEntityBaseline)
}
func (meta *Meta) IsEntityEnabled() bool {
	return meta.readFlag(bitPositionEntityEnabled)
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
		// call exemple elemf.Init(21, "field ", "hellkzae", fieldtype.Double, 5, "", true, false, true, true)
		field.Init(meta.id, meta.name, meta.description,
			meta.GetFieldType(), meta.GetFieldSize(), meta.value,
			meta.IsEntityBaseline(), meta.IsFieldNotNull(), meta.IsFieldCaseSensitive(),
			meta.IsFieldMultilingual(), meta.IsEntityEnabled())
		return field
	}
	return nil
}

func (meta *Meta) ToRelation() *Relation {
	if entitytype.EntityType(meta.objectType) == entitytype.Relation {
		var relation = new(Relation)
		// call exemple 	elemr.Init(21, "rel test", "hellkzae", "hell1", "52", nil, relationtype.Mto, false, true, false)
		relation.Init(meta.id, meta.name, meta.description,
			meta.value, meta.value, nil, meta.GetRelationType(),
			meta.IsRelationNotNull(), meta.IsEntityBaseline(), meta.IsEntityEnabled())
		return relation
	}
	return nil
}

func (meta *Meta) ToIndex() *Index {
	if entitytype.EntityType(meta.objectType) == entitytype.Index {
		var index = new(Index)
		var arr = strings.Split(meta.value, metaIndexSeparator)
		// call exemple 	elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
		index.Init(meta.id, meta.name, meta.description, arr, meta.IsIndexBitmap(), meta.IsIndexUnique(),
			meta.IsEntityBaseline(), meta.IsEntityEnabled())
		return index
	}
	return nil
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
	var temp uint64 = uint64(size & uint32(maxInt32))
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
func (meta *Meta) setEntityEnabled(value bool) {
	meta.writeFlag(bitPositionEntityEnabled, value)
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

func (meta *Meta) setRelationType(relationType relationtype.RelationType) {
	var temp uint64 = uint64(uint32(relationType) & uint32(maxInt8))
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
		mask <<= (bitPosition - 1)
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
