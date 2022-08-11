package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/acm"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
)

func NewCloudfrontOriginAccessIdentity(scope constructs.Construct, id string) cloudfront.CloudfrontOriginAccessIdentity {
	return cloudfront.NewCloudfrontOriginAccessIdentity(scope, &id, &cloudfront.CloudfrontOriginAccessIdentityConfig{})
}

type CloudfrontFrontendConfig struct {
	Domain               string
	Bucket               s3.S3Bucket
	OriginAccessIdentity cloudfront.CloudfrontOriginAccessIdentity
	Certificate          acm.AcmCertificate
}

func NewCloudfrontFrontend(scope constructs.Construct, cfg *CloudfrontFrontendConfig) cloudfront.CloudfrontDistribution {
	return cloudfront.NewCloudfrontDistribution(scope, jsii.String("cloudfront-distribution-frontend"), &cloudfront.CloudfrontDistributionConfig{
		Aliases:           jsii.Strings(cfg.Domain),
		Enabled:           jsii.Bool(true),
		DefaultRootObject: jsii.String("index.html"),
		Origin: []*cloudfront.CloudfrontDistributionOrigin{{
			OriginId:   cfg.Bucket.Id(),
			DomainName: cfg.Bucket.BucketRegionalDomainName(),
			S3OriginConfig: &cloudfront.CloudfrontDistributionOriginS3OriginConfig{
				OriginAccessIdentity: cfg.OriginAccessIdentity.CloudfrontAccessIdentityPath(),
			},
		}},
		DefaultCacheBehavior: &cloudfront.CloudfrontDistributionDefaultCacheBehavior{
			TargetOriginId:       cfg.Bucket.Id(),
			AllowedMethods:       jsii.Strings("GET", "HEAD"),
			CachedMethods:        jsii.Strings("GET", "HEAD"),
			ViewerProtocolPolicy: jsii.String("redirect-to-https"),
			Compress:             jsii.Bool(true),
			MinTtl:               jsii.Number(0),
			DefaultTtl:           jsii.Number(3600),
			MaxTtl:               jsii.Number(86400),
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
			AcmCertificateArn:            cfg.Certificate.Arn(),
			CloudfrontDefaultCertificate: jsii.Bool(false),
			MinimumProtocolVersion:       jsii.String("TLSv1.2_2021"),
			SslSupportMethod:             jsii.String("sni-only"),
		},
	})
}
