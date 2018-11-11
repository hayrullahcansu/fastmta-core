package exchange

import (
	".."
	"../../entity"
	"../smtp"
)

type GeneralSender struct {
	MessageChannel chan *entity.Message
	ParentRouter   *Router
	virtualMta     *smtp.VirtualMta
}

func NewGeneralSender(router *Router) *GeneralSender {
	return &GeneralSender{
		MessageChannel: make(chan *entity.Message, 100),
		ParentRouter:   router,
	}
}

func (sender *GeneralSender) AssignVirtualMta() {
	mta, ok := sender.ParentRouter.GetVirtualMta()
	if ok {
		sender.virtualMta = mta
	}
}

func (sender *GeneralSender) Run() {
	for {
		select {
		case msg, ok := <-sender.MessageChannel:
			if ok {
				_, err := core.NewDomain(msg.Host)
				if err != nil {
					//TODO: this is bounce domain not found
					_ = smtp.NewOutboundClient()
					//transactionResult := client.SendMessage(msg, nil, domain)
					//fmt.Println(transactionResult)
				}

			}
		}
	}
}
