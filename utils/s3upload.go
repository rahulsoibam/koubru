package utils

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3PictureDetails struct {
	FileName        string
	FileType        string
	FileContentType string
}

func S3UploadProfilePicture(uploader *s3manager.Uploader, file io.Reader, details S3PictureDetails) (string, error) {
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("PROFILE_PICTURE_BUCKET")),
		Key:         aws.String(details.FileName + details.FileType),
		Body:        file,
		ContentType: aws.String(details.FileContentType),
	})

	if err != nil {
		return "", err
	}
	return output.Location, nil
}
