package schema

import (
	"ring/schema"
)

const unknowFieldDataType string = ""

type Record struct {
	data       []string
	recordType *schema.Table
}

func (record *Record) SetRecordType(recordType string) {

}
