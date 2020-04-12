package inbound

import "time"

type InboundHandler interface {
	SetReadDeadline(time time.Time)
	SetReadDeadlineDefault()
	SetWriteDeadline(time time.Time)
	SetWriteDeadlineDefault()
}
