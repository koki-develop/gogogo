package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
)

type HostzoneMainConfig struct {
	Name string
}

func NewHostzoneMain(scope constructs.Construct, cfg *HostzoneMainConfig) route53.DataAwsRoute53Zone {
	return route53.NewDataAwsRoute53Zone(scope, jsii.String("route53-zone-default"), &route53.DataAwsRoute53ZoneConfig{
		Name:        &cfg.Name,
		PrivateZone: jsii.Bool(false),
	})
}
