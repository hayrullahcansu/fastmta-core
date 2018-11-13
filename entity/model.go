package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

type InboundMessage struct {
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

type OutboundMessageTransaction struct {
	gorm.Model
	MessageID          string
	Headers            map[string]*string
	RcptTo             string
	MailFrom           string
	Data               string
	Status             string
	MimeMode           string
	MessageDestination string
}

type Message struct {
	gorm.Model
	MessageID          string `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	Host               string `gorm:"type:varchar(200);not null;Column:host"`
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

type Domain struct {
	gorm.Model
	DomainName string     `gorm:"type:varchar(250);UNIQUE_INDEX;Column:domain_name"`
	MXRecords  []MXRecord `gorm:"foreignkey:DomainID"`
}

type MXRecord struct {
	gorm.Model
	DomainID       uint
	BaseDomainName string `gorm:"type:varchar(500);INDEX;Column:base_domain_name"`
	Host           string `gorm:"type:varchar(500);Column:mx_host_name"`
	Pref           uint16 `gorm:"Column:mx_preference"`
}
