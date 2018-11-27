package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"./boss"
	OS "./cross"
	"./global"
	"./logger"
)

var ops int64

const (
	Timeout   = time.Second * time.Duration(30)
	KeepAlive = time.Second * time.Duration(30)
	//MtaName       = "ZetaMail"
	//MaxErrorLimit = 10
)

func main() {
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
		logger.Info.Printf("RECEIVED SIGNAL: %s%s", s, OS.NewLine)
		AppCleanup()
		os.Exit(1)
	}()
	global.Run()
	boss := boss.New()
	boss.Run()
	// infinite print loop
	select {}

}
func AppCleanup() {
	time.Sleep(time.Millisecond * time.Duration(1000))
	logger.Info.Println("CLEANUP APP BEFORE EXIT!!!")
}
