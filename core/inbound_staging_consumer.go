package core

import (
	"fmt"

	OS "../cross"
	"../queue"
)

type InboundStagingConsumer struct {
	RabbitMqClient *queue.RabbitMqClient
}

func NewInboundStagingConsumer() *InboundStagingConsumer {
	return &InboundStagingConsumer{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *InboundStagingConsumer) Run() {
	consumer.RabbitMqClient.Connect(true)
	messageChannel, err := consumer.RabbitMqClient.Consume(queue.InboundStagingQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.InboundStagingQueueName, err, OS.NewLine))
	}

	for {
		select {
		case message, ok := <-messageChannel:
			if ok {
				//TODO: 1.) DKİM
			}
		}
	}
}
