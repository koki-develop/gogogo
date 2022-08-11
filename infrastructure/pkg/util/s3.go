package util

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
)

func NewS3Bucket(scope constructs.Construct, id, name string) s3.S3Bucket {
	return s3.NewS3Bucket(scope, &id, &s3.S3BucketConfig{Bucket: &name})
}

func NewS3PublicAccessBlock(scope constructs.Construct, id string, bucket *string) s3.S3BucketPublicAccessBlock {
	return s3.NewS3BucketPublicAccessBlock(scope, &id, &s3.S3BucketPublicAccessBlockConfig{
		Bucket:                bucket,
		BlockPublicAcls:       jsii.Bool(true),
		BlockPublicPolicy:     jsii.Bool(true),
		IgnorePublicAcls:      jsii.Bool(true),
		RestrictPublicBuckets: jsii.Bool(true),
	})
}
