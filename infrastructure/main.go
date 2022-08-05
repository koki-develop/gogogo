package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/ecr"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	aws.NewAwsProvider(stack, jsii.String("AWS"), &aws.AwsProviderConfig{
		Region: jsii.String("us-east-1"),
		DefaultTags: &aws.AwsProviderDefaultTags{
			Tags: &map[string]*string{
				"App": jsii.String("gogogo"),
			},
		},
	})

	s3bucketfrontend := s3.NewS3Bucket(stack, jsii.String("s3-bucket-frontend"), &s3.S3BucketConfig{
		Bucket: jsii.String("gogogo-frontend"),
	})

	cloudfrontoriginaccessidentity := cloudfront.NewCloudfrontOriginAccessIdentity(stack, jsii.String("cloudfront-origin-access-identity-frontend"), &cloudfront.CloudfrontOriginAccessIdentityConfig{})

	frontendbucketpolicy := iam.NewDataAwsIamPolicyDocument(stack, jsii.String("data-iam-policy-document-frontend-bucket-policy"), &iam.DataAwsIamPolicyDocumentConfig{
		Statement: []*iam.DataAwsIamPolicyDocumentStatement{{
			Effect:    jsii.String("Allow"),
			Actions:   jsii.Strings("s3:GetObject"),
			Resources: jsii.Strings(fmt.Sprintf("%s/*", *s3bucketfrontend.Arn())),
			Principals: []*iam.DataAwsIamPolicyDocumentStatementPrincipals{{
				Type:        jsii.String("AWS"),
				Identifiers: jsii.Strings(*cloudfrontoriginaccessidentity.IamArn()),
			}},
		}},
	})

	s3.NewS3BucketPolicy(stack, jsii.String("s3-bucket-policy-frontend"), &s3.S3BucketPolicyConfig{
		Bucket: s3bucketfrontend.Id(),
		Policy: frontendbucketpolicy.Json(),
	})

	s3.NewS3BucketPublicAccessBlock(stack, jsii.String("s3-public-access-block-frontend"), &s3.S3BucketPublicAccessBlockConfig{
		Bucket:                s3bucketfrontend.Bucket(),
		BlockPublicAcls:       jsii.Bool(true),
		BlockPublicPolicy:     jsii.Bool(true),
		IgnorePublicAcls:      jsii.Bool(true),
		RestrictPublicBuckets: jsii.Bool(true),
	})

	cloudfront.NewCloudfrontDistribution(stack, jsii.String("cloudfront-distribution-frontend"), &cloudfront.CloudfrontDistributionConfig{
		Enabled:           jsii.Bool(true),
		DefaultRootObject: jsii.String("index.html"),
		Origin: []*cloudfront.CloudfrontDistributionOrigin{{
			OriginId:   s3bucketfrontend.Id(),
			DomainName: s3bucketfrontend.BucketRegionalDomainName(),
			S3OriginConfig: &cloudfront.CloudfrontDistributionOriginS3OriginConfig{
				OriginAccessIdentity: cloudfrontoriginaccessidentity.CloudfrontAccessIdentityPath(),
			},
		}},
		DefaultCacheBehavior: &cloudfront.CloudfrontDistributionDefaultCacheBehavior{
			TargetOriginId:       s3bucketfrontend.Id(),
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

	apiecrrepository := ecr.NewEcrRepository(stack, jsii.String("ecr-repository-api"), &ecr.EcrRepositoryConfig{
		Name: jsii.String("gogogo-api"),
	})

	apilambdafunctioniamroleassumepolicy := iam.NewDataAwsIamPolicyDocument(stack, jsii.String("data-iam-policy-document-api-assume-policy"), &iam.DataAwsIamPolicyDocumentConfig{
		Statement: []*iam.DataAwsIamPolicyDocumentStatement{{
			Effect:  jsii.String("Allow"),
			Actions: jsii.Strings("sts:AssumeRole"),
			Principals: []*iam.DataAwsIamPolicyDocumentStatementPrincipals{{
				Type:        jsii.String("Service"),
				Identifiers: jsii.Strings("lambda.amazonaws.com"),
			}},
		}},
	})

	apilambdafunctioniamrole := iam.NewIamRole(stack, jsii.String("iam-role-api"), &iam.IamRoleConfig{
		Name:             jsii.String("gogogo-api-role"),
		AssumeRolePolicy: apilambdafunctioniamroleassumepolicy.Json(),
	})

	administratoraccesspolicy := iam.NewDataAwsIamPolicy(stack, jsii.String("iam-policy-adoministrator-access"), &iam.DataAwsIamPolicyConfig{
		Arn: jsii.String("arn:aws:iam::aws:policy/AdministratorAccess"),
	})
	iam.NewIamRolePolicyAttachment(stack, jsii.String("iam-role-policy-attachment-api-administorator-access-to-lambda-function-iam-role"), &iam.IamRolePolicyAttachmentConfig{
		Role:      apilambdafunctioniamrole.Name(),
		PolicyArn: administratoraccesspolicy.Arn(),
	})

	lambdafunction.NewLambdaFunction(stack, jsii.String("lambda-function-api"), &lambdafunction.LambdaFunctionConfig{
		FunctionName: jsii.String("gogogo-api"),
		Role:         apilambdafunctioniamrole.Arn(),
		PackageType:  jsii.String("Image"),
		ImageUri:     jsii.String(fmt.Sprintf("%s:latest", *apiecrrepository.RepositoryUrl())),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "infrastructure")

	cdktf.NewS3Backend(stack, &cdktf.S3BackendProps{
		Region:  jsii.String("us-east-1"),
		Bucket:  jsii.String("gogogo-tfstates"),
		Key:     jsii.String("terraform.tfstate"),
		Encrypt: jsii.Bool(true),
	})

	app.Synth()
}
