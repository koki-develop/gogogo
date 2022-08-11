package frontend

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
	"github.com/koki-develop/gogogo/infrastructure/pkg/util"
)

type s3MainInput struct {
	OriginAccessIdentity cloudfront.CloudfrontOriginAccessIdentity
}

func newS3Main(scope constructs.Construct, ipt *s3MainInput) s3.S3Bucket {
	bucket := util.NewS3Bucket(scope, "s3-bucket-frontend", "gogogo-frontend-files")

	util.NewS3PublicAccessBlock(scope, "s3-public-access-block-frontend", bucket.Bucket())

	frontendbucketpolicy := iam.NewDataAwsIamPolicyDocument(scope, jsii.String("data-iam-policy-document-frontend-bucket-policy"), &iam.DataAwsIamPolicyDocumentConfig{
		Statement: []*iam.DataAwsIamPolicyDocumentStatement{{
			Effect:    jsii.String("Allow"),
			Actions:   jsii.Strings("s3:GetObject"),
			Resources: jsii.Strings(fmt.Sprintf("%s/*", *bucket.Arn())),
			Principals: []*iam.DataAwsIamPolicyDocumentStatementPrincipals{{
				Type:        jsii.String("AWS"),
				Identifiers: &[]*string{ipt.OriginAccessIdentity.IamArn()},
			}},
		}},
	})
	s3.NewS3BucketPolicy(scope, jsii.String("s3-bucket-policy-frontend"), &s3.S3BucketPolicyConfig{
		Bucket: bucket.Id(),
		Policy: frontendbucketpolicy.Json(),
	})

	return bucket
}
