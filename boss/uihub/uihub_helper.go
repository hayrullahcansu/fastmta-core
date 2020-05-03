package uihub

import (
	"github.com/hayrullahcansu/fastmta-core/boss/brokerhub"
	"github.com/hayrullahcansu/fastmta-core/boss/hub"
)

func (m *LobbyManager) processBroker(c *UIClient, message *hub.Envelope) {
	switch message.Command {
	case "List":
		getBrokerList(m, c, message)
	}
}

func getBrokerList(m *LobbyManager, c *UIClient, message *hub.Envelope) {
	list := brokerhub.Manager().GetBrokerList()
	c.sendMessage(message, list)
}
