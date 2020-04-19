package transaction

type TransactionResult int

const (
	Success TransactionResult = iota + 1
	FailedToConnect
	ServiceNotAvailable
	RejectedByRemoteServer
	RetryRequired
	HostNotFound
	MaxConnections
	MaxMessages
	ClientAlreadyInUse
	Timeout
)

type TransactionGroupResult struct {
	TransactionResult TransactionResult
	ResultMessage     string
}
