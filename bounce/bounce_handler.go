package bounce

import (
	"math"
	"sync"
	"time"

	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/queue/priority"

	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/queue"
	"github.com/hayrullahcansu/fastmta-core/transaction"
)

type BounceHandler struct {
}

var instanceHandler *BounceHandler
var once sync.Once

func Handler() *BounceHandler {
	once.Do(func() {
		instanceHandler = newHandler()
	})
	return instanceHandler
}

func newHandler() *BounceHandler {
	handler := &BounceHandler{}
	return handler
}

// HandleSuccessToSend handles successful delivery.
// Logs success
// Deletes queued data
func (h *BounceHandler) HandleSuccessToSend(message *entity.Message, result *transaction.TransactionGroupResult) {
	logger.Infof("Sent message to %s", message.RcptTo)
}

// HandleFailedToSend handles failure of delivery.
// Logs failure
// Deletes queued data
func (h *BounceHandler) HandleFailedToSend(message *entity.Message, result *transaction.TransactionGroupResult) {
	logger.Errorf("Failed to send messsage to %s", message.RcptTo)
}

// HandleFailedToConnect handles failure of connection temprorely.
// Logs failure of the connection
// Defers the message
func (h *BounceHandler) HandleFailedToConnect(message *entity.Message, result *transaction.TransactionGroupResult) {
	logger.Errorf("Failed to connect messsage to %s", message.RcptTo)
	//TODO: Check if there was no MX record in DNS, so using A, we should fail and not retry.
	h.HandleDeferralToSend(message, result, 15)
}

// HandleDeferralToSend handles message deferal.
// Logs deferral
// Sets the next rety date time
func (h *BounceHandler) HandleDeferralToSend(message *entity.Message, result *transaction.TransactionGroupResult, nextRetryIntervalMinutes int) {
	logger.Warningf("Deliver deferral to send messsage to %s", message.RcptTo)
	message.DeferredCount++
	var nextRetryInterval int = constant.DefaultRetryInterval
	if nextRetryIntervalMinutes > 0 {
		nextRetryInterval = nextRetryIntervalMinutes
	} else {
		// Increase the deferred wait interval by doubling for each retry.
		nextRetryInterval = int(math.Pow(2, float64(message.DeferredCount)) * float64(nextRetryInterval))
		// If we reach over the max interval then set to the max interval value.
		if nextRetryInterval > constant.MaxRetryInterval {
			nextRetryInterval = constant.MaxRetryInterval
		}
	}
	message.AttemptSendTime = time.Now().Add(time.Duration(nextRetryInterval) * time.Minute)
	queue.Instance().EnqueueOutboundNormal(message)

}

// HandleThrottleToSend handles message throttle.
// Logs throttle
// Sets the next rety date time
func (h *BounceHandler) HandleThrottleToSend(message *entity.Message, result *transaction.TransactionGroupResult) {
	logger.Errorf("Delivery throttle to send messsage to %s", message.RcptTo)
	message.Priority = priority.LOW
	message.AttemptSendTime = time.Now().Add(time.Minute * 1)
	queue.Instance().EnqueueOutboundNormal(message)
}

// HandleEnqueueToSend handles maximum connection limit.
// It enqueues the message immediately.
func (h *BounceHandler) HandleEnqueueToSend(message *entity.Message, result *transaction.TransactionGroupResult) {
	logger.Errorf("Enqueue to send messsage to %s", message.RcptTo)
}

// HandleUnavailableToSend handles a service unavailable event, should be same as defer but only wait 1 minute before next retry.
func (h *BounceHandler) HandleUnavailableToSend(message *entity.Message, result *transaction.TransactionGroupResult) {
	logger.Errorf("Service unavailable to send messsage to %s", message.RcptTo)
	message.AttemptSendTime = time.Now().Add(time.Minute * 1)
	queue.Instance().EnqueueOutboundNormal(message)

}

// HandleTemporaryToSend handles something weird happening with this message, get it out of the way for a bit.
// Wait 5 minutes before next retry.
func (h *BounceHandler) HandleTemporaryToSend(message *entity.Message, result *transaction.TransactionGroupResult) {
	message.AttemptSendTime = time.Now().Add(time.Minute * 5)
	queue.Instance().EnqueueOutboundNormal(message)
}
