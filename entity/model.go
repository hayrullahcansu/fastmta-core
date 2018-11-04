package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MessageTransaction struct {
	gorm.Model
	MessageID          string
	Headers            map[string]*string
	RcptTo             []string
	MailFrom           string
	Data               string
	Status             string
	MimeMode           string
	MessageDestination string
}

type Message struct {
	gorm.Model
	MessageID          string `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	RcptTo             string `gorm:"type:varchar(200);not null;Column:rcpt_to"`
	MailFrom           string `gorm:"type:varchar(250);not null;Column:mail_from"`
	Data               string `gorm:"not null;Column:data"`
	Status             string `gorm:"type:varchar(10);not null;Column:status;INDEX"`
	MimeMode           string `gorm:"type:varchar(10);not null;Column:mimemod"`
	MessageDestination string `gorm:"type:varchar(10);not null;Column:message_destination"`
}

type Transaction struct {
	gorm.Model
	MessageID         string     `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	TransactionTime   *time.Time `gorm:"Column:transaction_timestamp;INDEX"`
	ServerHostname    string     `gorm:"type:varchar(250);not null;Column:server_hostname"`
	ServerResponse    string     `gorm:"Column:server_response"`
	TransactionStatus int        `gorm:"Column:transaction_status"`
}
