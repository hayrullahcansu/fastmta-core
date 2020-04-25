package entity

import "time"

/*
CreateNewTransactionLog insert new row into transaction_log table
If it cannot access database, the process will panic

messageID identity of the message that assigned in inbound_consumer.go

writer should be destination host or virtualMTA host

IO should be IN or OUT. it means which direction
(INBOUND, OUTBOUND) will produce the log

Message should be pure string. Put "original data..."
rather than original data after the 354 command.

command can be EHLO, HELO, DATA, QUIT...
*/
func CreateNewTransactionLog(messageID string, IO string, writer string, message string, command string) {
	db, err := GetDbContext()
	PanicOnError(err)
	log := &TransactionLog{
		MessageID:       messageID,
		Writer:          writer,
		IO:              IO,
		Message:         message,
		Command:         command,
		TransactionTime: time.Now(),
	}
	db.Create(&log)
}
