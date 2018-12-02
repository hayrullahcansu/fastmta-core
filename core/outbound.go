package core

import "../entity"

func SendMessages(message []*entity.Message, virtualMta *VirtualMta, domain *Domain) {
	if virtualMta.TLS {
		client := NewOutboundClientTLS()
		client.SendMessageTLS(message, virtualMta, domain)
	} else {
		client := NewOutboundClient()
		client.SendMessageNoTLS(message, virtualMta, domain)
	}
}
