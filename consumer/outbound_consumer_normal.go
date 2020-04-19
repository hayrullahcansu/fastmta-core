package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/hayrullahcansu/fastmta-core/relay"
)

type OutboundConsumerNormalSender struct {
	RabbitMqClient *rabbit.RabbitMqClient
}

func NewOutboundConsumerNormalSender() *OutboundConsumerNormalSender {
	return &OutboundConsumerNormalSender{
		RabbitMqClient: rabbit.New(),
	}
}

func (consumer *OutboundConsumerNormalSender) Run() {
	consumer.RabbitMqClient.Connect(true)
	ch, err := consumer.RabbitMqClient.Consume(constant.OutboundNormalQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", constant.OutboundNormalQueueName, err, cross.NewLine))
	}
	for {
		select {
		case outboundMessage, ok := <-ch:
			if ok {
				json.Unmarshal(outboundMessage.Body, pureMessage)
				logger.Infof("Recieved message From %s", constant.OutboundNormalQueueName)
				// if _, ok := caching.InstanceDomain().C.Get(pureMessage.Host); !ok {
				//exchange.InstanceRouter().
				// }
				go relay.InstanceManager().SendMessage(outboundMessage)
				logger.Infof("Sended message to %s", pureMessage.RcptTo)

				outboundMessage.Ack(false)
			}
		}
	}
}
func (consumer *OutboundConsumerNormalSender) Stop() {
	consumer.RabbitMqClient.Close()
}
