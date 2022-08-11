package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/koki-develop/gogogo/infrastructure/pkg/backend"
	"github.com/koki-develop/gogogo/infrastructure/pkg/frontend"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	domain := "go55.dev"

	stack := cdktf.NewTerraformStack(scope, &id)

	NewAwsProvider(stack)
	NewArchiveProvider(stack)

	hostzone := NewHostzoneMain(stack, &HostzoneMainConfig{Name: domain})

	frontend.Apply(stack, &frontend.Input{
		Domain:   domain,
		Hostzone: hostzone,
	})

	backend.Apply(stack)

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
