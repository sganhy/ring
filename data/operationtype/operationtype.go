package operationtype

type OperationType int8

//!!! reserved value for unit testing {4, 5, 6} !!!

const (
	Equal          OperationType = 1
	NotEqual       OperationType = 2
	Greater        OperationType = 3
	GreaterOrEqual OperationType = 10
	Less           OperationType = 11
	LessOrEqual    OperationType = 12
	Like           OperationType = 13
	NotLike        OperationType = 14
	//SoundsLike = 15, -- one day ...
	In OperationType = 17
)
