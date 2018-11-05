package boss

import (
	"../conf"
	ZMSmtp "../core/smtp"
)

type Boss struct {
	Config      *conf.Config
	VirtualMtas []*ZMSmtp.VirtualMta
	InboundMtas []*ZMSmtp.SmtpServer
}

func New(config *conf.Config) *Boss {
	boss := &Boss{
		Config:      config,
		VirtualMtas: make([]*ZMSmtp.VirtualMta, 0),
		InboundMtas: make([]*ZMSmtp.SmtpServer, 0),
	}
	return boss
}

func (boss *Boss) Run() {
	for _, vmta := range boss.Config.IPAddresses {
		vm := ZMSmtp.CreateNewVirtualMta(vmta.IP, "vmta1.localhost", 25, vmta.Inbound, vmta.Outbound)
		boss.VirtualMtas = append(boss.VirtualMtas, vm)
		inboundServer := ZMSmtp.CreateNewSmtpServer(vm, boss.Config)
		boss.InboundMtas = append(boss.InboundMtas, inboundServer)
		//TODO: pass InboundRabbitMqClient to inboundServer and inboundHandler
		go inboundServer.Run()
	}
}
