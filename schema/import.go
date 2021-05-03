package schema

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/sourcetype"
	"strconv"
	"strings"
)

type Import struct {
	id          int
	fileName    string
	source      sourcetype.SourceType
	initialized bool
	schema      *Schema
	tables      map[string]int32
	errorList   []error
	metaList    []*Meta
}

const (
	errorImportFileNotInitialized = "Import File object is not initialized."
	importTableTag1               = "object"
	importTableTag2               = "table"
	importFieldTag                = "field"
	importDescriptionTag          = "description"
	importRelationTag             = "relation"
	importIndexTag                = "index"
	importXmlAttributeNameTag     = "name"
	importXmlAttributeSubjectTag  = "suject"
	importXmlAttributeTypeTag     = "type"
	importXmlAttributeDefaultTag  = "default"
	importXmlAttributeIdTag1      = "id"
	importXmlAttributeIdTag2      = "type_id"
	importXmlAttributeToTag       = "to"
	importXmlAttributeInverseTag  = "inverse_relation"
	importXmlAttributeBaseline    = "baseline"
	importXmlAttributeCached      = "cached"
	importXmlBoolTrueValue1       = "true"
	importXmlBoolTrueValue2       = "1"
	importXmlBoolFalseValue1      = "false"
	importXmlBoolFalseValue2      = "0"
	importXmlAttributeSize        = "size"
	importXmlAttributeReadonly    = "readonly"
	importXmlAttributeNotNull     = "not_null"
	importXmlAttributeSensitive   = "case_sensitive"
	importXmlAttributeMultiLang   = "multilingual"
)

var (
	currentSchemaImportId = 1
)

func (importFile *Import) Init(source sourcetype.SourceType, fileName string) {
	currentSchemaImportId++
	importFile.id = currentSchemaImportId
	importFile.fileName = fileName
	importFile.initialized = true
	importFile.source = source
}

func (importFile *Import) Load() {
	if importFile.initialized == false {
		errors.New(errorImportFileNotInitialized)
	}
	if importFile.source == sourcetype.XmlDocument {
		importFile.loadXml()
		for i := 0; i < len(importFile.metaList); i++ {
			meta := importFile.metaList[i]
			if meta.objectType == 0 {
				fmt.Println("TABLE ==> " + meta.name)
			}
			if meta.objectType == 1 {
				var field = meta.toField()
				fmt.Println(field.String())
			}
		}
	}
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
	var meta *Meta
	var f *os.File

	*referenceId = 0
	*fieldId = 0
	*relationId = 0
	*indexId = 0
	importFile.errorList = make([]error, 2)
	importFile.metaList = make([]*Meta, 0, 20)

	// pass 1 build --> dico<lower(table_name),table_id>
	err := importFile.loadTableDico()
	if err != nil {
		importFile.errorList = append(importFile.errorList, err)
		return
	}

	// pass 2 build --> get []Meta
	f, err = os.Open(importFile.fileName)
	if err != nil {
		importFile.errorList = append(importFile.errorList, err)
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
			importFile.errorList = append(importFile.errorList, err)
			return
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			meta = importFile.manageElement(d, &ty, fieldId, relationId, indexId, referenceId)
			if meta != nil {
				importFile.metaList = append(importFile.metaList, meta)
			}
			break
		}
	}
	// de-allow table map
	importFile.tables = nil
}

func (importFile *Import) manageElement(d *xml.Decoder, ty *xml.StartElement, fieldId *int32, relationId *int32, indexId *int32, referenceId *int32) *Meta {
	var meta *Meta
	// TABLE
	if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
		*fieldId = 0
		*relationId = 0
		*indexId = 0
		meta = importFile.getXmlMeta(&ty.Attr, entitytype.Table, 0, 0, reflect.ValueOf(d).Elem().FieldByName("line").Int())
		meta.flags = importFile.getTableFlags(&ty.Attr)
		meta.description = importFile.getDescription(&ty.Attr)
		*referenceId = meta.id
	}
	// FIELDS
	if strings.ToLower(ty.Name.Local) == importFieldTag {
		(*fieldId)++
		line := reflect.ValueOf(d).Elem().FieldByName("line").Int()
		meta = importFile.getXmlMeta(&ty.Attr, entitytype.Field, *referenceId, *fieldId, line)
		meta.flags = importFile.getFieldFlags(&ty.Attr)
		meta.description = importFile.getDescription(&ty.Attr)
	}
	// RELATIONS
	if strings.ToLower(ty.Name.Local) == importRelationTag {
		(*relationId)++
		meta = importFile.getXmlMeta(&ty.Attr, entitytype.Relation, *referenceId, *relationId,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		meta.description = importFile.getDescription(&ty.Attr)
	}
	// INDEXES
	if strings.ToLower(ty.Name.Local) == importIndexTag {
		(*indexId)++
		meta = importFile.getXmlMeta(&ty.Attr, entitytype.Relation, *referenceId, *indexId,
			reflect.ValueOf(d).Elem().FieldByName("line").Int())
		meta.description = importFile.getDescription(&ty.Attr)
	}
	// DESCRIPTION
	if strings.ToLower(ty.Name.Local) == importDescriptionTag {
		importFile.addDescription(d, ty)
	}
	return meta
}

func (importFile *Import) getXmlMeta(attributes *[]xml.Attr, entityType entitytype.EntityType, referenceId int32, defaultId int32, line int64) *Meta {
	var result = new(Meta)
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
			result.id = importFile.getXmlAttributeId(attribute.Value)
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

func (importFile *Import) getXmlAttributeId(attributeValue string) int32 {
	result, err := strconv.ParseInt(attributeValue, 10, 32)
	if err == nil {
		return int32(result)
	}
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
	if entityType == entitytype.Field || entityType == entitytype.Relation || entityType == entitytype.Table {
		return value
	}
	return ""
}

func (importFile *Import) getTableFlags(attributes *[]xml.Attr) uint64 {
	meta := Meta{}
	count := len(*attributes)
	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		if strings.EqualFold(attribute.Value, importXmlBoolTrueValue1) || strings.EqualFold(attribute.Value, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeBaseline:
				// BASELINE
				meta.setEntityBaseline(true)
				break
			case importXmlAttributeReadonly:
				meta.setTableReadonly(true)
				break
			case importXmlAttributeCached:
				meta.setTableCached(true)
				break
			}
		}
	}
	return meta.flags
}

func (importFile *Import) getFieldFlags(attributes *[]xml.Attr) uint64 {
	meta := Meta{}
	count := len(*attributes)
	meta.setFieldCaseSensitive(true)
	meta.setFieldSize(0)
	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		var attributeValue = strings.Trim(attribute.Value, " ")

		if strings.EqualFold(attributeValue, importXmlBoolTrueValue1) || strings.EqualFold(attributeValue, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeBaseline:
				// BASELINE
				meta.setEntityBaseline(true)
				break
			case importXmlAttributeNotNull:
				meta.setFieldNotNull(true)
				break
			case importXmlAttributeMultiLang:
				meta.setFieldMultilingual(true)
				break
			}
		}
		if strings.EqualFold(attributeValue, importXmlBoolTrueValue1) || strings.EqualFold(attributeValue, importXmlBoolTrueValue2) {
			switch strings.ToLower(attribute.Name.Local) {
			case importXmlAttributeSensitive:
				meta.setFieldCaseSensitive(false)
				break
			}
		}
		if strings.ToLower(attribute.Name.Local) == importXmlAttributeSize {
			meta.setFieldSize(importFile.getFieldSize(attributeValue))
		}
	}
	return meta.flags
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

func (importFile *Import) getFieldSize(value string) uint32 {
	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		importFile.errorList = append(importFile.errorList, err)
	}
	return uint32(result)
}

func (importFile *Import) addDescription(d *xml.Decoder, ty *xml.StartElement) {
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			importFile.errorList = append(importFile.errorList, err)
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

func (importFile *Import) loadTableDico() error {
	importFile.tables = make(map[string]int32)
	f, err := os.Open(importFile.fileName)
	var tableName string
	var tableId int32

	if err != nil {
		importFile.errorList = append(importFile.errorList, err)
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
			importFile.errorList = append(importFile.errorList, err)
			return err
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
				tableName = importFile.getXmlAttribute(&ty.Attr, importXmlAttributeNameTag)
				tableId = importFile.getXmlAttributeId(importFile.getXmlAttribute(&ty.Attr, importXmlAttributeIdTag1))
				importFile.tables[strings.ToLower(tableName)] = tableId
			}
			continue
		}
	}
	return nil
}
