package mime

import "strings"

type HeaderCollection struct {
	M map[string]string
}

func (h *HeaderCollection) Add(key, value string) {
	(*h).M[key] = value
}

func (h *HeaderCollection) Get(key string) string {
	return h.M[key]
}

func NewHeaderCollection() *HeaderCollection {
	return &HeaderCollection{
		M: make(map[string]string),
	}
}

//GetData searchs case sensitive
func (h *HeaderCollection) GetData(key string) (d string, ok bool) {
	ok = false
	k := ""
	for k, d = range (*h).M {
		if k == key {
			ok = true
			return
		}
	}
	return
}

//GetDataI searchs case insensitive
func (h *HeaderCollection) GetDataI(key string) (d string, ok bool) {
	ok = false
	k := ""
	for k, d = range (*h).M {
		if strings.EqualFold(k, key) {
			ok = true
			return
		}
	}
	return
}
