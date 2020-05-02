package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var cli struct {
	Debug    bool
	Endpoint string `default:"http://127.0.0.1:8080"`
	Profile  string `default:"default"`
	Queue    string `default:"tmq"`
	Region   string
}

func main() {
	_ = kong.Parse(&cli)

	cfg := aws.NewConfig().WithCredentials(
		credentials.NewSharedCredentials("", cli.Profile),
	)
	if cli.Debug {
		cfg = cfg.WithLogLevel(aws.LogDebugWithHTTPBody)
	}
	if cli.Endpoint != "" {
		cfg = cfg.WithEndpoint(cli.Endpoint)
	}
	if cli.Region != "" {
		cfg = cfg.WithRegion(cli.Region)
	}

	sess := session.Must(session.NewSession())
	svc := sqs.New(sess, cfg)

	qURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(cli.Queue),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:     qURL.QueueUrl,
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"speaker": {
				DataType:    aws.String("String"),
				StringValue: aws.String("The Tick"),
			},
		},
		MessageBody: aws.String("Spoon!"),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println(*result.MessageId)
}
