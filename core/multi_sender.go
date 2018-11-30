package core

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	OS "../cross"
	"../entity"
	"../logger"
)

const (
	messageLimit       = 10
	workerConsumeLimit = 2
)

type domainMessageStack struct {
	Domain          string
	MessageStack    chan *entity.Message
	m               sync.Mutex
	consumerCounter int
}

type worker struct {
	parent      *MultipleSender
	messages    []*entity.Message
	ttl         time.Duration
	timeout     chan bool
	send        chan bool
	stop        chan bool
	cursorIndex int
}

func (d *domainMessageStack) handle(lastTry bool) (bool, chan *entity.Message) {
	d.m.Lock()
	defer d.m.Unlock()
	if d.isHandlableQuery(lastTry) {
		d.consumerCounter++
		return true, d.MessageStack
	}
	return false, nil
}
func (d *domainMessageStack) isHandlable(lastTry bool) bool {
	d.m.Lock()
	defer d.m.Unlock()
	return d.isHandlableQuery(lastTry)
}
func (d *domainMessageStack) isHandlableQuery(lastTry bool) bool {
	if d.consumerCounter < workerConsumeLimit && len(d.MessageStack) > 0 &&
		(lastTry || (!lastTry && len(d.MessageStack) > messageLimit)) {
		return true
	}
	return false
}

func (d *domainMessageStack) release() bool {
	d.m.Lock()
	defer d.m.Unlock()
	if d.consumerCounter > 0 {
		d.consumerCounter--
		return true
	}
	return false
}

func newWorker(parent *MultipleSender) *worker {
	return &worker{
		parent:      parent,
		ttl:         time.Second * 1,
		timeout:     make(chan bool),
		send:        make(chan bool),
		stop:        make(chan bool),
		cursorIndex: 0,
	}
}

func NewDomainMessageStack(domain string) *domainMessageStack {
	return &domainMessageStack{
		Domain:       domain,
		MessageStack: make(chan *entity.Message, 100),
	}
}

type MultipleSender struct {
	domainMessageStacks map[string]*domainMessageStack
	m                   sync.Mutex
	pool                []*worker
}

var bulkSender *MultipleSender
var onceBulkSender sync.Once

func InstanceBulkSender() *MultipleSender {
	onceBulkSender.Do(func() {
		workerLimit := runtime.NumCPU()
		bulkSender = &MultipleSender{
			domainMessageStacks: make(map[string]*domainMessageStack),
			pool:                make([]*worker, workerLimit),
		}
		for index := 0; index < workerLimit; index++ {
			bulkSender.pool[index] = newWorker(bulkSender)
		}
	})
	return bulkSender
}

func (b *MultipleSender) AppendMessage(host string, message *entity.Message) {
	b.m.Lock()
	defer b.m.Unlock()
	channel, ok := b.domainMessageStacks[host]
	if !ok {
		channel = NewDomainMessageStack(host)
		b.domainMessageStacks[host] = channel
	}
	channel.MessageStack <- message
}

func (b *MultipleSender) Run() {
	b.m.Lock()
	defer func() {
		fmt.Println("unlocked")
		b.m.Unlock()
	}()
	for index := 0; index < len(b.pool); index++ {
		go b.pool[index].run()
	}
}

func (b *MultipleSender) Stop() {
	b.m.Lock()
	defer b.m.Unlock()
	for index := 0; index < len(b.pool); index++ {
		b.pool[index].stop <- true
	}
}

func (b *MultipleSender) getDomainMessageStack() (bool, chan *entity.Message) {
	for _, stack := range b.domainMessageStacks {
		if stack.isHandlable(false) {
			return stack.handle(false)
		}
	}
loop:
	for _, stack := range b.domainMessageStacks {
		if stack.isHandlable(true) {
			return stack.handle(true)
		}
	}
	time.Sleep(time.Second * 5)
	goto loop
}
func (w *worker) run() {
	go func() {
		for {
			ok, channel := w.parent.getDomainMessageStack()
			if ok {
				for len(w.messages) < messageLimit {
					msg := <-channel
					w.messages = append(w.messages, msg)
				}
				w.send <- true
			}
		}
	}()
	w.setTtl()
	for {
		select {
		case <-w.timeout:
			w.send <- true
		case <-w.send:
			w.sendAllMessage()
		case <-w.stop:
			w.sendAllMessage()
		}
	}
}

func (w *worker) setTtl() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			w.timeout <- true
		}

	}()
	// for w.{
	// }

}

func (w *worker) sendAllMessage() {
	if len(w.messages) > 0 {
		logger.Info.Printf("Sended all message %d%s", len(w.messages), OS.NewLine)
	} else {
		logger.Info.Printf("Worker message array empty%s", len(w.messages), OS.NewLine)

	}
}
