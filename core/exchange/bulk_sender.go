package exchange

import (
	"runtime"
	"sync"
	"time"

	"../../entity"
)

type DomainMessageStack struct {
	Domain       string
	MessageStack chan *entity.Message
}

type worker struct {
	messages []*entity.Message
	ttl      time.Duration
	timeout  chan bool
}

func NewDomainMessageStack(domain string) *DomainMessageStack {
	return &DomainMessageStack{
		Domain:       domain,
		MessageStack: make(chan *entity.Message, 100),
	}
}

type BulkSender struct {
	DomainMessageStacks map[string]*DomainMessageStack
	m                   sync.Mutex
	pool                []*worker
}

var bulkSender *BulkSender
var onceBulkSender sync.Once

func InstanceBulkSender() *BulkSender {
	onceBulkSender.Do(func() {
		bulkSender = &BulkSender{
			DomainMessageStacks: make(map[string]*DomainMessageStack),
			pool:                make([]*worker, runtime.NumCPU()),
		}
	})
	return bulkSender
}

func (b *BulkSender) AppendMessage(host string, message *entity.Message) {
	b.m.Lock()
	defer b.m.Unlock()
	channel, ok := b.DomainMessageStacks[host]
	if !ok {
		channel = NewDomainMessageStack(host)
		b.DomainMessageStacks[host] = channel
	}
	channel.MessageStack <- message
}

func (b *BulkSender) Run() {
	b.m.Lock()
	defer b.m.Unlock()
	for index := 0; index < len(b.pool); index++ {
		go b.pool[index].run()
	}
}

func (w *worker) run() {
	for {
		select {
		case <-w.timeout:
			w.send()
		}
	}
}

func (w *worker) setTtl() {
	go func() {
		time.Sleep(1 * time.Second)
		w.timeout <- true
	}()
}

func (w *worker) send() {

}
