package schema

import (
	"ring/schema/entitytype"
	"strings"
)

type Index struct {
	id          int32
	name        string
	description string
	fields      []string
	bitmap      bool
	unique      bool
	baseline    bool
	active      bool
}

// call exemple elemi.Init(21, "rel test", "hellkzae", aarr, 52, false, true, true, true)
func (index *Index) Init(id int32, name string, description string, fields []string, bitmap bool,
	unique bool, baseline bool, active bool) {
	index.id = id
	index.name = name
	index.description = description
	index.loadFields(fields)
	index.bitmap = bitmap
	index.unique = unique
	index.baseline = baseline
	index.active = active
}

//******************************
// getters
//******************************
func (index *Index) GetId() int32 {
	return index.id
}

func (index *Index) GetName() string {
	return index.name
}

func (index *Index) GetDescription() string {
	return index.description
}

func (index *Index) GetFields() []string {
	return index.fields
}

func (index *Index) IsUnique() bool {
	return index.unique
}

func (index *Index) IsBitmap() bool {
	return index.bitmap
}

func (index *Index) IsBaseline() bool {
	return index.baseline
}

func (index *Index) IsActive() bool {
	return index.active
}

func (index *Index) ToMeta(tableId int32) *Meta {
	// we cannot have error here
	var result = new(Meta)

	// key
	result.id = index.id
	result.refId = tableId
	result.objectType = int8(entitytype.Index)

	// others
	result.dataType = 0
	result.name = index.name // max lenght 30 !! must be valided before
	result.description = index.description
	result.value = strings.Join(index.fields, metaIndexSeparator)

	// flags
	result.flags = 0
	result.setIndexBitmap(index.bitmap)
	result.setIndexUnique(index.unique)
	result.setEntityEnabled(index.active)
	result.setEntityBaseline(index.baseline)

	return result
}

//******************************
// public methods
//******************************
func (index *Index) Clone() *Index {
	newIndex := new(Index)
	/*
		id int32, name string, description string, fields []string, bitmap bool,
			unique bool, baseline bool, active bool
	*/
	fields := []string{}

	for i := 0; i < len(index.fields); i++ {
		fields = append(fields, index.fields[i])
	}
	newIndex.Init(index.id, index.name, index.description,
		fields, index.bitmap, index.unique, index.baseline, index.active)
	return newIndex
}

//******************************
// private methods
//******************************
func (index *Index) loadFields(fields []string) {
	// copy slice -- func make([]T, len, cap) []T
	if fields != nil {
		index.fields = make([]string, 0, len(fields))
		for i := 0; i < len(fields); i++ {
			index.fields = append(index.fields, fields[i])
		}
	} else {
		//TODO throw an error
		index.fields = make([]string, 0, 1)
	}
}
