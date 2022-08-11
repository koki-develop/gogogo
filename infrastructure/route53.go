package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
)

func NewHostzoneMain(scope constructs.Construct) route53.DataAwsRoute53Zone {
	return route53.NewDataAwsRoute53Zone(scope, jsii.String("route53-zone-default"), &route53.DataAwsRoute53ZoneConfig{
		Name:        jsii.String(Domain),
		PrivateZone: jsii.Bool(false),
	})
}
