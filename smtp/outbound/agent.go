package outbound

import (
	"net"
	"strings"

	"github.com/hayrullahcansu/fastmta-core/dns"
	"github.com/hayrullahcansu/fastmta-core/mta"
	"github.com/hayrullahcansu/fastmta-core/smtp/cmd"
	"github.com/hayrullahcansu/fastmta-core/smtp/rw"

	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/transaction"
)

type Agent struct {
	dialer rw.Dialer
	// virtualMta    *VirtualMta
	// messages      []*entity.Message
	virtualMta    *mta.VirtualMta
	_cmd          *cmd.SmtpCommander
	canPipeLining bool
	canTLS        bool
}

//TODO: fill params
func NewAgent(vmt *mta.VirtualMta) *Agent {
	agent := &Agent{
		virtualMta: vmt,
	}
	dialer := &net.Dialer{
		Timeout:   Timeout,
		KeepAlive: KeepAlive,
		LocalAddr: &net.TCPAddr{
			IP: []byte(agent.virtualMta.IPAddressString),
		},
	}
	if agent.canTLS {
		agent.dialer = rw.NewTLSDialer(dialer)
	} else {
		agent.dialer = rw.NewNoTLSDialer(dialer)
	}
	return agent
}

func (c *Agent) SendMessage(message *entity.Message) (bool, *transaction.TransactionGroupResult) {
	var result = &transaction.TransactionGroupResult{}
	if false {
		//if anyrule blocks to SendMessage
		c.ExecQuit()
	}
	domain, err := dns.NewDomain(message.Host)
	if err != nil {
		result.TransactionResult = transaction.FailedToConnect
		result.ResultMessage = err.Error()
		return true, result
	}

	ok, r, msg := c.Connect(domain)

	if !ok || r != transaction.Success {
		result.TransactionResult = r
		result.ResultMessage = msg
		return true, result
	}

	// Read the Server greeting.
	lines, err := c._cmd.ReadAllLine()
	if err != nil {
		result.TransactionResult = transaction.FailedToConnect
		result.ResultMessage = err.Error()
		return true, result
	}
	// Check we get a valid banner.
	if !strings.HasPrefix(lines, "2") {
		if strings.HasPrefix(lines, "421") {
			result.TransactionResult = transaction.ServiceNotAvalible
			result.ResultMessage = lines
			return true, result
		}
		result.TransactionResult = transaction.FailedToConnect
		result.ResultMessage = lines
		return true, result
	}

	// We have connected, so say helo
	r = c.ExecHelo()
	if r != transaction.Success {
		//TODO: add this rule like "this host not valid or unable to connect"
		result.TransactionResult = r
		result.ResultMessage = "service not avaliable"
		return true, result
	}

	r = c.ExecMailFrom(message.MailFrom)
	if r != transaction.Success {
		result.TransactionResult = r
		result.ResultMessage = "service not avaliable"
		return true, result
	}

	r = c.ExecRcptTo(message.RcptTo)
	if r != transaction.Success {
		result.TransactionResult = r
		result.ResultMessage = "service not avaliable"
		return true, result
	}

	mimeKit := false

	if mimeKit {
		//TODO: add dkim
	} else {
		r = c.ExecData(message.Data)
		if r != transaction.Success {
			result.TransactionResult = r
			result.ResultMessage = "service not avaliable"
			return true, result
		}
	}
	result.TransactionResult = transaction.Success
	result.ResultMessage = ""

	return false, result
}

func (c *Agent) Connect(domain *dns.Domain) (bool, transaction.TransactionResult, string) {
	err := c.dialer.Deal(domain.MXRecords[0].Host, Port)
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

func (c *Agent) ExecHelo() transaction.TransactionResult {
	// We have connected to the MX, Say EHLO.
	c._cmd.ExecEhlo(c.virtualMta.VmtaHostName)
	lines, _ := c._cmd.ReadAllLine()
	if strings.HasPrefix(lines, "421") {
		return transaction.ServiceNotAvalible
	}
	if !strings.HasPrefix(lines, "2") {
		// If server didn't respond with a success code on EHLO then we should retry with HELO
		_ = c._cmd.ExecHelo(c.virtualMta.VmtaHostName)
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

func (c *Agent) ExecMailFrom(mailFrom string) transaction.TransactionResult {
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

func (c *Agent) ExecRcptTo(rcptTo string) transaction.TransactionResult {
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

func (c *Agent) ExecData(data string) transaction.TransactionResult {
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

func (c *Agent) ExecRset() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (c *Agent) ExecQuit() transaction.TransactionResult {
	return transaction.RetryRequired
}
