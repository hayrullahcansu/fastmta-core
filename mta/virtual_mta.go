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

type VirtualMta struct {
	lock                       *sync.Mutex
	isInUsage                  bool
	IPAddressString            string
	VmtaHostName               string
	VmtaIPAddr                 *net.IPAddr
	GroupId                    int
	Port                       int
	IsSmtpInbound              bool
	IsSmtpOutbound             bool
	LocalPort                  int
	TLS                        bool
	ConcurrentConnectionNumber int
}

//CreateNewVirtualMta creates new dto
func CreateNewVirtualMta(ip string, hostname string, port int, groupId int, isInbound bool, isOutbound bool, tls bool) *VirtualMta {
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
		GroupId:         groupId,
	}
	if port == 25 {
		vm.TLS = false
	}
	return vm
}

func (virtualMta *VirtualMta) HandleLock() {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	if !virtualMta.isInUsage {
		virtualMta.isInUsage = true
	}
	virtualMta.ConcurrentConnectionNumber++
}

func (virtualMta *VirtualMta) IsInUsage() bool {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	return virtualMta.isInUsage
}

func (virtualMta *VirtualMta) ReleaseLock() {
	virtualMta.lock.Lock()
	defer virtualMta.lock.Unlock()
	if virtualMta.isInUsage && virtualMta.ConcurrentConnectionNumber < 2 {
		virtualMta.isInUsage = false
	}
	virtualMta.ConcurrentConnectionNumber--
}
