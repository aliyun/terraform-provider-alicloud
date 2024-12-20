package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudWafv3Domain_basic2308(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DomainMap2308)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandString(10)
	// once one domain has been set, it will not be set again for the wafv3 instance
	name := fmt.Sprintf("tftest%s.tftest.top", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DomainBasicDependence2308)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckForCleanUpInstances(t, string(connectivity.APSouthEast1), "waf", "waf", "waf", "waf")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_wafv3_instance.default.id}",
					"listen": []map[string]interface{}{
						{
							"https_ports":         []string{"443"},
							"http_ports":          []string{"80"},
							"cert_id":             "${local.certificate_id}",
							"cipher_suite":        "99",
							"xff_header_mode":     "2",
							"protection_resource": "share",
							"tls_version":         "tlsv1",
							"enable_tlsv3":        "true",
							"http2_enabled":       "true",
							"custom_ciphers":      []string{"ECDHE-ECDSA-AES128-GCM-SHA256"},
							"focus_https":         "false",
							"exclusive_ip":        "false",
							"xff_headers":         []string{"A", "B"},
							"ipv6_enabled":        "false",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends":           []string{"1.1.1.1"},
							"loadbalance":        "iphash",
							"sni_enabled":        "true",
							"sni_host":           "www.aliyun.com",
							"focus_http_backend": "false",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "1000",
							"keepalive_timeout":  "15",
							"connect_timeout":    "5",
							"read_timeout":       "30",
							"write_timeout":      "30",
							"request_headers": []map[string]interface{}{
								{
									"key":   "A",
									"value": "B",
								},
							},
						},
					},
					"domain":      "${var.name}",
					"access_type": "share",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":                   CHECKSET,
						"domain":                        name,
						"status":                        CHECKSET,
						"access_type":                   "share",
						"listen.#":                      "1",
						"listen.0.https_ports.#":        "1",
						"listen.0.http_ports.#":         "1",
						"listen.0.cert_id":              CHECKSET,
						"listen.0.cipher_suite":         "99",
						"listen.0.xff_header_mode":      "2",
						"listen.0.protection_resource":  "share",
						"listen.0.tls_version":          "tlsv1",
						"listen.0.enable_tlsv3":         "true",
						"listen.0.http2_enabled":        "true",
						"listen.0.custom_ciphers.#":     "1",
						"listen.0.focus_https":          "false",
						"listen.0.exclusive_ip":         "false",
						"listen.0.xff_headers.#":        "2",
						"listen.0.ipv6_enabled":         "false",
						"redirect.#":                    "1",
						"redirect.0.backends.#":         "1",
						"redirect.0.loadbalance":        "iphash",
						"redirect.0.sni_enabled":        "true",
						"redirect.0.sni_host":           "www.aliyun.com",
						"redirect.0.focus_http_backend": "false",
						"redirect.0.keepalive":          "true",
						"redirect.0.retry":              "true",
						"redirect.0.keepalive_requests": "1000",
						"redirect.0.keepalive_timeout":  "15",
						"redirect.0.connect_timeout":    "5",
						"redirect.0.read_timeout":       "30",
						"redirect.0.write_timeout":      "30",
						"redirect.0.request_headers.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listen": []map[string]interface{}{
						{
							"https_ports":         []string{"443"},
							"http_ports":          []string{"80", "88"},
							"cert_id":             "${local.certificate_id}",
							"cipher_suite":        "99",
							"xff_header_mode":     "2",
							"protection_resource": "share",
							"tls_version":         "tlsv1",
							"enable_tlsv3":        "true",
							"http2_enabled":       "true",
							"custom_ciphers":      []string{"ECDHE-ECDSA-AES128-GCM-SHA256"},
							"focus_https":         "false",
							"exclusive_ip":        "false",
							"xff_headers":         []string{"A", "B", "C"},
							"ipv6_enabled":        "false",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends":           []string{"1.1.1.1"},
							"loadbalance":        "iphash",
							"sni_enabled":        "true",
							"sni_host":           "www.aliyun.com",
							"focus_http_backend": "false",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "1000",
							"keepalive_timeout":  "15",
							"connect_timeout":    "5",
							"read_timeout":       "60",
							"write_timeout":      "60",
							"request_headers": []map[string]interface{}{
								{
									"key":   "A",
									"value": "B",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listen.#":                      "1",
						"listen.0.https_ports.#":        "1",
						"listen.0.http_ports.#":         "2",
						"listen.0.cert_id":              CHECKSET,
						"listen.0.cipher_suite":         "99",
						"listen.0.xff_header_mode":      "2",
						"listen.0.protection_resource":  "share",
						"listen.0.tls_version":          "tlsv1",
						"listen.0.enable_tlsv3":         "true",
						"listen.0.http2_enabled":        "true",
						"listen.0.custom_ciphers.#":     "1",
						"listen.0.focus_https":          "false",
						"listen.0.exclusive_ip":         "false",
						"listen.0.xff_headers.#":        "3",
						"listen.0.ipv6_enabled":         "false",
						"redirect.#":                    "1",
						"redirect.0.backends.#":         "1",
						"redirect.0.loadbalance":        "iphash",
						"redirect.0.sni_enabled":        "true",
						"redirect.0.sni_host":           "www.aliyun.com",
						"redirect.0.focus_http_backend": "false",
						"redirect.0.keepalive":          "true",
						"redirect.0.retry":              "true",
						"redirect.0.keepalive_requests": "1000",
						"redirect.0.keepalive_timeout":  "15",
						"redirect.0.connect_timeout":    "5",
						"redirect.0.read_timeout":       "60",
						"redirect.0.write_timeout":      "60",
						"redirect.0.request_headers.#":  "1",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_type"},
			},
		},
	})
}

var AlicloudWafv3DomainMap2308 = map[string]string{}

func AlicloudWafv3DomainBasicDependence2308(name string) string {
	casRegion := "cn-hangzhou"
	if strings.ToLower(os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")) == "international" {
		casRegion = "ap-southeast-1"
	}
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
// certificate can be generated in https://zh.rakko.tools/tools/46/ and the common name should be .tftest.top
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIIDeDCCAmCgAwIBAgIEN3ZT6zANBgkqhkiG9w0BAQsFADBVMQswCQYDVQQGEwJD
TjEVMBMGA1UEAwwMKi50ZnRlc3QudG9wMRAwDgYDVQQIDAdCZWlKaW5nMRAwDgYD
VQQHDAdCZWlKaW5nMQswCQYDVQQKDAJBQTAeFw0yMzA4MjgwNjQ5NDNaFw0yNTA4
MjcwNjQ5NDNaMFUxCzAJBgNVBAYTAkNOMRUwEwYDVQQDDAwqLnRmdGVzdC50b3Ax
EDAOBgNVBAgMB0JlaUppbmcxEDAOBgNVBAcMB0JlaUppbmcxCzAJBgNVBAoMAkFB
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzkk9NJUH7PLSQK4RRrGQ
Y5dVsftkhnKh9HhI6yrnlowWIDPS1PZHOU/5gQ7xPUPGdKQV5S7x8wROnAaXEimx
N4GdQw25pGhRJvlwme9fzJJiSe6lG49NCxmuBiEdJAzPKaTPpK1cG1f1TqdgCfHR
HAL6Jxb3ylHG2LlTNFLXikubUi5RT6/9C8psr713Zm4HveCI/cx0WdgZ+fmsc9ft
rkIB1DdyV1kQ51m8r2rLi3J7aC5ggGOiex/VlGSd4e6SOQLpdQEdDbodtOJ4LgVk
+arFNCMinUWIOPGFzXhdm6lssPbh4MOBrz8c/M9TcF4hoMn5/o/9johZIZ/DOvXt
ZQIDAQABo1AwTjAdBgNVHQ4EFgQUOnWiddgeZj17IeysatqhE361o5YwHwYDVR0j
BBgwFoAUOnWiddgeZj17IeysatqhE361o5YwDAYDVR0TBAUwAwEB/zANBgkqhkiG
9w0BAQsFAAOCAQEAfh3cnOszHM/5wXjY7BIkmgDOReksS+87ibhBz7T2ddZj+yCF
9GdIBzXCiHpQFXpW8a3kc3I7l3nGfMTkmF6ld3ot/6SXP17QKJwxtvUA4ib8QkWD
S7FT+UcHCUHv42Sh1e5uAlQ5pMSul7iKcR7jwlwZGZ0905HOqrmdyUGJ+Ud2uZWD
AC0dJF6Bv9VhNtci8Imp05PaPH09deXLZu8LRrKRZFy9qLW5R6Swv7nzxckOAqDk
TTc40xwvQROekWUyxeJL7xaHuylUHE0bxsiIfx5bZsBizRjprIwGlj85CSPuTZyP
DPfaiZAN/61h5HNAnxLltOZfqabKYYw7l9LBDg==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDOST00lQfs8tJA
rhFGsZBjl1Wx+2SGcqH0eEjrKueWjBYgM9LU9kc5T/mBDvE9Q8Z0pBXlLvHzBE6c
BpcSKbE3gZ1DDbmkaFEm+XCZ71/MkmJJ7qUbj00LGa4GIR0kDM8ppM+krVwbV/VO
p2AJ8dEcAvonFvfKUcbYuVM0UteKS5tSLlFPr/0LymyvvXdmbge94Ij9zHRZ2Bn5
+axz1+2uQgHUN3JXWRDnWbyvasuLcntoLmCAY6J7H9WUZJ3h7pI5Aul1AR0Nuh20
4nguBWT5qsU0IyKdRYg48YXNeF2bqWyw9uHgw4GvPxz8z1NwXiGgyfn+j/2OiFkh
n8M69e1lAgMBAAECggEAevPgTTT+0lYwx2h416ACJboP09O5KQGuUl5XaAPcoTjB
/1OkOFbKQPjQCAJ1+0QoR2F9w2plv6kziX/MD4FWJXVV3J+TpNCgfhBy8u1gNjiR
6Osa8gBJtXIK7ZBTJCeWWoXnVYoWuh2FEupkLck6D+4eV6oy6x4u3QIo+6jc24n9
dIXQG6/v/Iao34kB0LUdp/4WNaUDvfI6NDhEwchpKE95dtWIDlIN/YhfiYAdjrnl
YmH2VDbAGgsdEiHP4wLZfjgsGPPDGS0+qBHoSiJGH0E6wWEZdAE4TsYGRFsO86n3
LfjEPFGfPlcnZe2cTTe3kmyKq/DTjxtu2rh3I8o2CQKBgQD/5Xe7cenaOBefzPlx
GOEsB+qv49UmzEPOXDNZe9hmAawuuuxPUM+xlE++P+mEgQm1LPT4WWgtFLPVuwmx
ncxt4CJNZh+ZGFyAZ4dm4M4ZhIBXNonyIP+yGyAJUUVF9Iy3TYcJNiGzv2Rx9JRQ
XWJMQnTDILmZbmU+ltTea7/zqwKBgQDOXqCqb17MuLt7OcKWSgthm79OlaOdzZvl
i9qU6VzZKG7Axc5gA9yq6tHp3vWPI4bNdvwqIIa/nzVILjGA5fcYFbRN+7gHwo8s
rNAgi5PAoKWqQRovyJRY9Eq/sn6l1jbJZAOUAMZMWDm8z89OqK7PNQSIAtfFSneo
2QxJkGeTLwKBgGJkafBB8af9b1/7YWISLepPNPbihH/BhMThAMGEdAVs2TaymtA4
g1OFck/1pSVUtFXcbmjbf8ntruQcYbLQuNz6lFXsUXP9QPwCUrbE85ouL2bZSps2
AvsJoPzUKe2nBUAp6CUrkjPaAJYsc6ae8X/fAaRRfeu33ef9+OV4yrq3AoGAYFZo
ZmfrN2Kdkt7Z6dLTEVPlsMfGQ6pyNmxdM9rkzzNC0JcGymfDIb7RE35T3+hTy6La
AMiCXv3xn6qAzY2NFh87tpPlyymWMOLTnf3Kkcfszlfp45idOBGCu46V9NDVbppT
2UmrSIR/H5dbTXsNcAlt/hhlpeInjhkU1VqmH10CgYEA7Kk+QhWq705SczpWjm5J
9kHqfFzJLwAWNBduiia0WypgPhLe/4wT1rYQkBtKMVKrgFo7Cvi4YKlrtlDnXyeU
CIFqfEL5NriQelqrFsvgHsmD+MpvDoSWm5C8IrTubtlNyWUzXSVT4OIwzPobzPqG
LILJ+e7bLw8RrM0HfgFnl8c=
-----END PRIVATE KEY-----
EOF
}

resource "alicloud_wafv3_instance" "default" {}

locals {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
}


`, name, casRegion)
}
