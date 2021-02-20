package bulkquerytype

type BulkQueryType int8

//!!! reserved value for unitesting {4, 5, 6} !!!

const (
	Undefined          BulkQueryType = 0
	SimpleQuery        BulkQueryType = 1
	SetRoot            BulkQueryType = 2
	TraverseFromParent BulkQueryType = 3
	TraverseFromRoot   BulkQueryType = 8
)
