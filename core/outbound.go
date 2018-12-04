package core

import (
	OS "../cross"
	"../entity"
	"../logger"
	"./transaction"
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
		logger.Info.Printf("%s %s %s", resultSet[0].TransactionResult, resultSet[0].ResultMessage, OS.NewLine)
	} else {
		//different transaction result eachother
		for _, row := range resultSet {
			logger.Info.Printf("%s %s %s", row.TransactionResult, row.ResultMessage, OS.NewLine)
		}
	}
}
