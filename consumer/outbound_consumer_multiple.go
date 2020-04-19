package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/hayrullahcansu/fastmta-core/caching"
	"github.com/hayrullahcansu/fastmta-core/constant"
	OS "github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
)

type OutboundConsumerMultipleSender struct {
	RabbitMqClient *rabbit.RabbitMqClient
}

func NewOutboundConsumerMultipleSender() *OutboundConsumerMultipleSender {
	return &OutboundConsumerMultipleSender{
		RabbitMqClient: rabbit.New(),
	}
}

func (consumer *OutboundConsumerMultipleSender) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(constant.OutboundMultipleQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", constant.OutboundMultipleQueueName, err, OS.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Infof("Recieved message From %s", constant.OutboundMultipleQueueName)
				if _, ok := caching.InstanceDomain().C.Get(pureMessage.Host); !ok {
					//exchange.InstanceRouter().
				}
				// core.InstanceBulkSender().AppendMessage(pureMessage.Host, pureMessage)
				logger.Infof("queued message to send %s", pureMessage.RcptTo)

				outboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *OutboundConsumerMultipleSender) Stop() {
	consumer.RabbitMqClient.Close()
}
