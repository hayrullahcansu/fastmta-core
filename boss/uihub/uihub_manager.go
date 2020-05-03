package uihub

import (
	"sync"

	"github.com/hayrullahcansu/fastmta-core/boss/hub"
)

type LobbyManager struct {
	*hub.BaseRoomManager
	m         sync.Mutex
	UIClients map[*UIClient]bool
}

var _instance *LobbyManager

var _once sync.Once

func Manager() *LobbyManager {
	_once.Do(initialManagerInstance)
	return _instance
}

func initialManagerInstance() {
	_instance = &LobbyManager{
		BaseRoomManager: hub.NewBaseRoomManager(),
		UIClients:       make(map[*UIClient]bool),
	}
	go _instance.ListenEvents()
}

func (s *LobbyManager) ListenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case _ = <-s.Broadcast:
		case notify := <-s.Notify:
			s.OnNotify(notify)
			// default:
		}
	}
}

func (m *LobbyManager) ConnectLobby(c *UIClient) {
	m.UIClients[c] = true
	c.Notify = m.Notify
	m.Register <- c
}

func (m *LobbyManager) OnConnect(c interface{}) {
	client, ok := c.(*UIClient)
	if ok {
		client.Unregister = m.Unregister
	}
}

func (s *LobbyManager) OnNotify(notify *hub.Notify) {
	d := notify.Message
	if client, ok := notify.SentBy.(*UIClient); ok {
		s.onMessage(client, d)
	}
}

func (m *LobbyManager) onMessage(c *UIClient, message *hub.Envelope) {
	switch message.Target {
	case "Broker":
		m.processBroker(c, message)
	}
}
