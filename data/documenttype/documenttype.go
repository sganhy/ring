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
	XlsxSufix string = ".XLSX"
)

func GetDocumentType(fileName string) DocumentType {
	if strings.HasSuffix(strings.ToUpper(fileName), XlsxSufix) {
		return Xlsx
	}
	return Undefined
}
