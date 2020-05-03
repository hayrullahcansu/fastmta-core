package broker

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/hayrullahcansu/fastmta-core/util"

	"github.com/hayrullahcansu/fastmta-core/netw"

	"github.com/gorilla/websocket"
	"github.com/hayrullahcansu/fastmta-core/consumer"
	"github.com/hayrullahcansu/fastmta-core/core"
	"github.com/hayrullahcansu/fastmta-core/global"
	"github.com/hayrullahcansu/fastmta-core/in"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/mta"
)

//Broker keeps configrations, virtual MTAs, instance of all consumers in memory.
type Broker struct {
	c                        *websocket.Conn
	VirtualMtas              []*mta.VirtualMta
	InboundMtas              []*in.SmtpServer
	InboundConsumer          *consumer.InboundConsumer
	InboundStagingConsumer   *consumer.InboundStagingConsumer
	OutboundConsumerMultiple *consumer.OutboundConsumerMultipleSender
	OutboundConsumerNormal   *consumer.OutboundConsumerNormalSender
	Router                   *core.Router
	ID                       string
	Name                     string
	IsEnabled                bool
}

// New creates new instance of Broker
func New(id, name string, isEnabled bool) *Broker {
	Broker := &Broker{
		VirtualMtas:              make([]*mta.VirtualMta, 0),
		InboundMtas:              make([]*in.SmtpServer, 0),
		InboundConsumer:          consumer.NewInboundConsumer(),
		InboundStagingConsumer:   consumer.NewInboundStagingConsumer(),
		OutboundConsumerNormal:   consumer.NewOutboundConsumerNormalSender(),
		OutboundConsumerMultiple: consumer.NewOutboundConsumerMultipleSender(),
		Router:                   core.InstanceRouter(),
		ID:                       id,
		Name:                     name,
		IsEnabled:                isEnabled,
	}
	return Broker
}

//Run initializes configrations and virtual MTAs , starts consumers.
func (Broker *Broker) Run() {
	var addr = flag.String("addr", fmt.Sprintf("%s:%d", global.StaticConfig.Boss.Host, global.StaticConfig.Boss.Port), "http service address")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	logger.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Fatalf("dial:", err)
	}
	defer c.Close()
	Broker.c = c
	Broker.register()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logger.Infof("read:", err)
			return
		}
		fmt.Println(string(message))
		var msg json.RawMessage
		env := netw.Envelope{
			Message: &msg,
		}
		if err := json.Unmarshal(message, &env); err != nil {
			log.Fatal(err)
		}
		switch env.MessageCode {
		case netw.EEvent:
			var event netw.Event
			if err := json.Unmarshal(msg, &event); err != nil {
				log.Fatal(err)
			}
			env.Message = event

		case netw.ERegister:
			var register netw.Register
			if err := json.Unmarshal(msg, &register); err != nil {
				log.Fatal(err)
			}
			if register.Result == "ok" {
				Broker.init()
			} else {
				logger.Fatalf(register.Result)
			}
		default:
			continue
		}
	}

}

func (Broker *Broker) register() {
	dd := &netw.Envelope{
		Client:      "broker",
		MessageCode: netw.ERegister,
		Message: netw.Register{
			Id:        Broker.ID,
			Name:      Broker.Name,
			IsEnabled: Broker.IsEnabled,
		},
	}
	data := util.ToJson(dd)
	err := Broker.c.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		logger.Infof("write:", err)
		return
	}
}

func (Broker *Broker) init() {
	for _, vmta := range global.StaticConfig.IPAddresses {
		for _, port := range global.StaticConfig.Ports {
			vm := mta.CreateNewVirtualMta(vmta.IP, vmta.HostName, port, vmta.GroupID, vmta.Inbound, vmta.Outbound, false)
			Broker.VirtualMtas = append(Broker.VirtualMtas, vm)
			inboundServer := in.CreateNewSmtpServer(vm)
			Broker.InboundMtas = append(Broker.InboundMtas, inboundServer)
			go inboundServer.Run()
		}
	}
	go core.InstanceBulkSender().Run()
	go Broker.InboundConsumer.Run()
	go Broker.InboundStagingConsumer.Run()
	go Broker.OutboundConsumerNormal.Run()
	go Broker.OutboundConsumerMultiple.Run()

}
