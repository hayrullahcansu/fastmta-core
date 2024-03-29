package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hayrullahcansu/fastmta-core/global"

	"github.com/hayrullahcansu/fastmta-core/broker"
	OS "github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/logger"
)

var ops int64

const (
	Timeout   = time.Second * time.Duration(30)
	KeepAlive = time.Second * time.Duration(30)
	//MtaName       = "ZetaMail"
	//MaxErrorLimit = 10
)

func main() {
	// err := smtp.SendMail("smtp.gmail.com:587",
	// 	smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
	// 	from, []string{to}, []byte(msg))
	// host = "gmail-smtp-in.l.google.COM:25"
	// dialer := &net.Dialer{
	// 	Timeout:   Timeout,
	// 	KeepAlive: KeepAlive,
	// 	LocalAddr: &net.TCPAddr{
	// 		IP: "127.0.0.1",
	// 	},
	// }
	// conn, err := dialer.Dial("tcp", host)
	// if err != nil {
	// 	if opError, ok := err.(*net.OpError); ok {
	// 		if dnsError, ok := opError.Err.(*net.DNSError); ok {
	// 			return false, transaction.HostNotFound, dnsError.Error()
	// 		}
	// 	}
	// 	//TODO: define all error like dnsError
	// 	return false, transaction.ServiceNotAvailable, "service not avaliable"
	// }

	// load command line arguments
	start := time.Now()
	//m := &sync.Mutex{}

	//runtime.Goexit()
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	name := flag.String("name", "FastMta", "name to print")
	flag.Parse()
	log.Printf("Starting service for %s%s", *name, OS.NewLine)
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	//signal.Notify(sigs,syscall.SIGQUIT)
	// method invoked upon seeing signal

	go func() {
		s := <-sigs
		logger.Infof("RECEIVED SIGNAL: %s%s", s, OS.NewLine)
		AppCleanup()
		os.Exit(1)
	}()

	broker.InitSystem()
	broker := broker.New(
		global.StaticConfig.ID,
		global.StaticConfig.Name,
		global.StaticConfig.IsEnabled,
	)
	// rabbitClient := queue.New()
	// rabbitClient.Connect(true)
	// _, _ = rabbitClient.Consume(queue.InboundQueueName, "", false, false, true, nil)

	// rabbitClient2 := queue.New()
	// rabbitClient2.Connect(true)
	// _, _ = rabbitClient2.Consume(queue.InboundQueueName, "", false, false, true, nil)

	broker.Run()

	select {}

}
func AppCleanup() {
	time.Sleep(time.Millisecond * time.Duration(1000))
	logger.Infof("CLEANUP APP BEFORE EXIT!!!")
}
