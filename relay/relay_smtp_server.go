package relay

import (
	"net"

	"github.com/hayrullahcansu/fastmta-core/mta"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/hayrullahcansu/zetamail/queue"
)

type RelaySmtpServer struct {
	ID             string
	VMta           *mta.VirtualMta
	VmtaHostName   string
	VmtaIPAddr     *net.IPAddr
	Port           int
	rabbitMqClient *rabbit.RabbitMqClient
}

func RelaySmtpServer(vmta *mta.VirtualMta) *RelaySmtpServer {
	client := queue.New()
	client.Connect(true)

	return &RelaySmtpServer{
		ID:             "",
		VMta:           vmta,
		VmtaHostName:   vmta.VmtaHostName,
		VmtaIPAddr:     vmta.VmtaIPAddr,
		Port:           vmta.Port,
		rabbitMqClient: client,
	}
}
