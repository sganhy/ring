package schema

import (
	"errors"
	"fmt"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"strconv"
	"strings"
	"unicode"
)

const (
	validatorCharSeparator    string = ", "
	validatorAtLine           string = "\n at line %d"
	invalidEntityName         string = "Invalid %s name"
	invalidIndexValue         string = "Invalid index definition"
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
// getters and setters
//******************************

//******************************
// public methods
//******************************

func (valid *validator) ValidateImport(importFile *Import) bool {

	importFile.logInfo("Start validation", fmt.Sprintf("meta items count: %d", len(importFile.metaList)))
	if importFile.metaList == nil || len(importFile.metaList) <= 0 {
		importFile.logError(errors.New("empty metaList"))
		importFile.errorCount++
		valid.errorCount++
		return false
	}

	//{1} - step1
	valid.tableIdUnique(importFile)
	valid.tableNameUnique(importFile)
	valid.entityNameValid(importFile)
	valid.languageCodeValid(importFile)
	valid.entityTypeValid(importFile)
	valid.indexValueValid(importFile)

	//{2} - step2
	if importFile.errorCount == 0 {
		valid.tableValueValid(importFile)
		valid.entityNameUnique(importFile)
		// duplicate fields & relations into tables
	}

	//{3} - step3
	if importFile.errorCount == 0 {
		valid.inverseRelationValid(importFile)
		valid.indexValid(importFile)
		valid.duplicateIndex(importFile)
	}

	//{3} final checks
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
			val, ok = dico[meta.id]
			if !ok {
				// new slice
				val = make([]*Meta, 0, 2)
			}
			val = append(val, meta)
			dico[meta.id] = val
		}
	}
	for key, arr := range dico {
		if len(arr) > 1 {
			var message = fmt.Sprintf("Duplicate table Id '%d'", key)
			var description = fmt.Sprintf("table id '%d' is already in use in tables %s.",
				key, valid.joinMeta(arr, joinOperationByName))
			description += "\nat lines " + valid.joinMeta(arr, joinOperationByLineNumber)
			importFile.logErrorStr(703, message, description)
		}
	}
}

// Check if table name are unique
func (valid *validator) tableNameUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dico map[string]*Meta
	dico = make(map[string]*Meta)

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table {
			name := strings.ToUpper(meta.name)
			if val, ok := dico[name]; ok {
				var message = "Duplicate table name '%s'"
				var description = fmt.Sprintf("table name '%s' is already in use.\n at line %d and %d",
					meta.name, val.lineNumber, meta.lineNumber)
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
			var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
			var description = fmt.Sprintf("%s name cannot be empty."+validatorAtLine,
				strings.ToLower(metaType.String()), meta.lineNumber)
			importFile.logErrorStr(502, message, description)
		} else {
			valid.checkEntityName(importFile, meta, metaType)
		}
	}
}

func (valid *validator) checkEntityName(importFile *Import, meta *Meta, metaType entitytype.EntityType) {
	const descMaxLengthEntityName string = "%s name '%s' is too long (max length=%d)." + validatorAtLine

	// is name valid
	if valid.isValidName(meta.name) == false &&
		(metaType == entitytype.Table || metaType == entitytype.Field || metaType == entitytype.Index ||
			metaType == entitytype.Relation || metaType == entitytype.Tablespace) {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf("invalid %s name '%s'. A name can consist of any combination of letters(A to Z a to z), decimal digits(0 to 9) or underscore (_)."+validatorAtLine,
			strings.ToLower(metaType.String()), meta.name, meta.lineNumber)
		importFile.logErrorStr(501, message, description)
	}
	// is meta.name len > 28
	if len(meta.name) > prefixedEntityMaxLength && (metaType == entitytype.Table || metaType == entitytype.Field) {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descMaxLengthEntityName, strings.ToLower(metaType.String()), meta.name,
			prefixedEntityMaxLength, meta.lineNumber)
		importFile.logErrorStr(503, message, description)
	}
	// is meta.name len > 30
	if len(meta.name) > unPrefixedEntityMaxLength && metaType != entitytype.Table && metaType != entitytype.Field {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descMaxLengthEntityName, strings.ToLower(metaType.String()), meta.name,
			unPrefixedEntityMaxLength, meta.lineNumber)
		importFile.logErrorStr(505, message, description)
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

// Check for relation and field if there unique for the same table
func (valid *validator) entityNameUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dicoTable map[int32]map[string]bool
	var dicoTableName map[int32]string
	dicoTable = make(map[int32]map[string]bool)
	dicoTableName = make(map[int32]string)

	// (1) build table dictionary
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table {
			dicoTable[meta.id] = make(map[string]bool)
			dicoTableName[meta.id] = meta.name
		}
	}

	// (2) load field & relations
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		metaType := meta.GetEntityType()
		if metaType == entitytype.Field || metaType == entitytype.Relation {
			entityName := strings.ToUpper(meta.name)
			if _, ok := dicoTable[meta.refId][entityName]; ok {
				var message = "Duplicate relation or field"
				var description = fmt.Sprintf("relation or field '%s' for table '%s'"+validatorAtLine,
					meta.name, dicoTableName[meta.refId], meta.lineNumber)
				importFile.logErrorStr(811, message, description)
			} else {
				dicoTable[meta.refId][entityName] = true
			}
		}
	}

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

// Check if A name consist of any combination of letters(A to Z a to z), decimal digits(0 to 9) or underscore (_).
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

func (valid *validator) entityTypeValid(importFile *Import) {
	var metaList = importFile.metaList

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		metaType := meta.GetEntityType()
		if metaType == entitytype.Relation && meta.GetRelationType() == relationtype.NotDefined {
			var description = fmt.Sprintf("wrong relation type (must be OTOP, OTM, MTM, MTO, or OTOF)"+validatorAtLine, meta.lineNumber)
			importFile.logErrorStr(852, "Invalid relation type", description)
		}

		if metaType == entitytype.Field && meta.GetFieldType() == fieldtype.NotDefined {
			var description = fmt.Sprintf("wrong field type "+validatorAtLine, meta.lineNumber)
			importFile.logErrorStr(853, "Invalid field type", description)
		}
	}
}

// Check Relation inverse relation + inverse Type
func (valid *validator) inverseRelationValid(importFile *Import) {
	//	var dicoTable map[int32]map[string]bool

	var metaList = importFile.metaList
	var dicoRelation map[int32]map[string]relationtype.RelationType
	var relations []*Meta

	dicoRelation = make(map[int32]map[string]relationtype.RelationType)
	relations = make([]*Meta, 0, 10)

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]

		if meta.GetEntityType() == entitytype.Relation {
			var relationName = strings.ToUpper(meta.name)

			relations = append(relations, meta)
			if _, ok := dicoRelation[meta.refId]; !ok {
				dicoRelation[meta.refId] = make(map[string]relationtype.RelationType)
			}
			dicoRelation[meta.refId][relationName] = meta.GetRelationType()
		}
	}

	// (2) check relations
	for i := 0; i < len(relations); i++ {
		meta := metaList[i]
		relationType := meta.GetRelationType().InverseRelationType()

		fmt.Println(relationType)
	}
}

func (valid *validator) indexValueValid(importFile *Import) {
	var metaList = importFile.metaList

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]

		if meta.GetEntityType() == entitytype.Index && len(strings.Trim(meta.value, " ")) == 0 {
			var description = fmt.Sprintf("empty index definition"+validatorAtLine, meta.lineNumber)
			importFile.logErrorStr(860, invalidIndexValue, description)
		}
	}
}

func (valid *validator) tableValueValid(importFile *Import) {
	var metaList = importFile.metaList
	var fieldDico map[int32]bool

	fieldDico = make(map[int32]bool)

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		metaType := meta.GetEntityType()
		if metaType == entitytype.Field || metaType == entitytype.Relation {
			fieldDico[meta.refId] = true
		}
	}

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		if meta.GetEntityType() == entitytype.Table {
			if _, ok := fieldDico[meta.id]; !ok {
				var description = fmt.Sprintf("empty table definition"+validatorAtLine, meta.lineNumber)
				importFile.logErrorStr(864, "Invalid table definition", description)
			}
		}
	}

}

// Check if indexes reference existing fields
func (valid *validator) indexValid(importFile *Import) {

	var metaList = importFile.metaList
	// Upper(field)+Refid string
	var dicoField map[string]bool
	dicoField = make(map[string]bool)

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]
		metaType := meta.GetEntityType()

		if metaType == entitytype.Field || metaType == entitytype.Relation {
			key := strings.ToUpper(meta.name) + strconv.Itoa(int(meta.refId))
			dicoField[key] = true
		}
	}

	for i := 0; i < len(metaList); i++ {
		meta := metaList[i]

		if meta.GetEntityType() == entitytype.Index {
			strArr := strings.Split(meta.value, metaIndexSeparator)
			for j := 0; j < len(strArr); j++ {
				key := strings.Trim(strings.ToUpper(strArr[j]), "") + strconv.Itoa(int(meta.refId))
				if _, ok := dicoField[key]; !ok {
					var description = fmt.Sprintf("invalid indexed field or relation '%s' "+validatorAtLine, strArr[j], meta.lineNumber)
					importFile.logErrorStr(861, invalidIndexValue, description)
				}
			}
		}

	}

}

/*TODO Detect duplicate index definition ()
func (valid *validator) duplicateIndex(importFile *Import) {

}
*/
