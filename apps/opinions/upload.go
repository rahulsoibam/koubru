package opinions

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (a *App) S3UploadOpinion(file io.Reader, filename string) (string, error) {
	output, err := a.Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String("video/mp4"),
	})

	if err != nil {
		return "", err
	}
	return output.Location, nil
}
