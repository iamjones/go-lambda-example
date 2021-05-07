package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {

	t.Run("Successful request", func(t *testing.T) {
		_, err := handler(events.APIGatewayProxyRequest{})

		if err != nil {
			t.Fatal("Error not expected")
		}
	})
}