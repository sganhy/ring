package schema

import (
	"errors"
	"fmt"
	"os"
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
	invalidTablespaceValue    string = "Invalid tablespace definition"
	invalidRelationValue      string = "Invalid relation definition"
	invalidTableValue         string = "Invalid table definition"
	joinOperationByName       int    = 1
	joinOperationByLineNumber int    = 2
	prefixedEntityMaxLength   int    = 28
	unPrefixedEntityMaxLength int    = 30
)

type validator struct {
	errorCount    int
	fieldCount    int
	relationCount int
	tableCount    int
	indexCount    int
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

	//{0} - step0
	valid.loadStats(importFile)

	//{1} - step1
	valid.tableIdUnique(importFile)
	valid.tableNameUnique(importFile)
	valid.entityNameValid(importFile)
	valid.languageCodeValid(importFile)
	valid.entityTypeValid(importFile)
	valid.indexValueValid(importFile)
	valid.tableSpaceValueValid(importFile)

	//{2} - step2
	if importFile.errorCount == 0 {
		valid.tableValueValid(importFile)
		valid.fieldNameUnique(importFile)
		// duplicate fields & relations into tables
	}

	//{3} - step3
	if importFile.errorCount == 0 {
		valid.inverseRelationValid(importFile)
		valid.indexValid(importFile)
		//valid.duplicateIndex(importFile)
	}

	//{4} - step4 - compare with previous schema (if exist)
	if GetSchemaById(importFile.schemaId) != nil {
		valid.isTableIdReserved(importFile)

	}

	//{5} final checks
	if importFile.errorCount == 0 {
		valid.duplicateMetaKey(importFile)
	}

	return true
}

//******************************
// private methods
//******************************
func (valid *validator) loadStats(importFile *Import) {
	var metaList = importFile.metaList
	valid.tableCount = 0
	valid.fieldCount = 0
	valid.relationCount = 0
	valid.indexCount = 0

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		switch metaData.GetEntityType() {
		case entitytype.Table:
			valid.tableCount++
			break
		case entitytype.Field:
			valid.fieldCount++
			break
		case entitytype.Relation:
			valid.relationCount++
			break
		case entitytype.Index:
			valid.indexCount++
			break
		}
	}
}

func (valid *validator) tableIdUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dico map[int32][]*meta
	var ok bool
	var val []*meta

	dico = make(map[int32][]*meta, valid.tableCount)
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			val, ok = dico[metaData.id]
			if !ok {
				// new slice
				val = make([]*meta, 0, 2)
			}
			val = append(val, metaData)
			dico[metaData.id] = val
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
	var dico map[string]*meta
	dico = make(map[string]*meta, valid.tableCount)

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			name := strings.ToUpper(metaData.name)
			if val, ok := dico[name]; ok {
				var message = "Duplicate table name '%s'"
				var description = fmt.Sprintf("table name '%s' is already in use.\n at line %d and %d",
					metaData.name, val.lineNumber, metaData.lineNumber)
				message = fmt.Sprintf(message, metaData.name)
				importFile.logErrorStr(704, message, description)
			} else {
				dico[name] = metaData
			}
		}
	}
}

func (valid *validator) entityNameValid(importFile *Import) {
	var metaList = importFile.metaList

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		metaType := metaData.GetEntityType()

		// is metaData.name empty?
		if len(strings.TrimSpace(metaData.name)) == 0 {
			var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
			var description = fmt.Sprintf("%s name cannot be empty."+validatorAtLine,
				strings.ToLower(metaType.String()), metaData.lineNumber)
			importFile.logErrorStr(502, message, description)
		} else {
			valid.checkEntityName(importFile, metaData, metaType)
		}
	}
}

func (valid *validator) checkEntityName(importFile *Import, metaData *meta, metaType entitytype.EntityType) {
	const descMaxLengthEntityName string = "%s name '%s' is too long (max length=%d)." + validatorAtLine

	// is name valid
	if valid.isValidName(metaData.name) == false &&
		(metaType == entitytype.Table || metaType == entitytype.Field || metaType == entitytype.Index ||
			metaType == entitytype.Relation || metaType == entitytype.Tablespace) {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf("invalid %s name '%s'. A name can consist of any combination of letters(A to Z a to z), decimal digits(0 to 9) or underscore (_)."+validatorAtLine,
			strings.ToLower(metaType.String()), metaData.name, metaData.lineNumber)
		importFile.logErrorStr(501, message, description)
	}
	// is metaData.name len > 28
	if len(metaData.name) > prefixedEntityMaxLength && (metaType == entitytype.Table || metaType == entitytype.Field) {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descMaxLengthEntityName, strings.ToLower(metaType.String()), metaData.name,
			prefixedEntityMaxLength, metaData.lineNumber)
		importFile.logErrorStr(503, message, description)
	}
	// is metaData.name len > 30
	if len(metaData.name) > unPrefixedEntityMaxLength && metaType != entitytype.Table && metaType != entitytype.Field {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf(descMaxLengthEntityName, strings.ToLower(metaType.String()), metaData.name,
			unPrefixedEntityMaxLength, metaData.lineNumber)
		importFile.logErrorStr(505, message, description)
	}
	// schema name cannot be equal to "@meta"
	if metaType == entitytype.Schema && strings.ToLower(strings.Trim(metaData.name, " ")) == strings.ToLower(metaSchemaName) {
		var message = fmt.Sprintf(invalidEntityName, strings.ToLower(metaType.String()))
		var description = fmt.Sprintf("invalid schema name '%s'. A schema name cannot be equal to '@meta'."+validatorAtLine,
			metaData.name, metaData.lineNumber)
		importFile.logErrorStr(507, message, description)
	}

}

func (valid *validator) duplicateMetaKey(importFile *Import) {
	var metaList = importFile.metaList
	var dicoEntities map[string]bool
	dicoEntities = make(map[string]bool, len(metaList))

	// check on db unique key (pk_@meta) ==> id|schema_id|object_type|reference_id
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]

		metaKey := strconv.FormatInt(int64(metaData.id), 16) + "-" +
			strconv.FormatInt(int64(metaData.refId), 16) + "-" +
			strconv.FormatInt(int64(metaData.objectType), 16)

		if _, ok := dicoEntities[metaKey]; ok {
			// error duplicate meta key
			var message = fmt.Sprintf("Duplicate meta key")
			var description = fmt.Sprintf("duplicate meta key (type=%s): refid=%d, id=%d", strings.ToLower(metaData.GetEntityType().String()),
				metaData.refId, metaData.id)

			importFile.logErrorStr(527, message, description)
		} else {
			dicoEntities[metaKey] = true
		}
	}
}

// Check for relation and field if there unique for the same table
func (valid *validator) fieldNameUnique(importFile *Import) {
	var metaList = importFile.metaList
	var dicoTable map[int32]map[string]bool
	var dicoTableName map[int32]string
	dicoTable = make(map[int32]map[string]bool, valid.tableCount)
	dicoTableName = make(map[int32]string, valid.tableCount)

	// (1) build table dictionary
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			dicoTable[metaData.id] = make(map[string]bool)
			dicoTableName[metaData.id] = metaData.name
		}
	}

	// (2) load field & relations
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		metaType := metaData.GetEntityType()
		if metaType == entitytype.Field || metaType == entitytype.Relation {
			entityName := strings.ToUpper(metaData.name)
			if _, ok := dicoTable[metaData.refId][entityName]; ok {
				var message = "Duplicate relation or field"
				var description = fmt.Sprintf("relation or field '%s' for table '%s'"+validatorAtLine,
					metaData.name, dicoTableName[metaData.refId], metaData.lineNumber)
				importFile.logErrorStr(811, message, description)
			} else {
				dicoTable[metaData.refId][entityName] = true
			}
		}
	}

}

func (valid *validator) joinMeta(metaList []*meta, operation int) string {
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
		metaData := metaList[i]
		metaType := metaData.GetEntityType()

		if metaType == entitytype.Language {
			_, err := lang.IsCodeValid(metaData.value)
			if err != nil {
				importFile.logErrorStr(549, "Invalid language code", err.Error())
			}
		}
	}
}

func (valid *validator) entityTypeValid(importFile *Import) {
	var metaList = importFile.metaList

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		metaType := metaData.GetEntityType()
		if metaType == entitytype.Relation && metaData.GetRelationType() == relationtype.NotDefined {
			var description = fmt.Sprintf("wrong relation type (must be OTOP, OTM, MTM, MTO, or OTOF)"+validatorAtLine, metaData.lineNumber)
			importFile.logErrorStr(852, "Invalid relation type", description)
		}

		if metaType == entitytype.Field && metaData.GetFieldType() == fieldtype.NotDefined {
			var description = fmt.Sprintf("wrong field type "+validatorAtLine, metaData.lineNumber)
			importFile.logErrorStr(853, "Invalid field type", description)
		}
	}
}

// Check Relation inverse relation + inverse Type (case sensitif)
func (valid *validator) inverseRelationValid(importFile *Import) {
	//	var dicoTable map[int32]map[string]bool

	var metaList = importFile.metaList

	// <relation.RefId, <relation.Name, RelationType>>
	var dicoRelation map[int32]map[string]relationtype.RelationType
	var relations []*meta
	var ok bool
	var val relationtype.RelationType

	dicoRelation = make(map[int32]map[string]relationtype.RelationType, valid.relationCount)
	relations = make([]*meta, 0, 10)

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]

		if metaData.GetEntityType() == entitytype.Relation {
			var relationName = metaData.name

			relations = append(relations, metaData)
			if _, ok = dicoRelation[metaData.refId]; !ok {
				dicoRelation[metaData.refId] = make(map[string]relationtype.RelationType)
			}
			dicoRelation[metaData.refId][relationName] = metaData.GetRelationType()
		}
	}

	// (2) check relations
	for i := 0; i < len(relations); i++ {
		metaData := relations[i]

		if metaData.value == "" {
			var description = fmt.Sprintf("empty inverse relation definition"+validatorAtLine, metaData.lineNumber)
			importFile.logErrorStr(955, invalidRelationValue, description)
			continue
		}
		relationName := metaData.value
		if val, ok = dicoRelation[metaData.dataType][relationName]; !ok {
			var description = fmt.Sprintf("invalid inverse relation definition '%s'"+validatorAtLine, metaData.value, metaData.lineNumber)
			importFile.logErrorStr(956, invalidRelationValue, description)
			continue
		}

		if val.InverseRelationType() != metaData.GetRelationType() {
			var description = fmt.Sprintf("invalid relation type '%s'"+validatorAtLine, metaData.GetRelationType().String(),
				metaData.lineNumber)
			importFile.logErrorStr(957, invalidRelationValue, description)
		}
	}
}

func (valid *validator) indexValueValid(importFile *Import) {
	var metaList = importFile.metaList

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Index && len(strings.TrimSpace(metaData.value)) == 0 {
			var description = fmt.Sprintf("empty index definition"+validatorAtLine, metaData.lineNumber)
			importFile.logErrorStr(860, invalidIndexValue, description)
		}
	}
}

// Check table without fields or relations
func (valid *validator) tableValueValid(importFile *Import) {
	var metaList = importFile.metaList
	var fieldDico map[int32]bool

	fieldDico = make(map[int32]bool, valid.fieldCount+valid.relationCount)

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		metaType := metaData.GetEntityType()
		if metaType == entitytype.Field || metaType == entitytype.Relation {
			fieldDico[metaData.refId] = true
		}
	}

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			if _, ok := fieldDico[metaData.id]; !ok {
				var description = fmt.Sprintf("empty table definition"+validatorAtLine, metaData.lineNumber)
				importFile.logErrorStr(864, "Invalid table definition", description)
			}
		}
	}

}

// Check if indexes reference existing fields
func (valid *validator) indexValid(importFile *Import) {
	var metaList = importFile.metaList
	// field name+Refid string
	var dicoField map[string]bool
	dicoField = make(map[string]bool, valid.fieldCount)

	// (1) generate dictionary
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		metaType := metaData.GetEntityType()

		// make check case sensitive
		if metaType == entitytype.Field || metaType == entitytype.Relation {
			key := metaData.name + strconv.Itoa(int(metaData.refId))
			dicoField[key] = true
		}
	}
	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Index {
			strArr := strings.Split(metaData.value, metaIndexSeparator)
			for j := 0; j < len(strArr); j++ {
				key := strArr[j] + strconv.Itoa(int(metaData.refId))
				if _, ok := dicoField[key]; !ok {
					var description = fmt.Sprintf("invalid indexed field or relation '%s' "+validatorAtLine, strArr[j], metaData.lineNumber)
					importFile.logErrorStr(861, invalidIndexValue, description)
				}
			}
		}
	}
}

// check if path valid
func (valid *validator) tableSpaceValueValid(importFile *Import) {
	var metaList = importFile.metaList

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		metaType := metaData.GetEntityType()
		if metaType == entitytype.Tablespace {
			valueUpp := strings.ToUpper(metaData.value)

			if strings.HasSuffix(valueUpp, "/NUL") || strings.HasSuffix(valueUpp, "/NUL/") {
				var description = fmt.Sprintf("invalid tablespace '%s' contains value 'NUL'"+validatorAtLine,
					metaData.name, metaData.lineNumber)
				importFile.logErrorStr(884, invalidTablespaceValue, description)
				continue
			}

			err := os.MkdirAll(metaData.value, os.ModePerm)
			if err != nil {
				var description = fmt.Sprintf("invalid tablespace '%s' value '%s' error: %s"+validatorAtLine,
					metaData.name, metaData.value, err.Error(), metaData.lineNumber)
				importFile.logErrorStr(885, invalidTablespaceValue, description)

			}
		}
	}
}

// we cannot re-use a reserved {table_id}
func (valid *validator) isTableIdReserved(importFile *Import) {
	metaList := getMetaList(importFile.schemaId, entitytype.Table, true)

	if len(metaList) <= 0 {
		return
	}
	dicoTable := make(map[int32]string, len(metaList))

	// create dico
	for i := 0; i < len(metaList); i++ {
		dicoTable[metaList[i].id] = metaList[i].name
	}

	// compare tables
	for i := 0; i < len(importFile.metaList); i++ {
		metaData := importFile.metaList[i]
		metaType := metaData.GetEntityType()

		if metaType == entitytype.Table {
			if name, ok := dicoTable[metaData.id]; ok {
				if strings.EqualFold(name, metaData.name) == false {
					var description = fmt.Sprintf("table id %d already reserved for '%s' entity",
						metaData.id, metaData.name)
					importFile.logErrorStr(909, invalidTableValue, description)
				}
			}
		}
	}
}

/*TODO Detect duplicate index definition ()
func (valid *validator) duplicateIndex(importFile *Import) {

}
*/
