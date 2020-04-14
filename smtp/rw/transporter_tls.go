package rw

import (
	"bufio"
	"crypto/tls"

	"github.com/hayrullahcansu/fastmta-core/smtp/tcptimer"

	"github.com/hayrullahcansu/fastmta-core/cross"
)

type TLSTransporter struct {
	*BaseTransporter
	conn *tls.Conn
}

func NewTLSTransporter(connection *tls.Conn) *TLSTransporter {
	return &TLSTransporter{
		conn: connection,
		BaseTransporter: NewBaseTransporter(
			bufio.NewReader(connection),
			tcptimer.NewTLSTCPTimer(connection),
		),
	}
}

func (t *TLSTransporter) WriteLine(data string) error {
	t.BaseTransporter.timer.SetWriteDeadlineDefault()
	_, err := t.conn.Write([]byte(data + cross.NewLine))
	return err
}

func (t *TLSTransporter) Close() {
	t.conn.Close()
}
