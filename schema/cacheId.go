package schema

type CacheId struct {
	CurrentId     int64
	MaxId         int64
	ReservedRange int32
}
