package tcptimer

import (
	"net"
	"time"
)

type NoTLSTCPTimer struct {
	TcpTimer
	conn net.Conn
}

func NewNoTLSTCPTimer(conn net.Conn) *NoTLSTCPTimer {
	return &NoTLSTCPTimer{
		conn: conn,
	}
}

func (t *NoTLSTCPTimer) SetReadDeadline(time time.Time) {
	t.conn.SetReadDeadline(time)
}

func (t *NoTLSTCPTimer) SetReadDeadlineDefault() {
	t.conn.SetReadDeadline(time.Now().Add(ReadDeadLine))
}

func (t *NoTLSTCPTimer) SetWriteDeadline(time time.Time) {
	t.conn.SetWriteDeadline(time)
}

func (t *NoTLSTCPTimer) SetWriteDeadlineDefault() {
	t.conn.SetWriteDeadline(time.Now().Add(ReadDeadLine))
}
