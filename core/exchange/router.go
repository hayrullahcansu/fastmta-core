package exchange

import (
	"../../entity"
)

type Router struct {
	Domains        map[string]*Domain
	MessageChannel chan *entity.Message
	StopChannel    chan bool
}

func NewRouter() *Router {
	return &Router{
		Domains:        make(map[string]*Domain),
		MessageChannel: make(chan *entity.Message),
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
}

func (router *Router) RelayMessage(message *entity.Message) {
	router.MessageChannel <- message
}

func (router *Router) progressMessage(message *entity.Message) {
	domain, ok := router.Domains[message.Host]
	if !ok {
		domain = NewDomain(message.Host, router)
	}
	domain.ParentRouter = router
}

func (router *Router) Stop() {
	router.StopChannel <- true
}
