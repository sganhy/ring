package sortordertype

type SortOrderType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Ascending  SortOrderType = 1
	Descending SortOrderType = -1
)
