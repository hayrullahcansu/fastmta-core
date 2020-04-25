package consumer

import (
	"fmt"

	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/hayrullahcansu/fastmta-core/relay"
)

// OutboundConsumerNormalSender struct that includes RabbitMQClient
type OutboundConsumerNormalSender struct {
	RabbitMqClient *rabbit.RabbitMqClient
}

// NewOutboundConsumerNormalSender provides to get Outbound messages from the queue.
// It checks rules the messages and relays to target via agent
func NewOutboundConsumerNormalSender() *OutboundConsumerNormalSender {
	return &OutboundConsumerNormalSender{
		RabbitMqClient: rabbit.New(),
	}
}

// Run starts consuming from the queue
func (consumer *OutboundConsumerNormalSender) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(constant.OutboundNormalQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", constant.OutboundNormalQueueName, err, cross.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				go relay.InstanceManager().SendMessage(&outboundMessage)
			}
		}
	}
}

// Stop consuming from the queue
func (consumer *OutboundConsumerNormalSender) Stop() {
	consumer.RabbitMqClient.Close()
}
