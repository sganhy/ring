package documenttype

import (
	"strings"
)

type DocumentType int8

const (
	Xlsx      DocumentType = 1
	Json      DocumentType = 2
	Undefined DocumentType = 127
)

var (
	XlsxSuffix string = ".XLSX"
)

func GetDocumentType(fileName string) DocumentType {
	if strings.HasSuffix(strings.ToUpper(fileName), XlsxSuffix) {
		return Xlsx
	}
	return Undefined
}
