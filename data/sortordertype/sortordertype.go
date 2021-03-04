package sortordertype

type SortOrderType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	// check collision with operationType
	Ascending  SortOrderType = 101
	Descending SortOrderType = 102
)
