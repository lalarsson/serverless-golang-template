package main

import (
	"context"

	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"os"
	"time"
)

// Station representation of the sensor table
type Station struct {
	Datettime time.Time `json:"station_datetime"`
	StationId string    `json:"station_id"`
}

// Client definition to gather all persistant sessions
type Client struct {
	Session      *session.Session
	SnsClient    *sns.SNS
	DynamoClient *dynamodb.DynamoDB
	S3Client     *s3.S3
}

/** Initializers **/

// SetupClient intialize basic session
func SetupClient() *Client {
	sess := session.Must(session.NewSession())
	return &Client{Session: sess}
}

// SetupSnsClient SNS configuration
func (c *Client) SetupSnsClient() *Client {
	c.SnsClient = sns.New(c.Session)
	return c
}

// SetupDynamoDB Dynamodb configuration
func (c *Client) SetupDynamoDB() *Client {
	c.DynamoClient = dynamodb.New(c.Session)
	return c

}

// SetupS3Client S3 configuration
func (c *Client) SetupS3Client() *Client {
	c.S3Client = s3.New(c.Session)
	return c
}

// Publish Publish to SNS //
func (c *Client) Publish(params *sns.PublishInput) error {
	req, _ := c.SnsClient.PublishRequest(params)
	return req.Send()
}

// GetStationsOnDynamoDB Get all Stations with datetime older than five minutes ago //
func (c *Client) GetStationsOnDynamoDB() ([]Station, error) {
	var items []Station
	borderTime := time.Now().UTC().Add(time.Duration(-5) * time.Minute)

	filt := expression.Name("station_datetime").LessThanEqual(expression.Value(borderTime))

	proj := expression.NamesList(expression.Name("station_id"), expression.Name("station_datetime"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("sensors-dev"),
	}

	// Make the DynamoDB Query API call
	result, err := c.DynamoClient.Scan(params)
	for _, i := range result.Items {
		item := Station{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, nil

}

// Handler Function to be executed
func Handler(ctx context.Context) error {
	client := SetupClient().SetupDynamoDB()
	topicArn := os.Getenv("SNS_ARN")

	stations, err := client.GetStationsOnDynamoDB()
	if err != nil {
		return err
	}
	if len(stations) == 0 {
		return nil
	}
	client.SetupSnsClient()
	for _, elem := range stations {
		message := fmt.Sprintf("Basestation with id %s has not been seen since %s", elem.StationId, elem.Datettime.Format(time.RFC3339))
		param := &sns.PublishInput{
			Message:  aws.String(message),
			TopicArn: aws.String(topicArn),
		}
		client.Publish(param)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
