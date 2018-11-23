package boss

import (
	"../core"
	"../core/exchange"
	ZMSmtp "../core/smtp"
	"../global"
)

type Boss struct {
	VirtualMtas            []*ZMSmtp.VirtualMta
	InboundMtas            []*ZMSmtp.InboundSmtpServer
	InboundConsumer        *core.InboundConsumer
	InboundStagingConsumer *core.InboundStagingConsumer
	Router                 *exchange.Router
}

func New() *Boss {
	boss := &Boss{
		VirtualMtas:            make([]*ZMSmtp.VirtualMta, 0),
		InboundMtas:            make([]*ZMSmtp.InboundSmtpServer, 0),
		InboundConsumer:        core.NewInboundConsumer(),
		InboundStagingConsumer: core.NewInboundStagingConsumer(),
		Router:                 exchange.NewRouter(),
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

	go boss.InboundConsumer.Run()
	go boss.InboundStagingConsumer.Run()

}
