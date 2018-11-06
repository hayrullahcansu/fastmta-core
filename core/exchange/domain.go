package exchange

import (
	"sync"

	"../../entity"
	"github.com/golang-collections/go-datastructures/queue"
)

type Domain struct {
	mutex        *sync.Mutex
	Name         string
	MXRecords    []string
	ParentRouter *Router
	MessageQueue *queue.Queue
}

func NewDomain(name string, router *Router) *Domain {
	return &Domain{
		mutex:        &sync.Mutex{},
		Name:         name,
		ParentRouter: router,
		MessageQueue:    queue.New()
	}
}
func (domain *Domain) AddMessage(message *entity.Message) {

}
