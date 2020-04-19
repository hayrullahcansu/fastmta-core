package transaction

// TransactionResult extended from int
type TransactionResult int

const (
	// Success result
	Success TransactionResult = iota + 1
	// FailedToConnect that can helps to send message in next try.
	FailedToConnect
	// ServiceNotAvailable that can helps to send message in next try.
	ServiceNotAvailable
	// RejectedByRemoteServer the recipient host didn't accept to send messageç
	RejectedByRemoteServer
	// RetryRequired that can helps to send message in next try.
	RetryRequired
	// HostNotFound means it's hard bounce. You cannot send a messge to kind of host.
	HostNotFound
	// MaxConnections that reaches over Maximum connection limit for specific virtualMta
	MaxConnections
	// MaxMessages that virtualMta reached to maximum send rate in a periodik time.
	MaxMessages
	// ClientAlreadyInUse that virtualMta is in a usage. You should try to send via another virtualMta/
	ClientAlreadyInUse
	// Timeout means you cannot send the message anymore.
	Timeout
)

// TransactionGroupResult includes TransactionResult and response message from the server.
type TransactionGroupResult struct {
	TransactionResult TransactionResult
	ResultMessage     string
}
