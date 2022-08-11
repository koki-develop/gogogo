package frontend

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/cloudfront"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
)

type recordFrontendInput struct {
	Domain       string
	Hostzone     route53.DataAwsRoute53Zone
	Distribution cloudfront.CloudfrontDistribution
}

func newRecordMain(scope constructs.Construct, ipt *recordFrontendInput) route53.Route53Record {
	return route53.NewRoute53Record(scope, jsii.String("route53-record-frontend"), &route53.Route53RecordConfig{
		ZoneId: ipt.Hostzone.Id(),
		Name:   &ipt.Domain,
		Type:   jsii.String("A"),
		Alias: []*route53.Route53RecordAlias{{
			Name:                 ipt.Distribution.DomainName(),
			ZoneId:               ipt.Distribution.HostedZoneId(),
			EvaluateTargetHealth: jsii.Bool(false),
		}},
	})
}
