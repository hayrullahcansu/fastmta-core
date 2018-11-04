package smtp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/mail"
	"strings"
	"time"

	OS "../../cross"
	"../../entity"
	"github.com/google/uuid"
)

const (
	ReadDeadLine  = time.Second * time.Duration(30)
	WriteDeadLine = time.Second * time.Duration(30)
	MtaName       = "ZetaMail"
	MaxErrorLimit = 10
)

func InboundHandler(server *SmtpServer, conn net.Conn) {
	defer conn.Close()
	err := WriteAll(conn, fmt.Sprintf("220 %s ESMTP %s Ready", server.VmtaHostName, MtaName))
	if err != nil {

	}
	errorCounter := 0
	hasHello := false
	var mtaMessage *entity.MessageTransaction
	for {
		cmdOrginal, err := ReadAll(conn)
		cmd := strings.ToUpper(cmdOrginal)
		if err != nil {
			//TODO: fix this error
		}
		if cmd == "" {
			//probably wrong command
			errorCounter++
			err = WriteAll(conn, "500 5.5.2 Error: bad syntax")
			continue
		}
		if errorCounter >= MaxErrorLimit {
			_ = WriteAll(conn, fmt.Sprintf("421 4.7.0 %s Error: too many errors", server.VmtaHostName))
			break
		}
		//SMTP Commands that can be run before HELO is issued by client.
		if cmd == "QUIT" {
			_ = WriteAll(conn, "221 Goodbye")
			break
		}
		if cmd == "RSET" {
			_ = WriteAll(conn, "250 Ok")
			//TODO: Reset the mail transaction state. Forget any mail from rcpt to data.
			mtaMessage = nil
			continue
		}
		if cmd == "NOOP" {
			_ = WriteAll(conn, "250 Ok")
			continue
		}

		if strings.HasPrefix(cmd, "HELO") || strings.HasPrefix(cmd, "EHLO") {

			if strings.Index(cmd, " ") < 0 {
				errorCounter++
				_ = WriteAll(conn, "501 Syntax error")
				continue
			}
			heloHost := strings.Trim(cmd[strings.Index(cmd, " "):], " ")
			if strings.Index(heloHost, " ") >= 0 {
				errorCounter++
				_ = WriteAll(conn, "501 Syntax error")
				continue
			}
			hasHello = true
			if strings.HasPrefix(cmd, "HELO") {
				_ = WriteAll(conn, fmt.Sprintf("250 %s Hello [%s]", server.VmtaHostName, conn.RemoteAddr().String()))
			} else {
				msg := ""
				msg += fmt.Sprintf("250-%s Hello [%s]%s", server.VmtaHostName, conn.RemoteAddr().String(), OS.NewLine)
				msg += fmt.Sprintf("250-8BITMIME%s", OS.NewLine)
				msg += fmt.Sprintf("250-PIPELINING%s", OS.NewLine)
				msg += fmt.Sprintf("250 Ok")
				_ = WriteAll(conn, msg)
			}
			continue
		}
		if !hasHello {
			errorCounter++
			_ = WriteAll(conn, "503 5.5.1 Error: send HELO/EHLO first")
			continue
		}
		if cmd == "AUTH LOGIN" {
			errorCounter++
			_ = WriteAll(conn, "503 5.5.1 Error: authentication not enabled")
			//TODO: Enable Authentication
			continue
		}
		if strings.HasPrefix(cmd, "MAIL FROM:") {
			mtaMessage = &entity.MessageTransaction{
				MessageID: uuid.New().String(),
				RcptTo:    make([]string, 0),
				Headers:   make(map[string]*string),
			}
			bodyParaIndex := strings.Index(cmd, " BODY=")
			mimeMode := ""
			fromData := ""
			if bodyParaIndex > -1 {
				mimeMode = strings.Trim(cmd[bodyParaIndex+6:], " ")
				fromData = strings.Trim(cmdOrginal[0:bodyParaIndex], " ")
				if mimeMode == "7BIT" {
					mtaMessage.MimeMode = mimeMode
				} else if mimeMode == "8BITMIME" {
					mtaMessage.MimeMode = mimeMode
				} else {
					errorCounter++
					_ = WriteAll(conn, "501 Syntax error")
					continue
				}
			}
			mailFrom := ""
			address := strings.Trim(fromData[strings.Index(fromData, ":")+1:], " ")
			if address != "<>" {
				mailUser, err := mail.ParseAddress(address)
				if err != nil {
					errorCounter++
					_ = WriteAll(conn, "501 Syntax error")
					continue
				}
				mailFrom = mailUser.Address
				mtaMessage.MailFrom = mailFrom
				_ = WriteAll(conn, "250 Ok")
				continue
			}
		}

		//CONTUNIE HERE
		if strings.HasPrefix(cmd, "RCPT TO:") {
			if mtaMessage == nil || mtaMessage.MailFrom == "" {
				_ = WriteAll(conn, "503 Bad sequence of commands")
				continue
			}
			mailUser, err := mail.ParseAddress(strings.Trim(cmd[strings.Index(cmd, ":")+1:], " "))
			if err != nil {
				errorCounter++
				_ = WriteAll(conn, "501 Syntax error")
				continue
			}
			mailUserComponents := strings.Split(mailUser.Address, "@")
			_, domain := mailUserComponents[0], mailUserComponents[1]
			if domain == "" {
				//TODO: LOCALDELIVERY
			}
			mtaMessage.MessageDestination = "Relay"
			mtaMessage.RcptTo = append(mtaMessage.RcptTo, mailUser.Address)
			_ = WriteAll(conn, "250 Ok")
			continue
		}
		if cmd == "DATA" {
			if mtaMessage == nil || mtaMessage.MailFrom == "" {
				errorCounter++
				_ = WriteAll(conn, "503 Bad sequence of commands")
				continue
			} else if len(mtaMessage.RcptTo) < 1 {
				errorCounter++
				_ = WriteAll(conn, "554 No valid recipients")
				continue
			}
			_ = WriteAll(conn, "354 Go ahead")
			data, err := ReadData(conn)
			if err != nil {
				//TODO: fix this error
			}
			//TODO: AddHeader -> Received
			mtaMessage.Data = data
			_ = WriteAll(conn, "250 Message queued for delivery")
			mtaMessage = nil
			continue
		}
		errorCounter++
		_ = WriteAll(conn, "502 5.5.2 Error: command not recognized")
		continue
	}

}

func ReadData(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(ReadDeadLine))
	var builder strings.Builder
	reader := bufio.NewReader(conn)
	var dataLine string
	for {
		readLine, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		dataLine = string(readLine)
		if dataLine == "." {
			break
		}
		builder.WriteString(dataLine)
	}
	return builder.String(), nil
}

func ReadAll(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(ReadDeadLine))
	data, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	return string(data), err
}

func WriteAll(conn net.Conn, data string) error {
	conn.SetWriteDeadline(time.Now().Add(WriteDeadLine))
	_, err := conn.Write([]byte(data + OS.NewLine))
	return err
}