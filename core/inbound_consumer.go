package core

import (
	"fmt"

	OS "../cross"
	"../queue"
)

type InboundConsumer struct {
	RabbitMqClient *queue.RabbitMqClient
}

func NewInboundConsumer() *InboundConsumer {
	return &InboundConsumer{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *InboundConsumer) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(queue.InboundStagingQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.InboundStagingQueueName, err, OS.NewLine))
	}
	for {
		select {
		case msg, ok := <-ch:
			if ok {
				//TODO: Process Message Here

				msg.Ack(false)
				_ = string(msg.Body)
			}
		}
	}
}
func (consumer *InboundConsumer) Stop() {
	consumer.RabbitMqClient.Close()
}
