package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbCACertificatesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbCACertificatesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_ca_certificates.slb_ca_certificates"),
					resource.TestCheckResourceAttr("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.name", "tf-testAccSlbCACertificatesDataSourceBasic"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.fingerprint"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.common_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.expired_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.expired_timestamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.created_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.created_timestamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.resource_group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_ca_certificates.slb_ca_certificates", "certificates.0.region_id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbCACertificatesDataSourceBasic = `
variable "name" {
	default = "tf-testAccSlbCACertificatesDataSourceBasic"
}


resource "alicloud_slb_ca_certificate" "foo" {
  name = "${var.name}"
  ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}

data "alicloud_slb_ca_certificates" "slb_ca_certificates" {
  ids = ["${alicloud_slb_ca_certificate.foo.id}"]
  name_regex = "${var.name}"
}
`
