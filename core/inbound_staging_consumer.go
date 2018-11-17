package core

import (
	"encoding/json"
	"fmt"

	OS "../cross"
	"../queue"
	dkim "github.com/emersion/go-dkim"
)

var options := &dkim.SignOptions{
	Domain: "example.org",
	Selector: "brisbane",
	Signer: privateKey,
}

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
		case messageDelivery, ok := <-messageChannel:
			if ok {
				msg := &entity.Message{}				
				r := strings.NewReader(messageDelivery.Body)
				json.Unmarshal(messageDelivery.Body,msg)
				var b bytes.Buffer
				r := strings.NewReader(msg.Data)
				if err := dkim.Sign(&b, r, options); err != nil {
					//TODO: fix or report dkim error
				}
				data := json.Marshal(msg)
				err := consumer.RabbitMqClient.Publish(
					queue.OutboundExchange,
					queue.RoutingKeyOutbound,
					false,
					false,
					data,
				)

				if err == nil{
					messageDelivery.Ack(true)
				}else {
					messageChannel.Reject(true)
				}
			}
		}
	}
}
