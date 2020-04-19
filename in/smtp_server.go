package in

import (
	"fmt"
	"net"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/mta"
	"github.com/hayrullahcansu/fastmta-core/queue"
	"github.com/hayrullahcansu/fastmta-core/smtp/cmd"
	"github.com/hayrullahcansu/fastmta-core/smtp/rw"
)

type SmtpServer struct {
	ID           string
	VMta         *mta.VirtualMta
	VmtaHostName string
	VmtaIPAddr   *net.IPAddr
	Port         int
}

func CreateNewSmtpServer(vmta *mta.VirtualMta) *SmtpServer {
	return &SmtpServer{
		ID:           "",
		VMta:         vmta,
		VmtaHostName: vmta.VmtaHostName,
		VmtaIPAddr:   vmta.VmtaIPAddr,
		Port:         vmta.Port,
	}
}

func (smtpServer *SmtpServer) Run() {
	mergedAddress := fmt.Sprintf("%s:%d", smtpServer.VMta.IPAddressString, smtpServer.Port)
	listener, err := net.Listen("tcp", mergedAddress)

	//TODO: implement TLS inbound
	// tls.Listen("tcp", mergedAddress,)

	if err != nil {
		panic(fmt.Sprintf("%s Can't listen inbound %s", mergedAddress, err))
		//LOG
	}
	logger.Infof("%s Listening%s", mergedAddress, cross.NewLine)
	defer listener.Close()
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			//LOG
		}
		// Handle inbound connections in a new goroutine.
		go smtpServer.InboundHandler(conn)
	}
}

func (server *SmtpServer) InboundHandler(conn net.Conn) {
	defer conn.Close()
	t := rw.NewNoTLSTransporter(conn)
	_cmd := cmd.NewSmtpCommander(t)

	err := _cmd.SmtpReady(server.VmtaHostName, server.VmtaHostName)
	if err != nil {
	}
	errorCounter := 0
	hasHello := false
	var mtaMessage *entity.InboundMessage
	for {
		cmdOrginal, err := _cmd.ReadCommand()
		cmd := strings.ToUpper(cmdOrginal)
		if err != nil {
			//TODO: fix this error
			logger.Errorf("%s"+cross.NewLine, err)
			_ = _cmd.Timeout420()
			break
		}
		if cmd == "" {
			//probably wrong command
			errorCounter++
			err = _cmd.EmptyCommand()
			continue
		}
		if errorCounter >= constant.MaxErrorLimit {
			_ = _cmd.ErrorLimitExceed(server.VmtaHostName)
			break
		}
		//SMTP Commands that can be run before HELO is issued by client.
		if cmd == "QUIT" {
			_ = _cmd.GoodBye221()
			break
		}
		if cmd == "RSET" {
			_ = _cmd.Ok250()
			mtaMessage = nil
			continue
		}
		if cmd == "NOOP" {
			_ = _cmd.Ok250()
			continue
		}

		if strings.HasPrefix(cmd, "HELO") || strings.HasPrefix(cmd, "EHLO") {
			if strings.Index(cmd, " ") < 0 {
				errorCounter++
				_ = _cmd.SyntaxError501()
				continue
			}
			heloHost := strings.Trim(cmd[strings.Index(cmd, " "):], " ")
			if strings.Index(heloHost, " ") >= 0 {
				errorCounter++
				_ = _cmd.SyntaxError501()
				continue
			}
			hasHello = true
			if strings.HasPrefix(cmd, "HELO") {
				_ = _cmd.Hello250(server.VmtaHostName, conn.RemoteAddr().String())
			} else {
				msg := ""
				msg += fmt.Sprintf("250-%s Hello [%s]%s", server.VmtaHostName, conn.RemoteAddr().String(), cross.NewLine)
				msg += fmt.Sprintf("250-8BITMIME%s", cross.NewLine)
				msg += fmt.Sprintf("250-PIPELINING%s", cross.NewLine)
				msg += fmt.Sprintf("250 Ok")
				_ = _cmd.WriteLine(msg)
			}
			continue
		}
		if !hasHello {
			errorCounter++
			_ = _cmd.HelloFirst503()
			continue
		}
		if cmd == "AUTH LOGIN" {
			errorCounter++
			_ = _cmd.AuthenticationNotEnabled503()
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
					_ = _cmd.SyntaxError501()
					continue
				}
			}
			mailFrom := ""
			address := strings.Trim(fromData[strings.Index(fromData, ":")+1:], " ")
			if address != "<>" {
				mailUser, err := mail.ParseAddress(address)
				if err != nil {
					errorCounter++
					_ = _cmd.SyntaxError501()
					continue
				}
				mailFrom = mailUser.Address
				mtaMessage.MailFrom = mailFrom
				_ = _cmd.Ok250()
				continue
			}
		}

		//CONTUNIE HERE
		if strings.HasPrefix(cmd, "RCPT TO:") {
			if mtaMessage == nil || mtaMessage.MailFrom == "" {
				_ = _cmd.BadSequenceCommands503()
				continue
			}
			mailUser, err := mail.ParseAddress(strings.Trim(cmd[strings.Index(cmd, ":")+1:], " "))
			if err != nil {
				errorCounter++
				_ = _cmd.SyntaxError501()
				continue
			}
			mailUserComponents := strings.Split(mailUser.Address, "@")
			_, domain := mailUserComponents[0], mailUserComponents[1]
			if domain == "" {
				//TODO: LOCALDELIVERY
			}
			mtaMessage.MessageDestination = "Relay"
			mtaMessage.RcptTo = append(mtaMessage.RcptTo, mailUser.Address)
			_ = _cmd.Ok250()
			continue
		}
		if cmd == "DATA" {
			if mtaMessage == nil || mtaMessage.MailFrom == "" {
				errorCounter++
				_ = _cmd.BadSequenceCommands503()
				continue
			} else if len(mtaMessage.RcptTo) < 1 {
				errorCounter++
				_ = _cmd.NoValidRecipients554()
				continue
			}
			_ = _cmd.GoAHead354()
			data, err := _cmd.ReadData()
			if err != nil {
				//TODO: fix this error
			}

			//TODO: Check validity
			mtaMessage.Data = fmt.Sprintf("Received: by %s;%s%s%s", server.VmtaHostName, cross.NewLine, time.Now().UTC(), cross.NewLine) + data

			ok, err := AppendMessage(server, mtaMessage)
			if ok {
				_ = _cmd.QueuedForDelivery250()
				mtaMessage = nil
				continue
			}
			_ = _cmd.ExchangeServerStopped432()
		}
		errorCounter++
		_ = _cmd.CommandNotRecognized502()
		continue
	}

}

func AppendMessage(server *SmtpServer, message *entity.InboundMessage) (bool, error) {
	err := queue.Instance().EnqueueInbound(message)
	if err == nil {
		return true, err
	}
	return false, err
}
