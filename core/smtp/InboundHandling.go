package smtp

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	OS "../../cross"
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
	for {
		cmd, err := ReadAll(conn)
		cmd = strings.ToUpper(cmd)
		if err != nil {
			//error
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

	}

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
