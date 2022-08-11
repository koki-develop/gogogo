package backend

import (
	"os"
	"path"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-archive-go/archive"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/koki-develop/gogogo/infrastructure/pkg/util"
)

func newLambdaAPI(scope constructs.Construct) lambdafunction.LambdaFunction {
	assumepolicy := util.NewAssumePolicy(scope, "data-iam-policy-document-api-assume-policy", "lambda.amazonaws.com")
	iamrole := iam.NewIamRole(scope, jsii.String("iam-role-api"), &iam.IamRoleConfig{
		Name:             jsii.String("gogogo-api-role"),
		AssumeRolePolicy: assumepolicy.Json(),
	})
	// TODO: ちゃんと絞り込む
	adminpolicy := iam.NewDataAwsIamPolicy(scope, jsii.String("iam-policy-adoministrator-access"), &iam.DataAwsIamPolicyConfig{
		Arn: jsii.String("arn:aws:iam::aws:policy/AdministratorAccess"),
	})
	iam.NewIamRolePolicyAttachment(scope, jsii.String("iam-role-policy-attachment-api-administorator-access-to-lambda-function-iam-role"), &iam.IamRolePolicyAttachmentConfig{
		Role:      iamrole.Name(),
		PolicyArn: adminpolicy.Arn(),
	})

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// NOTE: backend 側で go build 実行済みであることを前提としている
	apisourcearchive := archive.NewDataArchiveFile(scope, jsii.String("archive-file-api-source"), &archive.DataArchiveFileConfig{
		Type:       jsii.String("zip"),
		SourceFile: jsii.String(path.Join(cwd, "../backend/dist/api")),
		OutputPath: jsii.String(path.Join(cwd, "dist/api.zip")),
	})

	return lambdafunction.NewLambdaFunction(scope, jsii.String("lambda-function-api"), &lambdafunction.LambdaFunctionConfig{
		FunctionName: jsii.String("gogogo-api"),
		Role:         iamrole.Arn(),
		PackageType:  jsii.String("Zip"),

		Filename:       apisourcearchive.OutputPath(),
		Handler:        jsii.String("api"),
		SourceCodeHash: apisourcearchive.OutputBase64Sha256(),
		Runtime:        jsii.String("go1.x"),
	})
}
