package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	endpoint := os.Getenv("TMQ_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8080"
	}

	creds := credentials.NewStaticCredentials("id", "secret", "token")
	sess := session.Must(session.NewSession())

	svc := sqs.New(sess,
		aws.NewConfig().
			WithCredentials(creds).
			WithRegion("us-west-2").
			WithEndpoint(endpoint),
	)

	qURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("example"),
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
