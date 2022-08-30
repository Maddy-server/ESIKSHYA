package db

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

func (store *SQLStore) UploadToS3(fileName string, data []byte) error {
	buff := bytes.NewBuffer(data)
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(store.awsSession)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(store.Config.AWS.BucketName),
		Key:                aws.String(fileName),
		Body:               buff,
		ACL:                aws.String("public-read"),
		ContentType:        aws.String(http.DetectContentType(data)),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	_ = result
	return err
}

func (store *SQLStore) GetBookImageUrl(bookId int32) string {
	// Create S3 service client
	svc := s3.New(store.awsSession)
	var url1 string

	tr, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(store.Config.AWS.BucketName),
		Key:    aws.String(fmt.Sprintf("books/%d", bookId)),
	})
	url1, err := tr.Presign(time.Hour * 48)
	if err != nil {
		logrus.Error(err)
	}

	return url1
}

func (store *SQLStore) GetBookPdfUrl(bookId int32) string {
	// Create S3 service client
	svc := s3.New(store.awsSession)
	var url1 string

	tr, _ := svc.GetObjectRequest(&s3.GetObjectInput{

		Bucket: aws.String(store.Config.AWS.BucketName),
		Key:    aws.String(fmt.Sprintf("pdf/%d", bookId)),
	})
	url1, err := tr.Presign(time.Hour * 48)
	if err != nil {
		logrus.Error(err)
	}

	return url1
}
