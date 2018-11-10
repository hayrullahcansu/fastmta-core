package exchange

import (
	".."
)

type BulkSender struct {
	DomainMessageStacks map[string]*core.DomainMessageStack
}
