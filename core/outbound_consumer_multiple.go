package core

import (
	"encoding/json"
	"fmt"

	"github.com/hayrullahcansu/zetamail/caching"
	OS "github.com/hayrullahcansu/zetamail/cross"
	"github.com/hayrullahcansu/zetamail/entity"
	"github.com/hayrullahcansu/zetamail/logger"
	"github.com/hayrullahcansu/zetamail/queue"
)

type OutboundConsumerMultipleSender struct {
	RabbitMqClient *queue.RabbitMqClient
}

func NewOutboundConsumerMultipleSender() *OutboundConsumerMultipleSender {
	return &OutboundConsumerMultipleSender{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *OutboundConsumerMultipleSender) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(queue.OutboundMultipleQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.OutboundMultipleQueueName, err, OS.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Info.Printf("Recieved message From %s", queue.OutboundMultipleQueueName)
				if _, ok := caching.InstanceDomain().C.Get(pureMessage.Host); !ok {
					//exchange.InstanceRouter().
				}
				InstanceBulkSender().AppendMessage(pureMessage.Host, pureMessage)
				logger.Info.Printf("queued message to send %s", pureMessage.RcptTo)

				outboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *OutboundConsumerMultipleSender) Stop() {
	consumer.RabbitMqClient.Close()
}
