package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type Client struct {
	s3api    s3iface.S3API
	uploader s3manageriface.UploaderAPI
}

func New() *Client {
	awscfg := &aws.Config{Region: aws.String("us-east-1")}
	sess := session.Must(session.NewSession(awscfg))

	return &Client{
		s3api:    s3.New(sess),
		uploader: s3manager.NewUploader(sess),
	}
}

func (cl *Client) Upload(bucket, key, contentType string, body io.Reader) error {
	_, err := cl.uploader.Upload(&s3manager.UploadInput{
		Bucket:      &bucket,
		Key:         &key,
		ContentType: &contentType,
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (cl *Client) Download(bucket, key string) (io.Reader, error) {
	resp, err := cl.s3api.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf, nil
}
