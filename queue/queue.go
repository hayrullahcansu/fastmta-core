package queue

import (
	"sync"

	"github.com/hayrullahcansu/fastmta-core/rabbit"
)

var instanceClient *rabbit.RabbitMqClient
var once sync.Once

//Instance returns new or existing instance of RabbitMQClient
func Instance() *rabbit.RabbitMqClient {
	once.Do(func() {
		instanceClient = newClient()
	})
	return instanceClient
}

func newClient() *rabbit.RabbitMqClient {
	client := rabbit.New()
	client.Connect(true)
	return client
}
