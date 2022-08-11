package backend

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/apigateway"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/lambdafunction"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type apiGatewayAPIConfig struct {
	LambdaFunction lambdafunction.LambdaFunction
}

func newAPIGatewayAPI(scope constructs.Construct, ipt *apiGatewayAPIConfig) apigateway.ApiGatewayRestApi {
	api := apigateway.NewApiGatewayRestApi(scope, jsii.String("api-gateway-api"), &apigateway.ApiGatewayRestApiConfig{
		Name: jsii.String("gogogo-api"),
		EndpointConfiguration: &apigateway.ApiGatewayRestApiEndpointConfiguration{
			Types: jsii.Strings("EDGE"),
		},
	})

	rsrc := apigateway.NewApiGatewayResource(scope, jsii.String("api-gateway-api-resource"), &apigateway.ApiGatewayResourceConfig{
		ParentId:  api.RootResourceId(),
		PathPart:  jsii.String("{proxy+}"),
		RestApiId: api.Id(),
	})

	mth := apigateway.NewApiGatewayMethod(scope, jsii.String("api-gateway-api-resource-any"), &apigateway.ApiGatewayMethodConfig{
		HttpMethod:    jsii.String("ANY"),
		ResourceId:    rsrc.Id(),
		RestApiId:     api.Id(),
		Authorization: jsii.String("NONE"),
	})

	apigateway.NewApiGatewayIntegration(scope, jsii.String("api-gateway-api-integration"), &apigateway.ApiGatewayIntegrationConfig{
		RestApiId:             api.Id(),
		ResourceId:            mth.ResourceId(),
		HttpMethod:            mth.HttpMethod(),
		IntegrationHttpMethod: jsii.String("POST"),
		Type:                  jsii.String("AWS_PROXY"),
		Uri:                   ipt.LambdaFunction.InvokeArn(),
	})

	dep := apigateway.NewApiGatewayDeployment(scope, jsii.String("api-gateway-deployment"), &apigateway.ApiGatewayDeploymentConfig{
		RestApiId: api.Id(),
		Lifecycle: &cdktf.TerraformResourceLifecycle{
			CreateBeforeDestroy: jsii.Bool(true),
		},
	})

	apigateway.NewApiGatewayStage(scope, jsii.String("api-gateway-stage"), &apigateway.ApiGatewayStageConfig{
		RestApiId:    api.Id(),
		DeploymentId: dep.Id(),
		StageName:    jsii.String("prod"),
	})

	lambdafunction.NewLambdaPermission(scope, jsii.String("lambda-function-permission-api"), &lambdafunction.LambdaPermissionConfig{
		Action:       jsii.String("lambda:InvokeFunction"),
		FunctionName: ipt.LambdaFunction.FunctionName(),
		Principal:    jsii.String("apigateway.amazonaws.com"),
		SourceArn:    jsii.String(fmt.Sprintf("%s/*/*", *api.ExecutionArn())),
	})

	return api
}
