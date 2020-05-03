package netw

import (
	"fmt"
	"sync"
)

type BaseRoomManager struct {
	EnvelopeListener
	Register         chan interface{}
	Unregister       chan interface{}
	Update           chan *Update
	Notify           chan *Notify
	Broadcast        chan *Envelope
	BroadcastStop    chan bool
	ListenEventsStop chan bool
	L                *sync.Mutex
}

type IBaseRoomManager interface {
	OnConnect(baseClient *BaseClient)
	OnDisconnect(baseClient *BaseClient)
	PurgeRoom()
}

func NewBaseRoomManager() *BaseRoomManager {
	return &BaseRoomManager{
		Register:         make(chan interface{}, 1),
		Unregister:       make(chan interface{}, 1),
		Update:           make(chan *Update, 10),
		Notify:           make(chan *Notify, 1),
		Broadcast:        make(chan *Envelope, 10),
		BroadcastStop:    make(chan bool),
		ListenEventsStop: make(chan bool),
		L:                &sync.Mutex{},
	}
}

func (s *BaseRoomManager) ListenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case b := <-s.Broadcast:
			println(b)
			break
		case b := <-s.Update:
			println(b.Type)
			break
		case notify := <-s.Notify:
			s.OnNotify(notify)
			// default:
			// break
		}
	}
}

func (m *BaseRoomManager) OnConnect(c interface{}) {

}

func (m *BaseRoomManager) OnDisconnect(c interface{}) {
	// _, ok := c.(*AmericanSPClient)
	// if ok {
	fmt.Println("OnDisconnectBase")
	// }
}

func (m *BaseRoomManager) PurgeRoom() {

}

func (s *BaseRoomManager) OnNotify(notify *Notify) {
	d := notify.Message.Message
	switch v := notify.Message.Message.(type) {
	case Message:
		t := d.(Message)
		s.OnMessage(notify.SentBy, &t)
	case Event:
		t := d.(Event)
		s.OnEvent(notify.SentBy, &t)
	case Register:
		t := d.(Register)
		s.OnRegister(notify.SentBy, &t)
	default:
		fmt.Printf("unexpected type %T", v)
	}
}

func (s *BaseRoomManager) OnEvent(c interface{}, event *Event) {}

func (s *BaseRoomManager) OnRegister(c interface{}, register *Register) {}

func (s *BaseRoomManager) OnMessage(c interface{}, message *Message) {}
