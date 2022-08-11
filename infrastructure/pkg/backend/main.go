package backend

import "github.com/aws/constructs-go/constructs/v10"

func Apply(scope constructs.Construct) {
	newS3Cats(scope)
	funcapi := newLambdaAPI(scope)
	newAPIGatewayAPI(scope, &apiGatewayAPIConfig{LambdaFunction: funcapi})
}
