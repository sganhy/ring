package schema

import (
	"sync"
)

type CacheId struct {
	CurrentId     int64
	MaxId         int64
	ReservedRange int32
	syncRoot      sync.Mutex
}
