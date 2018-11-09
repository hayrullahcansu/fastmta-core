package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"sync/atomic"

	"./boss"
	"./initialize"
)

var ops int64

func main() {
	// load command line arguments
	atomic.AddInt64(&ops, 1)
	atomic.AddInt64(&ops, -1)
	fmt.Println(atomic.LoadInt64(&ops))
	name := flag.String("name", "world", "name to print")
	flag.Parse()
	log.Printf("Starting sleepservice for %s", *name)
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	//signal.Notify(sigs,syscall.SIGQUIT)
	// method invoked upon seeing signal

	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		AppCleanup()
		os.Exit(1)
	}()
	conf := initialize.Run()
	boss := boss.New(conf)
	boss.Run()
	// infinite print loop
	select {}

}
func AppCleanup() {
	time.Sleep(time.Millisecond * time.Duration(3000))
	log.Println("CLEANUP APP BEFORE EXIT!!!")
}
