package exchange

import (
	"net"

	"../../entity"
)

type Domain struct {
	Name      string
	MXRecords []*net.MX
}

func NewDomain(name string) *Domain {
	domain := &Domain{
		Name: name,
	}
	mx, err := net.LookupMX(name)
	if err == nil {
		domain.MXRecords = mx
	}
	return domain
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
