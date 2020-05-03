package uihub

import (
	"github.com/hayrullahcansu/fastmta-core/boss/hub"
	"github.com/hayrullahcansu/fastmta-core/util"
)

type UIClient struct {
	*hub.BaseClient
	Id   string
	Name string
}

func NewClient() *UIClient {
	client := &UIClient{}
	base := hub.NewBaseClient(client)
	client.BaseClient = base
	return client
}
func (c *UIClient) sendMessage(message *hub.Envelope, data interface{}) {
	c.sendMessageManuel(message.Target, message.Command, data)
}

func (c *UIClient) sendMessageManuel(target, command string, data interface{}) {
	c.Send <- &hub.Envelope{
		Target:  target,
		Command: command,
		Message: util.ToJson(data),
	}
}
