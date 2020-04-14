package rw

import (
	"bufio"
	"strings"

	"github.com/hayrullahcansu/fastmta-core/smtp/tcptimer"
)

type BaseTransporter struct {
	Transporter
	timer  tcptimer.TcpTimer
	reader *bufio.Reader
}

func NewBaseTransporter(reader *bufio.Reader, timer tcptimer.TcpTimer) *BaseTransporter {
	return &BaseTransporter{
		reader: reader,
		timer:  timer,
	}
}

func (t *BaseTransporter) ReadData() (string, error) {
	t.timer.SetReadDeadlineDefault()
	var builder strings.Builder
	var dataLine string
	for {
		readLine, _, err := t.reader.ReadLine()
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

func (t *BaseTransporter) ReadAllLine() (string, error) {
	t.timer.SetReadDeadlineDefault()
	var builder strings.Builder
	var dataLine string

	readLine, _, err := t.reader.ReadLine()
	if err != nil {
		return "421 Connection ended abruptly", err
	}
	if len(readLine) == 0 {
		return "421 Connection ended abruptly", nil
	}

	for readLine[3] != '-' {
		dataLine = string(readLine)
		builder.WriteString(dataLine)
		readLine, _, err = t.reader.ReadLine()
		if err != nil {
			if builder.Len() == 0 {
				return "421 Connection ended abruptly", err
			}
			dataLine = ""
			break
		}

	}
	dataLine = string(readLine)
	builder.WriteString(dataLine)
	return builder.String(), nil
}

func (t *BaseTransporter) ReadAll() (string, error) {
	t.timer.SetReadDeadlineDefault()
	readLine, _, err := t.reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(readLine), nil
}
