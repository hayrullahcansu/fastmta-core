package tcptimer

import "time"

type TcpTimer interface {
	SetReadDeadline(time time.Time)
	SetReadDeadlineDefault()
	SetWriteDeadline(time time.Time)
	SetWriteDeadlineDefault()
}
