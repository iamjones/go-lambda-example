package main

import (
	"testing"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.PutItemOutput
	Error error
}

func (d mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &d.Response, d.Error
}

func TestHandler(t *testing.T) {

	t.Run("should return 200 and response when request is successful", func(t *testing.T) {
		requestBody := "{\"sku\": \"1234567\", \"name\": \"Hats\"}"

		m := mockedPutItem{
			Response: dynamodb.PutItemOutput{},
		}

		d := Dependencies{
			ddb: m,
			table: "Product",
		}

		response, err := d.handler(events.APIGatewayProxyRequest{Body: requestBody})

		if err != nil {
			t.Fatalf("Error not expected. Error %s", err)
		}

		if response.StatusCode != 200 {
			t.Fatalf("Expected 200 status code '%d' returned", response.StatusCode)
		}

		if response.Body != requestBody {
			t.Fatalf("Expect response body %s but received %s", response.Body, requestBody)
		}
	})

	t.Run("should return a 500 id the JSON is malformed", func(t *testing.T) {
		d := Dependencies{
			ddb: nil,
			table: "Product",
		}

		response, err := d.handler(events.APIGatewayProxyRequest{Body: "{sku\": \"1234567\", \"name\": \"Hats\"}"})

		if err == nil {
			t.Fatal("Error is expected but not returned")
		}

		if response.StatusCode != 500 {
			t.Fatalf("Expected 500 status code '%d' returned", response.StatusCode)
		}
	})

	t.Run("should return a 502 if there is a downstream error with dynamodb", func(t *testing.T) {
		requestBody := "{\"sku\": \"1234567\", \"name\": \"Hats\"}"

		m := mockedPutItem{
			Error: errors.New("Connection failed"),
		}

		d := Dependencies{
			ddb: m,
			table: "Product",
		}

		response, err := d.handler(events.APIGatewayProxyRequest{Body: requestBody})

		if err == nil {
			t.Fatal("Error is expected but not returned")
		}

		if response.StatusCode != 502 {
			t.Fatalf("Expected 502 status code '%d' returned", response.StatusCode)
		}
	})
}
