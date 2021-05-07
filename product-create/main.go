package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Product struct {
	Sku string `json:"sku"`
	Name string `json:"name"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("Request body: ", request.Body)

	var product Product
	err := json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		fmt.Println("Error converting JSON to struct: ", err)
	}

	fmt.Println("Product: ", product)

	// session := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// svc := dynamodb.New(session)

	// prod, err := dynamodbattribute.MarshalMap()

	return events.APIGatewayProxyResponse{
		Body: request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}