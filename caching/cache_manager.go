package caching

type CacheManager struct {
	//Items map[string]*core.Domain
	Items map[string]interface{}
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Items: make(map[string]interface{}),
	}
}

func (cm *CacheManager) Init() {

}

func (cm *CacheManager) Get(key string) (interface{}, bool) {
	val, ok := cm.Items[key]
	return val, ok
}

func (cm *CacheManager) Add(key string, value interface{}) bool {
	_, ok := cm.Items[key]
	if !ok {
		cm.Items[key] = value
		return true
	}
	return false
}
