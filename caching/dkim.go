package caching

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type DkimCacher struct {
	C *cache.Cache
}

var instanceDkimCacher *DkimCacher
var onceDkimCacher sync.Once

func InstanceDkim() *DkimCacher {
	onceDkimCacher.Do(func() {
		instanceDkimCacher = &DkimCacher{
			C: cache.New(5*time.Minute, 10*time.Minute),
		}
	})
	return instanceDkimCacher
}
