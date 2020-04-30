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

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: qURL.QueueUrl,
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      qURL.QueueUrl,
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})

	if err != nil {
		fmt.Println("Delete Error", err)
		return
	}

	fmt.Println("Message Deleted", resultDelete)
}
