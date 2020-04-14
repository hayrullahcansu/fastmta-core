package rw

import "net"

type BaseDialer struct {
	Dialer
	dialer      *net.Dialer
	transporter Transporter
}

func NewBaseDialer(dialer *net.Dialer) *BaseDialer {
	return &BaseDialer{
		dialer: dialer,
	}
}
func (d *BaseDialer) GetTransporter() Transporter {
	return d.transporter
}
