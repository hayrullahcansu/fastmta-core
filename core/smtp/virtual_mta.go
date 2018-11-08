package smtp

import (
	"fmt"
	"net"
	"sync"
)

type VirtualMta struct {
	lock            *sync.Mutex
	isInUsage       bool
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
		lock:            &sync.Mutex{},
		isInUsage:       false,
		IPAddressString: ip,
		VmtaHostName:    hostname,
		VmtaIPAddr:      ipAddr,
		Port:            port,
		IsSmtpInbound:   isInbound,
		IsSmtpOutbound:  isOutbound,
		LocalPort:       0,
	}
}

func (virtualMta *VirtualMta) HandleLock() bool {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	if !virtualMta.isInUsage {
		virtualMta.isInUsage = true
		return true
	}
	return false
}

func (virtualMta *VirtualMta) IsInUsage() bool {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	return virtualMta.isInUsage
}

func (virtualMta *VirtualMta) ReleaseLock() bool {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	if virtualMta.isInUsage {
		virtualMta.isInUsage = false
		return true
	}
	return false
}
