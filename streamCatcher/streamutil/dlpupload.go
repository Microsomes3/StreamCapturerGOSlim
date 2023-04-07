package streamutil

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DLPUploader struct{}

const (
	accountid = "6d8f181feef622b528f2fc75fbce8754"
)

var bucket = aws.String("maeplet")
var ACCESSKEY = os.Getenv("R2_ACCESS_KEY")
var SECRETKEY = os.Getenv("R2_SECRET_KEY")

func (d *DLPUploader) UploadFile(file *os.File, key string) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	fmt.Println("Uploading file: ", stat.Name())
	fmt.Println("Uploading file: ", stat.Size()/1024/1024, "MB")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountid),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ACCESSKEY, SECRETKEY, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      bucket,
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("video/mp4"),
	})

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Println("Uploaded: ", "https://pub-cf9c58b47aaa413eadbc9d4fba77649a.r2.dev/"+key)

	return nil
}
