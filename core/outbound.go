package core

import (
	"github.com/hayrullahcansu/fastmta-core/core/transaction"
	OS "github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
)

func SendMessages(message []*entity.Message, host string) {
	var resultSet []*transaction.TransactionGroupResult
	var result bool = true
	domain, err := NewDomain(host)
	if err != nil {
		resultSet := make([]*transaction.TransactionGroupResult, 1)
		resultSet[0] = &transaction.TransactionGroupResult{}
		resultSet[0].TransactionResult = transaction.HostNotFound
		resultSet[0].ResultMessage = err.Error()
	} else {
		virtualMta := InstancePool().GetVMtA()

		if virtualMta.TLS {
			client := NewOutboundClientTLS()
			result, resultSet = client.SendMessageTLS(message, virtualMta, domain)
		} else {
			client := NewOutboundClient()
			result, resultSet = client.SendMessageNoTLS(message, virtualMta, domain)
		}
	}

	if result {
		//AllMessage have the same error
		logger.Infof("%s %s %s", resultSet[0].TransactionResult, resultSet[0].ResultMessage, OS.NewLine)
	} else {
		//different transaction result eachother
		for _, row := range resultSet {
			logger.Infof("%s %s %s", row.TransactionResult, row.ResultMessage, OS.NewLine)
		}
	}
}
