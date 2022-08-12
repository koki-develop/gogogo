package util

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
)

func NewAssumePolicy(scope constructs.Construct, id, service string) iam.DataAwsIamPolicyDocument {
	return iam.NewDataAwsIamPolicyDocument(scope, &id, &iam.DataAwsIamPolicyDocumentConfig{
		Statement: []*iam.DataAwsIamPolicyDocumentStatement{{
			Effect:  jsii.String("Allow"),
			Actions: jsii.Strings("sts:AssumeRole"),
			Principals: []*iam.DataAwsIamPolicyDocumentStatementPrincipals{{
				Type:        jsii.String("Service"),
				Identifiers: &[]*string{&service},
			}},
		}},
	})

}

func NewDataIamPolicy(scope constructs.Construct, id, arn string) iam.DataAwsIamPolicy {
	return iam.NewDataAwsIamPolicy(scope, &id, &iam.DataAwsIamPolicyConfig{
		Arn: &arn,
	})
}
