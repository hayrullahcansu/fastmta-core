package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hayrullahcansu/zetamail/core/transaction"
	OS "github.com/hayrullahcansu/zetamail/cross"
	"github.com/hayrullahcansu/zetamail/entity"
)

const (
	Timeout   = time.Second * time.Duration(30)
	KeepAlive = time.Second * time.Duration(30)
	Port      = 25
	//MtaName       = "ZetaMail"
	//MaxErrorLimit = 10
)

type OutboundClient struct {
	dialer        *net.Dialer
	conn          net.Conn
	virtualMta    *VirtualMta
	messages      []*entity.Message
	domain        *Domain
	canPipeLining bool
}

func NewOutboundClient() *OutboundClient {
	return &OutboundClient{}
}

func (c *OutboundClient) SendMessageNoTLS(messages []*entity.Message, virtualMta *VirtualMta, domain *Domain) (bool, []*transaction.TransactionGroupResult) {
	c.virtualMta = virtualMta
	c.messages = messages
	c.domain = domain
	ok, r := c.CreateTcpClientNoTLS()
	resultSet := make([]*transaction.TransactionGroupResult, len(messages))
	resultSet[0] = &transaction.TransactionGroupResult{}
	if false {
		//if anyrule blocks to send
		c.ExecQuit()
	}

	ok, r, msg := c.ConnectNoTLS()
	if !ok || r != transaction.Success {
		resultSet[0].TransactionResult = r
		resultSet[0].ResultMessage = msg
		return true, resultSet
	}

	// Read the Server greeting.
	lines, err := ReadAllLineNoTLS(c.conn)
	if err != nil {
		resultSet[0].TransactionResult = transaction.FailedToConnect
		resultSet[0].ResultMessage = err.Error()
		return true, resultSet
	}
	// Check we get a valid banner.
	if !strings.HasPrefix(lines, "2") {
		if strings.HasPrefix(lines, "421") {
			resultSet[0].TransactionResult = transaction.ServiceNotAvalible
			resultSet[0].ResultMessage = lines
			return true, resultSet
		}
		resultSet[0].TransactionResult = transaction.FailedToConnect
		resultSet[0].ResultMessage = lines
		return true, resultSet
	}

	// We have connected, so say helo
	r = c.ExecHeloNoTLS()
	if r != transaction.Success {

		//TODO: add this rule like "this host not valid or unable to connect"
		resultSet[0].TransactionResult = r
		resultSet[0].ResultMessage = "service not avaliable"
		return true, resultSet
	}
	for i := 0; i < len(c.messages); i++ {
		resultSet[i] = &transaction.TransactionGroupResult{}
		r = c.ExecMailFromNoTLS(c.messages[i].MailFrom)
		if r != transaction.Success {
			resultSet[i].TransactionResult = r
			resultSet[i].ResultMessage = "service not avaliable"
			continue
		}

		r = c.ExecRcptToNoTLS(c.messages[i].RcptTo)
		if r != transaction.Success {
			resultSet[i].TransactionResult = r
			resultSet[i].ResultMessage = "service not avaliable"
			continue
		}

		mimeKit := false

		if mimeKit {
			//TODO: add dkim
		} else {
			r = c.ExecDataNoTLS(c.messages[i].Data)
			if r != transaction.Success {
				resultSet[i].TransactionResult = r
				resultSet[i].ResultMessage = "service not avaliable"
				continue
			}
		}
		resultSet[i].TransactionResult = transaction.Success
		resultSet[i].ResultMessage = ""
		continue
	}
	return false, resultSet
}

func (client *OutboundClient) CreateTcpClientNoTLS() (bool, transaction.TransactionResult) {
	client.virtualMta.HandleLock()
	client.dialer = &net.Dialer{
		Timeout:   Timeout,
		KeepAlive: KeepAlive,
	}
	return true, transaction.Success
}

func (c *OutboundClient) ConnectNoTLS() (bool, transaction.TransactionResult, string) {
	host := fmt.Sprintf("%s:%d", c.domain.MXRecords[0].Host, Port)
	host = "gmail-smtp-in.l.google.COM:25"
	conn, err := c.dialer.Dial("tcp", host)
	if err != nil {
		fmt.Println(err.Error())
		if opError, ok := err.(*net.OpError); ok {
			if dnsError, ok := opError.Err.(*net.DNSError); ok {
				return false, transaction.HostNotFound, dnsError.Error()
			}
		}
		//TODO: define all error like dnsError
		return false, transaction.ServiceNotAvalible, "service not avaliable"
	}
	c.conn = conn
	return true, transaction.Success, "connected"
}

func (c *OutboundClient) ExecHeloNoTLS() transaction.TransactionResult {
	// We have connected to the MX, Say EHLO.
	WriteLineNoTLS(c.conn, fmt.Sprintf("EHLO %s", c.virtualMta.VmtaHostName))
	lines, _ := ReadAllLineNoTLS(c.conn)
	if strings.HasPrefix(lines, "421") {
		return transaction.ServiceNotAvalible
	}
	if !strings.HasPrefix(lines, "2") {
		// If server didn't respond with a success code on EHLO then we should retry with HELO
		_ = WriteLineNoTLS(c.conn, fmt.Sprintf("HELO %s", c.virtualMta.VmtaHostName))
		lines, _ := ReadAllLineNoTLS(c.conn)
		if !strings.HasPrefix(lines, "250") {
			c.conn.Close()
			return transaction.ServiceNotAvalible
		}
	} else {
		// Server responded to EHLO
		// Check to see if it supports 8BITMIME

		// Check to see if the server supports pipelining
		c.canPipeLining = strings.Index(strings.ToUpper(lines), "PIPELINING") > -1
	}
	return transaction.Success
}

func (c *OutboundClient) ExecMailFromNoTLS(mailFrom string) transaction.TransactionResult {
	err := WriteLineNoTLS(c.conn, fmt.Sprintf("MAIL FROM: <%s>", mailFrom))
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	if !c.canPipeLining {
		lines, _ := ReadAllLineNoTLS(c.conn)
		if !strings.HasPrefix(lines, "250") {
			if strings.HasPrefix(lines, "421") {
				return transaction.ServiceNotAvalible
			}
			return transaction.RejectedByRemoteServer
		}
	}
	return transaction.Success
}

func (c *OutboundClient) ExecRcptToNoTLS(rcptTo string) transaction.TransactionResult {
	err := WriteLineNoTLS(c.conn, fmt.Sprintf("RCPT TO: <%s>", rcptTo))
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	if !c.canPipeLining {
		lines, _ := ReadAllLineNoTLS(c.conn)
		if !strings.HasPrefix(lines, "250") {
			return transaction.RejectedByRemoteServer
		}
	}
	return transaction.Success
}

func (c *OutboundClient) ExecDataNoTLS(data string) transaction.TransactionResult {
	// Data response or Mail From if pipelining
	err := WriteLineNoTLS(c.conn, fmt.Sprintf("DATA"))
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	lines, _ := ReadAllLineNoTLS(c.conn)
	// If the remote MX supports pipelining then we need to check the MAIL FROM and RCPT to responses.
	if c.canPipeLining {
		// Check MAIL FROM OK.
		if !strings.HasPrefix(lines, "250") {
			_, _ = ReadAllLineNoTLS(c.conn) // RCPT TO
			_, _ = ReadAllLineNoTLS(c.conn) // DATA
			return transaction.RejectedByRemoteServer
		}

		// Check RCPT TO OK.
		lines, _ = ReadAllLineNoTLS(c.conn) // RCPT TO
		if !strings.HasPrefix(lines, "250") {
			_, _ = ReadAllLineNoTLS(c.conn) // DATA
			return transaction.RejectedByRemoteServer
		}

		// Get the Data Command response.
		lines, _ = ReadAllLineNoTLS(c.conn) // DATA
	}
	if !strings.HasPrefix(lines, "354") {
		_, _ = ReadAllLineNoTLS(c.conn) // DATA
		return transaction.RejectedByRemoteServer
	}
	err = WriteLineNoTLS(c.conn, fmt.Sprintf("%s%s.", data, OS.NewLine))
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	lines, _ = ReadAllLineNoTLS(c.conn)
	if !strings.HasPrefix(lines, "250") {
		return transaction.RejectedByRemoteServer
	}
	return transaction.Success
}

func (c *OutboundClient) ExecRset() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (c *OutboundClient) ExecQuit() transaction.TransactionResult {
	return transaction.RetryRequired
}
