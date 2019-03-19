package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCasCertificatesDataSource_certificates(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CasClassicSupportedRegions)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCasDataSourceCertificates(randInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.common"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.finger_print"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.issuer"),
					resource.TestCheckResourceAttr("data.alicloud_cas_certificates.certs", "certificates.0.org_name", ""),
					resource.TestCheckResourceAttr("data.alicloud_cas_certificates.certs", "certificates.0.province", ""),
					resource.TestCheckResourceAttr("data.alicloud_cas_certificates.certs", "certificates.0.city", ""),
					resource.TestCheckResourceAttr("data.alicloud_cas_certificates.certs", "certificates.0.country", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.start_date"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.end_date"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.sans"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.expired"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "certificates.0.buy_in_aliyun"),
				),
			},
		},
	})
}

func testAccCheckAlicloudCasDataSourceCertificates(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_cas_certificate" "cert" {
  name = "tf_testAcc%d"
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
data "alicloud_cas_certificates" "certs" {
  name_regex = "${alicloud_cas_certificate.cert.name}"
  output_file = "${path.module}/cas_certificates.json"
}
`, randInt)
}
