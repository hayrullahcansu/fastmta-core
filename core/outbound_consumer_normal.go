package core

import (
	"encoding/json"
	"fmt"

	OS "../cross"
	"../entity"
	"../logger"
	"../queue"
)

type OutboundConsumerNormalSender struct {
	RabbitMqClient *queue.RabbitMqClient
}

func NewOutboundConsumerNormalSender() *OutboundConsumerNormalSender {
	return &OutboundConsumerNormalSender{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *OutboundConsumerNormalSender) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(queue.OutboundNormalQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.OutboundNormalQueueName, err, OS.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Info.Printf("Recieved message From %s", queue.OutboundNormalQueueName)

				logger.Info.Printf("Sended message to %s", pureMessage.RcptTo)

				outboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *OutboundConsumerNormalSender) Stop() {
	consumer.RabbitMqClient.Close()
}
