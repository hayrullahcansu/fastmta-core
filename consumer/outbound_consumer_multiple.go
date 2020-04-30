package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/hayrullahcansu/fastmta-core/relay"
)

// OutboundConsumerMultipleSender struct that includes RabbitMQClient
// It provides to send multiple message to the same domain.
type OutboundConsumerMultipleSender struct {
	RabbitMqClient *rabbit.RabbitMqClient
}

// NewOutboundConsumerMultipleSender creates new instance of OutboundConsumerMultipleSender
func NewOutboundConsumerMultipleSender() *OutboundConsumerMultipleSender {
	return &OutboundConsumerMultipleSender{
		RabbitMqClient: rabbit.New(),
	}
}

// Run starts consuming from the queue
func (consumer *OutboundConsumerMultipleSender) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(constant.OutboundMultipleQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", constant.OutboundMultipleQueueName, err, cross.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Infof("Recieved message From %s", constant.OutboundMultipleQueueName)
				// if _, ok := caching.InstanceDomain().Get(pureMessage.Host); !ok {
				//exchange.InstanceRouter().
				// }
				go relay.InstanceManager().SendMessage(&outboundMessage)
				logger.Infof("queued message to send %s", pureMessage.RcptTo)

				outboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *OutboundConsumerMultipleSender) Stop() {
	consumer.RabbitMqClient.Close()
}
