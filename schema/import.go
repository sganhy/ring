package schema

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"ring/schema/entitytype"
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
}

var currentSchemaImportId = 1

const errorImportFileNotInitialized = "Import File object is not initialized."
const importTableTag1 = "object"
const importTableTag2 = "table"
const importFieldTag = "field"
const importRelationTag = "relation"
const importIndexTag = "index"
const importXmlAttributeNameTag = "name"
const importXmlAttributeIdTag1 = "id"
const importXmlAttributeIdTag2 = "type_id"

func (importFile *Import) Init(source sourcetype.SourceType, fileName string) {
	currentSchemaImportId++
	importFile.id = currentSchemaImportId
	importFile.fileName = fileName
	importFile.initialized = true
	importFile.source = source
}

func (importFile *Import) Load() error {
	if importFile.initialized == false {
		return errors.New(errorImportFileNotInitialized)
	}
	if importFile.source == sourcetype.XmlDocument {
		return importFile.loadXml()
	}
	return nil
}

//******************************
// private methods
//******************************
func (importFile *Import) loadXml() error {
	f, err := os.Open(importFile.fileName)

	if err != nil {
		return err
	}
	defer f.Close()
	d := xml.NewDecoder(f)
	var referenceId int32 = 0
	//var schemaId int32 = 0
	var fieldId int32 = 0
	var relationId int32 = 0
	var indexId int32 = 0
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			return err
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if strings.ToLower(ty.Name.Local) == importTableTag1 || strings.ToLower(ty.Name.Local) == importTableTag2 {
				fieldId = 0
				relationId = 0
				indexId = 0
				var meta = getXmlMeta(&ty.Attr, entitytype.Table, 0, 0)
				referenceId = meta.id
				fmt.Println(meta.ToString())
			}
			if strings.ToLower(ty.Name.Local) == importFieldTag {
				fieldId++
				var meta = getXmlMeta(&ty.Attr, entitytype.Field, referenceId, fieldId)
				fmt.Println(meta.ToString())
			}
			if strings.ToLower(ty.Name.Local) == importRelationTag {
				relationId++
				var meta = getXmlMeta(&ty.Attr, entitytype.Relation, referenceId, relationId)
				fmt.Println(meta.ToString())
			}
			if strings.ToLower(ty.Name.Local) == importIndexTag {
				indexId++
				var meta = getXmlMeta(&ty.Attr, entitytype.Relation, referenceId, indexId)
				fmt.Println(meta.ToString())
			}
		default:
		}
	}
	return nil
}

func getXmlMeta(attributes *[]xml.Attr, entityType entitytype.EntityType, referenceId int32, defaultId int32) *Meta {
	var result = new(Meta)
	count := len(*attributes)
	result.objectType = int8(entityType)
	result.id = defaultId
	result.refId = referenceId
	for i := 0; i < count; i++ {
		var attribute = (*attributes)[i]
		switch strings.ToLower(attribute.Name.Local) {
		case importXmlAttributeNameTag:
			result.name = attribute.Value
			break
		case importXmlAttributeIdTag1, importXmlAttributeIdTag2:
			result.id = getXmlAttributeId(attribute.Value)
			break
		}
	}
	return result
}

func getXmlAttribute(attributes *[]xml.Attr, attributeName string) string {
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

func getXmlAttributeId(attributeValue string) int32 {
	result, err := strconv.ParseInt(attributeValue, 10, 32)
	if err == nil {
		return int32(result)
	}
	return -1
}
