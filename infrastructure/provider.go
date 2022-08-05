package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9"
)

func NewAwsProvider(scope constructs.Construct) aws.AwsProvider {
	return aws.NewAwsProvider(scope, jsii.String("AWS"), &aws.AwsProviderConfig{
		Region: jsii.String("us-east-1"),
		DefaultTags: &aws.AwsProviderDefaultTags{
			Tags: &map[string]*string{
				"App": jsii.String("gogogo"),
			},
		},
	})
}
