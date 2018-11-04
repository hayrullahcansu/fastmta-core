package smtp

import (
	"fmt"
	"net"
)

type VirtualMta struct {
	IPAddressString string
	VmtaHostName    string
	VmtaIPAddr      *net.IPAddr
	Port            int
	IsSmtpInbound   bool
	IsSmtpOutbound  bool
	LocalPort       int
}

//CreateNewVirtualMta creates new dto
func CreateNewVirtualMta(ip string, hostname string, port int, isInbound bool, isOutbound bool) *VirtualMta {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		panic(fmt.Sprintf("%s given ip address cant parsing", ip))
	}
	ipAddr := &net.IPAddr{IP: parsedIP}
	return &VirtualMta{
		IPAddressString: ip,
		VmtaHostName:    hostname,
		VmtaIPAddr:      ipAddr,
		Port:            port,
		IsSmtpInbound:   isInbound,
		IsSmtpOutbound:  isOutbound,
		LocalPort:       0,
	}
}

type SmtpServer struct {
	ID           string
	VMta         *VirtualMta
	VmtaHostName string
	VmtaIPAddr   *net.IPAddr
	Port         int
}

func CreateNewSmtpServer(vmta *VirtualMta) *SmtpServer {
	return &SmtpServer{
		ID:           "",
		VMta:         vmta,
		VmtaHostName: vmta.VmtaHostName,
		VmtaIPAddr:   vmta.VmtaIPAddr,
		Port:         vmta.Port,
	}
}

func (smtpServer *SmtpServer) Run() {
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
