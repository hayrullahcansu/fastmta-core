package relay

import (
	"encoding/json"
	"sync"

	"github.com/streadway/amqp"

	"github.com/hayrullahcansu/fastmta-core/bounce"
	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/smtp/outbound"
	"github.com/hayrullahcansu/fastmta-core/transaction"
	"github.com/hayrullahcansu/fastmta-core/transaction/status"

	"github.com/hayrullahcansu/fastmta-core/mta"
)

var instanceManager *Manager
var once sync.Once

// Manager manages to outbound messages
type Manager struct {
}

// InstanceManager return new or existing instance of Manager
func InstanceManager() *Manager {
	once.Do(func() {
		instanceManager = newManager()
	})
	return instanceManager
}

func newManager() *Manager {
	instance := &Manager{}

	return instance
}

// SendMessage sends message to target via agent
// And also it informs the result to bounce manager.
func (r *Manager) SendMessage(outboundMessage *amqp.Delivery) {
	pureMessage := &entity.Message{}
	json.Unmarshal(outboundMessage.Body, pureMessage)
	logger.Infof("Received message From %s", constant.OutboundNormalQueueName)
	// if _, ok := caching.InstanceDomain().C.Get(pureMessage.Host); !ok {
	//exchange.InstanceRouter().
	// }
	db, err := entity.GetDbContext()
	entity.PanicOnError(err)
	db.Model(&pureMessage).Update("status", status.OUT_QUEUE)
	db.Close()

	_mta := mta.InstanceManager().GetVirtualMtaGroup(pureMessage.GroupID).GetNextVirtualMta()
	agent := outbound.NewAgent(_mta)

	// if _, ok := caching.InstanceDomain().C.Get(pureMessage.Host); !ok {
	//exchange.InstanceRouter().
	// }
	_, result := agent.SendMessage(pureMessage)
	switch result.TransactionResult {
	case transaction.Success:
		bounce.Instance().HandleSuccessToSend(pureMessage, result)
	case transaction.HostNotFound:
	case transaction.FailedToConnect:
		bounce.Instance().HandleFailedToSend(pureMessage, result)
	case transaction.RejectedByRemoteServer:
		if result.ResultMessage[0] == '5' {
			bounce.Instance().HandleFailedToSend(pureMessage, result)
		} else {
			bounce.Instance().HandleDeferralToSend(pureMessage, result, -1)
		}
	case transaction.MaxMessages:
		bounce.Instance().HandleThrottleToSend(pureMessage, result)
	case transaction.MaxConnections:
		bounce.Instance().HandleEnqueueToSend(pureMessage, result)
	case transaction.ServiceNotAvailable:
		bounce.Instance().HandleUnavailableToSend(pureMessage, result)
	default:
		// Something weird happening with this message, get it out of the way for a bit.
		bounce.Instance().HandleTemporaryToSend(pureMessage, result)
	}
	outboundMessage.Ack(false)
}
