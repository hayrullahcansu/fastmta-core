package core

import (
	"encoding/json"
	"fmt"

	"../caching"
	OS "../cross"
	"../entity"
	"../logger"
	"../queue"
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
	ch, err := consumer.RabbitMqClient.Consume(queue.OutboundMultipleSenderQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.OutboundMultipleSenderQueueName, err, OS.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Info.Printf("Recieved message From %s", queue.OutboundMultipleSenderQueueName)
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
