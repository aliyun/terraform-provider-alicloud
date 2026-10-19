package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// collection-order: bypass client_ca_certificate_hostnames -- pre-existing attribute outside this change; this PR only adds Computed serial_number/fingerprint_sha256 (excluded candidates), hostnames TypeList and its basic12003 membership steps predate this change.

// Test Esa ClientCaCertificate. >>> Resource test cases, automatically generated.
// Case resource_ClientCaCertificate_set_test 12003
func TestAccAliCloudEsaClientCaCertificate_basic12003(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_ca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaClientCaCertificateMap12003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaClientCaCertificateBasicDependence12003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"certificate": testEsaClientCaCertificate12003,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":            CHECKSET,
						"serial_number":      CHECKSET,
						"fingerprint_sha256": CHECKSET,
						"certificate":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": []string{
						"test1.${data.alicloud_esa_sites.default.sites.0.site_name}", "test2.${data.alicloud_esa_sites.default.sites.0.site_name}", "test3.${data.alicloud_esa_sites.default.sites.0.site_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": []string{
						"test3.${data.alicloud_esa_sites.default.sites.0.site_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_certificate_hostnames": []string{
						"test1.${data.alicloud_esa_sites.default.sites.0.site_name}", "test2.${data.alicloud_esa_sites.default.sites.0.site_name}", "test3.${data.alicloud_esa_sites.default.sites.0.site_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_certificate_hostnames.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEsaClientCaCertificate_basic12003_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_ca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaClientCaCertificateMap12003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaClientCaCertificateBasicDependence12003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ca_cert_name": name,
					"site_id":             "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"client_ca_certificate_hostnames": []string{
						"test1.${data.alicloud_esa_sites.default.sites.0.site_name}", "test2.${data.alicloud_esa_sites.default.sites.0.site_name}", "test3.${data.alicloud_esa_sites.default.sites.0.site_name}"},
					"certificate": testEsaClientCaCertificate12003,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ca_cert_name":               name,
						"site_id":                           CHECKSET,
						"client_ca_certificate_hostnames.#": "3",
						"certificate":                       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudEsaClientCaCertificateMap12003 = map[string]string{
	"status":             CHECKSET,
	"create_time":        CHECKSET,
	"client_ca_cert_id":  CHECKSET,
	"serial_number":      CHECKSET,
	"fingerprint_sha256": CHECKSET,
}

func AliCloudEsaClientCaCertificateBasicDependence12003(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}
`, name)
}

// Test Esa ClientCaCertificate. <<< Resource test cases, automatically generated.

// generated for basic12003: CN=example-test123.com, SAN=example-test123.com,*.example-test123.com, CA:TRUE, valid 2026-2036.
const testEsaClientCaCertificate12003 = `-----BEGIN CERTIFICATE-----\nMIIDezCCAmOgAwIBAgIUGZsjd5vsP2z+dRKpngDEGaEzkekwDQYJKoZIhvcNAQEL\nBQAwOjELMAkGA1UEBhMCQ04xDTALBgNVBAoMBDExMjMxHDAaBgNVBAMME2V4YW1w\nbGUtdGVzdDEyMy5jb20wHhcNMjYxMDE4MDk1NjQ3WhcNMzYxMDE1MDk1NjQ3WjA6\nMQswCQYDVQQGEwJDTjENMAsGA1UECgwEMTEyMzEcMBoGA1UEAwwTZXhhbXBsZS10\nZXN0MTIzLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALPHT3xx\nPudUQEgvEcOfqLWxonUoy/klAE8Cimxf7UDnwoCXY2qruK6IuHNC1CBQMWKRuY1H\nzEQNghdErCFn03oKRTCXh6nhiMIuXR7R4xDUZYmEVZCi5vLNVBLBpQkvmqMecYN8\n3nSrQFoPGGSK16shLB8Ra8p5b1wIrmLI3/qtVAt5mAw9vSc5gqXvhX47U7zU3uv0\nMUMWGE+s8fpHLD3gvfJP2aD6v+d9X+btFapv22med2ORxfSba7gYilK7Uk10lqWI\n5kPLuavo1emnjig2JqdaKx7sB2ulmIdPKBUwPoPa3/OlYDGb91B5qHGjW9xOFoJg\nB5iawd4qsoZjMpcCAwEAAaN5MHcwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8E\nBAMCAYYwHQYDVR0OBBYEFOubsY2It6LeSZ64Ml42bKLRGlTbMDUGA1UdEQQuMCyC\nE2V4YW1wbGUtdGVzdDEyMy5jb22CFSouZXhhbXBsZS10ZXN0MTIzLmNvbTANBgkq\nhkiG9w0BAQsFAAOCAQEAeaDA0N+8T1uBOj/fwOg/quZI8VWvN/s/gX7bcsCyNjam\nXZkW2rH7RuTQeihqaiaIk68DtxSlJdUUkylszxR9x1dzqoN2DQmzMJ4CXbalZONS\nqDy5kR0kXU+vWL5EqrJlThIzvqM9yUF0fCVCG2/y5zzRi1SnCu300+JAy+k91ZYs\n5GfiS2PnXr4/s0ixwp3zztLRLOjgT5ULXQ/bOG4Cg8KoiRBmot0vkXo795y+2WPa\nkkgJROmjaWtrmhSCWUKykxHY4X7Aeinx2vMEM45RP2htq6SWT+dMZmJDtw4evqsz\nhZ960U5lm7qQYqu6L5VpcQs2jFzPWndpVG/8dPbWeQ==\n-----END CERTIFICATE-----`

// generated by https://www.ssleye.com/ssltool/self_sign.html
const testEsaClientCaCertificate = `-----BEGIN CERTIFICATE-----\nMIIDLDCCAhSgAwIBAgIUbW7alHwt3HICHGVsLkqpW9Cd4oIwDQYJKoZIhvcNAQEL\nBQAwLjEQMA4GA1UEAwwHdGVzdC1jYTENMAsGA1UECgwEdGVzdDELMAkGA1UEBhMC\nQ04wHhcNMjYwNjAxMDMyMzEzWhcNMzYwNTI5MDMyMzEzWjAuMRAwDgYDVQQDDAd0\nZXN0LWNhMQ0wCwYDVQQKDAR0ZXN0MQswCQYDVQQGEwJDTjCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAJvUEYwZodySRvwKaHVFgJsGe85tXYDzcHMlBS/7\n904+hVg/zPn3W83vlGHg0S9TUjDDYz+mx9XtGqbVTLNsnAFw6jnPvg6LoK3BE6fN\nT/muS4tkVN7zfKekC6RIPqpDoZEAVfXNzIdg31WdRhfDsUZ7ZgseTsh/o+0ophq/\nz8MwlLd0PYJJNQx+lvmBsBaqXKZ15+TXN4WwZ4Y0euXrX5uCd/eU+KrtsultZb/5\nF1yGR7ON2QCeYe0F1Csi+mdfuniipBtt4evtGRFRkTewzLPSoo8R/IokO50vtfl7\nNTTlwVw+N777PJZF5j0XdjGWtlxL98XpxkugpXSzsWUAewMCAwEAAaNCMEAwDwYD\nVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0OBBYEFIPFTEjW2zDy\nKmmsULjCisAJyu9QMA0GCSqGSIb3DQEBCwUAA4IBAQCSbwv8pOa1Tmv3eUCju4lW\n09JIYAE+H6n3tYRCwao0UdJ69c8yyNVVzbYf/3IbYBPs73OwW8sk3kKlZAHaPcUF\nFsCixLskgWKrJyCcK+sCt0c8JC+GeBHmJcncaq0E7XzxNB5R0KY2w6S7KCjeK1G9\niNsxvAVqNpISuHByd1LNKOT5S1niUx9Yd9eQTCcXEzZnE8Y2tf+gkUPjG9+5kxKL\nXHT3Choi7HjUYmwduKpV7wOAJdlBuyHegBii2xzRhdSriqffbqfczy9D5XSC6Q+w\nyF2ugGWW1IXF10DkMH3ixzbdGTohIUHpRZea3ELKhshMEo6ZqIhCLG3dv1kyeRG+\n-----END CERTIFICATE-----`
