package rw

import (
	"bufio"
	"net"

	"github.com/hayrullahcansu/fastmta-core/smtp/tcptimer"

	"github.com/hayrullahcansu/fastmta-core/cross"
)

type NoTLSTransporter struct {
	*BaseTransporter
	conn net.Conn
}

func NewNoTLSTransporter(connection net.Conn) *NoTLSTransporter {
	return &NoTLSTransporter{
		conn: connection,
		BaseTransporter: NewBaseTransporter(
			bufio.NewReader(connection),
			tcptimer.NewNoTLSTCPTimer(connection),
		),
	}
}

func (t *NoTLSTransporter) WriteLine(data string) error {
	t.BaseTransporter.timer.SetWriteDeadlineDefault()
	_, err := t.conn.Write([]byte(data + cross.NewLine))
	return err
}

func (t *NoTLSTransporter) Close() {
	t.conn.Close()
}
