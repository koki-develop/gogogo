package frontend

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/cdktf-provider-aws-go/aws/v9/route53"
	"github.com/koki-develop/gogogo/infrastructure/pkg/util"
)

type Input struct {
	Domain   string
	Hostzone route53.DataAwsRoute53Zone
}

func Apply(scope constructs.Construct, ipt *Input) {
	cert := newCertificateMain(scope, &certificateMainInput{
		Domain:   ipt.Domain,
		Hostzone: ipt.Hostzone,
	})

	cfoai := util.NewCloudfrontOriginAccessIdentity(scope, "cloudfront-origin-access-identity-frontend")

	s3 := newS3Main(scope, &s3MainInput{OriginAccessIdentity: cfoai})

	cf := newCloudfrontMain(scope, &cloudfrontMainInput{
		Domain:               ipt.Domain,
		Bucket:               s3,
		OriginAccessIdentity: cfoai,
		Certificate:          cert,
	})

	newRecordMain(scope, &recordFrontendInput{
		Domain:       ipt.Domain,
		Hostzone:     ipt.Hostzone,
		Distribution: cf,
	})
}
