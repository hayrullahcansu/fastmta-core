package rabbit

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/hayrullahcansu/fastmta-core/conf"
	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/global"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/queue/priority"
	"github.com/streadway/amqp"
)

// RabbitMqClient is a wrapper the official rabbitmq client.
// it provides infinitive connection. makes easier to usage the client.
type RabbitMqClient struct {
	MakeSureConnection bool
	IsConnected        bool
	Conf               *conf.RabbitMqConfig
	Conn               *amqp.Connection
	Channel            *amqp.Channel
}

// NewRabbitMqDialString returns connection string from the global config
func NewRabbitMqDialString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", global.StaticRabbitMqConfig.UserName, global.StaticRabbitMqConfig.Password, global.StaticRabbitMqConfig.Host, global.StaticRabbitMqConfig.Port)
}

// NewRabbitMq dials to rabbitmq server.
func NewRabbitMq() (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", global.StaticRabbitMqConfig.UserName, global.StaticRabbitMqConfig.Password, global.StaticRabbitMqConfig.Host, global.StaticRabbitMqConfig.Port))
}

// New returns new instance of RabbitMqClient
func New() *RabbitMqClient {
	client := &RabbitMqClient{
		MakeSureConnection: false,
		IsConnected:        false,
		Conf:               global.StaticRabbitMqConfig,
	}
	return client
}

// Connect to rabbitmq server. param1 makes sure infinity connection.
func (client *RabbitMqClient) Connect(makeSure bool) (*amqp.Connection, *amqp.Channel) {
	if !client.IsConnected {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", client.Conf.UserName, client.Conf.Password, client.Conf.Host, client.Conf.Port))
		if err != nil {
			panic(fmt.Sprintf("Can't connect rabbitClient error:%s", err))
		}
		channel, err := conn.Channel()
		if err != nil {
			panic(fmt.Sprintf("Can't connect rabbitClient error:%s", err))
		}
		client.IsConnected = true
		client.Conn = conn
		client.Channel = channel
		if makeSure {
			client.MakeSureConnectionEndless()
		}
		return conn, channel
	}
	return client.Conn, client.Channel
}

// ConnectForInit just using for ensure rabbitmq settings.
func (client *RabbitMqClient) ConnectForInit() (*amqp.Connection, *amqp.Channel, error) {
	if !client.IsConnected {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", client.Conf.UserName, client.Conf.Password, client.Conf.Host, client.Conf.Port))
		if err != nil {
			fmt.Printf("Can't connect rabbitClient error:%s%s", err, cross.NewLine)
		}
		channel, err := conn.Channel()
		if err != nil {
			fmt.Printf("Can't connect rabbitClient error:%s%s", err, cross.NewLine)
		}
		client.IsConnected = true
		client.Conn = conn
		client.Channel = channel
		return conn, channel, nil
	}
	return client.Conn, client.Channel, nil
}

// MakeSureConnectionEndless creates a chan and a goroutine that is for infinity connection.
// if the connection is broken. It will try to connect again.
func (client *RabbitMqClient) MakeSureConnectionEndless() {
	if !client.MakeSureConnection {
		client.MakeSureConnection = true
		go func() {
			notify := client.Conn.NotifyClose(make(chan *amqp.Error)) //error channel
			err, ok := <-notify
			if ok {
				client.IsConnected = false
				client.MakeSureConnection = false
				logger.Infof("RabbitMqClient error handled: %s%s", err, cross.NewLine)
				defer client.Connect(true)
			}
		}()
		go func() {
			notifyReturn := client.Channel.NotifyReturn(make(chan amqp.Return)) //error channel
			message, ok := <-notifyReturn
			if ok && client.IsConnected {
				err := client.Publish(message.Exchange, message.RoutingKey, false, false, message.Body, priority.LOW)
				if err != nil {
					logger.Infof("RabbitMqClient returned Message cant publish: %s%s", err, cross.NewLine)
				}
			}
		}()
	}
}

// QueueDeclare for declaring a queue
func (client *RabbitMqClient) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return client.Channel.QueueDeclare(
		name,
		durable,    // durable
		autoDelete, // delete when usused
		exclusive,  // exclusive
		noWait,     // no-wait
		args,       // arguments
	)
}

// QueueBind for binding key in an exchange
func (client *RabbitMqClient) QueueBind(name string, exchangeName string, routingKey string, noWait bool, args amqp.Table) error {
	return client.Channel.QueueBind(
		name,
		routingKey,   // durable
		exchangeName, // delete when usused
		noWait,       // exclusive
		args,         // arguments
	)
}

// ExchangeDeclare for declaring an exchange.
func (client *RabbitMqClient) ExchangeDeclare(name string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error {
	return client.Channel.ExchangeDeclare(
		name,       // name
		"direct",   // type
		durable,    // durable
		autoDelete, // auto-deleted
		internal,   // internal
		noWait,     // no-wait
		args,       // arguments
	)
}

// Publish for sent a message to custom exchange and queue
func (client *RabbitMqClient) Publish(exchange string, routingKey string, mandatory bool, immediate bool, data []byte, priority priority.Priority) error {
	return client.Channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		mandatory,  // mandatory
		immediate,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			Priority:    uint8(priority),
		})
}

// Consume returns a channel that receives message from the queue.
func (client *RabbitMqClient) Consume(queue string, consumerTag string, autoAck bool, exclusive bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return client.Channel.Consume(
		queue,       // exchange
		consumerTag, // routing key
		autoAck,     // mandatory
		exclusive,   // immediate
		false,
		noWait,
		args)
}

// Close the connection.
func (client *RabbitMqClient) Close() {
	client.Channel.Close()
	client.Conn.Close()
}

// EnqueueInbound add a message that's type InboundMessage
func (client *RabbitMqClient) EnqueueInbound(message *entity.InboundMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return client.Publish(
		"",
		constant.InboundQueueName,
		false,
		false,
		data,
		message.Priority,
	)
}

// EnqueueInboundStaging add a message that's type Message
func (client *RabbitMqClient) EnqueueInboundStaging(message *entity.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return client.Publish(
		"",
		constant.InboundStagingQueueName,
		false,
		false,
		data,
		message.Priority,
	)
}

// EnqueueOutboundMultiple add a message that's type Message
func (client *RabbitMqClient) EnqueueOutboundMultiple(message *entity.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return client.Publish(
		"",
		constant.OutboundMultipleQueueName,
		false,
		false,
		data,
		message.Priority,
	)
}

// EnqueueOutboundNormal add a message that's type Message
func (client *RabbitMqClient) EnqueueOutboundNormal(message *entity.Message) error {
	que := getOutboundQueue(message)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return client.Publish(
		"",
		que,
		false,
		false,
		data,
		message.Priority,
	)
}

func getOutboundQueue(message *entity.Message) string {
	que := constant.OutboundNormalQueueName
	dd := message.AttemptSendTime.Sub(time.Now()) / time.Second
	diff := math.Ceil(float64(dd))
	fmt.Println(diff)
	if diff > 0 {
		if diff < 10 {
			que = constant.OutboundWaiting1
		} else if diff < 60 {
			que = constant.OutboundWaiting10
		} else if diff < 300 {
			que = constant.OutboundWaiting60
		} else {
			que = constant.OutboundWaiting300
		}
	}
	return que
}
