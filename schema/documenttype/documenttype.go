package documenttype

import (
	"strings"
)

type DocumentType int8

const (
	Xml       DocumentType = 1
	Json      DocumentType = 2
	Undefined DocumentType = 127
)

var (
	XmlSufix string = ".XML"
)

func GetDocumentType(fileName string) DocumentType {
	if strings.HasSuffix(strings.ToUpper(fileName), XmlSufix) {
		return Xml
	}
	return Undefined
}
