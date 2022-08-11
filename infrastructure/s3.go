package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
	"github.com/koki-develop/gogogo/infrastructure/pkg/util"
)

// cats.json 保存用の S3 バケット
func NewS3Cats(scope constructs.Construct) s3.S3Bucket {
	bucket := util.NewS3Bucket(scope, "s3-bucket-cats", "gogogo-cats")

	util.NewS3PublicAccessBlock(scope, "s3-public-access-block-cats", bucket.Bucket())

	return bucket
}
