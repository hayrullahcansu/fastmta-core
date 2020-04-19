package relay

import (
	"net"

	"github.com/hayrullahcansu/fastmta-core/mta"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/hayrullahcansu/zetamail/queue"
)

// SMTPServer connects to target and relays messages
type SMTPServer struct {
	ID             string
	VMta           *mta.VirtualMta
	VmtaHostName   string
	VmtaIPAddr     *net.IPAddr
	Port           int
	rabbitMqClient *rabbit.RabbitMqClient
}

// NewSMTPServer returns new instance of SMTPServer
func NewSMTPServer(vmta *mta.VirtualMta) *SMTPServer {
	client := queue.New()
	client.Connect(true)

	return &SMTPServer{
		ID:             "",
		VMta:           vmta,
		VmtaHostName:   vmta.VmtaHostName,
		VmtaIPAddr:     vmta.VmtaIPAddr,
		Port:           vmta.Port,
		rabbitMqClient: client,
	}
}
