package caching

import (
	"sync"
	"time"

	"github.com/hayrullahcansu/fastmta-core/dns"
	"github.com/patrickmn/go-cache"
)

// DomainCacher caches domain data in memory.
type DomainCacher struct {
	c *cache.Cache
}

var instanceDomainCacher *DomainCacher
var onceDomainCacher sync.Once

// InstanceDomain returns new or existing instance of DomainCacher.
func InstanceDomain() *DomainCacher {
	onceDomainCacher.Do(func() {
		instanceDomainCacher = newInstanceDomain(5 * time.Minute)
	})
	return instanceDomainCacher
	return instanceDomainCacher
}

func newInstanceDomain(d time.Duration) *DomainCacher {
	return &DomainCacher{
		c: cache.New(d, 10*time.Minute),
	}
}

func (c *DomainCacher) AddOrUpdate(k string, domain *dns.Domain) {
	c.c.SetDefault(k, domain)
}

func (c *DomainCacher) Get(k string) *dns.Domain {
	if dd, ok := c.c.Get(k); ok {
		return dd.(*dns.Domain)
	}
	return nil
}
