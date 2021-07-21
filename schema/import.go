package schema

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"reflect"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"ring/schema/sourcetype"
	"runtime"
	"strconv"
	"strings"
)

type Import struct {
	id              int
	fileName        string
	schemaId        int32
	schemaName      string
	source          sourcetype.SourceType
	initialized     bool
	loaded          bool
	jobId           int64
	tables          map[string]int32
	errorCount      int32
	metaList        []*meta
	newSchema       *Schema
	logger          *log
	tablespaceCount int32
	provider        databaseprovider.DatabaseProvider
}

const (
	errorImportFileNotInitialized   string = "Import File object is not initialized."
	errorImportInvalidId            string = "Invalid attribute {id}"
	errorImportFieldSize            string = "Invalid field size"
	importTableTag1                 string = "object"
	importTableTag2                 string = "table"
	importFieldTag                  string = "field"
	importTablespaceTag             string = "tablespace"
	importSchemaTag                 string = "schema"
	importDescriptionTag            string = "description"
	importRelationTag               string = "relation"
	importIndexTag                  string = "index"
	importIndexFieldTag             string = "index_field"
	importXmlAttributeVersionTag    string = "version"
	importXmlAttributeLangTag       string = "default_language"
	importXmlAttributeFileTag       string = "file"
	importXmlAttributeNameTag       string = "name"
	importXmlAttributeSubjectTag    string = "suject"
	importXmlAttributeTypeTag       string = "type"
	importXmlAttributeDefaultTag    string = "default"
	importXmlAttributeIdTag1        string = "id"
	importXmlAttributeIdTag2        string = "type_id"
	importXmlAttributeToTag         string = "to"
	importXmlAttributeInverseTag    string = "inverse_relation"
	importXmlAttributeBaseline      string = "baseline"
	importXmlAttributeConstraintTag string = "constraint"
	importXmlAttributeCached        string = "cached"
	importXmlBoolTrueValue1         string = "true"
	importXmlBoolTrueValue2         string = "1"
	importXmlBoolFalseValue1        string = "false"
	importXmlBoolFalseValue2        string = "0"
	importXmlAttributeSize          string = "size"
	importXmlAttributeReadonly      string = "readonly"
	importXmlAttributeUnique        string = "unique"
	importXmlAttributeBitmap        string = "bitmap"
	importXmlAttributeNotNull       string = "not_null"
	importXmlAttributeSensitive     string = "case_sensitive"
	importXmlAttributeMultiLang     string = "multilingual"
	importMinId                     int64  = -2147483648
	importMaxId                     int64  = 2147483647
	baseErrorId                     int32  = 11
)

var (
	currentSchemaImportId int = 101
)

func (importFile *Import) Init(source sourcetype.SourceType, fileName string) {
	currentSchemaImportId++
	importFile.id = currentSchemaImportId
	importFile.fileName = fileName
	importFile.initialized = true
	importFile.source = source
	importFile.logger = new(log)
	importFile.logger.Init(schemaNotDefined, 0, false)
	importFile.loaded = false
}

//******************************
// getters and setters
//******************************
func (importFile *Import) GetJobId() int64 {
	return importFile.jobId
}

func (importFile *Import) GetId() int {
	return importFile.id
}

func (importFile *Import) GetFile() string {
	return importFile.fileName
}

func (importFile *Import) GetSchemaName() string {
	return importFile.schemaName
}

func (importFile *Import) GetSchemaId() int32 {
	return importFile.schemaId
}

func (importFile *Import) GetSourceType() sourcetype.SourceType {
	return importFile.source
}

func (importFile *Import) GetSchema() *Schema {
	return importFile.newSchema
}

func (importFile *Import) GetDatabaseProvider() databaseprovider.DatabaseProvider {
	return importFile.provider
}

//******************************
// public methods
//******************************
func (importFile *Import) Load() {
	var metaSchema = GetSchemaByName(metaSchemaName)
	importFile.provider = getDefaultDbProvider()
	importFile.tablespaceCount = 0
	importFile.metaList = make([]*meta, 0, 20)
	importFile.errorCount = 0
	importFile.jobId = metaSchema.getJobIdNextValue()
	importFile.loadSchemaInfo()
	importFile.logInfo("Load Schema", "import_file: "+importFile.fileName)

	if importFile.schemaName == "" {
		// no schema information
		return
	}
	if importFile.initialized == false {
		errors.New(errorImportFileNotInitialized)
	}
	if importFile.source == sourcetype.XmlDocument {
		importFile.loadXml()
		valid := new(validator)
		valid.Init()
		valid.ValidateImport(importFile)
	}
	importFile.loaded = true
}

func (importFile *Import) Upgrade() {
	//const upgradeSchema string = "Upgrade schema"

	if importFile.IsValid() == true {
		var err error
		err = importFile.saveMetaList()
		if err != nil {
			importFile.logError(err)
			return
		}

		err = importFile.saveMetaIdList()
		if err != nil {
			importFile.logError(err)
			return
		}
		go importFile.analyzeMetaTables()
		//TODO rollback !!
		// no error
		importFile.loaded = false
		importFile.metaList = nil
		importFile.newSchema = getSchemaById(importFile.schemaId)
		// database.go
		upgradeSchema(importFile.jobId, importFile.newSchema)
		runtime.GC()

	}
}

func (importFile *Import) IsValid() bool {
	return importFile.errorCount == 0 && importFile.loaded
}

//******************************
// private methods
//******************************
func (importFile *Import) loadXml() {
	referenceId := new(int32)
	//var schemaId int32 = 0
	fieldId := new(int32)
	relationId := new(int32)
	indexId := new(int32)

	var metaData *meta
	var f *os.File

	*referenceId = 0
	*fieldId = 0
	*relationId = 0
	*indexId = 0

	// pass 1 build --> dico<lower(table_name),table_id>
	err := importFile.loadTableInfo()
	if err != nil {
		// unable to open file
		return
	}

	// pass 2 build --> get []Meta
	f, err = os.Open(importFile.fileName)
	if err != nil {
		importFile.logError(err)
		return
	}
	defer f.Close()
	d := xml.NewDecoder(f)
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			importFile.logError(err)
			return
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			metaData = importFile.manageElement(d, &ty, fieldId, relationId, indexId, referenceId)
			if metaData != nil {
				importFile.metaList = append(importFile.metaList, metaData)
			}
			break
		}
	}
	// de-allow table map
	importFile.tables = nil
}

func (importFile *Import) manageElement(d *xml.Decoder, ty *xml.StartElement, fieldId *int32, relationId *int32, indexId *int32, referenceId *int32) *meta {
	var metaData *meta

	// SCHEMA
	if strings.ToLower(ty.Name.Local) == importSchemaTag {
		metaData = importFile.getXmlMetaSchema(&ty.Attr, reflect.ValueOf(d).Elem().FieldByName("line").Int())
	}
	// TABLESPACE
	if strings.ToLower(ty.Name.Local) == importTablespaceTag {
		metaData = importFile.getXmlMetaTablespace(&ty.Attr, reflect.ValueOf(d).Elem().FieldByName("line").Int())
	}
	// TABLE
	if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
		*fieldId = 0
		*relationId = 0
		*indexId = 0
		metaData = importFile.getXmlMeta(&ty.Attr, entitytype.Table, importFile.schemaId, 0,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		metaData.flags = importFile.getTableFlags(&ty.Attr)
		metaData.description = importFile.getDescription(&ty.Attr)
		*referenceId = metaData.id
	}
	// FIELDS
	if strings.ToLower(ty.Name.Local) == importFieldTag {
		(*fieldId)++
		line := reflect.ValueOf(d).Elem().FieldByName("line").Int()
		metaData = importFile.getXmlMeta(&ty.Attr, entitytype.Field, *referenceId, *fieldId, line)
		metaData.flags = importFile.getFieldFlags(&ty.Attr, line)
		metaData.description = importFile.getDescription(&ty.Attr)
	}
	// RELATIONS
	if strings.ToLower(ty.Name.Local) == importRelationTag {
		(*relationId)++
		metaData = importFile.getXmlMeta(&ty.Attr, entitytype.Relation, *referenceId, *relationId,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		metaData.flags = importFile.getRelationFlags(&ty.Attr)
		metaData.description = importFile.getDescription(&ty.Attr)
	}
	// INDEXES
	if strings.ToLower(ty.Name.Local) == importIndexTag {
		(*indexId)++
		metaData = importFile.getXmlMeta(&ty.Attr, entitytype.Index, *referenceId, *indexId,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		metaData.flags = importFile.getIndexFlags(&ty.Attr)
		metaData.value = importFile.getIndexValue(d, ty)
		metaData.description = importFile.getDescription(&ty.Attr)
	}
	// DESCRIPTION
	if strings.ToLower(ty.Name.Local) == importDescriptionTag {
		importFile.addDescription(d, ty)
	}
	return metaData
}

func (importFile *Import) getXmlMetaSchema(attributes *[]xml.Attr, line int64) *meta {
	var result = new(meta)
	var schema = new(Schema)

	result.flags = uint64(importFile.provider)
	result.name = importFile.getXmlAttribute(attributes, importXmlAttributeNameTag)
	result.objectType = int8(entitytype.Schema)
	result.refId = 0
	result.id = 0
	result.lineNumber = line
	result.enabled = true
	result.value = schema.getPhysicalName(importFile.provider, importFile.schemaName)

	return result
}

func (importFile *Import) getXmlMetaTablespace(attributes *[]xml.Attr, line int64) *meta {
	var result = new(meta)
	importFile.tablespaceCount++

	result.name = importFile.getXmlAttribute(attributes, importXmlAttributeNameTag)
	result.objectType = int8(entitytype.Tablespace)
	result.refId = 0
	result.id = importFile.tablespaceCount
	result.lineNumber = line
	result.value = importFile.getXmlAttribute(attributes, importXmlAttributeFileTag)
	result.enabled = true
	var isIndex = strings.ToLower(importFile.getXmlAttribute(attributes, importIndexTag))
	if isIndex == importXmlBoolTrueValue1 || isIndex == importXmlBoolTrueValue2 {
		result.setTablespaceIndex(true)
	}
	var isTable = strings.ToLower(importFile.getXmlAttribute(attributes, importTableTag2))
	if isTable == importXmlBoolTrueValue1 || isTable == importXmlBoolTrueValue2 {
		result.setTablespaceTable(true)
	}

	return result
}

func (importFile *Import) getXmlMeta(attributes *[]xml.Attr, entityType entitytype.EntityType, referenceId int32, defaultId int32, line int64) *meta {
	var result = new(meta)
	count := len(*attributes)
	result.objectType = int8(entityType)
	result.id = defaultId
	result.refId = referenceId
	result.lineNumber = line
	result.enabled = true

	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		switch strings.ToLower(attribute.Name.Local) {
		case importXmlAttributeNameTag:
			// NAME
			result.name = attribute.Value
			break
		case importXmlAttributeIdTag1, importXmlAttributeIdTag2:
			// ID
			result.id = importFile.getXmlAttributeId(attribute.Value, line)
			break
		case importXmlAttributeTypeTag, importXmlAttributeToTag:
			// DATA_TYPE
			result.dataType = importFile.getDataType(entityType, attribute.Name.Local, attribute.Value)
			break
		case importXmlAttributeDefaultTag, importXmlAttributeSubjectTag, importXmlAttributeInverseTag:
			// VALUE
			result.value = importFile.getValue(entityType, attribute.Value)
			break
		}
	}
	return result
}

func (importFile *Import) getXmlAttribute(attributes *[]xml.Attr, attributeName string) string {
	if attributes != nil {
		count := len(*attributes)
		for i := 0; i < count; i++ {
			var attribute = (*attributes)[i]
			if strings.EqualFold(attribute.Name.Local, attributeName) {
				return attribute.Value
			}
		}
	}
	return ""
}

func (importFile *Import) getXmlAttributeId(attributeValue string, line int64) int32 {
	result, err := strconv.ParseInt(attributeValue, 10, 64)
	if err == nil {
		if result >= importMinId && result <= importMaxId {
			return int32(result)
		}
		importFile.logFileError(errorImportInvalidId, "", line)
		return -2
	}
	importFile.logFileError(errorImportInvalidId, err.Error(), line)
	return -1
}

func (importFile *Import) getDataType(entityType entitytype.EntityType, attributeName string, value string) int32 {
	if entityType == entitytype.Field && strings.ToLower(attributeName) == importXmlAttributeTypeTag {
		return int32(fieldtype.GetFieldType(value))
	}
	if entityType == entitytype.Relation && strings.ToLower(attributeName) == importXmlAttributeToTag {
		if val, ok := importFile.tables[strings.ToLower(value)]; ok {
			return val
		}
	}
	return 0
}

func (importFile *Import) getValue(entityType entitytype.EntityType, value string) string {
	if entityType == entitytype.Field || entityType == entitytype.Relation || entityType == entitytype.Table ||
		entityType == entitytype.Schema {
		// get Field: default, Relation: inverse relationship name,Table: Subject
		return value
	}
	return ""
}

func (importFile *Import) getTableFlags(attributes *[]xml.Attr) uint64 {
	metaData := meta{}
	count := len(*attributes)
	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		if strings.EqualFold(attribute.Value, importXmlBoolTrueValue1) || strings.EqualFold(attribute.Value, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeBaseline:
				// BASELINE
				metaData.setEntityBaseline(true)
				break
			case importXmlAttributeReadonly:
				metaData.setTableReadonly(true)
				break
			case importXmlAttributeCached:
				metaData.setTableCached(true)
				break
			}
		}
	}
	return metaData.flags
}

func (importFile *Import) getRelationFlags(attributes *[]xml.Attr) uint64 {
	metaData := meta{}
	count := len(*attributes)
	metaData.flags = 0
	metaData.setRelationConstraint(true)

	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		var attributeValue = strings.Trim(attribute.Value, " ")

		if strings.EqualFold(attributeValue, importXmlBoolTrueValue1) || strings.EqualFold(attributeValue, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeBaseline:
				// BASELINE
				metaData.setEntityBaseline(true)
				break
			case importXmlAttributeNotNull:
				metaData.setRelationNotNull(true)
				break
			}
		}
		if strings.ToLower(attribute.Name.Local) == importXmlAttributeTypeTag {
			metaData.setRelationType(relationtype.GetRelationType(attributeValue))
		}
		if strings.ToLower(attribute.Name.Local) == importXmlAttributeConstraintTag &&
			(strings.EqualFold(attributeValue, importXmlBoolFalseValue1) ||
				strings.EqualFold(attributeValue, importXmlBoolFalseValue2)) {
			metaData.setRelationConstraint(false)
		}
	}
	return metaData.flags
}

func (importFile *Import) getFieldFlags(attributes *[]xml.Attr, line int64) uint64 {
	metaData := meta{}
	count := len(*attributes)
	metaData.setFieldCaseSensitive(true)
	metaData.setFieldSize(0)
	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		var attributeValue = strings.Trim(attribute.Value, " ")

		if strings.EqualFold(attributeValue, importXmlBoolTrueValue1) || strings.EqualFold(attributeValue, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeBaseline:
				// BASELINE
				metaData.setEntityBaseline(true)
				break
			case importXmlAttributeNotNull:
				metaData.setFieldNotNull(true)
				break
			case importXmlAttributeMultiLang:
				metaData.setFieldMultilingual(true)
				break
			}
		}
		if strings.EqualFold(attributeValue, importXmlBoolFalseValue1) || strings.EqualFold(attributeValue, importXmlBoolFalseValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeSensitive:
				metaData.setFieldCaseSensitive(false)
				break
			}
		}
		if strings.ToLower(attribute.Name.Local) == importXmlAttributeSize {
			metaData.setFieldSize(importFile.getFieldSize(attributeValue, line))
		}
	}
	return metaData.flags
}

func (importFile *Import) getIndexFlags(attributes *[]xml.Attr) uint64 {
	metaData := meta{}
	count := len(*attributes)
	metaData.flags = 0

	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		var attributeValue = strings.Trim(attribute.Value, " ")

		if strings.EqualFold(attributeValue, importXmlBoolTrueValue1) || strings.EqualFold(attributeValue, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeBaseline:
				// BASELINE
				metaData.setEntityBaseline(true)
				break
			case importXmlAttributeUnique:
				metaData.setIndexUnique(true)
				break
			case importXmlAttributeBitmap:
				metaData.setIndexBitmap(true)
				break
			}
		}
		if strings.ToLower(attribute.Name.Local) == importXmlAttributeTypeTag {
			metaData.setRelationType(relationtype.GetRelationType(attributeValue))
		}
	}
	return metaData.flags
}

func (importFile *Import) getIndexValue(d *xml.Decoder, ty *xml.StartElement) string {
	//importIndexFieldTag
	var result strings.Builder

	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			importFile.logError(err)
			return ""
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importIndexFieldTag {
				fieldName := importFile.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
				result.WriteString(fieldName)
				result.WriteString(metaIndexSeparator)
			}
			break
		case xml.EndElement:
			if strings.ToLower(ty.Name.Local) == importIndexTag {
				return strings.Trim(result.String(), metaIndexSeparator)
			}
		}
	}
	return ""
}

// Get description from element attribute
func (importFile *Import) getDescription(attributes *[]xml.Attr) string {
	result := ""
	if attributes != nil {
		for i := 0; i < len(*attributes); i++ {
			var attribute = (*attributes)[i]
			if strings.ToLower(attribute.Name.Local) == importDescriptionTag {
				result = attribute.Value
			}
		}
	}
	return result
}

func (importFile *Import) getFieldSize(value string, line int64) uint32 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		importFile.logFileError(errorImportFieldSize, err.Error(), line)
		return 0
	}
	if result >= importMinId && result <= importMaxId {
		return uint32(result)
	}
	importFile.logFileError(errorImportFieldSize, "", line)
	return 0
}

func (importFile *Import) addDescription(d *xml.Decoder, ty *xml.StartElement) {
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			importFile.logError(err)
			return
		}
		switch ty := tok.(type) {
		case xml.CharData:
			if len(importFile.metaList) > 0 {
				lastMeta := importFile.metaList[len(importFile.metaList)-1]
				bytes := xml.CharData(ty)
				lastMeta.description = string([]byte(bytes))
			}
			break
		case xml.EndElement:
			return
		}
	}
}

//go:noinline
func (importFile *Import) logInfo(message string, description string) {
	importFile.logger.writePartialLog(9, levelInfo, importFile.jobId, message, description)
}

//go:noinline
func (importFile *Import) logError(err error) {
	importFile.logger.writePartialLog(baseErrorId+importFile.errorCount, levelError, importFile.jobId, err)
	importFile.errorCount++
}

//go:noinline
func (importFile *Import) logErrorStr(id int32, message string, description string) {
	importFile.logger.writePartialLog(id, levelError, importFile.jobId, message, description)
	importFile.errorCount++
}

//go:noinline
func (importFile *Import) logFileError(message string, description string, lineNumber int64) {
	description += "\n at line " + strconv.FormatInt(lineNumber, 10)
	importFile.logger.writePartialLog(baseErrorId+importFile.errorCount, levelError, importFile.jobId, message, description)
	importFile.errorCount++
}

func (importFile *Import) loadTableInfo() error {
	importFile.tables = make(map[string]int32)
	f, err := os.Open(importFile.fileName)
	var tableName string
	var tableId int32

	if err != nil {
		importFile.logError(err)
		return err
	}
	defer f.Close()
	d := xml.NewDecoder(f)
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			importFile.logError(err)
			return err
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
				tableName = importFile.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
				tableId = importFile.getXmlAttributeId(importFile.getXmlAttribute(&ty.Attr, importXmlAttributeIdTag1), -1)
				importFile.tables[strings.ToLower(tableName)] = tableId
			}
			continue
		}
	}
	return nil
}

func (importFile *Import) loadSchemaInfo() {
	f, err := os.Open(importFile.fileName)

	if err != nil {
		importFile.logError(err)
		return
	}
	defer f.Close()
	d := xml.NewDecoder(f)
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			importFile.logError(err)
			return
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importSchemaTag {
				importFile.schemaName = importFile.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
				importFile.schemaId = getSchemaId(importFile.schemaName)
				importFile.logger.setSchemaId(importFile.schemaId)
				// add @version parameter
				var metaParam = importFile.getSchemaVersion(importFile.getXmlAttribute(&ty.Attr, importXmlAttributeVersionTag))
				var metaLanguage = importFile.getSchemaLanguage(importFile.getXmlAttribute(&ty.Attr, importXmlAttributeLangTag))
				importFile.metaList = append(importFile.metaList, metaParam)
				importFile.metaList = append(importFile.metaList, metaLanguage)
				return
			}
			continue
		}
	}
}

func (importFile *Import) getSchemaVersion(value string) *meta {
	parameter := new(parameter)

	var schemaVersion = parameter.getVersionParameter(importFile.schemaId, importFile.schemaId, entitytype.Schema, value)
	var metaData = schemaVersion.toMeta(importFile.schemaId)
	metaData.refId = importFile.schemaId
	return metaData
}

func (importFile *Import) getSchemaLanguage(value string) *meta {
	lang := new(Language)
	lang.Init(1, value)
	return lang.toMeta()
}

func (importFile *Import) saveMetaList() error {
	existingMetaList := getMetaList(importFile.schemaId)
	var metaList = importFile.metaList
	var err error

	err = nil
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	// create new schema
	if len(existingMetaList) == 0 {
		for i := 0; i < len(metaList) && err == nil; i++ {
			err = query.insertMeta(metaList[i], importFile.schemaId)
		}
	}

	return err
}

func (importFile *Import) saveMetaIdList() error {
	queryExist := new(metaQuery)
	queryInsert := new(metaQuery)
	metaid := new(metaId)
	metaList := importFile.metaList

	metaid.schemaId = importFile.schemaId
	metaid.objectType = int8(entitytype.Table)
	metaid.value = 0

	queryExist.setSchema(metaSchemaName)
	queryExist.setTable(metaIdTableName)
	queryInsert.setSchema(metaSchemaName)
	queryInsert.setTable(metaIdTableName)

	queryExist.addFilter(metaFieldId, operatorEqual, 0)
	queryExist.addFilter(metaSchemaId, operatorEqual, importFile.schemaId)
	queryExist.addFilter(metaObjectType, operatorEqual, int8(entitytype.Table))

	for i := 0; i < len(metaList); i++ {
		metaData := metaList[i]
		if metaData.GetEntityType() == entitytype.Table {
			queryExist.setParamValue(metaData.id, 0)
			exist, err := queryExist.exists()
			if err != nil {
				return err
			}
			if exist == false {
				metaid.id = metaData.id
				queryInsert.insertMetaId(metaid)
			}
		}
	}

	return nil
}

func (importFile *Import) analyzeMetaTables() {
	var table *Table
	var metaSchema = GetSchemaByName(metaSchemaName)

	table = metaSchema.GetTableByName(metaIdTableName)
	table.Vacuum(importFile.jobId)
	table.Analyze(importFile.jobId)

	table = metaSchema.GetTableByName(metaTableName)
	table.Vacuum(importFile.jobId)
	table.Analyze(importFile.jobId)
}
