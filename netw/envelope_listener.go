package netw

type EnvelopeListener interface {
	OnNotify(notify *Notify)
	OnEvent(c interface{}, event *Event)
	OnRegister(c interface{}, register *Register)
	OnMessage(c interface{}, message *Message)
}
