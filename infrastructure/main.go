package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/apigateway"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/ecr"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/iam"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	NewAwsProvider(stack)

	cloudfrontoriginaccessidentity := cloudfront.NewCloudfrontOriginAccessIdentity(stack, jsii.String("cloudfront-origin-access-identity-frontend"), &cloudfront.CloudfrontOriginAccessIdentityConfig{})
	s3bucketfrontend := NewS3Frontend(stack, cloudfrontoriginaccessidentity)

	NewCloudfrontFrontend(stack, s3bucketfrontend, cloudfrontoriginaccessidentity)

	apiecrrepository := ecr.NewEcrRepository(stack, jsii.String("ecr-repository-api"), &ecr.EcrRepositoryConfig{
		Name: jsii.String("gogogo-api"),
	})

	apilambdafunctioniamroleassumepolicy := NewAssumePolicy(stack, "data-iam-policy-document-api-assume-policy", "lambda.amazonaws.com")

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

	apifunction := lambdafunction.NewLambdaFunction(stack, jsii.String("lambda-function-api"), &lambdafunction.LambdaFunctionConfig{
		FunctionName: jsii.String("gogogo-api"),
		Role:         apilambdafunctioniamrole.Arn(),
		PackageType:  jsii.String("Image"),
		ImageUri:     jsii.String(fmt.Sprintf("%s:latest", *apiecrrepository.RepositoryUrl())),
	})

	apigatewayapi := apigateway.NewApiGatewayRestApi(stack, jsii.String("api-gateway-api"), &apigateway.ApiGatewayRestApiConfig{
		Name: jsii.String("gogogo-api"),
		EndpointConfiguration: &apigateway.ApiGatewayRestApiEndpointConfiguration{
			Types: jsii.Strings("EDGE"),
		},
	})

	apigatewayapiresource := apigateway.NewApiGatewayResource(stack, jsii.String("api-gateway-api-resource"), &apigateway.ApiGatewayResourceConfig{
		ParentId:  apigatewayapi.RootResourceId(),
		PathPart:  jsii.String("{proxy+}"),
		RestApiId: apigatewayapi.Id(),
	})

	apigatewayapiresourceany := apigateway.NewApiGatewayMethod(stack, jsii.String("api-gateway-api-resource-any"), &apigateway.ApiGatewayMethodConfig{
		HttpMethod:    jsii.String("ANY"),
		ResourceId:    apigatewayapiresource.Id(),
		RestApiId:     apigatewayapi.Id(),
		Authorization: jsii.String("NONE"),
	})

	apigateway.NewApiGatewayIntegration(stack, jsii.String("api-gateway-api-integration"), &apigateway.ApiGatewayIntegrationConfig{
		RestApiId:             apigatewayapi.Id(),
		ResourceId:            apigatewayapiresourceany.ResourceId(),
		HttpMethod:            apigatewayapiresourceany.HttpMethod(),
		IntegrationHttpMethod: jsii.String("POST"),
		Type:                  jsii.String("AWS_PROXY"),
		Uri:                   apifunction.InvokeArn(),
	})

	apigatewaydeployment := apigateway.NewApiGatewayDeployment(stack, jsii.String("api-gateway-deployment"), &apigateway.ApiGatewayDeploymentConfig{
		RestApiId: apigatewayapi.Id(),
		Lifecycle: &cdktf.TerraformResourceLifecycle{
			CreateBeforeDestroy: jsii.Bool(true),
		},
	})

	apigateway.NewApiGatewayStage(stack, jsii.String("api-gateway-stage"), &apigateway.ApiGatewayStageConfig{
		RestApiId:    apigatewayapi.Id(),
		DeploymentId: apigatewaydeployment.Id(),
		StageName:    jsii.String("prod"),
	})

	lambdafunction.NewLambdaPermission(stack, jsii.String("lambda-function-permission-api"), &lambdafunction.LambdaPermissionConfig{
		Action:       jsii.String("lambda:InvokeFunction"),
		FunctionName: apifunction.FunctionName(),
		Principal:    jsii.String("apigateway.amazonaws.com"),
		SourceArn:    jsii.String(fmt.Sprintf("%s/*/*", *apigatewayapi.ExecutionArn())),
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
