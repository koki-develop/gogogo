package util

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
)

func NewCloudfrontOriginAccessIdentity(scope constructs.Construct, id string) cloudfront.CloudfrontOriginAccessIdentity {
	return cloudfront.NewCloudfrontOriginAccessIdentity(scope, &id, &cloudfront.CloudfrontOriginAccessIdentityConfig{})
}
