package boss

import (
	"../conf"
	ZMSmtp "../core/smtp"
	"../queue"
)

type Boss struct {
	Config                *conf.Config
	VirtualMtas           []*ZMSmtp.VirtualMta
	InboundMtas           []*ZMSmtp.SmtpServer
	InboundRabbitMqClient *queue.RabbitMqClient
}

func New(config *conf.Config) *Boss {
	boss := &Boss{
		Config:      config,
		VirtualMtas: make([]*ZMSmtp.VirtualMta, 0),
		InboundMtas: make([]*ZMSmtp.SmtpServer, 0),
	}
	client := queue.New()
	client.Connect(config, true)
	client.ExchangeDeclare(queue.InboundExchange, true, false, false, false, nil)
	que, err := client.QueueDeclare(queue.InboundStagingQueueName, true, false, false, false, nil)
	client.QueueBind(que.Name, queue.InboundExchange, "", false, nil)
	boss.InboundRabbitMqClient = client
	return
}

func (boss *Boss) Run() {
	for _, vmta := range boss.Config.IPAddresses {
		vm := ZMSmtp.CreateNewVirtualMta(vmta.IP, "vmta1.localhost", 25, vmta.Inbound, vmta.Outbound)
		boss.VirtualMtas = append(boss.VirtualMtas, vm)
		inboundServer := ZMSmtp.CreateNewSmtpServer(vm)
		boss.InboundMtas = append(boss.InboundMtas, inboundServer)
		//TODO: pass InboundRabbitMqClient to inboundServer and inboundHandler
		go inboundServer.Run()
	}
}
