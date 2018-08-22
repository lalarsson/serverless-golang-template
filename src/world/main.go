package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("hello")
	return events.APIGatewayProxyResponse{Body: "hello world", StatusCode: 200}, nil
}

func main() {
	log.Print("hello")
	lambda.Start(Handler)
}
