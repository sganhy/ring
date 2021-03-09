package dmlstatement

type DmlStatement int8

const (
	Insert DmlStatement = 1
	Update DmlStatement = 2
	Delete DmlStatement = 3
)
