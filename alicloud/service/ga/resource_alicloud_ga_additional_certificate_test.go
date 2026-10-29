package ga_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaAdditionalCertificate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_additional_certificate.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaAdditionalCertificateMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaAdditionalCertificate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaadditionalcertificate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaAdditionalCertificateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":    "${alicloud_ga_listener.default.id}",
					"domain":         "${local.domain}",
					"certificate_id": "${local.certificate_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id": acctest.CHECKSET,
						"listener_id":    acctest.CHECKSET,
						"domain":         acctest.CHECKSET,
						"certificate_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_id": "${local.certificate_id_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudGaAdditionalCertificateMap0 = map[string]string{}

func AliCloudGaAdditionalCertificateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Basic"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
  		auto_pay               = true
  		auto_use_coupon        = true
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		count            = 2
  		certificate_name = join("-", [var.name, count.index])
		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIIDOTCCAiGgAwIBAgIUdRfPhtKuM0RuDc0vYpuqbZ68m0UwDQYJKoZIhvcNAQEL
BQAwLDELMAkGA1UEBhMCQ04xHTAbBgNVBAMMFGFsaWNsb3VkLXByb3ZpZGVyLmNu
MB4XDTI2MTAyODE1NDE1OFoXDTMxMTAyNzE1NDE1OFowLDELMAkGA1UEBhMCQ04x
HTAbBgNVBAMMFGFsaWNsb3VkLXByb3ZpZGVyLmNuMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEA4Sr6K4qXPR7ZTfruMJ8RmQxOk/RvnK/8dYMjexBmY8WJ
WrMCpfbOsuiYYiaQYciITVOV1046zShgwLLoDfIMAjPgDk/sEu/jI8EUY1/3DYzw
+8DT0aumptZA3Sl/ZtJUcNRBt3fpIf/CJNYBQcH8Bjp7TT7WrBFNKZuAN7iS0ddT
eECImbBLYrIO838OAtPi6TEfbWKA48f2wU82Ut4Nq54OgVVszLvgH7ToUTave8pf
00YJScdc28QhUAIrrcLpHnaRXpRsEG2kgIameeYueyo533DDhtpZGhzDfmC8J2Wy
SdR9zG2SxMB+pVsmQAk9F9qpCyAjbsYezBla2AehuwIDAQABo1MwUTAdBgNVHQ4E
FgQUYrKkKT/IORw+C2TT5/lxSSNZvhkwHwYDVR0jBBgwFoAUYrKkKT/IORw+C2TT
5/lxSSNZvhkwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAtuO6
hrTbSuA/UZi2b/O5ABCNdU0hGAjbEhqMjtRd15Fr1+H0rDZjO139bcBsOwbx2aRE
KrKIeNuAmK6ZPRycI+JeyviWOFkhjJ+cLT12AgA8HZPLsY7WGEcBz4GCRAdKlhqf
3BbAH7jRQDomnjafzCWcUUqnIXRzDt85uqXPxvUf1RzsZl+Bm04fLeaUVChp1bkp
xeKvEgYJyzOu3L+PwJiEtZ1HykloyCPbpkrLbXxVf7sh4nLo2a7VQMjB0wHmC7WN
S2yUXGRGfEBlmNTpzY6O8JalzWJkVBZME1r6ILmxMvag4VS81p4VAlse5fMLN9/I
8YLoLo7lE3fZEdo5UQ==
-----END CERTIFICATE-----
EOF
		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA4Sr6K4qXPR7ZTfruMJ8RmQxOk/RvnK/8dYMjexBmY8WJWrMC
pfbOsuiYYiaQYciITVOV1046zShgwLLoDfIMAjPgDk/sEu/jI8EUY1/3DYzw+8DT
0aumptZA3Sl/ZtJUcNRBt3fpIf/CJNYBQcH8Bjp7TT7WrBFNKZuAN7iS0ddTeECI
mbBLYrIO838OAtPi6TEfbWKA48f2wU82Ut4Nq54OgVVszLvgH7ToUTave8pf00YJ
Scdc28QhUAIrrcLpHnaRXpRsEG2kgIameeYueyo533DDhtpZGhzDfmC8J2WySdR9
zG2SxMB+pVsmQAk9F9qpCyAjbsYezBla2AehuwIDAQABAoIBAHQwkytItsYpVwTw
SVz5mEEboqiiZ7246lNROkvUqKEvW7Nb0DAH7Z1QfOGY+MQckIkmapnRmnInNdX5
+9jGX//6+NEp2K3k75aRmY22BWwOoIOoeyUkqDhLMuls8nsCaeRXAxpDXB/ZNDuB
2uAt03j7UIkJtgjQXLrtrPoT2l9A/8dwfD/YolPdOLw+l6rZEXb/AzSQz+pJ6SMB
ehj2/EuQgOQvLrWdXDel/ykl5NUYrKBHv0rwjER4YTKtUwJa8jrqfQrCgMqqJlM7
ErEW+v0Lo7Z6xRwUQhiF8SiaAvJeHm3IQzz/XM4PuVHEzbRav4T5FwizlnjUa2FW
EoarGMECgYEA9rwXNnW/ssML73CxDGyaX5AO1iuOpW07CNFoHMng3jRIrDn5NSDe
NGto/C4f3sS+w3MX10Wk4R9umSRrw7A6sFsj3hEAwzOPQhzms+JvD7ofLYxTUFkx
ethOZ7ECNoudM8bgIEy/QJmvBDLKVLD6PSPdqkmDfmEL3A72IJx1xzECgYEA6Z+P
aU1HtaxazDcn6ErWc1hIUOQNFRO74eyc3FJADHVnm92Xsfr0vsH7PGZkmqSftguy
ecI/SQ4KAGsQpdqyFDtFUSwPgdLHRf3QUfD5hOruL5HCmUOKoe8U4ewUQ6gJ1PIf
EStNSoMG4pVXTmS1fm8CnKqrmBNaVTn4TpcX1KsCgYEAycxCa/9YHzvZaidQFupI
RDqNd+s+CtzQuMCUl18htHssqPftc9+C/XqP7hlHZaqinzZl3GKQnLNZgJVCMlXT
WmPgbQKZDg4AcEE1CxD6ZvGCCH3JsmCQ0HlTfwlN1xDPLQbOTnyWu542f4eO5j+t
Gw5XgKwJiPaUfZIbCccNqTECgYEAobjhtrxXcqzAIfPgfU4pC6nGTbDAWaHSgr+f
imipi7ZOJbiMxRj7CX4Q46vh4qcpJAUO5HiSBIliZQ24dYdhlN2+y0Ub6vu1FH9y
EPcmU1upzGYzdpMpyZVOOGKLWJZfT/jXhrQaCgWP45QL3umaNc602SzQIjLramn5
uLWhOocCgYEAstYQR5bgpmuexSqFYwK7FGyi4Sh6+NQWGez1Oijzz5V1NSGhdJ8j
a52a5oSXTQFLM2/DjjvDB+WuSEPhp5DTmJY2I1E81YNCHpT64T5E1EupgaiZkzp/
O5XDKhiybbWfTefLsMVczvGia69okf7GwXCCFm8wvUwiMtPZ4aqtYeU=
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ssl_certificates_service_certificate" "update" {
  		certificate_name = "${var.name}-update"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7jCCAtagAwIBAgIQUNnSVa/sQNeb9pBN9NhkwTANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yMzA4MDkwMzM4MThaFw0yODA4MDcwMzM4MThaMCwxCzAJBgNVBAYTAkNOMR0w
GwYDVQQDExRhbGljbG91ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAOgskr8dEfZYdjr0xaIqlCkmE802vABoj3SQNn3rLWnUj+1v
Wqbpsj6Bu61Scb8mtl/OZOOM7sgq0Q1hpdO8xvMGxTMuZ2bjX0EqCMqh4AvFofHL
a/iVD07hfoM1Jo8CEidh1uvcOuXP1TlaqU020x1TX3a3niJu4JVkmCkCOwAbWYuj
O8IsgBCsFaF9d4+C1JRYOtRbIHCNhd0sxG8AGovUDLvlkePeH5NF7DNvFXgGJ4iv
EQcY9pP08RBFUkaznOw/r64Up7zhLb+Ie4SyAvs1FulhMAmIXOcbsND39hJ+/WIP
8beWvIN1eCS8zcvgAvDgMkV8oqqVbQu1dqx5WuMCAwEAAaOB2TCB1jAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFCiBJgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEF
BQcwAYYVaHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8v
Y2EubXlzc2wuY29tL215c3NsdGVzdHJzYS5jcnQwHwYDVR0RBBgwFoIUYWxpY2xv
dWQtcHJvdmlkZXIuY24wDQYJKoZIhvcNAQELBQADggEBALd0hFZAd2XHJgETbHQs
h4YUBNKxrIy6JiWfxffhIL1ZK5pI443DC4VRGfxVi3zWqs01WbNtJ2b1KdfSoovH
Zwi3hdMF1IwoAB/Y2sS4zjqS0H1od7MN9KKHes6bl3yCgpmaYs5cHbyg0IJHmeq3
rCgbKsvHfUwtzBNNPHlpANakAYd/5O1pztmUskWMUVaExfpMoQLo/AX9Lqm8pVjw
xs921I703l/E5zEnd3PVSYagy/KQJrwVt+wQZS11HsAryfO9kct/9f+c85VDo6Ht
iRirW/EnNPQRSno4z0V2x1Rn5+ZaoJo8cWzPvKrdfCG9TUozt4AR/LIudNLb6NNW
n7g=
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6CySvx0R9lh2OvTFoiqUKSYTzTa8AGiPdJA2festadSP7W9a
pumyPoG7rVJxvya2X85k44zuyCrRDWGl07zG8wbFMy5nZuNfQSoIyqHgC8Wh8ctr
+JUPTuF+gzUmjwISJ2HW69w65c/VOVqpTTbTHVNfdreeIm7glWSYKQI7ABtZi6M7
wiyAEKwVoX13j4LUlFg61FsgcI2F3SzEbwAai9QMu+WR494fk0XsM28VeAYniK8R
Bxj2k/TxEEVSRrOc7D+vrhSnvOEtv4h7hLIC+zUW6WEwCYhc5xuw0Pf2En79Yg/x
t5a8g3V4JLzNy+AC8OAyRXyiqpVtC7V2rHla4wIDAQABAoIBABKGQ+sluaIrKrvH
feFTfmDOHfRYsqVhslh9jSt80THJePZb1SLOMJ+WIFBS7Kpwv0pjoF8bho3IBMgJ
i36aaFFJsABGao+mApqjbPIl+kdWLHarYWEDG6aSjVKQshPk+WfVAZ3uA3EEpSGf
XzS+9Bc56LsDKYXbzOV+kjlraSO35AMec3CpISdx4K1caEAhKX6it9bvPq4pSYXi
PQspba0Jv46VV7MaabVjLzsinz5/md4vxyYHNIJAukHUfwJIsVC9ZNxukwSw+CzE
MMO64ylq2DGokNerGsLetuViV8UWi7qmUmms2fAmchodW16olgNkYTz27+V/A42S
eex63pkCgYEA+CqKhqp3qPe2E9KVrycrwjoycxmhOn3Iz1xiN7uAEv+DzfKtfZVf
mcOIiqw4Z82RkgjHb9vJuTigKdDkB1zE2gSDnep44sDWJM/5nPjGlMgnkiJWJhci
CnD0P4d6cT5wyDt7Q0/tS6ql2UrCpW4ktw1AP0Rm/z/VBD8jGkVenjcCgYEA74DM
Z2Qmh3bPt1TykpOlw+H+sEuvlkYxqMlbtn3Rv3WgEPIBekOFrgP7n/uLW1Aizn8w
EhNBBAE8w5jvklqZWYbpFMJQc09eqUkI8aTbLooZbzYj1f3CrzBRKn1GoTPmN9V0
j9r+TbH3/5CEoqlsJdmeQPofuv5Qid2oEutZcrUCgYBuZ16hco0xmqJiRzlYZvDM
w99V3X0g7Hy947e+W6gqy4nzwZb1W9LgMWE5cEzXwViVw1oWpY0k3dBDSi9oJxlc
dM2pH3sQRgH+9pdyAis2XaVdGfGBmKEITCAdc0RBxSmfqva3h4NmOlD2TpAx0MJ8
vWRrwR6hR+CYtw4CzgG+GQKBgQDGmi5lugW9JUe/xeBUrbyyv0+cT1auLUz2ouq7
XIA23Mo74wJYqW9Lyp+4nTWFJeGHDK8G/hJWyNPjeomG+jvZombbQPrHc9SSWi7h
eowKfpfywZlb1M7AyTc1HacY+9l3CTlcJQPl16NHuEZUQFue02NIjGENhd+xQy4h
ainFVQKBgAoPs9ebtWbBOaqGpOnquzb7WibvW/ifOfzv5/aOhkbpp+KCjcSON6CB
QF3BEXMcNMGWlpPrd8PaxCAzR4MyU///ekJri2icS9lrQhGSz2TtYhdED4pv1Aag
7eTPl5L7xAwphCSwy8nfCKmvlqcX/MSJ7A+LHB/2hdbuuEOyhpbu
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		name           = var.name
  		protocol       = "HTTPS"
  		port_ranges {
    		from_port = 8080
    		to_port   = 8080
  		}
  		certificates {
    		id = join("-", [alicloud_ssl_certificates_service_certificate.default.0.id, "%s"])
  		}
	}

locals {
  domain                = "alicloud-provider.cn"
  certificate_id        = join("-", [alicloud_ssl_certificates_service_certificate.default.1.id, "%s"])
  certificate_id_update = join("-", [alicloud_ssl_certificates_service_certificate.update.id, "%s"])
}
`, name, acctest.DefaultRegionToTest, acctest.DefaultRegionToTest, acctest.DefaultRegionToTest)
}
