package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nlopes/slack"
)

func Handler(ctx context.Context, snsEvent events.SNSEvent) error {
	token := os.Getenv("SLACK_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL_ID")
	user := os.Getenv("SLACK_USER_NAME")
	api := slack.New(token)

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
		params := slack.PostMessageParameters{
			Username:  user,
			IconEmoji: ":angry:",
		}
		channelID, timestamp, err := api.PostMessage(channel, snsRecord.Message, params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}
		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
