package brokerhub

import (
	"fmt"
	"sync"

	"github.com/hayrullahcansu/fastmta-core/util"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/fastmta-core/netw"
)

type Manager struct {
	*netw.BaseRoomManager
	m                sync.Mutex
	BrokerClients    map[*BrokerClient]bool
	BrokerClientsMap map[string]*BrokerClient
}

var _instance *Manager

var _once sync.Once

func Instance() *Manager {
	_once.Do(initialManagerInstance)
	return _instance
}

func initialManagerInstance() {
	_instance = &Manager{
		BaseRoomManager:  netw.NewBaseRoomManager(),
		BrokerClients:    make(map[*BrokerClient]bool),
		BrokerClientsMap: make(map[string]*BrokerClient),
	}
	go _instance.ListenEvents()
}

func (s *Manager) ListenEvents() {
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

func (s *Manager) OnNotify(notify *netw.Notify) {
	d := notify.Message.Message
	switch v := notify.Message.Message.(type) {
	case netw.Event:
		t := d.(netw.Event)
		s.OnEvent(notify.SentBy, &t)
	case netw.Register:
		t := d.(netw.Register)
		s.OnRegister(notify.SentBy, &t)
	default:
		fmt.Printf("unexpected type %T", v)
	}
}

func (m *Manager) ConnectLobby(c *BrokerClient) {
	m.BrokerClients[c] = true
	c.Notify = m.Notify
	m.Register <- c
}

func (m *Manager) OnConnect(c interface{}) {
	client, ok := c.(*BrokerClient)
	if ok {
		client.Unregister = m.Unregister
	}
}
func (m *Manager) OnDisconnect(c interface{}) {
	client, ok := c.(*BrokerClient)
	if ok {
		m.m.Lock()
		defer m.m.Unlock()
		if _, ok := m.BrokerClientsMap[client.Id]; ok {
			delete(m.BrokerClientsMap, client.Id)
		}
	}
}
func (m *Manager) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
	client, ok := c.(*BrokerClient)
	if ok {
		//TODO: check player able to play?
		mode := playGame.Mode
		guid := uuid.New()
		playGame.Id = guid.String()
		playGame.Mode = mode
		client.Send <- &netw.Envelope{
			Client:      "client_id",
			Message:     playGame,
			MessageCode: netw.EPlayGame,
		}
	}
}
func (m *Manager) OnRegister(c interface{}, register *netw.Register) {
	m.m.Lock()
	defer m.m.Unlock()
	BrokerClient, ok := c.(*BrokerClient)
	if ok {
		if util.IsNullOrEmpty(register.Id) {
			register.Result = "Invalid Command"
			BrokerClient.Send <- &netw.Envelope{
				Client:      "boss",
				Message:     register,
				MessageCode: netw.ERegister,
			}
			return
		}
		if _, ok := m.BrokerClientsMap[register.Id]; ok {
			register.Result = "there is already BrokerClient with same id"
			BrokerClient.Send <- &netw.Envelope{
				Client:      "boss",
				Message:     register,
				MessageCode: netw.ERegister,
			}
			return
		}

		BrokerClient.Id = register.Id
		BrokerClient.Name = register.Name
		BrokerClient.IsEnabled = register.IsEnabled
		m.BrokerClientsMap[register.Id] = BrokerClient

		register.Result = "ok"
		BrokerClient.Send <- &netw.Envelope{
			Client:      "boss",
			Message:     register,
			MessageCode: netw.ERegister,
		}
	}
}
