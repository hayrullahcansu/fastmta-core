package core

import (
	"fmt"
	"sync"

	"github.com/hayrullahcansu/fastmta-core/entity"
)

var (
	RouterInstance Router
)

type Router struct {
	BulkChannel            chan *entity.Message
	GeneralChannel         chan *entity.Message
	MessageChannel         chan *entity.Message
	StopChannel            chan bool
	OutboundVirtualMtaPool []*VirtualMta
}

var instanceRouter *Router
var once sync.Once

func InstanceRouter() *Router {
	once.Do(func() {
		instanceRouter = newRouter()
	})
	return instanceRouter
}

func newRouter() *Router {

	return &Router{
		BulkChannel:    make(chan *entity.Message, 1000),
		GeneralChannel: make(chan *entity.Message, 500),
		MessageChannel: make(chan *entity.Message, 1000),
		StopChannel:    make(chan bool),
	}
}

func (router *Router) Init(virtualMtas *[]*VirtualMta) {
	router.OutboundVirtualMtaPool = *virtualMtas
}

func (router *Router) Run() {
	defer close(router.MessageChannel)
	defer close(router.StopChannel)
	for {
		select {
		case message, ok := <-router.MessageChannel:
			if ok {
				router.progressMessage(message)
			}
		case stop := <-router.StopChannel:
			if stop {
				break
			}
		}
	}
	go router.progressBulkMessage()
	go router.progressGeneralMessage()

}

func (router *Router) RelayMessage(message *entity.Message) {
	router.MessageChannel <- message
}

func (router *Router) progressMessage(message *entity.Message) {
	switch message.Host {
	case "gmail.com":
	case "yandex.com":
	case "yandex.com.tr":
		router.BulkChannel <- message
		break
	default:
		router.GeneralChannel <- message
	}

	// if !ok {
	// domain = NewDomain(message.Host, router)
	// domain.ParentRouter = router
	// router.Domains[message.Host] = domain
	// go domain.Run()
	// }
	// domain.AddMessage(message)
}

func (router *Router) Stop() {
	router.StopChannel <- true
}

func (router *Router) progressBulkMessage() {
	for {
		select {
		case bulk, ok := <-router.BulkChannel:
			if ok {

				fmt.Println(bulk.Host)
			}
		}
	}
}
func (router *Router) progressGeneralMessage() {
	for {
		select {
		case general, ok := <-router.GeneralChannel:
			if ok {
				fmt.Println(general.Host)
			}
		}
	}
}

func (router *Router) GetVirtualMta() (*VirtualMta, bool) {
	for _, value := range router.OutboundVirtualMtaPool {
		if !value.IsInUsage() {
			return value, true
		}
	}
	return nil, false
}
