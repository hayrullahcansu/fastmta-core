package tcptimer

import (
	"crypto/tls"
	"time"
)

type TLSTCPTimer struct {
	TcpTimer
	conn *tls.Conn
}

func NewTLSTCPTimer(conn *tls.Conn) *TLSTCPTimer {
	return &TLSTCPTimer{
		conn: conn,
	}
}

func (t *TLSTCPTimer) SetReadDeadline(time time.Time) {
	t.conn.SetReadDeadline(time)
}

func (t *TLSTCPTimer) SetReadDeadlineDefault() {
	t.conn.SetReadDeadline(time.Now().Add(ReadDeadLine))
}

func (t *TLSTCPTimer) SetWriteDeadline(time time.Time) {
	t.conn.SetWriteDeadline(time)
}

func (t *TLSTCPTimer) SetWriteDeadlineDefault() {
	t.conn.SetWriteDeadline(time.Now().Add(ReadDeadLine))
}
