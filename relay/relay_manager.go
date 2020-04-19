package relay

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/queue"
	"github.com/hayrullahcansu/fastmta-core/smtp/outbound"
	"github.com/hayrullahcansu/fastmta-core/transaction"

	"github.com/hayrullahcansu/fastmta-core/mta"
)

var instanceManager *RelayManager
var once sync.Once

type RelayManager struct {
}

func InstanceManager() *RelayManager {
	once.Do(func() {
		instanceManager = newManager()
	})
	return instanceManager
}

func newManager() *RelayManager {
	instance := &RelayManager{}

	return instance
}

func (r *RelayManager) SendMessage(outboundMessage *amqp.Delivery) {
	_mta := mta.InstanceManager().GetVirtualMtaGroup(1).GetNextVirtualMta()
	agent := outbound.NewAgent(_mta)

	pureMessage := &entity.Message{}
	json.Unmarshal(outboundMessage.Body, pureMessage)
	logger.Infof("Recieved message From %s", constant.OutboundNormalQueueName)
	// if _, ok := caching.InstanceDomain().C.Get(pureMessage.Host); !ok {
	//exchange.InstanceRouter().
	// }
	_, result := agent.SendMessage(pureMessage)
	switch result.TransactionResult {
	case transaction.Success:
		logger.Infof("Sended message to %s", pureMessage.RcptTo)
	case transaction.HostNotFound:
	case transaction.FailedToConnect:
		logger.Errorf("Failed to Send messsage to %s", pureMessage.RcptTo)
	case transaction.RejectedByRemoteServer:
		if result.ResultMessage[0] == '5' {
			logger.Errorf("Failed to Send messsage to %s", pureMessage.RcptTo)
		} else {
			logger.Warningf("Deliver Deferral to Send messsage to %s", pureMessage.RcptTo)
		}
	case transaction.MaxMessages:
		logger.Errorf("Delivery Throttle to Send messsage to %s", pureMessage.RcptTo)
	case transaction.MaxConnections:
		logger.Errorf("Enqueue to Send messsage to %s", pureMessage.RcptTo)
	case transaction.ServiceNotAvalible:
		logger.Errorf("Service Unavailable to Send messsage to %s", pureMessage.RcptTo)
	default:
		// Something weird happening with this message, get it out of the way for a bit.
		pureMessage.AttemptSendTime = time.Now().Add(time.Minute * 5)
		data, err := json.Marshal(message)
		if err == nil {
			queue.Instance().EnqueueOutboundNormal(data)
		} else {
			outboundMessage.Reject(true)
		}
	}
	outboundMessage.Ack(false)
}
