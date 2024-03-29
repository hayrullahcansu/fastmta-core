package entity

import (
	"time"

	"github.com/hayrullahcansu/fastmta-core/queue/priority"

	dkim "github.com/emersion/go-dkim"
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
	Priority           priority.Priority `gorm:"Column:priority;"`
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
	MessageID          string            `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	Host               string            `gorm:"type:varchar(200);not null;Column:host"`
	RcptTo             string            `gorm:"type:varchar(200);not null;Column:rcpt_to"`
	MailFrom           string            `gorm:"type:varchar(250);not null;Column:mail_from"`
	Data               string            `gorm:"not null;Column:data"`
	Status             string            `gorm:"type:varchar(10);not null;Column:status;INDEX"`
	MimeMode           string            `gorm:"type:varchar(10);not null;Column:mimemod"`
	MessageDestination string            `gorm:"type:varchar(10);not null;Column:message_destination"`
	AttemptSendTime    time.Time         `gorm:"Column:attempt_send_time;INDEX"`
	DeferredCount      int               `gorm:"Column:deffered_count;"`
	GroupID            int               `gorm:"Column:group_id;"`
	Priority           priority.Priority `gorm:"Column:priority;"`
}

type Transaction struct {
	gorm.Model
	MessageID         string    `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	TransactionTime   time.Time `gorm:"Column:transaction_timestamp;INDEX"`
	ServerHostname    string    `gorm:"type:varchar(250);not null;Column:server_hostname"`
	ServerResponse    string    `gorm:"Column:server_response"`
	TransactionStatus int       `gorm:"Column:transaction_status"`
}
type TransactionLog struct {
	gorm.Model
	MessageID       string    `gorm:"INDEX;not null;Column:message_id"`
	Writer          string    `gorm:"type:varchar(250);not null;Column:writer"`
	IO              string    `gorm:"Column:io"`
	Message         string    `gorm:"Column:message"`
	Command         string    `gorm:"Column:command"`
	TransactionTime time.Time `gorm:"Column:transaction_timestamp;INDEX"`
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

type Dkimmer struct {
	gorm.Model
	DomainName string            `gorm:"type:varchar(500);INDEX;Column:base_domain_name"`
	Selector   string            `gorm:"type:varchar(500);INDEX;Column:selector"`
	PrivateKey string            `gorm:"type:varchar(1000);Column:private_key"`
	Enabled    bool              `gorm:"Column:enabled"`
	Options    *dkim.SignOptions `gorm:"-"`
}

type Settings struct {
	DefaultSend string            `gorm:"type:varchar(500);INDEX;Column:base_domain_name"`
	Selector    string            `gorm:"type:varchar(500);INDEX;Column:selector"`
	PrivateKey  string            `gorm:"type:varchar(1000);Column:private_key"`
	Enabled     bool              `gorm:"Column:enabled"`
	Options     *dkim.SignOptions `gorm:"-"`
}
