package core

import (
	"net"

	"github.com/hayrullahcansu/zetamail/entity"
)

type Domain struct {
	Name      string
	MXRecords []*net.MX
	Dkimmer   entity.Dkimmer
}

func NewDomain(name string) (*Domain, error) {
	domain := &Domain{
		Name: name,
	}
	mx, err := net.LookupMX(name)
	if err == nil {
		domain.MXRecords = mx
		return domain, nil
	}
	return domain, err
}

func InitDomain(name string, domain *entity.Domain) *Domain {
	domainInstance := &Domain{
		Name: name,
	}
	if domain.MXRecords == nil || len(domain.MXRecords) < 1 {
		mx, err := net.LookupMX(name)
		if err == nil {
			domainInstance.MXRecords = mx
			//TODO: save db
		}
	} else {
		domainInstance.MXRecords = make([]*net.MX, len(domain.MXRecords))
		for i, mx := range domain.MXRecords {
			domainInstance.MXRecords[i] = &net.MX{
				Host: mx.Host,
				Pref: mx.Pref,
			}
		}
	}
	return domainInstance
}
