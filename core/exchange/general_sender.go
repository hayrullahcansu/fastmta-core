package exchange

type GeneralSender struct {
	domain       Domain
	ParentRouter *Router
}

func NewGeneralSender(router *Router) *GeneralSender {
	return &GeneralSender{
		ParentRouter: router,
	}
}

func (sender *GeneralSender) Send() {

}
