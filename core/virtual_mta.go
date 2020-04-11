package core

import (
	"fmt"
	"net"
	"sync"

	"github.com/hayrullahcansu/fastmta-core/global"
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
	Port                       int
	IsSmtpInbound              bool
	IsSmtpOutbound             bool
	LocalPort                  int
	TLS                        bool
	ConcurrentConnectionNumber int
}

func InstancePool() *VMtaPool {
	poolOnce.Do(func() {
		Pool = newVMtaPool()
	})
	return Pool
}

func (v *VMtaPool) GetVMtA() *VirtualMta {
	v.m.Lock()
	defer func() {
		counter++
		v.m.Unlock()
	}()
	return v.virtualMtas[counter%len(v.virtualMtas)]
}

func newVMtaPool() *VMtaPool {
	_pool := &VMtaPool{
		m:           &sync.Mutex{},
		virtualMtas: make([]*VirtualMta, 0),
		rules:       make([]*Rule, 0),
	}
	for _, vmta := range global.StaticConfig.IPAddresses {
		vm := CreateNewVirtualMta(vmta.IP, vmta.HostName, 25, vmta.Inbound, vmta.Outbound, false)
		_pool.virtualMtas = append(_pool.virtualMtas, vm)
	}
	return _pool
}

//CreateNewVirtualMta creates new dto
func CreateNewVirtualMta(ip string, hostname string, port int, isInbound bool, isOutbound bool, tls bool) *VirtualMta {
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
