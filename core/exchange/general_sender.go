package exchange

import "../smtp"

type GeneralSender struct {
	domain       Domain
	ParentRouter *Router
	virtualMta   *smtp.VirtualMta
}

func NewGeneralSender(router *Router) *GeneralSender {
	return &GeneralSender{
		ParentRouter: router,
	}
}

func (sender *GeneralSender) AssignVirtualMta() {
	mta, ok := sender.ParentRouter.GetVirtualMta()
	if ok {
		sender.virtualMta = mta
	}
}

func (sender *GeneralSender) Send() {

}
