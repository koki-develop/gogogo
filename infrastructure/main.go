package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
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
	NewRecordFrontend(stack, &RecordFrontendConfig{
		Domain:       domain,
		Hostzone:     hostzone,
		Distribution: cffrontend,
	})

	// バックエンド
	NewS3Cats(stack)
	apifunc := NewAPILambda(stack)
	NewAPIGatewayForAPI(stack, &APIGatewayForAPIConfig{LambdaFunction: apifunc})

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
