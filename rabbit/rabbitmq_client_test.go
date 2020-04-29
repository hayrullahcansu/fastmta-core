package rabbit

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetOutboundQueue(t *testing.T) {
	tests := map[string]struct {
		input  *entity.Message
		output string
		err    error
	}{
		"successful outbound normal": {
			input: &entity.Message{
				AttemptSendTime: time.Now(),
			},
			output: constant.OutboundNormalQueueName,
			err:    nil,
		},
		"successful outbound 1": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Second * 3),
			},
			output: constant.OutboundWaiting1,
			err:    nil,
		},
		"successful outbound 10 1": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Second * 10),
			},
			output: constant.OutboundWaiting10,
			err:    nil,
		},
		"successful outbound 10 2": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Second * 15),
			},
			output: constant.OutboundWaiting10,
			err:    nil,
		},
		"successful outbound 60": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Second * 60),
			},
			output: constant.OutboundWaiting60,
			err:    nil,
		},
		"successful outbound 60 2": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Second * 61),
			},
			output: constant.OutboundWaiting60,
			err:    nil,
		},
		"successful outbound 300 1": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Second * 300),
			},
			output: constant.OutboundWaiting300,
			err:    nil,
		},
		"successful outbound 300 2": {
			input: &entity.Message{
				AttemptSendTime: time.Now().Add(time.Minute * 100),
			},
			output: constant.OutboundWaiting300,
			err:    nil,
		},
	}

	for head, test := range tests {
		a, b := test.output, getOutboundQueue(test.input)
		if !assert.EqualValues(t, a, b) {
			t.Errorf(head)
		}

		if !reflect.DeepEqual(a, b) {
			t.Errorf(head)
		}
		if diff := deep.Equal(a, b); diff != nil {
			t.Error(diff)
		}
	}
}
