package consumer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/streadway/amqp"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/queue"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/hayrullahcansu/fastmta-core/smtp/mime"
)

// InboundConsumer struct that includes RabbitMQClient
type InboundConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    chan bool
}

// NewInboundConsumer provides to get Inbound messages from the queue.
// It checks recipients the messages and creates entity.Message for each recipient
func NewInboundConsumer() *InboundConsumer {
	return &InboundConsumer{
		q: make(chan bool),
	}
}

// Run starts consuming from the queue
func (consumer *InboundConsumer) Run() {
	conn, err := amqp.Dial(rabbit.NewRabbitMqDialString())
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", constant.InboundQueueName, err, cross.NewLine))
	}
	ch, err := conn.Channel()
	deliveries, err := ch.Consume(constant.InboundQueueName, "", false, false, true, false, nil)
	if err != nil {

	}
	for {
		select {
		case inboundMessage, ok := <-deliveries:
			if ok {
				pureMessage := &entity.InboundMessage{}
				json.Unmarshal(inboundMessage.Body, pureMessage)
				logger.Infof("Received message From %s", constant.InboundQueueName)
				headers := mime.ParseHeaders(pureMessage.Data)

				groupId := constant.DefaultGroupID
				if dd, err := strconv.Atoi(headers.Get(constant.HeaderGroupKey)); err == nil {
					groupId = dd
				}
				db, err := entity.GetDbContext()
				entity.PanicOnError(err)
				for i := 0; i < len(pureMessage.RcptTo); i++ {
					msg := &entity.Message{
						MessageID: uuid.New().String(),
						MailFrom:  string(pureMessage.MailFrom),
						Data:      string(pureMessage.Data),
						Status:    "w",
						RcptTo:    string(pureMessage.RcptTo[i]),
						Host:      string(pureMessage.RcptTo[i][strings.LastIndex(pureMessage.RcptTo[i], "@")+1:]),
						GroupID:   groupId,
					}
					db.Create(&msg)
					queue.Instance().EnqueueInboundStaging(msg)
				}
				db.Close()
				inboundMessage.Ack(false)
			}
		case <-consumer.q:
			break
		}
	}
}

// Stop consuming from the queue
func (consumer *InboundConsumer) Stop() {
	consumer.q <- true
}
