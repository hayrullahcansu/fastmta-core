package brokerhub

import "github.com/hayrullahcansu/fastmta-core/dto"

func (l *Manager) GetBrokerList() []dto.Broker {
	l.L.Lock()
	defer l.L.Unlock()
	list := make([]dto.Broker, 0, len(l.BrokerClientsMap))
	for _, v := range l.BrokerClientsMap {
		list = append(list, v.ToDto())
	}
	return list
}
