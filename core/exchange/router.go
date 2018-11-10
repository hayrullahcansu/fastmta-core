package exchange

import (
	"fmt"

	"../../entity"
	"../caching"
	"../smtp"
)

type Router struct {
	BulkChannel            chan *entity.Message
	GeneralChannel         chan *entity.Message
	MessageChannel         chan *entity.Message
	StopChannel            chan bool
	OutboundVirtualMtaPool []*smtp.VirtualMta
	CacheManager           *caching.CacheManager
}

func NewRouter() *Router {

	return &Router{
		BulkChannel:    make(chan *entity.Message, 1000),
		GeneralChannel: make(chan *entity.Message, 500),
		MessageChannel: make(chan *entity.Message, 1000),
		StopChannel:    make(chan bool),
	}
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

func (router *Router) SetCacheManager(cacheManager *caching.CacheManager) {
	router.CacheManager = cacheManager
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

func (router *Router) GetVirtualMta() (*smtp.VirtualMta, bool) {
	for _, value := range router.OutboundVirtualMtaPool {
		if !value.IsInUsage() {
			return value, true
		}
	}
	return nil, false
}
