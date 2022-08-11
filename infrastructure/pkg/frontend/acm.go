package frontend

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/acm"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type certificateMainInput struct {
	Domain   string
	Hostzone route53.DataAwsRoute53Zone
}

func newCertificateMain(scope constructs.Construct, ipt *certificateMainInput) acm.AcmCertificate {
	cert := acm.NewAcmCertificate(scope, jsii.String("acm-certificate-frontend"), &acm.AcmCertificateConfig{
		DomainName:       &ipt.Domain,
		ValidationMethod: jsii.String("DNS"),
		Lifecycle: &cdktf.TerraformResourceLifecycle{
			CreateBeforeDestroy: jsii.Bool(true),
		},
	})

	validrec := route53.NewRoute53Record(scope, jsii.String("route53-record-frontend-certificate-validation"), &route53.Route53RecordConfig{
		ZoneId:  ipt.Hostzone.ZoneId(),
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
