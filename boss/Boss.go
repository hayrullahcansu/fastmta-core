package boss

import (
	"../caching"
	"../core"
	"../core/exchange"
	ZMSmtp "../core/smtp"
	"../global"
)

type Boss struct {
	VirtualMtas        []*ZMSmtp.VirtualMta
	InboundMtas        []*ZMSmtp.InboundSmtpServer
	InboundConsumer    *core.InboundConsumer
	Router             *exchange.Router
	DomainCacheManager *caching.CacheManager
}

func New() *Boss {
	boss := &Boss{
		VirtualMtas:        make([]*ZMSmtp.VirtualMta, 0),
		InboundMtas:        make([]*ZMSmtp.InboundSmtpServer, 0),
		InboundConsumer:    core.NewInboundConsumer(),
		DomainCacheManager: caching.NewCacheManager(),
		Router:             exchange.NewRouter(),
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
	boss.DomainCacheManager.Init()
	boss.Router.SetDomainCacheManager(boss.DomainCacheManager)
	go boss.InboundConsumer.Run()

}
