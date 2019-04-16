package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cas_certificate", &resource.Sweeper{
		Name: "alicloud_cas_certificate",
		F:    testSweepCasCertificate,
	})
}

func testSweepCasCertificate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var allcerts []cas.Certificate
	req := cas.CreateDescribeUserCertificateListRequest()
	req.RegionId = client.RegionId
	req.ShowSize = requests.NewInteger(PageSizeLarge)
	req.CurrentPage = requests.NewInteger(1)
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
				return casClient.DescribeUserCertificateList(req)
			})
			raw = rsp
			return err
		}); err != nil {
			log.Printf("[ERROR] Error retrieving Certificates: %s", WrapError(err))
		}
		resp, _ := raw.(*cas.DescribeUserCertificateListResponse)
		if resp == nil || len(resp.CertificateList) < 1 {
			break
		}
		allcerts = append(allcerts, resp.CertificateList...)

		if len(resp.CertificateList) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.CurrentPage); err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		} else {
			req.CurrentPage = page
		}
	}

	for _, v := range allcerts {
		name := v.Name
		id := v.Id
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), "cert-"+strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Certificate: %s (%d)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Certificate: %s (%d)", name, id)
		req := cas.CreateDeleteUserCertificateRequest()
		req.CertId = requests.NewInteger(id)
		_, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
			return casClient.DeleteUserCertificate(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Certificate (%s (%d)): %s", name, id, err)
		}
	}

	return nil
}

func TestAccAlicloudCasCertificate_basic(t *testing.T) {
	var v *cas.Certificate

	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CasClassicSupportedRegions)
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cas_certificate.cert",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCasCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCasCertificateConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCasExists("alicloud_cas_certificate.cert", v),
					resource.TestCheckResourceAttr("alicloud_cas_certificate.cert", "name", fmt.Sprintf("tf_testAcc_%v", randInt)),
				),
			},
		},
	})
}

func testAccCheckCasExists(n string, cert *cas.Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		casService := &CasService{client: client}
		certInfo, err := casService.DescribeCas(rs.Primary.Attributes["id"])
		log.Printf("[WARN] Cert id %#v", rs.Primary.ID)

		if err == nil {
			cert = certInfo
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckCasCertificateDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cas_certificate" {
			continue
		}

		// Try to find the cert
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		casService := &CasService{client: client}
		cert, err := casService.DescribeCas(rs.Primary.Attributes["id"])

		if cert != nil {
			return WrapError(err)
		}
	}

	return nil
}

func testAccCasCertificateConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_cas_certificate" "cert" {
  name = "tf_testAcc_%d"
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIIFoTCCBImgAwIBAgIQA06FWjWYxa0rnA/tSO0JUTANBgkqhkiG9w0BAQsFADBu
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMS0wKwYDVQQDEyRFbmNyeXB0aW9uIEV2ZXJ5d2hlcmUg
RFYgVExTIENBIC0gRzEwHhcNMTgxMjE3MDAwMDAwWhcNMTkxMjE3MTIwMDAwWjAp
MScwJQYDVQQDEx5mcmVlLm9zcy5jZXJ0aWZpY2F0ZXN0ZXN0cy5jb20wggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCWA4caf9ZdcsIcDbebMktxyk4Ddoit
8EV0kNdywyNEOGukuZWU4pXUuHv0miZAZU5DBTKCz1bC/1Nb5Fc8xyat8OPX24DG
cWE/XeDCVvCHP/8KrG8Mnej2nn18IH601mADMm2NBArQt8yaBogCba0JX9elVZYW
tQwF8724ymvEt6xDWQs1uWFGLP4WSZWJduOHDuipccEf3zPohBfNsKIvuOICKNY8
93J4HN4jYZxLQdHnUXAxr6XK5pIizr1yd2IKUcPb/nNvD4ouboFJp121V56wxgXn
sRwQctG42DsloMojlAKF0r+Xf3YUXw8VD2UkCc9jUqe2Nmv+QerReJDxAgMBAAGj
ggJ+MIICejAfBgNVHSMEGDAWgBRVdE+yck/1YLpQ0dfmUVyaAYca1zAdBgNVHQ4E
FgQUXGJ22anPRc0VuJsZIt1R4sxJ2cQwKQYDVR0RBCIwIIIeZnJlZS5vc3MuY2Vy
dGlmaWNhdGVzdGVzdHMuY29tMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggr
BgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggr
BgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEw
fQYIKwYBBQUHAQEEcTBvMCEGCCsGAQUFBzABhhVodHRwOi8vb2NzcC5kY29jc3Au
Y24wSgYIKwYBBQUHMAKGPmh0dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9FbmNy
eXB0aW9uRXZlcnl3aGVyZURWVExTQ0EtRzEuY3J0MAkGA1UdEwQCMAAwggEEBgor
BgEEAdZ5AgQCBIH1BIHyAPAAdgC72d+8H4pxtZOUI5eqkntHOFeVCqtS6BqQlmQ2
jh7RhQAAAWe7Zn9bAAAEAwBHMEUCIQCHzPlasoPFAG2+3wFXbljXjvD0YoGKxYvg
FPDp1gLJrQIgLUSpZhhO1c78hIvX/CN9fClUFeDddAp0EFqD/bIJtwEAdgCHdb/n
WXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAWe7ZoA2AAAEAwBHMEUCIQCY
KcgsUvmTy1xGRD5Ai1lK17ncotqkxHGNImlu8s+KugIgN9nkEzG3aFGm3RrQkGX8
+/m/TAyldZTiy8x8BSHqKh0wDQYJKoZIhvcNAQELBQADggEBAKUWU9X01Y/87JAg
cbrP5xvxWbES8VsxOs5QcFmGpLIZZr1mdIYm+l0WfksxOb8xwRog/fWOjFE02tuf
SqnDnSiUwknFX1YAcf5Z9xei+UQo17W0U46wUjwnP5BRhoIu5pPt8+eTs7/IMkpR
gsXzKFv3wf+0CsqfkfTOB2eDk4SCQVWi3N0ESbV+bDpZ+4/yrlIK9VXIOcXKrZM4
U8JxqXbxQJTnKht4tXAtrAdrCAoLJsG7l+hGIRlNYiogB0uOwb/T9N4NBBHaqpKv
6YkVeMOtiJK/rzK+luuLQEnfeLvRTUdxPbTfEIOZTt3dYrG1QcCfPGKm12dzvx8e
zq8BC8Q=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEqjCCA5KgAwIBAgIQAnmsRYvBskWr+YBTzSybsTANBgkqhkiG9w0BAQsFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0xNzExMjcxMjQ2MTBaFw0yNzExMjcxMjQ2MTBaMG4xCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j
b20xLTArBgNVBAMTJEVuY3J5cHRpb24gRXZlcnl3aGVyZSBEViBUTFMgQ0EgLSBH
MTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALPeP6wkab41dyQh6mKc
oHqt3jRIxW5MDvf9QyiOR7VfFwK656es0UFiIb74N9pRntzF1UgYzDGu3ppZVMdo
lbxhm6dWS9OK/lFehKNT0OYI9aqk6F+U7cA6jxSC+iDBPXwdF4rs3KRyp3aQn6pj
pp1yr7IB6Y4zv72Ee/PlZ/6rK6InC6WpK0nPVOYR7n9iDuPe1E4IxUMBH/T33+3h
yuH3dvfgiWUOUkjdpMbyxX+XNle5uEIiyBsi4IvbcTCh8ruifCIi5mDXkZrnMT8n
wfYCV6v6kDdXkbgGRLKsR4pucbJtbKqIkUGxuZI2t7pfewKRc5nWecvDBZf3+p1M
pA8CAwEAAaOCAU8wggFLMB0GA1UdDgQWBBRVdE+yck/1YLpQ0dfmUVyaAYca1zAf
BgNVHSMEGDAWgBQD3lA1VtFMu2bwo+IbG8OXsj3RVTAOBgNVHQ8BAf8EBAMCAYYw
HQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQIMAYBAf8C
AQAwNAYIKwYBBQUHAQEEKDAmMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdp
Y2VydC5jb20wQgYDVR0fBDswOTA3oDWgM4YxaHR0cDovL2NybDMuZGlnaWNlcnQu
Y29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNybDBMBgNVHSAERTBDMDcGCWCGSAGG
/WwBAjAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5jb20vQ1BT
MAgGBmeBDAECATANBgkqhkiG9w0BAQsFAAOCAQEAK3Gp6/aGq7aBZsxf/oQ+TD/B
SwW3AU4ETK+GQf2kFzYZkby5SFrHdPomunx2HBzViUchGoofGgg7gHW0W3MlQAXW
M0r5LUvStcr82QDWYNPaUy4taCQmyaJ+VB+6wxHstSigOlSNF2a6vg4rgexixeiV
4YSB03Yqp2t3TeZHM9ESfkus74nQyW7pRGezj+TC44xCagCQQOzzNmzEAP2SnCrJ
sNE2DpRVMnL8J6xBRdjmOsC3N6cQuKuRXbzByVBjCqAA8t1L0I+9wXJerLPyErjy
rMKWaBFLmfK/AHNF4ZihwPGOc7w6UHczBZXH5RFzJNnww+WnKuTPI0HfnVH8lg==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAlgOHGn/WXXLCHA23mzJLccpOA3aIrfBFdJDXcsMjRDhrpLmV
lOKV1Lh79JomQGVOQwUygs9Wwv9TW+RXPMcmrfDj19uAxnFhP13gwlbwhz//Cqxv
DJ3o9p59fCB+tNZgAzJtjQQK0LfMmgaIAm2tCV/XpVWWFrUMBfO9uMprxLesQ1kL
NblhRiz+FkmViXbjhw7oqXHBH98z6IQXzbCiL7jiAijWPPdyeBzeI2GcS0HR51Fw
Ma+lyuaSIs69cndiClHD2/5zbw+KLm6BSaddtVeesMYF57EcEHLRuNg7JaDKI5QC
hdK/l392FF8PFQ9lJAnPY1KntjZr/kHq0XiQ8QIDAQABAoIBACiPm6AWoKdzt/hN
3S8hUjTaNm3JRvuA08bIwvhMuuRfPPu1EjTHbyutFhb09xLCUX7dkOK9nP/seWWH
P+83CcZOM8zRlOgTD/BKOdNSHobzTspcBUqsB6lnARbm0lui+yLiJ6zRQvtcNv4O
dgfyD69RMsWJdqN9IFsbpFiqoqj1dqAwT9c44bNq6c0kY4lOuPnRLU2sfJDJJlOU
GtHoO1U5mP6kJ4OdjoxSYZnhsox12fXzHashxd6mNxn4KORK2urybaTokVGUCGqx
qWmzEDWNQb4GYBWiR9BDHnNULFLpx08dI9Hb9HNyPS9EEVGBlxmP7+2SpvePr++c
e3ZyN+ECgYEAyhciKtYy0mdv7MF1MXCFVFJLkaPxChtYcoQerrKAoFkxYl/reY2p
CRIESMUabnopkuWBeTJ5CCoV47gbmHBQLdr43hRd/cpGVE+oiA679+mktEeCou5f
xScsOQsfJxUkPVJz9vmVbtmLvvf74Z0Cx7YWOtMmHVsZ57fACg7bKfcCgYEAvggN
9lbY1afjmDUSj5ZPfUi5/bK5i7da2Xqqx3HTPWVs04nHYtBJpFCzfBx3xbWyj/h+
M1pLhTxSKJPa+hDG8NYLrzzYp6gtCfhJF7GND3OquFHMt7+cgCKPFLazMGTEwLKM
hzpqTHlvQaoQpe8M1lzcNmevbAsJdDh1d4tJolcCgYEAxqx5dZ2A9yKjgSErgoA5
Q41oJ3UBqcr6aBKFS3/HPlyRVUIxcB2ZWYZx2cyUUJoetwCUCb9aB3HAdU/xKSr5
WCtW0JU7Vh5+h7KMX74EgxQaTPWkc2NfmaYKLsZFSRnat8KQqPPzObf7T7Hh2YqP
SiEzt38PkHqYfBpEXF8AjT8CgYB9lzMrGFCsPA4l/QVsUknsohESA3mvRhnL289c
ivSyAgM/dzKIMuJIr3E/2EysJR6DGhbF96orvycJXFZ/qHDioIQOZ6dEfthtW2Nr
PlPc33P350/mLMPQx4ZKiUi59g82z4oioU+5hRQrkKr6D5grYCnF5xa/0DeKUPoJ
bMvYdwKBgQCCp2hKvOsJ7k2Yy5hC7jqksoDXW4TEW1VMTBCvnJ9IKBNcgGhvNhZP
t9s3qsJd0RvW6aHXkLEt3Znti9Nki5vBsppP2QKc5vmkmTnRd/e2dqZeL7Orua6i
ENPkagd6EvjAMsbne1dTgIe7R2yh3pdRNOiLVy2uyVJaz8Qn8at/8Q==
-----END RSA PRIVATE KEY-----
EOF
}
`, randInt)
}
