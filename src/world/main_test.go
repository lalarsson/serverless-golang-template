package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	movies, err := Handler(events.APIGatewayProxyRequest{})
	assert.IsType(t, nil, err)
	assert.NotEqual(t, 0, len(movies.Body))
}
