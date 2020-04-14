package rw

import (
	"fmt"
	"net"
)

type NoTLSDialer struct {
	*BaseDialer
	conn net.Conn
}

func NewNoTLSDialer(dialer *net.Dialer) *NoTLSDialer {
	return &NoTLSDialer{
		BaseDialer: NewBaseDialer(dialer),
	}
}

func (d *NoTLSDialer) Deal(host string, port int) error {
	conn, err := d.dialer.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	d.conn = conn
	d.transporter = NewNoTLSTransporter(conn)
	return err
}
