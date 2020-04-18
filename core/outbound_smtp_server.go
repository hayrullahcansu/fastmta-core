package core

import (
	"net"

	"github.com/hayrullahcansu/fastmta-core/mta"
	"github.com/hayrullahcansu/fastmta-core/queue"
)

type OutboundSmtpServer struct {
	ID             string
	VMta           *mta.VirtualMta
	VmtaHostName   string
	VmtaIPAddr     *net.IPAddr
	Port           int
	RabbitMqClient *queue.RabbitMqClient
}

func CreateNewOutboundSmtpServer(vmta *mta.VirtualMta) *OutboundSmtpServer {
	client := queue.New()
	client.Connect(true)

	return &OutboundSmtpServer{
		ID:             "",
		VMta:           vmta,
		VmtaHostName:   vmta.VmtaHostName,
		VmtaIPAddr:     vmta.VmtaIPAddr,
		Port:           vmta.Port,
		RabbitMqClient: client,
	}
}
