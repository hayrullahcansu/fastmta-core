package smtp

import (
	".."
	"../../entity"
	"../transaction"
)

func SendMessage(message *entity.Message, virtualMta *VirtualMta, domain *core.Domain) transaction.TransactionResult {
	return transaction.RetryRequired
}
