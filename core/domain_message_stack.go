package core

import (
	"../entity"
)

type DomainMessageStack struct {
	Domain       *Domain
	MessageStack chan *entity.Message
}

func NewDomainMessageStack(domain *Domain) *DomainMessageStack {
	return &DomainMessageStack{
		Domain:       domain,
		MessageStack: make(chan *entity.Message, 1000),
	}
}
