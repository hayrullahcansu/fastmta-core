package smtp

import (
	"net"

	"../../global"
	"../../queue"
)

type OutboundSmtpServer struct {
	ID             string
	VMta           *VirtualMta
	VmtaHostName   string
	VmtaIPAddr     *net.IPAddr
	Port           int
	RabbitMqClient *queue.RabbitMqClient
}

func CreateNewOutboundSmtpServer(vmta *VirtualMta) *OutboundSmtpServer {
	client := queue.New(global.StaticRabbitMqConfig)
	client.Connect(true)
	client.ExchangeDeclare(queue.OutboundExchange, true, false, false, false, nil)
	//que, _ := client.QueueDeclare(queue.InboundStagingQueueName, true, false, false, false, nil)
	//client.QueueBind(que.Name, queue.InboundExchange, "", false, nil)
	return &OutboundSmtpServer{
		ID:             "",
		VMta:           vmta,
		VmtaHostName:   vmta.VmtaHostName,
		VmtaIPAddr:     vmta.VmtaIPAddr,
		Port:           vmta.Port,
		RabbitMqClient: client,
	}
}
