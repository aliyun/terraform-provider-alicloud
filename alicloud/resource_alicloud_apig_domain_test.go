// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Apig Domain. >>> Resource test cases, automatically generated.
// Case domain_basic_test 12907
func TestAccAliCloudApigDomain_basic12907(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudApigDomainMap12907)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf.apig%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigDomainBasicDependence12907)
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"domain_name":       name,
					"gateway_type":      "API",
					"protocol":          "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"domain_name":       name,
						"gateway_type":      "API",
						"protocol":          "HTTP",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gateway_type"},
			},
		},
	})
}

var AlicloudApigDomainMap12907 = map[string]string{}

func AlicloudApigDomainBasicDependence12907(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)
}

// Case domain_https_full_update 12906
func TestAccAliCloudApigDomain_basic12906(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudApigDomainMap12906)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf.apig%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigDomainBasicDependence12906)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Step 0: create an HTTP domain
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"domain_name":       name,
					"gateway_type":      "API",
					"domain_scope":      "Dedicated",
					"protocol":          "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"domain_name":       name,
						"gateway_type":      "API",
						"domain_scope":      "Dedicated",
						"protocol":          "HTTP",
					}),
				),
			},
			{
				// Step 1: switch to HTTPS with a certificate and TLS 1.2-1.3 cipher config
				Config: testAccConfig(map[string]interface{}{
					"protocol":        "HTTPS",
					"cert_identifier": "${alicloud_ssl_certificates_service_certificate.default.id}-cn-hangzhou",
					"force_https":     "true",
					"http2_option":    "Open",
					"tls_min":         "TLS 1.2",
					"tls_max":         "TLS 1.3",
					"tls_cipher_suites_config": []map[string]interface{}{
						{
							"config_type": "Custom",
							"tls_cipher_suite": []map[string]interface{}{
								{
									"name":             "ECDHE-RSA-AES128-GCM-SHA256",
									"support_versions": []string{"TLS 1.2", "TLS 1.3"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":                               "HTTPS",
						"cert_identifier":                        CHECKSET,
						"force_https":                            "true",
						"http2_option":                           "Open",
						"tls_min":                                "TLS 1.2",
						"tls_max":                                "TLS 1.3",
						"tls_cipher_suites_config.#":             "1",
						"tls_cipher_suites_config.0.config_type": "Custom",
						"tls_cipher_suites_config.0.tls_cipher_suite.#":      "1",
						"tls_cipher_suites_config.0.tls_cipher_suite.0.name": "ECDHE-RSA-AES128-GCM-SHA256",
					}),
				),
			},
			{
				// Step 2: raise TLS to 1.3, enable mTLS, switch off forced HTTPS/HTTP2
				Config: testAccConfig(map[string]interface{}{
					"force_https":        "false",
					"http2_option":       "Close",
					"tls_min":            "TLS 1.3",
					"tls_max":            "TLS 1.3",
					"m_tls_enabled":      "true",
					"client_ca_cert":     "-----BEGIN CERTIFICATE-----\\nMIIDIzCCAgugAwIBAgIUDgC59vYBh6W9wT3Wy/W2fDvKQEMwDQYJKoZIhvcNAQEL\\nBQAwITEfMB0GA1UEAwwWdGYtdGVzdGFjYy5leGFtcGxlLmNvbTAeFw0yNjA3MDEw\\nOTAzMDVaFw0zNjA2MjgwOTAzMDVaMCExHzAdBgNVBAMMFnRmLXRlc3RhY2MuZXhh\\nbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCV/puAPJBo\\ny6aIYj3bnKKobycrVLKlN90z06ArcC62bSQPX81Ab8D428uvmyYGlYlMG66DrpcK\\nnHlgsSoXBK17cl/o1tsrUapI7vTZOjLkfia4Rq3N91hgUfm6e41ckWoRUCwTsSC/\\n+iSWDvFsUqpzjSEKM1eqaW7p5LtzmaUNcWPpiyuhFK2uoCi1anacBB/8HaMJM1KJ\\nMLpbRwhE2Olas1ASkpVqczZS0OZ3h7Rgj9nP0t+YIxiul4hujhbgTI6MFBWrjdSt\\nfS5g1eOwFVuzZLKqBwjj2ASITXtrGndaPOutNoc7SB09GXaNkISLEeITkSjceGNr\\nw0KjZ6gLNYvpAgMBAAGjUzBRMB0GA1UdDgQWBBTqnyl6uqUSFizgnladn9aAjzbP\\nlTAfBgNVHSMEGDAWgBTqnyl6uqUSFizgnladn9aAjzbPlTAPBgNVHRMBAf8EBTAD\\nAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAMgUoy0GgIa6vZA9VAFCfhoqPxpw7Grosh\\nQGw0BZSE4AoEi/iljCWwjhKdrwr97guaZeDXYw5NIeFwfLbVuo0qK56VQb/td2RX\\narG5INSH7sVVZJruMgQBBOSYe3vr4bB+6/TM/FAQWEbejLbz95impgUhYgMTr6b5\\n+jcCY545RUsYnkH+qe6qCdi7l1WR+QWgxrrpdcuWDipFY5/81D6mx8crykq/X2KO\\nhWNin2Y7nBT3vwOFdBdeCO6/njmP+lJ3wG5A+1/5+UYwm2oOElt8Ib/jvqCuRkEi\\nw0NVDKP009+0XdKOcWl5bGawyg06Y4InvNLu/Q9U1TsilOjmYZkG\\n-----END CERTIFICATE-----",
					"ca_cert_identifier": "${alicloud_ssl_certificates_service_certificate.default.id}-cn-hangzhou",
					"tls_cipher_suites_config": []map[string]interface{}{
						{
							"config_type": "Custom",
							"tls_cipher_suite": []map[string]interface{}{
								{
									"name":             "ECDHE-RSA-AES256-GCM-SHA384",
									"support_versions": []string{"TLS 1.3"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_https":                            "false",
						"http2_option":                           "Close",
						"tls_min":                                "TLS 1.3",
						"tls_max":                                "TLS 1.3",
						"m_tls_enabled":                          "true",
						"client_ca_cert":                         CHECKSET,
						"ca_cert_identifier":                     CHECKSET,
						"tls_cipher_suites_config.0.config_type": "Custom",
						"tls_cipher_suites_config.0.tls_cipher_suite.0.name": "ECDHE-RSA-AES256-GCM-SHA384",
					}),
				),
			},
			{
				// Step 3: migrate to a different resource group (async ChangeResourceGroup)
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gateway_type"},
			},
		},
	})
}

var AlicloudApigDomainMap12906 = map[string]string{}

func AlicloudApigDomainBasicDependence12906(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIIDIzCCAgugAwIBAgIUDgC59vYBh6W9wT3Wy/W2fDvKQEMwDQYJKoZIhvcNAQEL
BQAwITEfMB0GA1UEAwwWdGYtdGVzdGFjYy5leGFtcGxlLmNvbTAeFw0yNjA3MDEw
OTAzMDVaFw0zNjA2MjgwOTAzMDVaMCExHzAdBgNVBAMMFnRmLXRlc3RhY2MuZXhh
bXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCV/puAPJBo
y6aIYj3bnKKobycrVLKlN90z06ArcC62bSQPX81Ab8D428uvmyYGlYlMG66DrpcK
nHlgsSoXBK17cl/o1tsrUapI7vTZOjLkfia4Rq3N91hgUfm6e41ckWoRUCwTsSC/
+iSWDvFsUqpzjSEKM1eqaW7p5LtzmaUNcWPpiyuhFK2uoCi1anacBB/8HaMJM1KJ
MLpbRwhE2Olas1ASkpVqczZS0OZ3h7Rgj9nP0t+YIxiul4hujhbgTI6MFBWrjdSt
fS5g1eOwFVuzZLKqBwjj2ASITXtrGndaPOutNoc7SB09GXaNkISLEeITkSjceGNr
w0KjZ6gLNYvpAgMBAAGjUzBRMB0GA1UdDgQWBBTqnyl6uqUSFizgnladn9aAjzbP
lTAfBgNVHSMEGDAWgBTqnyl6uqUSFizgnladn9aAjzbPlTAPBgNVHRMBAf8EBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAMgUoy0GgIa6vZA9VAFCfhoqPxpw7Grosh
QGw0BZSE4AoEi/iljCWwjhKdrwr97guaZeDXYw5NIeFwfLbVuo0qK56VQb/td2RX
arG5INSH7sVVZJruMgQBBOSYe3vr4bB+6/TM/FAQWEbejLbz95impgUhYgMTr6b5
+jcCY545RUsYnkH+qe6qCdi7l1WR+QWgxrrpdcuWDipFY5/81D6mx8crykq/X2KO
hWNin2Y7nBT3vwOFdBdeCO6/njmP+lJ3wG5A+1/5+UYwm2oOElt8Ib/jvqCuRkEi
w0NVDKP009+0XdKOcWl5bGawyg06Y4InvNLu/Q9U1TsilOjmYZkG
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAlf6bgDyQaMumiGI925yiqG8nK1SypTfdM9OgK3Autm0kD1/N
QG/A+NvLr5smBpWJTBuug66XCpx5YLEqFwSte3Jf6NbbK1GqSO702Toy5H4muEat
zfdYYFH5unuNXJFqEVAsE7Egv/oklg7xbFKqc40hCjNXqmlu6eS7c5mlDXFj6Ysr
oRStrqAotWp2nAQf/B2jCTNSiTC6W0cIRNjpWrNQEpKVanM2UtDmd4e0YI/Zz9Lf
mCMYrpeIbo4W4EyOjBQVq43UrX0uYNXjsBVbs2SyqgcI49gEiE17axp3WjzrrTaH
O0gdPRl2jZCEixHiE5Eo3Hhja8NCo2eoCzWL6QIDAQABAoIBAAVfH9yAzr8iA+3A
buyteFnF2UZA+0DVdlOD0amck9+umur+CFC1b9i5rlq0mLEFq+wQ1bgbiYc0wVgI
IDTA0yGnn+2rvB+aBhokjJo27lmmduaEiXbl08FnTiUyhYZ6Iq1KDLoLzttxLtw8
3sJ9V2NZ+4PtAMe2jOVNbrUeHH4Vsm/5nacRklp1EIrrF618nbG88+2dMscbPdmy
NItTv8QdNZ9Pwazs+fRoJJW0rzKyXZ+3bn+wmBw+hhZEPYHM1z5383oNW5v6rLI2
YKYoB3PIkfF0GsSw90JU/qcyGLf3Ccqg8Xy2ah482Ln2CqxZM9pf8iPRQvohBsXR
0Gs/5x0CgYEAxnaSJ05Hlnz6nhG5bUoQ87ZI17J9PkXoiStZIzEx9S/Qw7FPk+90
FSrH2YwrVaETCu3NfQYcwRJJedQbwWdsjsRGZLUn+GOeKOSdZMadtSXUuwpBPgG5
nyM6gH0Py7qOvCecps7+VGWTYqpczI+YS5XUMMqZbH6eWSbS7OvlIIcCgYEAwXrS
Ai0lSV2O9oRRT/ATxEkIbGGa4GszlOOSBvDvJssekY8INy4i61clMV34P9aMzqfg
tlO7ic2Jr2dBKmIKUMoujmKeZgOYfCvu4ox2DqGhOyOjvJ3d1C7WAjeV+BQzEOhz
y/RQ5JXBTz7ZxXcSKZIoPjsqr3B4eXuk9QjPPA8CgYBFWOExAtVY7ErWOPNGEP9j
aWqClEfXHq5mX9NBzMrcFd0oxCg+VQmG6+/xQF1UCniQ9Q88hIo/nJg4Dbm1FuKD
8Gl4fyR8UrLNLzUgJZat2Y4/3RF3DTtDNBgZFZoTYhjF/kFquCF+dA/QBh9vCy34
G16Nvf1mP8gs9rf1OWhSuQKBgBz85NgUoYCDdvbyTih25M9EzfFHEmhLR3goPGmz
0XDzf8n5LxbtX6f474ac+KO/5mrT9jP7CZ8U32sbQkUyWS9Pi3gjyG2qXj9Eac8h
klKQ3tI4fcC1ulWfCstcPqjjhd8jpK3LFg+ZbFQOK5yNQXhfAI6KWNPeOv6gis93
mWz7AoGASLaPPEZb0o8QSLpoKNwFNgBe/bu0165tJ7LLFdvdzaGB6TNGTC/L9CSd
sjzZjJlYcIuPims6QS0yYkHe5D0OtIPvHIHNa8JDp0aa8JzsAdJk7aoHj+sh5NU4
jd6R8k23Hryzrq9NMCdnFiB/38aHy8VYS5flruTWeHKqDmcHuTo=
-----END RSA PRIVATE KEY-----
EOF
}

`, name)
}

// Case domain_serverless_scope
// A Serverless-scope domain uses a managed certificate and does not accept an
// explicit protocol; it is HTTPS-only, driven by force_https.
func TestAccAliCloudApigDomain_serverless(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudApigDomainMapServerless)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf.apig%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigDomainBasicDependenceServerless)
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"domain_name":       name,
					"gateway_type":      "API",
					"domain_scope":      "Serverless",
					"force_https":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"domain_name":       name,
						"gateway_type":      "API",
						"domain_scope":      "Serverless",
						"force_https":       "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gateway_type"},
			},
		},
	})
}

var AlicloudApigDomainMapServerless = map[string]string{}

func AlicloudApigDomainBasicDependenceServerless(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)
}

// Test Apig Domain. <<< Resource test cases, automatically generated.
