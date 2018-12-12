package main

import (
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	snsRecord := events.SNSEntity{
		Message:   "test event",
		Timestamp: time.Now(),
	}
	ev := events.SNSEvent{
		Records: []events.SNSEventRecord{
			events.SNSEventRecord{
				EventSource: "test",
				SNS:         snsRecord,
			},
		},
	}

	err := Handler(nil, ev)
	assert.IsType(t, nil, err)
}
