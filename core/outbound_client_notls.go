package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	OS "../cross"
	"../entity"
	"./transaction"
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
	message       *entity.Message
	domain        *Domain
	canPipeLining bool
}

func NewOutboundClient() *OutboundClient {
	return &OutboundClient{}
}

func (c *OutboundClient) SendMessageNoTLS(message *entity.Message, virtualMta *VirtualMta, domain *Domain) (transaction.TransactionResult, string) {
	c.virtualMta = virtualMta
	c.message = message
	c.domain = domain
	ok, r := c.CreateTcpClientNoTLS()
	if !ok {
		return r, "Ip address already in usage"
	}

	if false {
		//if anyrule blocks to send
		c.ExecQuit()
	}

	ok, r, msg := c.ConnectNoTLS()
	if !ok || r != transaction.Success {
		return r, msg
	}

	// Read the Server greeting.
	lines, err := ReadAllLineNoTLS(c.conn)
	if err != nil {
		return transaction.FailedToConnect, err.Error()
	}
	// Check we get a valid banner.
	if !strings.HasPrefix(lines, "2") {
		if strings.HasPrefix(lines, "421") {
			return transaction.ServiceNotAvalible, lines
		}
		return transaction.FailedToConnect, lines
	}

	// We have connected, so say helo
	r = c.ExecHeloNoTLS()
	if r != transaction.Success {

		//TODO: add this rule like "this host not valid or unable to connect"
		return r, "service not avaliable"
	}

	r = c.ExecMailFromNoTLS(c.message.MailFrom)
	if r != transaction.Success {
		return r, "service not avaliable"
	}

	r = c.ExecRcptToNoTLS(c.message.RcptTo)
	if r != transaction.Success {
		return r, "service not avaliable"
	}

	mimeKit := false

	if mimeKit {
		//TODO: add dkim
	} else {
		r = c.ExecDataNoTLS(c.message.Data)
		if r != transaction.Success {
			return r, "service not avaliable"
		}
	}
	return transaction.Success, ""
}

func (client *OutboundClient) CreateTcpClientNoTLS() (bool, transaction.TransactionResult) {
	if client.virtualMta.HandleLock() {
		client.dialer = &net.Dialer{
			Timeout:   Timeout,
			KeepAlive: KeepAlive,
			LocalAddr: &net.TCPAddr{
				IP:   client.virtualMta.VmtaIPAddr.IP,
				Port: 25,
			},
		}
		return true, transaction.Success
	} else {
		return false, transaction.ClientAlreadyInUse
	}

	return false, transaction.RetryRequired
}

func (c *OutboundClient) ConnectNoTLS() (bool, transaction.TransactionResult, string) {
	conn, err := c.dialer.Dial("tcp", fmt.Sprintf("%s:%s", c.domain.MXRecords[0].Host, Port))
	if err != nil {
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
