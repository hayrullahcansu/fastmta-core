package mta

import (
	"sync/atomic"
)

type VirtualMtaGroup struct {
	virtualMtas []*VirtualMta
	next        uint32
}

func InitVirtualMtaGroup(virtualMtas []*VirtualMta) *VirtualMtaGroup {
	return &VirtualMtaGroup{
		virtualMtas: virtualMtas,
	}
}
func NewVirtualMtaGroup() *VirtualMtaGroup {
	return &VirtualMtaGroup{
		virtualMtas: make([]*VirtualMta, 0),
	}
}

func (g *VirtualMtaGroup) GetNextVirtualMta() *VirtualMta {
	n := atomic.AddUint32(&g.next, 1)
	return g.virtualMtas[(int(n)-1)%len(g.virtualMtas)]
}

func (g *VirtualMtaGroup) AddVirtualMta(mta *VirtualMta) {
	g.virtualMtas = append(g.virtualMtas, mta)
}
