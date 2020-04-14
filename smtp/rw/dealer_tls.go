package rw

import (
	"crypto/tls"
	"fmt"
	"net"
)

type TLSDialer struct {
	*BaseDialer
	conn *tls.Conn
}

func NewTLSDialer(dialer *net.Dialer) *TLSDialer {
	return &TLSDialer{
		BaseDialer: NewBaseDialer(dialer),
	}
}

func (d *TLSDialer) Deal(host string, port int) error {
	conn, err := tls.DialWithDialer(d.dialer, "tcp", fmt.Sprintf("%s:%d", host, port), nil)
	d.conn = conn
	d.transporter = NewTLSTransporter(conn)
	return err
}
