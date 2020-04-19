package boss

import (
	"github.com/hayrullahcansu/fastmta-core/consumer"
	"github.com/hayrullahcansu/fastmta-core/core"
	"github.com/hayrullahcansu/fastmta-core/global"
	"github.com/hayrullahcansu/fastmta-core/in"
	"github.com/hayrullahcansu/fastmta-core/mta"
)

type Boss struct {
	VirtualMtas              []*mta.VirtualMta
	InboundMtas              []*in.SmtpServer
	InboundConsumer          *consumer.InboundConsumer
	InboundStagingConsumer   *consumer.InboundStagingConsumer
	OutboundConsumerMultiple *consumer.OutboundConsumerMultipleSender
	OutboundConsumerNormal   *consumer.OutboundConsumerNormalSender
	Router                   *core.Router
}

func New() *Boss {
	boss := &Boss{
		VirtualMtas:              make([]*mta.VirtualMta, 0),
		InboundMtas:              make([]*in.SmtpServer, 0),
		InboundConsumer:          consumer.NewInboundConsumer(),
		InboundStagingConsumer:   consumer.NewInboundStagingConsumer(),
		OutboundConsumerNormal:   consumer.NewOutboundConsumerNormalSender(),
		OutboundConsumerMultiple: consumer.NewOutboundConsumerMultipleSender(),
		Router:                   core.InstanceRouter(),
	}
	return boss
}

func (boss *Boss) Run() {
	for _, vmta := range global.StaticConfig.IPAddresses {
		for _, port := range global.StaticConfig.Ports {
			vm := mta.CreateNewVirtualMta(vmta.IP, vmta.HostName, port, vmta.GroupId, vmta.Inbound, vmta.Outbound, false)
			boss.VirtualMtas = append(boss.VirtualMtas, vm)
			inboundServer := in.CreateNewSmtpServer(vm)
			boss.InboundMtas = append(boss.InboundMtas, inboundServer)
			go inboundServer.Run()
		}
	}
	go core.InstanceBulkSender().Run()
	go boss.InboundConsumer.Run()
	go boss.InboundStagingConsumer.Run()
	go boss.OutboundConsumerNormal.Run()
	go boss.OutboundConsumerMultiple.Run()

}
