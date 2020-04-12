package tcptimer

import "time"

const (
	ReadDeadLine  = time.Second * time.Duration(30)
	WriteDeadLine = time.Second * time.Duration(30)
	MtaName       = "FastMTA"
	MaxErrorLimit = 10
)
