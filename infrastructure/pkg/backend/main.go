package backend

import "github.com/aws/constructs-go/constructs/v10"

func Apply(scope constructs.Construct) {
	cats := newS3Cats(scope)

	f := newLambdaAPI(scope, &lambdaAPIInput{
		CatsBucket: cats,
	})

	newAPIGatewayAPI(scope, &apiGatewayAPIInput{
		LambdaFunction: f,
	})
}
