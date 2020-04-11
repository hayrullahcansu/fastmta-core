package boss

import (
	"github.com/hayrullahcansu/fastmta-core/core"
	"github.com/hayrullahcansu/fastmta-core/global"
)

type Boss struct {
	VirtualMtas              []*core.VirtualMta
	InboundMtas              []*core.InboundSmtpServer
	InboundConsumer          *core.InboundConsumer
	InboundStagingConsumer   *core.InboundStagingConsumer
	OutboundConsumerMultiple *core.OutboundConsumerMultipleSender
	OutboundConsumerNormal   *core.OutboundConsumerNormalSender
	Router                   *core.Router
}

func New() *Boss {
	boss := &Boss{
		VirtualMtas:              make([]*core.VirtualMta, 0),
		InboundMtas:              make([]*core.InboundSmtpServer, 0),
		InboundConsumer:          core.NewInboundConsumer(),
		InboundStagingConsumer:   core.NewInboundStagingConsumer(),
		OutboundConsumerNormal:   core.NewOutboundConsumerNormalSender(),
		OutboundConsumerMultiple: core.NewOutboundConsumerMultipleSender(),
		Router:                   core.InstanceRouter(),
	}
	return boss
}

func (boss *Boss) Run() {
	for _, vmta := range global.StaticConfig.IPAddresses {
		vm := core.CreateNewVirtualMta(vmta.IP, vmta.HostName, 25, vmta.Inbound, vmta.Outbound, false)
		boss.VirtualMtas = append(boss.VirtualMtas, vm)
		inboundServer := core.CreateNewInboundSmtpServer(vm)
		boss.InboundMtas = append(boss.InboundMtas, inboundServer)
		go inboundServer.Run()
	}
	go core.InstanceBulkSender().Run()
	go boss.InboundConsumer.Run()
	go boss.InboundStagingConsumer.Run()
	go boss.OutboundConsumerNormal.Run()
	go boss.OutboundConsumerMultiple.Run()

}
