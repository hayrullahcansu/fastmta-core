package cmd

import (
	"fmt"

	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/smtp/rw"
)

type SmtpCommander struct {
	t rw.Transporter
}

func NewSmtpCommander(t rw.Transporter) *SmtpCommander {
	return &SmtpCommander{
		t: t,
	}
}

/*			SMTP COMMANDS			*/
func (s *SmtpCommander) ExecEhlo(hostName string) error {
	return s.t.WriteLine(fmt.Sprintf("EHLO %s", hostName))
}

func (s *SmtpCommander) ExecHelo(hostName string) error {
	return s.t.WriteLine(fmt.Sprintf("HELO %s", hostName))
}

func (s *SmtpCommander) ExecMailFrom(mailFrom string) error {
	return s.t.WriteLine(fmt.Sprintf("MAIL FROM: <%s>", mailFrom))
}

func (s *SmtpCommander) ExecRcptTo(rcptTo string) error {
	return s.t.WriteLine(fmt.Sprintf("RCPT TO: <%s>", rcptTo))
}

func (s *SmtpCommander) ExecDataCommand() error {
	return s.t.WriteLine("DATA")
}

func (s *SmtpCommander) ExecData(data string) error {
	return s.t.WriteLine(fmt.Sprintf("%s%s.", data, cross.NewLine))
}

/*			WRITING COMMANDS			*/
func (s *SmtpCommander) WriteLine(data string) error {
	return s.t.WriteLine(data)
}

func (s *SmtpCommander) SmtpReady(hostName string, mtaName string) error {
	return s.t.WriteLine(fmt.Sprintf("220 %s ESMTP %s Ready", hostName, mtaName))
}

func (s *SmtpCommander) Hello250(hostName string, ipAddr string) error {
	return s.t.WriteLine(fmt.Sprintf("250 %s Hello [%s]", hostName, ipAddr))
}

/*			SUCCESS COMMANDS			*/

func (s *SmtpCommander) GoodBye221() error {
	return s.t.WriteLine("221 Goodbye")
}

func (s *SmtpCommander) Ok250() error {
	return s.t.WriteLine("250 Ok")
}

func (s *SmtpCommander) QueuedForDelivery250() error {
	return s.t.WriteLine("250 Message queued for delivery")
}

func (s *SmtpCommander) GoAHead354() error {
	return s.t.WriteLine("354 Go ahead")
}

/*			ERROR COMMANDS			*/
func (s *SmtpCommander) Timeout420() error {
	return s.t.WriteLine("420 Timeout connection problem")
}

func (s *SmtpCommander) ErrorLimitExceed(hostName string) error {
	return s.t.WriteLine(fmt.Sprintf("421 4.7.0 %s Error: too many errors", hostName))
}

func (s *SmtpCommander) ExchangeServerStopped432() error {
	return s.t.WriteLine("432 The recipient’s Exchange Server incoming mail queue has been stopped")
}

func (s *SmtpCommander) EmptyCommand() error {
	return s.t.WriteLine("500 5.5.2 Error: bad syntax")
}

func (s *SmtpCommander) SyntaxError501() error {
	return s.t.WriteLine("501 Syntax error")
}

func (s *SmtpCommander) CommandNotRecognized502() error {
	return s.t.WriteLine("502 5.5.2 Error: command not recognized")
}

func (s *SmtpCommander) HelloFirst503() error {
	return s.t.WriteLine("503 5.5.1 Error: send HELO/EHLO first")
}

func (s *SmtpCommander) AuthenticationNotEnabled503() error {
	return s.t.WriteLine("503 5.5.1 Error: authentication not enabled")
}

func (s *SmtpCommander) BadSequenceCommands503() error {
	return s.t.WriteLine("503 Bad sequence of commands")
}

func (s *SmtpCommander) NoValidRecipients554() error {
	return s.t.WriteLine("554 No valid recipients")
}

/*			READING COMMANDS			*/
func (s *SmtpCommander) ReadCommand() (string, error) {
	return s.t.ReadAll()
}

func (s *SmtpCommander) ReadAllLine() (string, error) {
	return s.t.ReadAllLine()
}
func (s *SmtpCommander) ReadData() (string, error) {
	return s.t.ReadData()
}

func (s *SmtpCommander) Close() {
	s.t.Close()
}
