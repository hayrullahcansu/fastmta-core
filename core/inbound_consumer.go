package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/streadway/amqp"

	"github.com/google/uuid"
	OS "github.com/hayrullahcansu/zetamail/cross"
	"github.com/hayrullahcansu/zetamail/entity"
	"github.com/hayrullahcansu/zetamail/logger"
	"github.com/hayrullahcansu/zetamail/queue"
)

type InboundConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    chan bool
}

func NewInboundConsumer() *InboundConsumer {
	return &InboundConsumer{
		q: make(chan bool),
	}
}

func (consumer *InboundConsumer) Run() {
	conn, err := amqp.Dial(queue.NewRabbitMqDialString())
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.InboundQueueName, err, OS.NewLine))
	}
	ch, err := conn.Channel()
	deliveries, err := ch.Consume(queue.InboundQueueName, "", false, false, true, false, nil)
	if err != nil {

	}
	for {
		select {
		case inboundMessage, ok := <-deliveries:
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
						ch.Publish(
							queue.InboundStagingExchange,
							queue.RoutingKeyInboundStaging,
							false,
							false,
							amqp.Publishing{
								ContentType: "text/plain",
								Body:        data,
							},
						)
					}

				}
				inboundMessage.Ack(false)
			}
		case <-consumer.q:
			break
		}
	}
}

func (consumer *InboundConsumer) Stop() {
	consumer.q <- true
}
