package sourcetype

type SourceType int8

const (
	XmlDocument    SourceType = 5
	NativeDataBase SourceType = 6
	Undefined      SourceType = 127
)
