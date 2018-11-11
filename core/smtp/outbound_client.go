package smtp

import (
	"fmt"
	"net"
	"regexp"
	"time"

	".."
	"../../entity"
	"../transaction"
)

const (
	Timeout   = time.Second * time.Duration(30)
	KeepAlive = time.Second * time.Duration(30)
	Port      = 25
	//MtaName       = "ZetaMail"
	//MaxErrorLimit = 10
)

type OutboundClient struct {
	dialer     *net.Dialer
	conn       *net.Conn
	virtualMta *VirtualMta
	message    *entity.Message
	domain     *core.Domain
}

func NewOutboundClient() *OutboundClient {
	return &OutboundClient{}
}

func (client *OutboundClient) SendMessage(message *entity.Message, virtualMta *VirtualMta, domain *core.Domain) (transaction.TransactionResult, string) {
	client.virtualMta = virtualMta
	client.message = message
	client.domain = domain
	ok, r := client.CreateTcpClient()
	if !ok {
		return r, "Ip address already in usage"
	}

	if false {
		//if anyrule blocks to send
		client.ExecQuit()
	}

	ok, r, msg := client.Connect()
	if !ok || r != transaction.Success {
		return r, msg
	}

	return transaction.RetryRequired, ""
}

func (client *OutboundClient) CreateTcpClient() (bool, transaction.TransactionResult) {
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

func (c *OutboundClient) Connect() (bool, transaction.TransactionResult, string) {
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
	c.conn = &conn
	return true, transaction.Success, "connected"
}

func (client *OutboundClient) ExecHelo() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (client *OutboundClient) ExecMailFrom() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (client *OutboundClient) ExecRcptTo() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (client *OutboundClient) ExecRset() transaction.TransactionResult {
	return transaction.RetryRequired
}

func (client *OutboundClient) ExecQuit() transaction.TransactionResult {
	return transaction.RetryRequired
}

func checkErr(err error) (transaction.TransactionResult, string) {
	opError, _ := err.(*net.OpError)
	_, _ = opError.Err.(*net.DNSError)
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return transaction.Timeout, "timeout"
	} else if match, _ := regexp.MatchString(".*lookup.*", err.Error()); match {
		return transaction.HostNotFound, "host not found"
	} else if match, _ := regexp.MatchString(".*connection refused.*", err.Error()); match {
		return transaction.RejectedByRemoteServer, "connection refused"
	} else {
		return transaction.ServiceNotAvalible, "service not avaliable"
	}
}
