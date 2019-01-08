package core

import (
	"bufio"
	"net"
	"strings"
	"time"

	OS "github.com/hayrullahcansu/zetamail/cross"
)

func ReadDataNoTLS(conn net.Conn) (string, error) {
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

func ReadAllLineNoTLS(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(ReadDeadLine))
	var builder strings.Builder
	reader := bufio.NewReader(conn)
	var dataLine string

	readLine, _, err := reader.ReadLine()
	if err != nil {
		return "421 Connection ended abruptly", err
	}
	if len(readLine) == 0 {
		return "421 Connection ended abruptly", nil
	}

	for readLine[3] != '-' {
		dataLine = string(readLine)
		builder.WriteString(dataLine)
		readLine, _, err = reader.ReadLine()
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

func ReadAllNoTLS(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(ReadDeadLine))
	reader := bufio.NewReader(conn)

	readLine, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(readLine), nil
}

func WriteLineNoTLS(conn net.Conn, data string) error {
	conn.SetWriteDeadline(time.Now().Add(WriteDeadLine))
	_, err := conn.Write([]byte(data + OS.NewLine))
	return err
}
