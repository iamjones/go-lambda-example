package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Product struct {
	Sku string `json:"sku"`
	Name string `json:"name"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("Request body: ", request.Body)

	var product Product
	err := json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
		fmt.Println("Error converting JSON to struct: ", err)
	}

	fmt.Println("Product: ", product)

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(session)

	prod, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
		fmt.Println("Error creating DynamoDB map of product: ", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      prod,
		TableName: aws.String("Product"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
		fmt.Println("Error writing product to DynamoDB: ", err)
	}

	return events.APIGatewayProxyResponse{
		Body: request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}