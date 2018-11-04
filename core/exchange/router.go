package exchange

import (
	"../../entity"
)

type Router struct {
	RelayMessage chan *entity.MessageTransaction
	StopChannel  chan bool
}

func New() *Router {
	return &Router{
		RelayMessage: make(chan *entity.MessageTransaction),
		StopChannel:  make(chan bool),
	}
}

func (router *Router) Run() {
	for {
		select {
		case message, ok := <-router.RelayMessage:
			if ok {
				RelayMessage(message)
			}
		case stop := <-router.StopChannel:
			if stop {
				break
			}
		}
	}
}

func RelayMessage(message *entity.MessageTransaction) {
	//SaveModel
	//
}

func (router *Router) Stop() {
	router.StopChannel <- true
}
