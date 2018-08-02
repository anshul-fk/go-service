package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	sent_string := "dummy string"
	if path == "send_message" {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := sqs.New(sess)

		// URL to our queue
		qURL := "https://sqs.us-east-1.amazonaws.com/723008196684/dev-integrations-worker"
		result, err := svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(10),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"Title": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("The Whistler"),
				},
				"Author": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("John Grisham"),
				},
				"WeeksOn": &sqs.MessageAttributeValue{
					DataType:    aws.String("Number"),
					StringValue: aws.String("6"),
				},
			},
			MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
			QueueUrl:    &qURL,
		})

		if err != nil {
			fmt.Println("Error", err)
			return
		}
		sent_string = *result.MessageId
		fmt.Println("Success", *result.MessageId)
	}

	fmt.Fprintf(w, "Hi there, I love %s! %s", r.URL.Path[1:], sent_string)
}
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
