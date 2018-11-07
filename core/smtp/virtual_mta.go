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
