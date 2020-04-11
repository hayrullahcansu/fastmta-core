package core

import (
	"encoding/json"
	"fmt"
	"net"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	OS "github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/queue"
)

const (
	ReadDeadLine  = time.Second * time.Duration(30)
	WriteDeadLine = time.Second * time.Duration(30)
	MtaName       = "ZetaMail"
	MaxErrorLimit = 10
)

func InboundHandler(server *InboundSmtpServer, conn net.Conn) {
	defer conn.Close()
	err := WriteLineNoTLS(conn, fmt.Sprintf("220 %s ESMTP %s Ready", server.VmtaHostName, MtaName))
	if err != nil {
	}
	errorCounter := 0
	hasHello := false
	var mtaMessage *entity.InboundMessage
	for {
		cmdOrginal, err := ReadAllNoTLS(conn)
		cmd := strings.ToUpper(cmdOrginal)
		if err != nil {
			//TODO: fix this error
			fmt.Printf("%s"+OS.NewLine, err)
			_ = WriteLineNoTLS(conn, "420 Timeout connection problem")
			break
		}
		if cmd == "" {
			//probably wrong command
			errorCounter++
			err = WriteLineNoTLS(conn, "500 5.5.2 Error: bad syntax")
			continue
		}
		if errorCounter >= MaxErrorLimit {
			_ = WriteLineNoTLS(conn, fmt.Sprintf("421 4.7.0 %s Error: too many errors", server.VmtaHostName))
			break
		}
		//SMTP Commands that can be run before HELO is issued by client.
		if cmd == "QUIT" {
			_ = WriteLineNoTLS(conn, "221 Goodbye")
			break
		}
		if cmd == "RSET" {
			_ = WriteLineNoTLS(conn, "250 Ok")
			mtaMessage = nil
			continue
		}
		if cmd == "NOOP" {
			_ = WriteLineNoTLS(conn, "250 Ok")
			continue
		}

		if strings.HasPrefix(cmd, "HELO") || strings.HasPrefix(cmd, "EHLO") {

			if strings.Index(cmd, " ") < 0 {
				errorCounter++
				_ = WriteLineNoTLS(conn, "501 Syntax error")
				continue
			}
			heloHost := strings.Trim(cmd[strings.Index(cmd, " "):], " ")
			if strings.Index(heloHost, " ") >= 0 {
				errorCounter++
				_ = WriteLineNoTLS(conn, "501 Syntax error")
				continue
			}
			hasHello = true
			if strings.HasPrefix(cmd, "HELO") {
				_ = WriteLineNoTLS(conn, fmt.Sprintf("250 %s Hello [%s]", server.VmtaHostName, conn.RemoteAddr().String()))
			} else {
				msg := ""
				msg += fmt.Sprintf("250-%s Hello [%s]%s", server.VmtaHostName, conn.RemoteAddr().String(), OS.NewLine)
				msg += fmt.Sprintf("250-8BITMIME%s", OS.NewLine)
				msg += fmt.Sprintf("250-PIPELINING%s", OS.NewLine)
				msg += fmt.Sprintf("250 Ok")
				_ = WriteLineNoTLS(conn, msg)
			}
			continue
		}
		if !hasHello {
			errorCounter++
			_ = WriteLineNoTLS(conn, "503 5.5.1 Error: send HELO/EHLO first")
			continue
		}
		if cmd == "AUTH LOGIN" {
			errorCounter++
			_ = WriteLineNoTLS(conn, "503 5.5.1 Error: authentication not enabled")
			//TODO: Enable Authentication
			continue
		}
		if strings.HasPrefix(cmd, "MAIL FROM:") {
			mtaMessage = &entity.InboundMessage{
				MessageID: uuid.New().String(),
				RcptTo:    make([]string, 0),
				Headers:   make(map[string]*string),
			}
			bodyParaIndex := strings.Index(cmd, " BODY=")
			mimeMode := ""
			fromData := cmdOrginal
			if bodyParaIndex > -1 {
				mimeMode = strings.Trim(cmd[bodyParaIndex+6:], " ")
				fromData = strings.Trim(cmdOrginal[0:bodyParaIndex], " ")
				if mimeMode == "7BIT" {
					mtaMessage.MimeMode = mimeMode
				} else if mimeMode == "8BITMIME" {
					mtaMessage.MimeMode = mimeMode
				} else {
					errorCounter++
					_ = WriteLineNoTLS(conn, "501 Syntax error")
					continue
				}
			}
			mailFrom := ""
			address := strings.Trim(fromData[strings.Index(fromData, ":")+1:], " ")
			if address != "<>" {
				mailUser, err := mail.ParseAddress(address)
				if err != nil {
					errorCounter++
					_ = WriteLineNoTLS(conn, "501 Syntax error")
					continue
				}
				mailFrom = mailUser.Address
				mtaMessage.MailFrom = mailFrom
				_ = WriteLineNoTLS(conn, "250 Ok")
				continue
			}
		}

		//CONTUNIE HERE
		if strings.HasPrefix(cmd, "RCPT TO:") {
			if mtaMessage == nil || mtaMessage.MailFrom == "" {
				_ = WriteLineNoTLS(conn, "503 Bad sequence of commands")
				continue
			}
			mailUser, err := mail.ParseAddress(strings.Trim(cmd[strings.Index(cmd, ":")+1:], " "))
			if err != nil {
				errorCounter++
				_ = WriteLineNoTLS(conn, "501 Syntax error")
				continue
			}
			mailUserComponents := strings.Split(mailUser.Address, "@")
			_, domain := mailUserComponents[0], mailUserComponents[1]
			if domain == "" {
				//TODO: LOCALDELIVERY
			}
			mtaMessage.MessageDestination = "Relay"
			mtaMessage.RcptTo = append(mtaMessage.RcptTo, mailUser.Address)
			_ = WriteLineNoTLS(conn, "250 Ok")
			continue
		}
		if cmd == "DATA" {
			if mtaMessage == nil || mtaMessage.MailFrom == "" {
				errorCounter++
				_ = WriteLineNoTLS(conn, "503 Bad sequence of commands")
				continue
			} else if len(mtaMessage.RcptTo) < 1 {
				errorCounter++
				_ = WriteLineNoTLS(conn, "554 No valid recipients")
				continue
			}
			_ = WriteLineNoTLS(conn, "354 Go ahead")
			data, err := ReadDataNoTLS(conn)
			if err != nil {
				//TODO: fix this error
			}

			//TODO: Check validity
			mtaMessage.Data = fmt.Sprintf("Received: by %s;%s%s%s", server.VmtaHostName, OS.NewLine, time.Now().UTC(), OS.NewLine) + data

			ok, err := AppendMessage(server, mtaMessage)
			if ok {
				_ = WriteLineNoTLS(conn, "250 Message queued for delivery")
				mtaMessage = nil
				continue
			}
			_ = WriteLineNoTLS(conn, "432 The recipient’s Exchange Server incoming mail queue has been stopped")
		}
		errorCounter++
		_ = WriteLineNoTLS(conn, "502 5.5.2 Error: command not recognized")
		continue
	}

}

func AppendMessage(server *InboundSmtpServer, message *entity.InboundMessage) (bool, error) {
	data, err := json.Marshal(message)
	if err == nil {
		err = server.RabbitMqClient.Publish(queue.InboundExchange, queue.RoutingKeyInbound, false, false, data)
		if err == nil {
			return true, err
		}
	}
	return false, err
}
