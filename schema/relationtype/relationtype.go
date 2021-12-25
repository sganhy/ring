package relationtype

import (
	"strings"
)

type RelationType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	strOtop string = "otop"
	strOtm  string = "otm"
	strMtm  string = "mtm"
	strMto  string = "mto"
	strOtof string = "otof"
)

const (
	Otop      RelationType = 1
	Otm       RelationType = 2
	Mtm       RelationType = 3
	Mto       RelationType = 11
	Otof      RelationType = 12
	Undefined RelationType = 123
)

func (relationType RelationType) String() string {
	switch relationType {
	case Otop:
		return strings.ToUpper(strOtop)
	case Otm:
		return strings.ToUpper(strOtm)
	case Mtm:
		return strings.ToUpper(strMtm)
	case Mto:
		return strings.ToUpper(strMto)
	case Otof:
		return strings.ToUpper(strOtof)
	}
	return ""
}

func (relationType RelationType) InverseRelationType() RelationType {
	switch relationType {
	case Otop:
		return Otof
	case Otm:
		return Mto
	case Mtm:
		return Mtm
	case Mto:
		return Otm
	case Otof:
		return Otop
	}
	return Undefined
}

func GetRelationType(value string) RelationType {
	switch strings.ToLower(value) {
	case strings.ToLower(strOtop):
		return Otop
	case strings.ToLower(strOtm):
		return Otm
	case strings.ToLower(strMtm):
		return Mtm
	case strings.ToLower(strMto):
		return Mto
	case strings.ToLower(strOtof):
		return Otof
	}
	return Undefined
}

func GetRelationTypeById(entityId int) RelationType {
	if entityId <= 127 && entityId >= -128 {
		var newId = RelationType(entityId)
		if newId == Otop || newId == Otm || newId == Mtm || newId == Mto || newId == Otof {
			return newId
		}
	}
	return Undefined

}
