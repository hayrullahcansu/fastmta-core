package brokerhub

import (
	"time"

	"github.com/hayrullahcansu/fastmta-core/dto"
	"github.com/hayrullahcansu/fastmta-core/netw"
)

type BrokerClient struct {
	*netw.BaseClient
	Id          string
	Name        string
	IsEnabled   bool
	StartedDate time.Time
}

func NewClient() *BrokerClient {
	client := &BrokerClient{}
	base := netw.NewBaseClient(client)
	client.BaseClient = base
	client.StartedDate = time.Now()
	return client
}

func (b *BrokerClient) ToDto() dto.Broker {
	status := "active"
	if !b.IsEnabled {
		status = "passive"
	}
	return dto.Broker{
		ID:          b.Id,
		Name:        b.Name,
		Status:      status,
		StartedDate: b.StartedDate,
	}
}
