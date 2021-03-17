package bulkretrievetype

type BulkRetrieveType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Undefined          BulkRetrieveType = 0
	SimpleQuery        BulkRetrieveType = 1
	SetRoot            BulkRetrieveType = 2
	TraverseFromParent BulkRetrieveType = 3
	TraverseFromRoot   BulkRetrieveType = 8
)
