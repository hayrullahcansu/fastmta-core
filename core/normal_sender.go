package core

import (
	"github.com/hayrullahcansu/zetamail/entity"
)

type NormalSender struct {
	MessageChannel chan *entity.Message
	ParentRouter   *Router
	virtualMta     *VirtualMta
}

func NewGeneralSender(router *Router) *NormalSender {
	return &NormalSender{
		MessageChannel: make(chan *entity.Message, 100),
		ParentRouter:   router,
	}
}

func (sender *NormalSender) AssignVirtualMta() {
	mta, ok := sender.ParentRouter.GetVirtualMta()
	if ok {
		sender.virtualMta = mta
	}
}

func (sender *NormalSender) Run() {
	for {
		select {
		case msg, ok := <-sender.MessageChannel:
			if ok {
				_, err := NewDomain(msg.Host)
				if err != nil {
					//TODO: this is bounce domain not found
					_ = NewOutboundClient()
					//transactionResult := client.SendMessage(msg, nil, domain)
					//fmt.Println(transactionResult)
				}

			}
		}
	}
}
