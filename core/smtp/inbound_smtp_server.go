package smtp

import (
	"fmt"
	"net"

	"../../queue"
)

type InboundSmtpServer struct {
	ID             string
	VMta           *VirtualMta
	VmtaHostName   string
	VmtaIPAddr     *net.IPAddr
	Port           int
	RabbitMqClient *queue.RabbitMqClient
}

func CreateNewInboundSmtpServer(vmta *VirtualMta) *InboundSmtpServer {
	client := queue.New()
	client.Connect(true)
	client.ExchangeDeclare(queue.InboundExchange, true, false, false, false, nil)
	que, _ := client.QueueDeclare(queue.InboundQueueName, true, false, false, false, nil)
	client.QueueBind(que.Name, queue.InboundExchange, queue.RoutingKeyInbound, false, nil)
	return &InboundSmtpServer{
		ID:             "",
		VMta:           vmta,
		VmtaHostName:   vmta.VmtaHostName,
		VmtaIPAddr:     vmta.VmtaIPAddr,
		Port:           vmta.Port,
		RabbitMqClient: client,
	}
}

func (smtpServer *InboundSmtpServer) Run() {
	mergedAddress := fmt.Sprintf("%s:%d", smtpServer.VMta.IPAddressString, smtpServer.Port)
	listener, err := net.Listen("tcp", mergedAddress)

	if err != nil {
		panic(fmt.Sprintf("%s Can't listen inbound %s", mergedAddress, err))
		//LOG
	}
	fmt.Printf("%s Listening", mergedAddress)
	defer listener.Close()

	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			//LOG
		}
		// Handle inbound connections in a new goroutine.
		go InboundHandler(smtpServer, conn)
	}
}
