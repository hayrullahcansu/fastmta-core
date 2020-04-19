package constant

import "time"

const (
	InboundExchange            string = "fastmta_ex_inbound"
	InboundStagingExchange     string = "fastmta_ex_inbound_staging"
	OutboundExchange           string = "fastmta_ex_outbound"
	WaitingExchange            string = "fastmta_ex_waiting"
	RoutingKeyInbound          string = "inbound"
	RoutingKeyInboundStaging   string = "inbound_staging"
	RoutingKeyOutboundMultiple string = "outbound_multiple"
	RoutingKeyOutboundNormal   string = "outbound_normal"
	RoutingKeyWaiting          string = "waiting_route"
	InboundQueueName           string = "fastmta_inbound"
	InboundStagingQueueName    string = "fastmta_inbound_staging"
	OutboundMultipleQueueName  string = "fastmta_outbound_multiple"
	OutboundNormalQueueName    string = "fastmta_outbound_normal"
	OutboundWaiting1           string = "fastmta_outbound_waiting_1"
	OutboundWaiting10          string = "fastmta_outbound_waiting_10"
	OutboundWaiting60          string = "fastmta_outbound_waiting_60"
	OutboundWaiting300         string = "fastmta_outbound_waiting_300"
)

const (
	ReadDeadLine  = time.Second * time.Duration(30)
	WriteDeadLine = time.Second * time.Duration(30)
	MtaName       = "ZetaMail"
	MaxErrorLimit = 10
)
