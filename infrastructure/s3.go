package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
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

type S3FrontendConfig struct {
	OriginAccessIdentity cloudfront.CloudfrontOriginAccessIdentity
}

// フロントエンドの静的ファイル配置用の S3 バケット
func NewS3Frontend(scope constructs.Construct, cfg *S3FrontendConfig) s3.S3Bucket {
	bucket := NewS3Bucket(scope, "s3-bucket-frontend", "gogogo-frontend-files")

	NewS3PublicAccessBlock(scope, "s3-public-access-block-frontend", bucket.Bucket())

	frontendbucketpolicy := iam.NewDataAwsIamPolicyDocument(scope, jsii.String("data-iam-policy-document-frontend-bucket-policy"), &iam.DataAwsIamPolicyDocumentConfig{
		Statement: []*iam.DataAwsIamPolicyDocumentStatement{{
			Effect:    jsii.String("Allow"),
			Actions:   jsii.Strings("s3:GetObject"),
			Resources: jsii.Strings(fmt.Sprintf("%s/*", *bucket.Arn())),
			Principals: []*iam.DataAwsIamPolicyDocumentStatementPrincipals{{
				Type:        jsii.String("AWS"),
				Identifiers: &[]*string{cfg.OriginAccessIdentity.IamArn()},
			}},
		}},
	})
	s3.NewS3BucketPolicy(scope, jsii.String("s3-bucket-policy-frontend"), &s3.S3BucketPolicyConfig{
		Bucket: bucket.Id(),
		Policy: frontendbucketpolicy.Json(),
	})

	return bucket
}

// cats.json 保存用の S3 バケット
func NewS3Cats(scope constructs.Construct) s3.S3Bucket {
	bucket := NewS3Bucket(scope, "s3-bucket-cats", "gogogo-cats")

	NewS3PublicAccessBlock(scope, "s3-public-access-block-cats", bucket.Bucket())

	return bucket
}
