package mta

import (
	"sync"

	"github.com/hayrullahcansu/fastmta-core/global"
)

var instanceManager *VirtualMtaManager
var once sync.Once

type VirtualMtaManager struct {
	virtualMtaGroups map[int]*VirtualMtaGroup
}

func InstanceManager() *VirtualMtaManager {
	once.Do(func() {
		instanceManager = newManager()
	})
	return instanceManager
}

func newManager() *VirtualMtaManager {
	instance := &VirtualMtaManager{
		virtualMtaGroups: make(map[int]*VirtualMtaGroup, 0),
	}
	for _, vmta := range global.StaticConfig.IPAddresses {
		for _, port := range global.StaticConfig.Ports {
			vm := CreateNewVirtualMta(vmta.IP, vmta.HostName, port, vmta.GroupId, vmta.Inbound, vmta.Outbound, false)
			group, ok := instance.virtualMtaGroups[vm.GroupId]
			if !ok {
				group = NewVirtualMtaGroup()
				instance.virtualMtaGroups[vm.GroupId] = group
			}
			group.AddVirtualMta(vm)
		}
	}
	return instance
}
func (m *VirtualMtaManager) GetVirtualMtaGroup(groupId int) *VirtualMtaGroup {
	return m.virtualMtaGroups[groupId]
}
