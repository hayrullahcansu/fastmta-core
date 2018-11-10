package transaction

type TransactionResult int

const (
	Success TransactionResult = iota + 1
	FailedToConnect
	ServiceNotAvalible
	RejectedByRemoteServer
	RetryRequired
	HostNotFound
	MaxConnections
	MaxMessages
	ClientAlreadyInUse
)
