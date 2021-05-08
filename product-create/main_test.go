package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {

	t.Run("Successful request", func(t *testing.T) {
		requestBody := "{\"sku\": \"1234567\", \"name\": \"Hats\"}"
		response, err := handler(events.APIGatewayProxyRequest{Body: requestBody})

		if err != nil {
			t.Fatal("Error not expected")
		}

		if response.StatusCode != 200 {
			t.Fatalf("Expected 200 status code '%d' returned", response.StatusCode)
		}

		if response.Body != requestBody {
			t.Fatalf("Expect response body %s but received %s", response.Body, requestBody)
		}
	})

	t.Run("JSON format error in request body", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{Body: "{sku\": \"1234567\", \"name\": \"Hats\"}"})

		if err == nil {
			t.Fatal("Error is expected but not returned")
		}

		if response.StatusCode != 500 {
			t.Fatalf("Expected 500 status code '%d' returned", response.StatusCode)
		}
	})
}
