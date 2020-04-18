package core

import (
	"github.com/hayrullahcansu/fastmta-core/dns"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/mta"
	"github.com/hayrullahcansu/fastmta-core/smtp/outbound"
)

type NormalSender struct {
	MessageChannel chan *entity.Message
	ParentRouter   *Router
	virtualMta     *mta.VirtualMta
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
				_, err := dns.NewDomain(msg.Host)
				if err != nil {
					//TODO: this is bounce domain not found
					agent := outbound.NewAgent(sender.virtualMta)
					_, transactionResult := agent.SendMessage(msg)
					logger.Info(transactionResult)
				}

			}
		}
	}
}
