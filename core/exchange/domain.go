package exchange

import (
	"sync"
)

type Domain struct {
	mutex        *sync.Mutex
	Name         string
	MXRecords    []string
	ParentRouter *Router
}

func New(name string, router *Router) *Domain {
	return &Domain{
		mutex:        &sync.Mutex{},
		Name:         name,
		ParentRouter: router,
	}
}
