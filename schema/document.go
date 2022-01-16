package schema

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"reflect"
	"ring/schema/databaseprovider"
	"ring/schema/documenttype"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/relationtype"
	"strconv"
	"strings"
	"time"
)

const (
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
)

type document struct {
	schemaId        int32
	filePath        string
	creator         string
	creationTime    *time.Time
	updateTime      *time.Time
	metaList        []*meta
	documentType    documenttype.DocumentType
	logger          *log
	currentJobId    int64
	tablespaceCount int32
	tables          map[string]int32
	provider        databaseprovider.DatabaseProvider
	schemaName      string
	isNewSchema     bool
}

func (doc *document) Init(filePath string, docType documenttype.DocumentType,
	currentJobId int64, provider databaseprovider.DatabaseProvider) {

	doc.logger = new(log)
	doc.currentJobId = currentJobId
	doc.provider = provider
	doc.documentType = docType
	doc.filePath = filePath

}

//******************************
// getters / setters
//******************************
func (doc *document) GetCreationTime() *time.Time {
	return doc.creationTime
}
func (doc *document) GetUpdateTime() *time.Time {
	return doc.updateTime
}
func (doc *document) GetDocumentType() documenttype.DocumentType {
	return doc.documentType
}

//******************************
// public methods
//******************************
func (doc *document) ParseFile(filePath string) error {
	doc.filePath = filePath
	doc.tablespaceCount = 0
	return nil
}

//******************************
// private methods
//******************************

//go:noinline
func (doc *document) logError(err error) {
	doc.logger.writePartialLog(13, levelError, doc.currentJobId, err)
}

//go:noinline
func (doc *document) logFileError(message string, description string, lineNumber int64) {
	description += "\n at line " + strconv.FormatInt(lineNumber, 10)
	doc.logger.writePartialLog(14, levelError, doc.currentJobId, message, description)
}

func (doc *document) isXmlValid() (bool, error) {
	fileInfo, errStat := os.Stat(doc.filePath)
	if errStat != nil {
		doc.logError(errStat)
		return false, errStat
	}
	// get the size
	if fileInfo.Size() > maxXmlFileSize {
		doc.logError(errors.New("File size is bigger than " +
			strconv.Itoa(int(maxXmlFileSize)) + " bytes"))
		return false, errStat
	}
	filerc, err := os.Open(doc.filePath)
	if err != nil {
		doc.logError(err)
		return false, err
	}
	defer filerc.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(filerc)
	contents := buf.Bytes()
	err = xml.Unmarshal(contents, new(interface{}))
	if err != nil {
		doc.logError(err)
	}
	return err == nil, err
}

func (doc *document) loadXml() {
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
	err := doc.loadTableInfo()
	if err != nil {
		// unable to open file
		return
	}

	// pass 2 build --> get []Meta
	f, err = os.Open(doc.filePath)
	if err != nil {
		doc.logError(err)
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
			doc.logError(err)
			return
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			metaData = doc.manageElement(d, &ty, fieldId, relationId, indexId, referenceId)
			if metaData != nil {
				doc.metaList = append(doc.metaList, metaData)
			}
			break
		}
	}
	// de-allow table map
	doc.tables = nil
}

func (doc *document) manageElement(d *xml.Decoder, ty *xml.StartElement, fieldId *int32, relationId *int32, indexId *int32, referenceId *int32) *meta {
	var metaData *meta

	// SCHEMA
	if strings.ToLower(ty.Name.Local) == importSchemaTag {
		metaData = doc.getXmlMetaSchema(&ty.Attr, reflect.ValueOf(d).Elem().FieldByName("line").Int())
	}
	// TABLESPACE
	if strings.ToLower(ty.Name.Local) == importTablespaceTag {
		metaData = doc.getXmlMetaTablespace(&ty.Attr, reflect.ValueOf(d).Elem().FieldByName("line").Int())
	}
	// TABLE
	if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
		*fieldId = 0
		*relationId = 0
		*indexId = 0
		metaData = doc.getXmlMeta(&ty.Attr, entitytype.Table, doc.schemaId, 0,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		metaData.flags = doc.getTableFlags(&ty.Attr)
		metaData.description = doc.getDescription(&ty.Attr)
		*referenceId = metaData.id
	}
	// FIELDS
	if strings.ToLower(ty.Name.Local) == importFieldTag {
		(*fieldId)++
		line := reflect.ValueOf(d).Elem().FieldByName("line").Int()
		metaData = doc.getXmlMeta(&ty.Attr, entitytype.Field, *referenceId, *fieldId, line)
		metaData.flags = doc.getFieldFlags(&ty.Attr, line)
		metaData.description = doc.getDescription(&ty.Attr)
	}
	// RELATIONS
	if strings.ToLower(ty.Name.Local) == importRelationTag {
		(*relationId)++
		metaData = doc.getXmlMeta(&ty.Attr, entitytype.Relation, *referenceId, *relationId,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		metaData.flags = doc.getRelationFlags(&ty.Attr)
		metaData.description = doc.getDescription(&ty.Attr)
	}
	// INDEXES
	if strings.ToLower(ty.Name.Local) == importIndexTag {
		(*indexId)++
		metaData = doc.getXmlMeta(&ty.Attr, entitytype.Index, *referenceId, *indexId,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		metaData.flags = doc.getIndexFlags(&ty.Attr)
		metaData.value = doc.getIndexValue(d, ty)
		metaData.description = doc.getDescription(&ty.Attr)
	}
	// DESCRIPTION
	if strings.ToLower(ty.Name.Local) == importDescriptionTag {
		doc.addDescription(d, ty)
	}
	return metaData
}

func (doc *document) getXmlMetaSchema(attributes *[]xml.Attr, line int64) *meta {
	var result = new(meta)
	var schema = new(Schema)

	result.flags = uint64(doc.provider)
	result.name = doc.getXmlAttribute(attributes, importXmlAttributeNameTag)
	result.objectType = int8(entitytype.Schema)
	result.refId = 0
	result.id = doc.schemaId // copy schema_id
	result.lineNumber = line
	result.enabled = true
	result.value = schema.getPhysicalName(doc.provider, doc.schemaName)

	return result
}

func (doc *document) getXmlMetaTablespace(attributes *[]xml.Attr, line int64) *meta {
	var result = new(meta)
	doc.tablespaceCount++

	result.name = doc.getXmlAttribute(attributes, importXmlAttributeNameTag)
	result.objectType = int8(entitytype.Tablespace)
	result.refId = 0
	result.id = doc.tablespaceCount
	result.lineNumber = line
	result.value = doc.getXmlAttribute(attributes, importXmlAttributeFileTag)
	result.enabled = true
	var isIndex = strings.ToLower(doc.getXmlAttribute(attributes, importIndexTag))
	if isIndex == importXmlBoolTrueValue1 || isIndex == importXmlBoolTrueValue2 {
		result.setTablespaceIndex(true)
	}
	var isTable = strings.ToLower(doc.getXmlAttribute(attributes, importTableTag2))
	if isTable == importXmlBoolTrueValue1 || isTable == importXmlBoolTrueValue2 {
		result.setTablespaceTable(true)
	}

	return result
}

func (doc *document) getXmlMeta(attributes *[]xml.Attr, entityType entitytype.EntityType, referenceId int32, defaultId int32, line int64) *meta {
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
			result.id = doc.getXmlAttributeId(attribute.Value, line)
			break
		case importXmlAttributeTypeTag, importXmlAttributeToTag:
			// DATA_TYPE
			result.dataType = doc.getDataType(entityType, attribute.Name.Local, attribute.Value)
			break
		case importXmlAttributeDefaultTag, importXmlAttributeSubjectTag, importXmlAttributeInverseTag:
			// VALUE
			result.value = doc.getValue(entityType, attribute.Value)
			break
		}
	}
	return result
}

func (doc *document) getXmlAttribute(attributes *[]xml.Attr, attributeName string) string {
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

func (doc *document) getXmlAttributeId(attributeValue string, line int64) int32 {
	result, err := strconv.ParseInt(attributeValue, 10, 64)
	if err == nil {
		if result >= importMinId && result <= importMaxId {
			return int32(result)
		}
		doc.logFileError(errorImportInvalidId, "", line)
		return -2
	}
	doc.logFileError(errorImportInvalidId, err.Error(), line)
	return -1
}

func (doc *document) getDataType(entityType entitytype.EntityType, attributeName string, value string) int32 {
	if entityType == entitytype.Field && strings.ToLower(attributeName) == importXmlAttributeTypeTag {
		return int32(fieldtype.GetFieldType(value))
	}
	if entityType == entitytype.Relation && strings.ToLower(attributeName) == importXmlAttributeToTag {
		if val, ok := doc.tables[strings.ToLower(value)]; ok {
			return val
		}
	}
	return 0
}

func (doc *document) getValue(entityType entitytype.EntityType, value string) string {
	if entityType == entitytype.Field || entityType == entitytype.Relation || entityType == entitytype.Table ||
		entityType == entitytype.Schema {
		// get Field: default, Relation: inverse relationship name,Table: Subject
		return value
	}
	return ""
}

func (doc *document) getTableFlags(attributes *[]xml.Attr) uint64 {
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

func (doc *document) getRelationFlags(attributes *[]xml.Attr) uint64 {
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

func (doc *document) getFieldFlags(attributes *[]xml.Attr, line int64) uint64 {
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
			metaData.setFieldSize(doc.getFieldSize(attributeValue, line))
		}
	}
	return metaData.flags
}

func (doc *document) getIndexFlags(attributes *[]xml.Attr) uint64 {
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

func (doc *document) getIndexValue(d *xml.Decoder, ty *xml.StartElement) string {
	//importIndexFieldTag
	var result strings.Builder

	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			doc.logError(err)
			return ""
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importIndexFieldTag {
				fieldName := doc.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
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
func (doc *document) getDescription(attributes *[]xml.Attr) string {
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

func (doc *document) getFieldSize(value string, line int64) uint32 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		doc.logFileError(errorImportFieldSize, err.Error(), line)
		return 0
	}
	if result >= importMinId && result <= importMaxId {
		return uint32(result)
	}
	doc.logFileError(errorImportFieldSize, "", line)
	return 0
}

func (doc *document) addDescription(d *xml.Decoder, ty *xml.StartElement) {
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			doc.logError(err)
			return
		}
		switch ty := tok.(type) {
		case xml.CharData:
			if len(doc.metaList) > 0 {
				lastMeta := doc.metaList[len(doc.metaList)-1]
				bytes := xml.CharData(ty)
				lastMeta.description = string([]byte(bytes))
			}
			break
		case xml.EndElement:
			return
		}
	}
}

func (doc *document) loadTableInfo() error {
	doc.tables = make(map[string]int32)
	f, err := os.Open(doc.filePath)
	var tableName string
	var tableId int32

	if err != nil {
		doc.logError(err)
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
			doc.logError(err)
			return err
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
				tableName = doc.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
				tableId = doc.getXmlAttributeId(doc.getXmlAttribute(&ty.Attr, importXmlAttributeIdTag1), -1)
				doc.tables[strings.ToLower(tableName)] = tableId
			}
			continue
		}
	}
	return nil
}

func (doc *document) loadSchemaInfo() {
	f, err := os.Open(doc.filePath)

	if err != nil {
		doc.logError(err)
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
			doc.logError(err)
			return
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importSchemaTag {
				doc.schemaName = doc.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
				doc.schemaId, doc.isNewSchema = getSchemaId(doc.schemaName)
				doc.logger.setSchemaId(doc.schemaId)

				// add @version parameter
				var metaParam = doc.getSchemaVersion(doc.getXmlAttribute(&ty.Attr, importXmlAttributeVersionTag))
				var metaLanguage = doc.getSchemaLanguage(doc.getXmlAttribute(&ty.Attr, importXmlAttributeLangTag))
				var metaLastUpgrade = doc.getLastUpgrade()
				var metaCreationTime = doc.getCreationTime()

				doc.metaList = append(doc.metaList, metaParam)
				doc.metaList = append(doc.metaList, metaLastUpgrade)
				doc.metaList = append(doc.metaList, metaLanguage)
				doc.metaList = append(doc.metaList, metaCreationTime)

				return
			}
			continue
		}
	}
}

func (doc *document) getSchemaVersion(value string) *meta {
	parameter := new(parameter)

	var schemaVersion = parameter.getVersionParameter(doc.schemaId, doc.schemaId, entitytype.Schema, value)
	var metaData = schemaVersion.toMeta(0)

	return metaData
}

func (doc *document) getSchemaLanguage(value string) *meta {
	lang := new(parameter)
	var schema = GetSchemaById(doc.schemaId)
	if schema != nil && schema.getParameterByName(parameterDefaultLanguage) != nil {
		lang = schema.getParameterByName(parameterDefaultLanguage)
	} else {
		lang = lang.getLanguageParameter(doc.schemaId, value)
	}
	return lang.toMeta(0)
}

func (doc *document) getLastUpgrade() *meta {
	param := new(parameter)
	param = param.getLastUpgradeParameter(doc.schemaId, 0, entitytype.Schema)
	return param.toMeta(0)
}

func (doc *document) getCreationTime() *meta {
	param := new(parameter)
	var schema = GetSchemaById(doc.schemaId)
	if schema != nil && schema.getParameterByName(parameterCreationTime) != nil {
		param = schema.getParameterByName(parameterCreationTime)
	} else {
		param = param.getCreationTimeParameter(doc.schemaId, 0, entitytype.Schema)
	}
	return param.toMeta(0)
}
