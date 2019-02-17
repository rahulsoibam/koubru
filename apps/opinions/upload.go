package opinions

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/rahulsoibam/koubru/errs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type SQSResponse struct {
	Event  SQSEvent   `json:"event"`
	Error  *SQSError  `json:"error"`
	Output *SQSOutput `json:"output"`
}

type SQSError struct {
	Error string `json:"Error"`
	Cause string `json:"Cause"`
}

type SQSEvent struct {
	TargetBucket string `json:"target_bucket"`
	SourceBucket string `json:"source_bucket"`
	Key          string `json:"key"`
}

type SQSOutput struct {
	Source    string `json:"source"`
	Hls       string `json:"hls"`
	Thumbnail string `json:"thumbnail"`
}

func (a *App) S3UploadOpinion(file io.Reader, filename string) error {
	_, err := a.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		return err
	}
	return nil
}

func (a *App) PollSQSAndGetLinks(svc sqsiface.SQSAPI, bucket, filename, resultQueueURL string) (*SQSOutput, error) {
	log.Println("Client: ", svc)
	log.Println("Bucket: ", bucket)
	log.Println("Filename: ", filename)
	log.Println("Result Queue URL: ", resultQueueURL)

	for {
		resp, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:          aws.String(resultQueueURL),
			VisibilityTimeout: aws.Int64(0),
			WaitTimeSeconds:   aws.Int64(20),
		})
		if err != nil {
			log.Println("Failed to receive message", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, msg := range resp.Messages {
			result := &SQSResponse{}
			if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), result); err != nil {
				log.Println("Failed to unmarshal message", err)
				continue
			}
			if result.Event.SourceBucket == bucket && result.Event.Key == filename {
				log.Println(msg)
				if result.Error != nil {
					log.Println(result.Error)
					return nil, errs.OpinionBadPayload
				}

				svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(resultQueueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})
				return result.Output, nil
			}
		}
	}
}
