package backend

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-archive-go/archive"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/s3"
	"github.com/koki-develop/gogogo/infrastructure/pkg/util"
)

type lambdaAPIInput struct {
	CatsBucket s3.S3Bucket
}

func newLambdaAPI(scope constructs.Construct, ipt *lambdaAPIInput) lambdafunction.LambdaFunction {
	assumepolicy := util.NewAssumePolicy(scope, "data-iam-policy-document-api-assume-policy", "lambda.amazonaws.com")

	policy := iam.NewDataAwsIamPolicyDocument(scope, jsii.String("data-iam-policy-document-api-policy"), &iam.DataAwsIamPolicyDocumentConfig{
		Statement: []*iam.DataAwsIamPolicyDocumentStatement{{
			Effect:    jsii.String("Allow"),
			Actions:   jsii.Strings("s3:GetObject"),
			Resources: jsii.Strings(fmt.Sprintf("%s/*", *ipt.CatsBucket.Arn())),
		}},
	})
	iamrole := iam.NewIamRole(scope, jsii.String("iam-role-api"), &iam.IamRoleConfig{
		Name:             jsii.String("gogogo-api-role"),
		AssumeRolePolicy: assumepolicy.Json(),
		InlinePolicy: []*iam.IamRoleInlinePolicy{{
			Name:   jsii.String("allow-access-to-cats-bucket"),
			Policy: policy.Json(),
		}},
	})

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// NOTE: backend 側で go build 実行済みであることを前提としている
	src := archive.NewDataArchiveFile(scope, jsii.String("archive-file-api-source"), &archive.DataArchiveFileConfig{
		Type:       jsii.String("zip"),
		SourceFile: jsii.String(path.Join(cwd, "../backend/dist/api")),
		OutputPath: jsii.String(path.Join(cwd, "dist/api.zip")),
	})

	return lambdafunction.NewLambdaFunction(scope, jsii.String("lambda-function-api"), &lambdafunction.LambdaFunctionConfig{
		FunctionName: jsii.String("gogogo-api"),
		Role:         iamrole.Arn(),
		PackageType:  jsii.String("Zip"),

		Filename:       src.OutputPath(),
		Handler:        jsii.String("api"),
		SourceCodeHash: src.OutputBase64Sha256(),
		Runtime:        jsii.String("go1.x"),
	})
}
