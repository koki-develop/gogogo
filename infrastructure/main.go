package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/apigateway"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	domain := "go55.dev"

	stack := cdktf.NewTerraformStack(scope, &id)

	NewAwsProvider(stack)
	NewArchiveProvider(stack)

	hostzone := NewHostzoneMain(stack, &HostzoneMainConfig{Name: domain})

	// フロントエンド
	certfrontend := NewCertificateFrontend(stack, &CertificateFrontendConfig{Hostzone: hostzone})
	cfoai := NewCloudfrontOriginAccessIdentity(stack, "cloudfront-origin-access-identity-frontend")
	s3frontend := NewS3Frontend(stack, &S3FrontendConfig{OriginAccessIdentity: cfoai})
	cffrontend := NewCloudfrontFrontend(stack, &CloudfrontFrontendConfig{
		Domain:               domain,
		Bucket:               s3frontend,
		OriginAccessIdentity: cfoai,
		Certificate:          certfrontend,
	})

	route53.NewRoute53Record(stack, jsii.String("route53-record-frontend"), &route53.Route53RecordConfig{
		ZoneId: hostzone.Id(),
		Name:   jsii.String("go55.dev"),
		Type:   jsii.String("A"),
		Alias: []*route53.Route53RecordAlias{{
			Name:                 cffrontend.DomainName(),
			ZoneId:               cffrontend.HostedZoneId(),
			EvaluateTargetHealth: jsii.Bool(false),
		}},
	})

	NewS3Cats(stack)

	apifunction := NewAPILambda(stack)

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
