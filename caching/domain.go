package caching

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// DomainCacher caches domain data in memory.
type DomainCacher struct {
	C *cache.Cache
}

var instanceDomainCacher *DomainCacher
var onceDomainCacher sync.Once

// InstanceDomain returns new or existing instance of DomainCacher.
func InstanceDomain() *DomainCacher {
	onceDomainCacher.Do(func() {
		instanceDomainCacher = &DomainCacher{
			C: cache.New(5*time.Minute, 10*time.Minute),
		}
	})
	return instanceDomainCacher
}
