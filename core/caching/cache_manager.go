package caching

import (
	".."
)

type CacheManager struct {
	Domains map[string]*core.Domain
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Domains: make(map[string]*core.Domain),
	}
}

func (cm *CacheManager) Init() {

}

func (cm *CacheManager) Get(key string) (*core.Domain, bool) {
	val, ok := cm.Domains[key]
	return val, ok
}

func (cm *CacheManager) Add(key string, value *core.Domain) bool {
	_, ok := cm.Domains[key]
	if !ok {
		cm.Domains[key] = value
		return true
	}
	return false
}
