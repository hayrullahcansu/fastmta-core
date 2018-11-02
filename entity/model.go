package entity

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	MessageID string `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	RcptTo    string `gorm:"type:varchar(200);not null;Column:rcpt_to"`
	MailFrom  string `gorm:"type:varchar(250);not null;Column:mail_from"`
	Data      string `gorm:"not null;Column:data"`
	Status    string `gorm:"type:varchar(10);not null;Column:status;INDEX"`
}

type Transaction struct {
	gorm.Model
	MessageID         string `gorm:"UNIQUE_INDEX;not null;Column:message_id"`
	Timestamp         int    `gorm:"Column:transaction_timestamp;INDEX"`
	ServerHostname    string `gorm:"type:varchar(250);not null;Column:server_hostname"`
	ServerResponse    string `gorm:"Column:server_response"`
	TransactionStatus int    `gorm:"Column:transaction_status"`
}

type User struct {
	gorm.Model
	Name         string
	Age          sql.NullInt64
	Birthday     *time.Time
	Email        string  `gorm:"type:varchar(100);unique_index;"`
	Role         string  `gorm:"size:255"`        // set field size to 255
	MemberNumber *string `gorm:"unique;not null"` // set member number to unique and not null
	Num          int     `gorm:"AUTO_INCREMENT"`  // set num to auto incrementable
	Address      string  `gorm:"index:addr"`      // create index with name `addr` for address
	IgnoreMe     int     `gorm:"-"`               // ignore this field
}
