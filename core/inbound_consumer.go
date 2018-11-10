package core

import (
	"encoding/json"
	"fmt"

	OS "../cross"
	"../entity"
	"../queue"
)

type InboundConsumer struct {
	RabbitMqClient        *queue.RabbitMqClient
	MessageProcessChannel chan *entity.Message
}

func NewInboundConsumer() *InboundConsumer {
	return &InboundConsumer{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *InboundConsumer) SetChannel(channel chan *entity.Message) {
	if channel == nil {
		panic(fmt.Sprintf("MessageChannel can not be nil%s", OS.NewLine))
	}
	consumer.MessageProcessChannel = channel
}

func (consumer *InboundConsumer) Run() {
	if consumer.MessageProcessChannel == nil {
		panic(fmt.Sprintf("MessageChannel can not be nil%s", OS.NewLine))
	}
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(queue.InboundStagingQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.InboundStagingQueueName, err, OS.NewLine))
	}
	for {
		select {
		case msg, ok := <-ch:
			if ok {
				pureMessage := &entity.Message{}
				json.Unmarshal(msg.Body, pureMessage)
				consumer.MessageProcessChannel <- pureMessage
				msg.Ack(false)
				_ = string(msg.Body)
			}
		}
	}
}
func (consumer *InboundConsumer) Stop() {
	consumer.RabbitMqClient.Close()
}
