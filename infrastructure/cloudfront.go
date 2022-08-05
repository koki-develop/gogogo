package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
)

func NewCloudfrontFrontend(scope constructs.Construct, bucket s3.S3Bucket, identity cloudfront.CloudfrontOriginAccessIdentity) cloudfront.CloudfrontDistribution {
	return cloudfront.NewCloudfrontDistribution(scope, jsii.String("cloudfront-distribution-frontend"), &cloudfront.CloudfrontDistributionConfig{
		Enabled:           jsii.Bool(true),
		DefaultRootObject: jsii.String("index.html"),
		Origin: []*cloudfront.CloudfrontDistributionOrigin{{
			OriginId:   bucket.Id(),
			DomainName: bucket.BucketRegionalDomainName(),
			S3OriginConfig: &cloudfront.CloudfrontDistributionOriginS3OriginConfig{
				OriginAccessIdentity: identity.CloudfrontAccessIdentityPath(),
			},
		}},
		DefaultCacheBehavior: &cloudfront.CloudfrontDistributionDefaultCacheBehavior{
			TargetOriginId:       bucket.Id(),
			AllowedMethods:       jsii.Strings("GET", "HEAD"),
			CachedMethods:        jsii.Strings("GET", "HEAD"),
			ViewerProtocolPolicy: jsii.String("redirect-to-https"),
			Compress:             jsii.Bool(true),
			MinTtl:               jsii.Number(0),
			DefaultTtl:           jsii.Number(0), // TODO: 適切に設定する
			MaxTtl:               jsii.Number(0), // TODO: 適切に設定する
			ForwardedValues: &cloudfront.CloudfrontDistributionDefaultCacheBehaviorForwardedValues{
				QueryString: jsii.Bool(false),
				Cookies: &cloudfront.CloudfrontDistributionDefaultCacheBehaviorForwardedValuesCookies{
					Forward: jsii.String("none"),
				},
			},
		},
		Restrictions: &cloudfront.CloudfrontDistributionRestrictions{
			GeoRestriction: &cloudfront.CloudfrontDistributionRestrictionsGeoRestriction{
				RestrictionType: jsii.String("none"),
			},
		},
		ViewerCertificate: &cloudfront.CloudfrontDistributionViewerCertificate{
			CloudfrontDefaultCertificate: jsii.Bool(true),
		},
	})
}
