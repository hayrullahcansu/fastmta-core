package queue

import (
	"fmt"

	"../conf"
	OS "../cross"
	"github.com/streadway/amqp"
)

type RabbitMqClient struct {
	MakeSureConnection bool
	IsConnected        bool
	Conf               *conf.RabbitMqConfig
	Conn               *amqp.Connection
	Channel            *amqp.Channel
}

func New() *RabbitMqClient {
	client := &RabbitMqClient{}
	return client
}

func (client *RabbitMqClient) Connect(conf *conf.RabbitMqConfig, makeSure bool) (*amqp.Connection, *amqp.Channel) {
	if !client.IsConnected {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.UserName, conf.Password, conf.Host, conf.Port))
		if err != nil {
			panic(fmt.Sprintf("Can't connect rabbitClient error:%s", err))
		}
		channel, err := conn.Channel()
		if err != nil {
			panic(fmt.Sprintf("Can't connect rabbitClient error:%s", err))
		}
		client.IsConnected = true
		if makeSure {
			client.MakeSureConnectionEndless()
		}
		return conn, channel
	}
	return client.Conn, client.Channel
}

func (client *RabbitMqClient) MakeSureConnectionEndless() {
	//TODO: makesure connection is endless
	if !client.MakeSureConnection {
		client.MakeSureConnection = true
		go func() {
			notify := client.Conn.NotifyClose(make(chan *amqp.Error))           //error channel
			notifyReturn := client.Channel.NotifyReturn(make(chan amqp.Return)) //error channel
			for {
				select {
				case err := <-notify:
					client.IsConnected = false
					client.MakeSureConnection = false
					fmt.Printf("RabbitMqClient error handled: %s%s", err, OS.NewLine)
					defer client.Connect(client.Conf, true)
					break
				case message := <-notifyReturn:
					err := client.Publish(message.Exchange, message.RoutingKey, false, false, message.Body)
					if err != nil {
						fmt.Printf("RabbitMqClient returned Message cant publish: %s%s", err, OS.NewLine)
					}
				}
			}
		}()
	}
}

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
func (client *RabbitMqClient) ExchangeDeclare(name string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error {
	return client.Channel.ExchangeDeclare(
		"logs",     // name
		"fanout",   // type
		durable,    // durable
		autoDelete, // auto-deleted
		internal,   // internal
		noWait,     // no-wait
		args,       // arguments
	)
}

func (client *RabbitMqClient) Publish(exchange string, routingKey string, mandatory bool, immediate bool, data []byte) error {
	return client.Channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		mandatory,  // mandatory
		immediate,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
}
