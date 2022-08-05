package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
)

func NewS3Bucket(scope constructs.Construct, id string, bucket string) s3.S3Bucket {
	return s3.NewS3Bucket(scope, &id, &s3.S3BucketConfig{Bucket: &bucket})
}
