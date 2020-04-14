package outbound

import "time"

const (
	Timeout   = time.Second * time.Duration(30)
	KeepAlive = time.Second * time.Duration(30)
	Port      = 25
	//MtaName       = "ZetaMail"
	//MaxErrorLimit = 10
)
