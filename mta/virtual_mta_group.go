package mta

import (
	"sync/atomic"
)

// VirtualMtaGroup includes virtual MTAs which has same groupID.
// Main idea is saparate IP addresses according to hot or cold, priority, relaying message rate.
type VirtualMtaGroup struct {
	virtualMtas []*VirtualMta
	next        uint32
}

// InitVirtualMtaGroup initialize bulk virtual MTAs
func InitVirtualMtaGroup(virtualMtas []*VirtualMta) *VirtualMtaGroup {
	return &VirtualMtaGroup{
		virtualMtas: virtualMtas,
	}
}

// NewVirtualMtaGroup return new instance of VirtualMtaGroup
func NewVirtualMtaGroup() *VirtualMtaGroup {
	return &VirtualMtaGroup{
		virtualMtas: make([]*VirtualMta, 0),
	}
}

// GetNextVirtualMta return a virtual MTA from the pool. This method works as round robin.
func (g *VirtualMtaGroup) GetNextVirtualMta() *VirtualMta {
	n := atomic.AddUint32(&g.next, 1)
	return g.virtualMtas[(int(n)-1)%len(g.virtualMtas)]
}

// AddVirtualMta adding for a virtual MTA into existing pool.
func (g *VirtualMtaGroup) AddVirtualMta(mta *VirtualMta) {
	g.virtualMtas = append(g.virtualMtas, mta)
}
