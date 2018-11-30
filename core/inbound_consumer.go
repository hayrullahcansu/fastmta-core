package core

import (
	"encoding/json"
	"fmt"
	"strings"

	OS "../cross"
	"../entity"
	"../logger"
	"../queue"
	"github.com/google/uuid"
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
	ch, err := consumer.RabbitMqClient.Consume(queue.InboundQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.InboundQueueName, err, OS.NewLine))
	}
	for {
		select {
		case inboundMessage, ok := <-ch:
			if ok {
				pureMessage := &entity.InboundMessage{}
				json.Unmarshal(inboundMessage.Body, pureMessage)
				logger.Info.Printf("Recieved message From %s", queue.InboundQueueName)
				for i := 0; i < len(pureMessage.RcptTo); i++ {
					msg := &entity.Message{
						MessageID: uuid.New().String(),
						MailFrom:  string(pureMessage.MailFrom),
						Data:      string(pureMessage.Data),
						Status:    "w",
						RcptTo:    string(pureMessage.RcptTo[i]),
						Host:      string(pureMessage.RcptTo[i][strings.LastIndex(pureMessage.RcptTo[i], "@")+1:]),
					}
					data, err := json.Marshal(msg)
					if err == nil {
						consumer.RabbitMqClient.Publish(
							queue.InboundStagingExchange,
							queue.RoutingKeyInboundStaging,
							false,
							false,
							data,
						)
					}

				}
				inboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *InboundConsumer) Stop() {
	consumer.RabbitMqClient.Close()
}
