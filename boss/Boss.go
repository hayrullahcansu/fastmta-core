package boss

import (
	ZMSmtp "../core/smtp"
	"../global"
)

type Boss struct {
	VirtualMtas []*ZMSmtp.VirtualMta
	InboundMtas []*ZMSmtp.InboundSmtpServer
}

func New() *Boss {
	boss := &Boss{
		VirtualMtas: make([]*ZMSmtp.VirtualMta, 0),
		InboundMtas: make([]*ZMSmtp.InboundSmtpServer, 0),
	}
	return boss
}

func (boss *Boss) Run() {
	for _, vmta := range global.StaticConfig.IPAddresses {
		vm := ZMSmtp.CreateNewVirtualMta(vmta.IP, vmta.HostName, 25, vmta.Inbound, vmta.Outbound)
		boss.VirtualMtas = append(boss.VirtualMtas, vm)
		inboundServer := ZMSmtp.CreateNewInboundSmtpServer(vm)
		boss.InboundMtas = append(boss.InboundMtas, inboundServer)
		go inboundServer.Run()
	}
}
