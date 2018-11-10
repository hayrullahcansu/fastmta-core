package caching

import (
	"../exchange"
)

type CacheManager struct {
	Domains map[string]*exchange.Domain
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Domains: make(map[string]*exchange.Domain),
	}
}

func (cm *CacheManager) Init() {

}

func (cm *CacheManager) Get(key string) (*exchange.Domain, bool) {
	val, ok := cm.Domains[key]
	return val, ok
}

func (cm *CacheManager) Add(key string, value *exchange.Domain) bool {
	_, ok := cm.Domains[key]
	if !ok {
		cm.Domains[key] = value
		return true
	}
	return false
}
