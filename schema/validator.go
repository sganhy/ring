package schema

import (
	"errors"
	"fmt"
	"ring/schema/entitytype"
	"strconv"
	"strings"
	"unicode"
)

const (
	mssgDuplicateTableId      string = "Duplicate table Id '%d'"
	descDuplicateTableId      string = "table id '%d' is already in use in tables %s."
	mssgDuplicateTableName    string = "Duplicate table name '%s'"
	descDuplicateTableName    string = "table name '%s' is already in use.\n at line %d and %d"
	mssgEmptyEntityName       string = "Invalid %s name"
	descEmptyEntityName       string = "%s name cannot be empty." + validatorAtLine
	mssgInvalidEntityName     string = mssgEmptyEntityName
	descInvalidEntityName     string = "invalid %s name '%s'. A name can consist of any combination of letters(A to Z a to z), decimal digits(0 to 9) or underscore (_)." + validatorAtLine
	mssgMaxLengthEntityName   string = mssgEmptyEntityName
	descMaxLengthEntityName   string = "%s name '%s' is too long (max length=%d)." + validatorAtLine
	validatorCharSeparator    string = ", "
	validatorAtLine           string = "\n at line %d"
	joinOperationByName       int    = 1
	joinOperationByLineNumber int    = 2
	prefixedEntityMaxLength   int    = 28
	unPrefixedEntityMaxLength int    = 30
)

type validator struct {
	errorCount int
}

func (valid *validator) Init() {
	valid.errorCount = 0
}

//******************************
// getters
//******************************

//******************************
// public methods
//******************************

func (valid *validator) ValidateImport(importFile *Import) bool {
	fmt.Println("Start validation import ==> ")
	if importFile.metaList == nil || len(importFile.metaList) <= 0 {
		importFile.logError(errors.New("empty metaList"))
		importFile.errorCount++
		valid.errorCount++
		return false
	}

	valid.tableIdUnique(importFile)
	valid.tableNameUnique(importFile)
	valid.entityNameValid(importFile)
	valid.entityNameUnique(importFile)
	valid.languageCodeValid(importFile)

	// final check if all is ok
	if importFile.errorCount == 0 {
		valid.duplicateMetaKey(importFile)
	}

	return true
}

//******************************
// private methods
//******************************
func (valid *validator) tableIdUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dico map[int32][]*Meta
	var ok bool
	var val []*Meta

	dico = make(map[int32][]*Meta)
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table {
			if val, ok = dico[meta.id]; ok {
			} else {
				// new slice
				val = make([]*Meta, 0, 2)
			}
			val = append(val, meta)
			dico[meta.id] = val
		}
	}
	for key, arr := range dico {
		if len(arr) > 1 {
			var message = fmt.Sprintf(mssgDuplicateTableId, key)
			var description = fmt.Sprintf(descDuplicateTableId, key, valid.joinMeta(arr, joinOperationByName))
			description += validatorAtLine + valid.joinMeta(arr, joinOperationByLineNumber)
			importFile.logErrorStr(703, message, description)
		}
	}
}

func (valid *validator) tableNameUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dico map[string]*Meta
	dico = make(map[string]*Meta)

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table {
			name := strings.ToUpper(meta.name)
			if val, ok := dico[name]; ok {
				var message = fmt.Sprintf(mssgDuplicateTableName, meta.name)
				var description = fmt.Sprintf(descDuplicateTableName, meta.name, val.lineNumber, meta.lineNumber)
				importFile.logErrorStr(704, message, description)
			} else {
				dico[name] = meta
			}
		}
	}
}

func (valid *validator) entityNameValid(importFile *Import) {
	var metaList = importFile.metaList

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		metaType := meta.GetEntityType()

		// is meta.name empty?
		if len(meta.name) == 0 {
			var message = fmt.Sprintf(mssgEmptyEntityName, strings.ToLower(metaType.String()))
			var description = fmt.Sprintf(descEmptyEntityName, strings.ToLower(metaType.String()), meta.lineNumber)
			importFile.logErrorStr(502, message, description)
		} else {
			valid.checkEntityName(importFile, meta, metaType)
		}
	}
}

func (valid *validator) checkEntityName(importFile *Import, meta *Meta, metaType entitytype.EntityType) {
	// is name valid
	if valid.isValidName(meta.name) == false &&
		(metaType == entitytype.Table || metaType == entitytype.Field || metaType == entitytype.Index ||
			metaType == entitytype.Relation || metaType == entitytype.Tablespace) {
		var message = fmt.Sprintf(mssgInvalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descInvalidEntityName, strings.ToLower(metaType.String()), meta.name, meta.lineNumber)
		importFile.logErrorStr(501, message, description)
	}
	// is meta.name len > 28
	if len(meta.name) > prefixedEntityMaxLength && (metaType == entitytype.Table || metaType == entitytype.Field) {
		var message = fmt.Sprintf(mssgMaxLengthEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descMaxLengthEntityName, strings.ToLower(metaType.String()), meta.name,
			prefixedEntityMaxLength, meta.lineNumber)
		importFile.logErrorStr(503, message, description)
	}
	// is meta.name len > 30
	if len(meta.name) > unPrefixedEntityMaxLength && metaType != entitytype.Table && metaType != entitytype.Field {
		var message = fmt.Sprintf(mssgMaxLengthEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descMaxLengthEntityName, strings.ToLower(metaType.String()), meta.name,
			unPrefixedEntityMaxLength, meta.lineNumber)
		importFile.logErrorStr(504, message, description)
	}
}

func (valid *validator) duplicateMetaKey(importFile *Import) {
	var metaList = importFile.metaList
	var dicoEntities map[string]bool
	dicoEntities = make(map[string]bool, len(metaList))

	// check on db unique key (pk_@meta) ==> id|schema_id|object_type|reference_id
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]

		metaKey := strconv.FormatInt(int64(meta.id), 16) + "-" +
			strconv.FormatInt(int64(meta.refId), 16) + "-" +
			strconv.FormatInt(int64(meta.objectType), 16)

		if _, ok := dicoEntities[metaKey]; ok {
			// error duplicate meta key
			var message = fmt.Sprintf("Duplicate meta key")
			var description = fmt.Sprintf("duplicate meta key (type=%s): refid=%d, id=%d", strings.ToLower(meta.GetEntityType().String()),
				meta.refId, meta.id)

			importFile.logErrorStr(527, message, description)
		} else {
			dicoEntities[metaKey] = true
		}
	}
}

func (valid *validator) entityNameUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dicoTable map[int32]string
	dicoTable = make(map[int32]string)
	// (1) build table dictionary
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table {
			dicoTable[meta.id] = meta.name
		}
	}
	// (2) build dictionary entity & check

}

func (valid *validator) joinMeta(metaList []*Meta, operation int) string {
	var result strings.Builder
	for i := 0; i < len(metaList); i++ {
		switch operation {
		case joinOperationByName:
			result.WriteString(metaList[i].name)
			break
		case joinOperationByLineNumber:
			result.WriteString(strconv.FormatInt(metaList[i].lineNumber, 10))
			break
		}
		if i < len(metaList)-1 {
			result.WriteString(validatorCharSeparator)
		}
	}
	return result.String()
}

func (valid *validator) isValidName(name string) bool {
	for _, c := range name {
		if unicode.IsLetter(c) == false && unicode.IsDigit(c) == false && c != '_' {
			return false
		}
	}
	return true
}

func (valid *validator) languageCodeValid(importFile *Import) {
	var metaList = importFile.metaList
	lang := new(Language)

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		metaType := meta.GetEntityType()

		if metaType == entitytype.Language {
			_, err := lang.IsCodeValid(meta.value)
			if err != nil {
				importFile.logErrorStr(549, "Invalid language code", err.Error())
			}
		}
	}
}
