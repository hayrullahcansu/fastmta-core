package mta

import (
	"sync"

	"github.com/hayrullahcansu/fastmta-core/global"
)

var instanceManager *VirtualMtaManager
var once sync.Once

// VirtualMtaManager keeps all virtual MTA group in mapped table.
type VirtualMtaManager struct {
	virtualMtaGroups map[int]*VirtualMtaGroup
}

// InstanceManager returns new or existing instance of VirtualMtaManager
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
			vm := CreateNewVirtualMta(vmta.IP, vmta.HostName, port, vmta.GroupID, vmta.Inbound, vmta.Outbound, false)
			group, ok := instance.virtualMtaGroups[vm.GroupID]
			if !ok {
				group = NewVirtualMtaGroup()
				instance.virtualMtaGroups[vm.GroupID] = group
			}
			group.AddVirtualMta(vm)
		}
	}
	return instance
}

// GetVirtualMtaGroup returns VirtualMtaGroup by GroupId
func (m *VirtualMtaManager) GetVirtualMtaGroup(GroupID int) *VirtualMtaGroup {
	return m.virtualMtaGroups[GroupID]
}
