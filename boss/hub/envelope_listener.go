package hub

type EnvelopeListener interface {
	OnNotify(notify *Notify)
}
