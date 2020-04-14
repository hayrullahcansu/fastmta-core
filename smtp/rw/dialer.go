package rw

type Dialer interface {
	Deal(host string, port int) error
	GetTransporter() Transporter
}
