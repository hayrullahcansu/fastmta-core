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

func CreateBoss(config *conf.Config) *Boss {
	return &Boss{
		Config:      config,
		VirtualMtas: make([]*ZMSmtp.VirtualMta, 0),
		InboundMtas: make([]*ZMSmtp.SmtpServer, 0),
	}
}

func (boss *Boss) Run() {
	for _, vmta := range boss.Config.IPAddresses {
		vm := ZMSmtp.CreateNewVirtualMta(vmta.IP, "vmta1.localhost", 25, vmta.Inbound, vmta.Outbound)
		boss.VirtualMtas = append(boss.VirtualMtas, vm)
		inboundServer := ZMSmtp.CreateNewSmtpServer(vm)
		boss.InboundMtas = append(boss.InboundMtas, inboundServer)
		go inboundServer.Run()
	}
}
