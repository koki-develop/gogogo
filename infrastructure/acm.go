package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/acm"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CertificateFrontendConfig struct {
	Hostzone route53.DataAwsRoute53Zone
}

func NewCertificateFrontend(scope constructs.Construct, cfg *CertificateFrontendConfig) acm.AcmCertificate {
	cert := acm.NewAcmCertificate(scope, jsii.String("acm-certificate-frontend"), &acm.AcmCertificateConfig{
		DomainName:       jsii.String("go55.dev"),
		ValidationMethod: jsii.String("DNS"),
		Lifecycle: &cdktf.TerraformResourceLifecycle{
			CreateBeforeDestroy: jsii.Bool(true),
		},
	})

	validrec := route53.NewRoute53Record(scope, jsii.String("route53-record-api-certificate-validation"), &route53.Route53RecordConfig{
		ZoneId:  cfg.Hostzone.ZoneId(),
		Name:    cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordName(),
		Type:    cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordType(),
		Records: &[]*string{cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordValue()},
		Ttl:     jsii.Number(60),
	})

	acm.NewAcmCertificateValidation(scope, jsii.String("acm-certificate-validation-frontend"), &acm.AcmCertificateValidationConfig{
		CertificateArn:        cert.Arn(),
		ValidationRecordFqdns: &[]*string{validrec.Fqdn()},
	})

	return cert
}
