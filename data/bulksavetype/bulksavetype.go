package bulksavetype

type BulkSaveType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Undefined           BulkSaveType = 0
	DeleteRecord        BulkSaveType = 1
	InsertRecord        BulkSaveType = 2
	UpdateRecord        BulkSaveType = 3
	RelateRecords       BulkSaveType = 7
	BindRelation        BulkSaveType = 8
	InsertMtm           BulkSaveType = 9
	InsertMtmIfNotExist BulkSaveType = 10
	Cancelled           BulkSaveType = 17
)
