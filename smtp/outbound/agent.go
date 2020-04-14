package outbound

import (
	"net"
	"strings"

	"github.com/hayrullahcansu/fastmta-core/conf"
	"github.com/hayrullahcansu/fastmta-core/smtp/cmd"
	"github.com/hayrullahcansu/fastmta-core/smtp/rw"

	"github.com/hayrullahcansu/fastmta-core/core/transaction"
	"github.com/hayrullahcansu/fastmta-core/dns"
	"github.com/hayrullahcansu/fastmta-core/entity"
)

type Agent struct {
	dialer rw.Dialer
	// virtualMta    *VirtualMta
	messages      []*entity.Message
	virtualMta    *conf.VirtualMta
	domain        *dns.Domain
	_cmd          *cmd.SmtpCommander
	canPipeLining bool
	canTLS        bool
}

//TODO: fill params
func NewAgent() *Agent {
	agent := &Agent{}
	dialer := &net.Dialer{
		Timeout:   Timeout,
		KeepAlive: KeepAlive,
		LocalAddr: &net.TCPAddr{
			IP: []byte(agent.virtualMta.IP),
		},
	}
	if agent.canTLS {
		agent.dialer = rw.NewTLSDialer(dialer)
	} else {
		agent.dialer = rw.NewNoTLSDialer(dialer)
	}
	return agent
}

func (c *Agent) SendMessages(messages []*entity.Message, domain *dns.Domain) (bool, []*transaction.TransactionGroupResult) {
	c.messages = messages
	c.domain = domain
	resultSet := make([]*transaction.TransactionGroupResult, len(messages))
	resultSet[0] = &transaction.TransactionGroupResult{}
	if false {
		//if anyrule blocks to SendMessageTLS
		c.ExecQuitTLS()
	}

	ok, r, msg := c.ConnectTLS()

	if !ok || r != transaction.Success {
		resultSet[0].TransactionResult = r
		resultSet[0].ResultMessage = msg
		return true, resultSet
	}

	// Read the Server greeting.
	lines, err := c._cmd.ReadAllLine()
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
	r = c.ExecHeloTLS()
	if r != transaction.Success {
		//TODO: add this rule like "this host not valid or unable to connect"
		resultSet[0].TransactionResult = r
		resultSet[0].ResultMessage = "service not avaliable"
		return true, resultSet
	}
	for i := 0; i < len(c.messages); i++ {
		resultSet[i] = &transaction.TransactionGroupResult{}

		r = c.ExecMailFromTLS(c.messages[i].MailFrom)
		if r != transaction.Success {
			resultSet[i].TransactionResult = r
			resultSet[i].ResultMessage = "service not avaliable"
			continue
		}

		r = c.ExecRcptToTLS(c.messages[i].RcptTo)
		if r != transaction.Success {
			resultSet[i].TransactionResult = r
			resultSet[i].ResultMessage = "service not avaliable"
			continue
		}

		mimeKit := false

		if mimeKit {
			//TODO: add dkim
		} else {
			r = c.ExecDataTLS(c.messages[i].Data)
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

func (c *Agent) ConnectTLS() (bool, transaction.TransactionResult, string) {
	err := c.dialer.Deal(c.domain.MXRecords[0].Host, Port)
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if dnsError, ok := opError.Err.(*net.DNSError); ok {
				return false, transaction.HostNotFound, dnsError.Error()
			}
		}
		//TODO: define all error like dnsError
		return false, transaction.ServiceNotAvalible, "service not avaliable"
	}

	c._cmd = cmd.NewSmtpCommander(c.dialer.GetTransporter())
	return true, transaction.Success, "connected"
}

func (c *Agent) ExecHeloTLS() transaction.TransactionResult {
	// We have connected to the MX, Say EHLO.
	c._cmd.ExecEhlo(c.virtualMta.HostName)
	lines, _ := c._cmd.ReadAllLine()
	if strings.HasPrefix(lines, "421") {
		return transaction.ServiceNotAvalible
	}
	if !strings.HasPrefix(lines, "2") {
		// If server didn't respond with a success code on EHLO then we should retry with HELO
		_ = c._cmd.ExecHelo(c.virtualMta.HostName)
		lines, _ := c._cmd.ReadAllLine()
		if !strings.HasPrefix(lines, "250") {
			c._cmd.Close()
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

func (c *Agent) ExecMailFromTLS(mailFrom string) transaction.TransactionResult {
	err := c._cmd.ExecMailFrom(mailFrom)
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	if !c.canPipeLining {
		lines, _ := c._cmd.ReadAllLine()
		if !strings.HasPrefix(lines, "250") {
			if strings.HasPrefix(lines, "421") {
				return transaction.ServiceNotAvalible
			}
			return transaction.RejectedByRemoteServer
		}
	}
	return transaction.Success
}

func (c *Agent) ExecRcptToTLS(rcptTo string) transaction.TransactionResult {
	err := c._cmd.ExecRcptTo(rcptTo)
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	if !c.canPipeLining {
		lines, _ := c._cmd.ReadAllLine()
		if !strings.HasPrefix(lines, "250") {
			return transaction.RejectedByRemoteServer
		}
	}
	return transaction.Success
}

func (c *Agent) ExecDataTLS(data string) transaction.TransactionResult {
	// Data response or Mail From if pipelining
	err := c._cmd.ExecDataCommand()
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	lines, _ := c._cmd.ReadAllLine()
	// If the remote MX supports pipelining then we need to check the MAIL FROM and RCPT to responses.
	if c.canPipeLining {
		// Check MAIL FROM OK.
		if !strings.HasPrefix(lines, "250") {
			_, _ = c._cmd.ReadAllLine() // RCPT TO
			_, _ = c._cmd.ReadAllLine() // DATA
			return transaction.RejectedByRemoteServer
		}

		// Check RCPT TO OK.
		lines, _ = c._cmd.ReadAllLine() // RCPT TO
		if !strings.HasPrefix(lines, "250") {
			_, _ = c._cmd.ReadAllLine() // DATA
			return transaction.RejectedByRemoteServer
		}

		// Get the Data Command response.
		lines, _ = c._cmd.ReadAllLine() // DATA
	}
	if !strings.HasPrefix(lines, "354") {
		_, _ = c._cmd.ReadAllLine() // DATA
		return transaction.RejectedByRemoteServer
	}
	err = c._cmd.ExecData(data)
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {
			if opError.Timeout() {
				return transaction.Timeout
			}
		}
		return transaction.ServiceNotAvalible
	}
	lines, _ = c._cmd.ReadAllLine()
	if !strings.HasPrefix(lines, "250") {
		return transaction.RejectedByRemoteServer
	}
	return transaction.Success
}

func (c *Agent) ExecRsetTLS() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (c *Agent) ExecQuitTLS() transaction.TransactionResult {
	return transaction.RetryRequired
}
