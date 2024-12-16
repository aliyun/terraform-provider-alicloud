package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudWafv3DomainDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	// cn region needs icp domain, and using intl region instead
	checkoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_domain.default.id}_fake"]`,
		}),
	}
	DomainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"domain": `"${alicloud_wafv3_domain.default.domain}"`,
		}),
		fakeConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"domain": `"${alicloud_wafv3_domain.default.domain}_fake"`,
		}),
	}
	backendConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"backend": `"1.1.1.1"`,
		}),
		fakeConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"backend": `"1.1.1.2"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_wafv3_domain.default.id}"]`,
			"domain":  `"${alicloud_wafv3_domain.default.domain}"`,
			"backend": `"1.1.1.1"`,
		}),
		fakeConfig: testAccCheckAliCloudWafv3DomainSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_wafv3_domain.default.id}_fake"]`,
			"domain":  `"${alicloud_wafv3_domain.default.domain}_fake"`,
			"backend": `"1.1.1.2"`,
		}),
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckForCleanUpInstances(t, string(connectivity.APSouthEast1), "waf", "waf", "waf", "waf")
	}

	Wafv3DomainCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, DomainConf, backendConf, allConf)
}

var existWafv3DomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":            "1",
		"domains.#":        "1",
		"domains.0.id":     CHECKSET,
		"domains.0.domain": CHECKSET,
		"domains.0.resource_manager_resource_group_id": CHECKSET,
		"domains.0.status":                             CHECKSET,
		"domains.0.listen.#":                           "1",
		"domains.0.listen.0.cert_id":                   CHECKSET,
		"domains.0.listen.0.cipher_suite":              "99",
		"domains.0.listen.0.enable_tlsv3":              "true",
		"domains.0.listen.0.exclusive_ip":              "false",
		"domains.0.listen.0.focus_https":               "false",
		"domains.0.listen.0.http2_enabled":             "true",
		"domains.0.listen.0.ipv6_enabled":              "false",
		"domains.0.listen.0.protection_resource":       "share",
		"domains.0.listen.0.tls_version":               "tlsv1",
		"domains.0.listen.0.xff_header_mode":           "2",
		"domains.0.listen.0.http_ports.#":              "1",
		"domains.0.listen.0.https_ports.#":             "1",
		"domains.0.listen.0.custom_ciphers.#":          "1",
		"domains.0.listen.0.xff_headers.#":             "2",
		"domains.0.redirect.#":                         "1",
		"domains.0.redirect.0.loadbalance":             "iphash",
		"domains.0.redirect.0.focus_http_backend":      "false",
		"domains.0.redirect.0.keepalive":               "true",
		"domains.0.redirect.0.keepalive_requests":      "80",
		"domains.0.redirect.0.retry":                   "true",
		"domains.0.redirect.0.sni_enabled":             "true",
		"domains.0.redirect.0.sni_host":                "www.aliyun.com",
		"domains.0.redirect.0.connect_timeout":         "30",
		"domains.0.redirect.0.keepalive_timeout":       "30",
		"domains.0.redirect.0.read_timeout":            "30",
		"domains.0.redirect.0.write_timeout":           "30",
		"domains.0.redirect.0.backends.#":              "1",
		"domains.0.redirect.0.request_headers.#":       "1",
		"domains.0.redirect.0.request_headers.0.key":   CHECKSET,
		"domains.0.redirect.0.request_headers.0.value": CHECKSET,
	}
}

var fakeWafv3DomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"domains.#": "0",
	}
}

var Wafv3DomainCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_wafv3_domains.default",
	existMapFunc: existWafv3DomainMapFunc,
	fakeMapFunc:  fakeWafv3DomainMapFunc,
}

func testAccCheckAliCloudWafv3DomainSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	casRegion := "cn-hangzhou"
	if strings.ToLower(os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")) == "international" {
		casRegion = "ap-southeast-1"
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tftest%d.tftest.top"
	}

	resource "alicloud_wafv3_instance" "default" {
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
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
  		key              = <<EOF
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

	resource "alicloud_wafv3_domain" "default" {
  		instance_id = alicloud_wafv3_instance.default.id
  		domain      = var.name
  		access_type = "share"
  		listen {
    		http_ports          = [80]
    		https_ports         = [443]
    		cert_id             = local.certificate_id
    		cipher_suite        = 99
    		xff_header_mode     = 2
    		protection_resource = "share"
			tls_version         = "tlsv1"
			enable_tlsv3        = true
			http2_enabled       = true
			custom_ciphers      = ["ECDHE-ECDSA-AES128-GCM-SHA256"]
			focus_https         = false
			ipv6_enabled        = false
			exclusive_ip        = false
			xff_headers         = ["A", "B"]
  		}
  		redirect {
			backends           = ["1.1.1.1"]
			loadbalance        = "iphash"
			sni_enabled        = true
			sni_host           = "www.aliyun.com"
			focus_http_backend = false
			keepalive          = true
			retry              = true
			keepalive_requests = 80
			keepalive_timeout  = 30
			connect_timeout    = 30
			read_timeout       = 30
			write_timeout      = 30
    		request_headers {
      			key   = "A"
      			value = "B"
    		}
  		}

	}

	locals {
  		certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
	}

	data "alicloud_wafv3_domains" "default" {
  		instance_id    = alicloud_wafv3_instance.default.id
  		enable_details = true
		%s
	}
`, rand, casRegion, strings.Join(pairs, "\n   "))
	return config
}
