package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Dependencies struct {
	ddb dynamodbiface.DynamoDBAPI
	table string
}

type Product struct {
	Sku string `json:"sku"`
	Name string `json:"name"`
}

func (d *Dependencies) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Request body: ", request.Body)
	log.Println("Environment: ", os.Getenv("AWS_SAM_LOCAL"))

	var product Product
	err := json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		log.Println("Error converting JSON to struct: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	log.Println("Product: ", product)

	prod, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		log.Println("Error creating DynamoDB map of product: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      prod,
		TableName: aws.String(d.table),
	}

	_, err = d.ddb.PutItem(input)
	if err != nil {
		log.Println("Error writing product to DynamoDB: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body: request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	d := Dependencies{
		ddb: createDynamoDbSession(),
		table: "Product",
	}

	lambda.Start(d.handler)
}

func createDynamoDbSession() *dynamodb.DynamoDB {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		log.Println("Configuring DynamoDB for local use")
		return dynamodb.New(session, &aws.Config{Endpoint: aws.String("http://localhost:8000/")})
	} else {
		log.Println("Configuring DynamoDB for prod use")
		return dynamodb.New(session)
	}
}