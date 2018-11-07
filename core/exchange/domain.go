package exchange

import (
	"net"
	"sync"

	"../../entity"
	"../exchange"
	"github.com/emnl/goods/queue"
)

type Domain struct {
	mutex        *sync.Mutex
	Name         string
	MXRecords    []*net.MX
	ParentRouter *Router
	Messages     *queue.Queue
}

func NewDomain(name string, router *Router) *Domain {
	domain := &Domain{
		mutex:        &sync.Mutex{},
		Name:         name,
		ParentRouter: router,
		Messages:     queue.New(),
	}
	mx, err := net.LookupMX(name)
	if err == nil {
		domain.MXRecords = mx
	}
	return domain
}

func (domain *Domain) Run(outboundClient *exchange.OutboundSmtpServer) {
	outboundClient.ConsumeMessage(domain.)
}

func (domain *Domain) AddMessage(message *entity.Message) {
	domain.Messages.Enqueue(message)
}

func (domain *Domain) RemoveFromRouter() {

}
