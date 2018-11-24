package core

import (
	"encoding/json"
	"fmt"

	OS "../cross"
	"../entity"
	"../logger"
	"../queue"
)

type OutboundConsumer struct {
	RabbitMqClient *queue.RabbitMqClient
}

func NewOutboundConsumer() *OutboundConsumer {
	return &OutboundConsumer{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *OutboundConsumer) Run() {
	consumer.RabbitMqClient.Connect(true)
	consumer.RabbitMqClient.ExchangeDeclare(queue.OutboundExchange, true, false, false, false, nil)
	que, _ := consumer.RabbitMqClient.QueueDeclare(queue.OutboundQueueName, true, false, false, false, nil)
	consumer.RabbitMqClient.QueueBind(que.Name, queue.OutboundExchange, queue.RoutingKeyOutbound, false, nil)
	ch, err := consumer.RabbitMqClient.Consume(queue.OutboundQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.OutboundQueueName, err, OS.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Info.Printf("Recieved message From %s", queue.OutboundQueueName)
				logger.Info.Printf("Sended message to %s", pureMessage.RcptTo)

				outboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *OutboundConsumer) Stop() {
	consumer.RabbitMqClient.Close()
}
