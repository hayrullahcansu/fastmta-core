package core

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	OS "../cross"
	"../entity"
	"../queue"
	dkim "github.com/emersion/go-dkim"
)

const testPrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIICXwIBAAKBgQDwIRP/UC3SBsEmGqZ9ZJW3/DkMoGeLnQg1fWn7/zYtIxN2SnFC
jxOCKG9v3b4jYfcTNh5ijSsq631uBItLa7od+v/RtdC2UzJ1lWT947qR+Rcac2gb
to/NMqJ0fzfVjH4OuKhitdY9tf6mcwGjaNBcWToIMmPSPDdQPNUYckcQ2QIDAQAB
AoGBALmn+XwWk7akvkUlqb+dOxyLB9i5VBVfje89Teolwc9YJT36BGN/l4e0l6QX
/1//6DWUTB3KI6wFcm7TWJcxbS0tcKZX7FsJvUz1SbQnkS54DJck1EZO/BLa5ckJ
gAYIaqlA9C0ZwM6i58lLlPadX/rtHb7pWzeNcZHjKrjM461ZAkEA+itss2nRlmyO
n1/5yDyCluST4dQfO8kAB3toSEVc7DeFeDhnC1mZdjASZNvdHS4gbLIA1hUGEF9m
3hKsGUMMPwJBAPW5v/U+AWTADFCS22t72NUurgzeAbzb1HWMqO4y4+9Hpjk5wvL/
eVYizyuce3/fGke7aRYw/ADKygMJdW8H/OcCQQDz5OQb4j2QDpPZc0Nc4QlbvMsj
7p7otWRO5xRa6SzXqqV3+F0VpqvDmshEBkoCydaYwc2o6WQ5EBmExeV8124XAkEA
qZzGsIxVP+sEVRWZmW6KNFSdVUpk3qzK0Tz/WjQMe5z0UunY9Ax9/4PVhp/j61bf
eAYXunajbBSOLlx4D+TunwJBANkPI5S9iylsbLs6NkaMHV6k5ioHBBmgCak95JGX
GMot/L2x0IYyMLAz6oLWh2hm7zwtb0CgOrPo1ke44hFYnfc=
-----END RSA PRIVATE KEY-----
`

var (
	testPrivateKey *rsa.PrivateKey
)

func init() {
	block, _ := pem.Decode([]byte(testPrivateKeyPEM))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	testPrivateKey = key
}

type InboundStagingConsumer struct {
	RabbitMqClient *queue.RabbitMqClient
}

func NewInboundStagingConsumer() *InboundStagingConsumer {
	return &InboundStagingConsumer{
		RabbitMqClient: queue.New(),
	}
}

func (consumer *InboundStagingConsumer) Run() {
	consumer.RabbitMqClient.Connect(true)
	messageChannel, err := consumer.RabbitMqClient.Consume(queue.InboundStagingQueueName, "", false, false, true, nil)
	if err != nil {
		panic(fmt.Sprintf("error handled in %s queue: %s%s", queue.InboundStagingQueueName, err, OS.NewLine))
	}

	for {
		select {
		case messageDelivery, ok := <-messageChannel:
			if ok {
				msg := &entity.Message{}
				options := &dkim.SignOptions{
					Domain:   "example.org",
					Selector: "brisbane",
					Signer:   testPrivateKey,
				}
				json.Unmarshal(messageDelivery.Body, msg)
				var b bytes.Buffer
				r := strings.NewReader(msg.Data)
				if err := dkim.Sign(&b, r, options); err != nil {
					//TODO: fix or report dkim error
				}
				msg.Data = string(b.Bytes())
				data, err := json.Marshal(msg)
				err = consumer.RabbitMqClient.Publish(
					queue.OutboundExchange,
					queue.RoutingKeyOutbound,
					false,
					false,
					data,
				)

				if err == nil {
					messageDelivery.Ack(true)
				} else {
					messageDelivery.Reject(true)
				}
			}
		}
	}
}