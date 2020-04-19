package mta

import (
	"fmt"
	"net"
	"sync"
)

var Pool *VMtaPool
var poolOnce sync.Once
var counter int

type VMtaPool struct {
	m           *sync.Mutex
	virtualMtas []*VirtualMta
	rules       []*Rule
}

type Rule struct {
}

// VirtualMta includes IPaddress, port, hostname, ability inbound processing or outbound processing.
type VirtualMta struct {
	lock                       *sync.Mutex
	isInUsage                  bool
	IPAddressString            string
	VmtaHostName               string
	VmtaIPAddr                 *net.IPAddr
	GroupID                    int
	Port                       int
	IsSmtpInbound              bool
	IsSmtpOutbound             bool
	LocalPort                  int
	TLS                        bool
	ConcurrentConnectionNumber int
}

//CreateNewVirtualMta creates new instance of VirtualMta
func CreateNewVirtualMta(ip string, hostname string, port int, GroupID int, isInbound bool, isOutbound bool, tls bool) *VirtualMta {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		panic(fmt.Sprintf("%s given ip address cant parsing", ip))
	}
	ipAddr := &net.IPAddr{IP: parsedIP}
	vm := &VirtualMta{
		lock:            &sync.Mutex{},
		isInUsage:       false,
		IPAddressString: ip,
		VmtaHostName:    hostname,
		VmtaIPAddr:      ipAddr,
		Port:            port,
		IsSmtpInbound:   isInbound,
		IsSmtpOutbound:  isOutbound,
		LocalPort:       0,
		TLS:             tls,
		GroupID:         GroupID,
	}
	if port == 25 {
		vm.TLS = false
	}
	return vm
}

// HandleLock locks the current virtual mta. that provides to block multiple usages.
func (virtualMta *VirtualMta) HandleLock() {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	if !virtualMta.isInUsage {
		virtualMta.isInUsage = true
	}
	virtualMta.ConcurrentConnectionNumber++
}

// IsInUsage returns it's in usage or not.
func (virtualMta *VirtualMta) IsInUsage() bool {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	return virtualMta.isInUsage
}

// ReleaseLock unlocks the current virtual mta.
func (virtualMta *VirtualMta) ReleaseLock() {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	if virtualMta.isInUsage && virtualMta.ConcurrentConnectionNumber < 2 {
		virtualMta.isInUsage = false
	}
	virtualMta.ConcurrentConnectionNumber--
}
