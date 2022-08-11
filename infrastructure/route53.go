package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
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

type RecordFrontendConfig struct {
	Domain       string
	Hostzone     route53.DataAwsRoute53Zone
	Distribution cloudfront.CloudfrontDistribution
}

func NewRecordFrontend(scope constructs.Construct, cfg *RecordFrontendConfig) route53.Route53Record {
	return route53.NewRoute53Record(scope, jsii.String("route53-record-frontend"), &route53.Route53RecordConfig{
		ZoneId: cfg.Hostzone.Id(),
		Name:   &cfg.Domain,
		Type:   jsii.String("A"),
		Alias: []*route53.Route53RecordAlias{{
			Name:                 cfg.Distribution.DomainName(),
			ZoneId:               cfg.Distribution.HostedZoneId(),
			EvaluateTargetHealth: jsii.Bool(false),
		}},
	})
}
