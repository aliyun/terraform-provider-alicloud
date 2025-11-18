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
	ra := resourceAttrInit(resourceId, AliCloudWafv3DomainMap2308)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandString(10)
	// once one domain has been set, it will not be set again for the wafv3 instance
	name := fmt.Sprintf("tftest%s.tftest.top", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudWafv3DomainBasicDependence2308)
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
							"cipher_suite":        "1",
							"xff_header_mode":     "2",
							"protection_resource": "share",
							"tls_version":         "tlsv1",
							"enable_tlsv3":        "false",
							"http2_enabled":       "false",
							"custom_ciphers":      []string{},
							"focus_https":         "true",
							"exclusive_ip":        "true",
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
						"listen.0.cipher_suite":         "1",
						"listen.0.xff_header_mode":      "2",
						"listen.0.protection_resource":  "share",
						"listen.0.tls_version":          "tlsv1",
						"listen.0.enable_tlsv3":         "false",
						"listen.0.http2_enabled":        "false",
						"listen.0.custom_ciphers.#":     "0",
						"listen.0.focus_https":          "true",
						"listen.0.exclusive_ip":         "true",
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
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AliCloudWafv3DomainMap2308 = map[string]string{
	"status": CHECKSET,
}

func AliCloudWafv3DomainBasicDependence2308(name string) string {
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
MIIDojCCAoqgAwIBAgIEW02xdjANBgkqhkiG9w0BAQsFADBqMQswCQYDVQQGEwJD
TjEVMBMGA1UEAwwMKi50ZnRlc3QudG9wMQ8wDQYDVQQIDAbljJfkuqwxFTATBgNV
BAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UECgwTRGVmYXVsdCBDb21wYW55IEx0ZDAe
Fw0yNTExMTcxMDEwMDRaFw0yNzExMTcxMDEwMDRaMGoxCzAJBgNVBAYTAkNOMRUw
EwYDVQQDDAwqLnRmdGVzdC50b3AxDzANBgNVBAgMBuWMl+S6rDEVMBMGA1UEBwwM
RGVmYXVsdCBDaXR5MRwwGgYDVQQKDBNEZWZhdWx0IENvbXBhbnkgTHRkMIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyxPGODI9JxKzsB9cYkWSkY075aFD
tKvD8XdGztsOtd9IXuQqqPciTUbvaisYV/wcH9FhTea14nDC2IO05OX6EBDdrX2P
eHGaXAcahnfHzqc5Sd+JEfZMLdzhkKyIg3YOAC24wV5gqIjAuAfoMbLXuglRlSCe
pA7wBBn6zUCCBGzjwSB1RtmgZeSEf5q95ZYRZhU3mTWd158LlwImvOcVG5aJIlq8
Sf+Yj1Wr9hYD0s73M6x5YqgHtw+XwVVs+EJA+oepUdj74CejYWAYyD+OMA6KcLcY
N1Scg3VLs3mS9YhM3aZqrt04eR1ouLRuqfG/oN7kP0G5nQAov/rN25XmNQIDAQAB
o1AwTjAdBgNVHQ4EFgQUvJe0lF2jdEw4KpuaLIRIXtparo8wHwYDVR0jBBgwFoAU
vJe0lF2jdEw4KpuaLIRIXtparo8wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsF
AAOCAQEAcssSZQBjVdLb3USF2GRrKstRcaKqCBbno1QrtLc4WGDeOGzCjLq9I92l
0aGJjB8N4Aw1oTz3z9TFM0eTJbylNNz1U/PgHOhh6rkxq2pMPdT2bwStzA7Yw0oA
slmXMQkLQ4mT3ANDpXBm1MEOCAQcVS9LQzjpPslCLjXZzgPZd+RyEgToXPCh0gSY
cX30HpnPDYEccgShhWdNIuW3ptBm2h7sp5rbyLcIrnj7b8gnjR1wfCeMokhMDlQU
DbyajZrnCjM1ErmcGesfYyNXGut4V/ZD6fRI6umsaHy2ojoc1PLNL0fFx1vpcrv2
ecFIr8QsgvjmPBaf36mE4L9R/p+bhA==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDLE8Y4Mj0nErOw
H1xiRZKRjTvloUO0q8Pxd0bO2w6130he5Cqo9yJNRu9qKxhX/Bwf0WFN5rXicMLY
g7Tk5foQEN2tfY94cZpcBxqGd8fOpzlJ34kR9kwt3OGQrIiDdg4ALbjBXmCoiMC4
B+gxste6CVGVIJ6kDvAEGfrNQIIEbOPBIHVG2aBl5IR/mr3llhFmFTeZNZ3XnwuX
Aia85xUblokiWrxJ/5iPVav2FgPSzvczrHliqAe3D5fBVWz4QkD6h6lR2PvgJ6Nh
YBjIP44wDopwtxg3VJyDdUuzeZL1iEzdpmqu3Th5HWi4tG6p8b+g3uQ/QbmdACi/
+s3bleY1AgMBAAECggEAHmXdW/gZM0oXX5uyJnunjuYHOz5Cdj7n27MxBDHCJ/M+
UAXzZMtpdQvTNp8wy6rPQuwClricLUMXx3UflMvf/JupsxiCa1MF+hSOIea4H/Tk
HGy7hdbKXGsd0Jwi6xq2ycwORdOswE2IG4QDe7z7kbtnXN0/Ieuifdqtsh+f1q15
v3KBy+xNvrLezt+tEW/MKec4+z4w9a/t9vlZaFdy8pmbVoJ5R8R7JOEIb7jVU9bQ
dw0pR74dT9t+Wy/9IZ8Zg5fPtghHVtsfqaPI2y+PZ+GLIE2zDCkkk9IBTj+1UTyz
MYXqiV1F5xM4A4SvxjdnigjI1kQe4cyGt4Nb76IzmQKBgQDkCTW5JfVNI2+Gvnah
g26yzWbmpwajZDj4XKasZeKpO+sbs5ffMWiQXrel4pd4k8VEFLDYzkorivYFe47D
0MhQJ0srAXS5cxNsWxpNnCOdskU7ATw+IyduD/YXkVsO0CQ6rJYu/Lj+Pm6Z3HwB
anMJluAMh5pihUhpNMOEH6i0nwKBgQDj+wdbCnJmSIwf/7TJoXCNeOmH1uirv7td
hnh8ij9CiWHkHQizmO95GN2PDzvQzNbAscFsyCgTPDTB9Tw9fjZ6v+QL91pcBceV
hg1cjYpB0g0UuU88kW50NMytyQ3GJND72QoSsSejmFZrT8EilyVq5dF3Iqx6KuM1
3HnzfNLAqwKBgQDKw62zfh8KgjHO3Fjb8ORjtOSEv4vViW5m2OuTa90Joi/CKAUo
/uP9S1t882jAXJURnlxJy9SDt0JfSah+UY1sSCQ8j0TdCYgB11giHm5E8JlCiCyr
C9MpuKOX/TW8jDNKwN+h1DIiUB1ETpstHxw/MJr0STdr1xu2AKbBOu1l/QKBgCZs
2xKHMjz3IVcLXEdXxIazyHiyykiYalYbIhernXnzeeJe2maa8lHw1PcV1DkfLVsR
Gs7g9ZA4Z5QBZ4Pd0ATkbuVtbXdxKfCTxZDB8nmhk77YdPh6cql3dMAd0QqCjg7E
yCPaZBn4xSgVKzJPU2kvDx0LZRK6Q82CObPGaCc9AoGBANVlk+qQvp1Gi+IDnx7N
/i3yoxzoMXllpBvzYLN7xEMUJOcoDVhD1PzoBmr4Gu4iiLqmsVk9rsssUWRVzFkC
MLsSpBrUluqMOczU/b3XmpWkagSGBNTQ9LdwZJmX9e8tlUhK498wCZTWCWETJtdb
ehoKzHPSsE7bB1GDRjbHXC8F
-----END PRIVATE KEY-----
EOF
}

resource "alicloud_wafv3_instance" "default" {}

locals {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.id, "%s"])
}


`, name, casRegion)
}

// Test Wafv3 Domain. >>> Resource test cases, automatically generated.
// Case 企业级能力 7652
func TestAccAliCloudWafv3Domain_basic7652(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DomainMap7652)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DomainBasicDependence7652)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_wafv3_instances.default.ids.0}",
					"listen": []map[string]interface{}{
						{
							"http2_enabled":       "false",
							"enable_tlsv3":        "false",
							"ipv6_enabled":        "false",
							"protection_resource": "share",
							"http_ports": []string{
								"80"},
						},
					},
					"redirect": []map[string]interface{}{
						{
							"loadbalance":        "iphash",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"backends": []string{
								"1.1.1.1"},
							"xff_proto":          "true",
							"connect_timeout":    "5",
							"read_timeout":       "5",
							"write_timeout":      "5",
							"keepalive":          "false",
							"retry":              "false",
							"keepalive_requests": "60",
							"keepalive_timeout":  "15",
						},
					},
					"domain":      "cwaf-zctest-0731.wafqax.top",
					"access_type": "share",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"domain":      "cwaf-zctest-0731.wafqax.top",
						"access_type": "share",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listen": []map[string]interface{}{
						{
							"ipv6_enabled":        "false",
							"protection_resource": "share",
							"http_ports": []string{
								"81"},
							"xff_header_mode": "1",
							"http2_enabled":   "false",
							"enable_tlsv3":    "false",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends": []string{
								"123.56.107.188"},
							"loadbalance":        "iphash",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"xff_proto":          "true",
							"connect_timeout":    "6",
							"read_timeout":       "6",
							"write_timeout":      "6",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "100",
							"keepalive_timeout":  "16",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudWafv3DomainMap7652 = map[string]string{
	"status": CHECKSET,
}

func AlicloudWafv3DomainBasicDependence7652(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Case 更新Tags属性 9852
func TestAccAliCloudWafv3Domain_basic9852(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DomainMap9852)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DomainBasicDependence9852)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_wafv3_instances.default.ids.0}",
					"listen": []map[string]interface{}{
						{
							"http2_enabled":       "false",
							"enable_tlsv3":        "false",
							"ipv6_enabled":        "false",
							"protection_resource": "share",
							"http_ports": []string{
								"80"},
						},
					},
					"redirect": []map[string]interface{}{
						{
							"loadbalance":        "iphash",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"backends": []string{
								"1.1.1.1"},
							"xff_proto":          "true",
							"connect_timeout":    "5",
							"read_timeout":       "5",
							"write_timeout":      "5",
							"keepalive":          "false",
							"retry":              "false",
							"keepalive_requests": "60",
							"keepalive_timeout":  "15",
						},
					},
					"domain":      "qiyezhili-1118.wafqax.top",
					"access_type": "share",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"domain":      "qiyezhili-1118.wafqax.top",
						"access_type": "share",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listen": []map[string]interface{}{
						{
							"ipv6_enabled":        "false",
							"protection_resource": "share",
							"http_ports": []string{
								"81"},
							"xff_header_mode": "1",
							"http2_enabled":   "false",
							"enable_tlsv3":    "false",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends": []string{
								"123.56.107.188"},
							"loadbalance":        "iphash",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"xff_proto":          "true",
							"connect_timeout":    "6",
							"read_timeout":       "6",
							"write_timeout":      "6",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "100",
							"keepalive_timeout":  "16",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudWafv3DomainMap9852 = map[string]string{
	"status": CHECKSET,
	"cname":  CHECKSET,
}

func AlicloudWafv3DomainBasicDependence9852(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Case 企业级能力qiyezhili-1118.wafqax.top账号_线上_换账号_测试通过_待发布 11009
func TestAccAliCloudWafv3Domain_basic11009(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DomainMap11009)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DomainBasicDependence11009)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_wafv3_instances.default.ids.0}",
					"listen": []map[string]interface{}{
						{
							"protection_resource": "share",
							"http_ports": []string{
								"81", "82", "83"},
							"https_ports":     []string{},
							"xff_header_mode": "2",
							"xff_headers": []string{
								"testa", "testb", "testc"},
							"custom_ciphers": []string{},
							"ipv6_enabled":   "true",
						},
					},
					"redirect": []map[string]interface{}{
						{
							"loadbalance":        "iphash",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"backends": []string{
								"1.1.1.1", "3.3.3.3", "2.2.2.2"},
							"xff_proto":          "true",
							"connect_timeout":    "5",
							"read_timeout":       "5",
							"write_timeout":      "5",
							"keepalive":          "true",
							"retry":              "true",
							"keepalive_requests": "1000",
							"keepalive_timeout":  "15",
							"request_headers": []map[string]interface{}{
								{
									"key":   "testkey1",
									"value": "testValue1",
								},
								{
									"key":   "key1",
									"value": "value1",
								},
								{
									"key":   "key22",
									"value": "value22",
								},
							},
						},
					},
					"domain":                             "qiyezhili-1118.wafqax.top",
					"access_type":                        "share",
					"resource_manager_resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":                        CHECKSET,
						"domain":                             "qiyezhili-1118.wafqax.top",
						"access_type":                        "share",
						"resource_manager_resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listen": []map[string]interface{}{
						{
							"ipv6_enabled":        "false",
							"protection_resource": "share",
							"xff_header_mode":     "1",
							"https_ports":         []string{},
							"http_ports": []string{
								"84", "86"},
							"xff_headers":    []string{},
							"custom_ciphers": []string{},
						},
					},
					"redirect": []map[string]interface{}{
						{
							"backends": []string{
								"123.56.107.188", "1.1.1.1"},
							"loadbalance":        "leastTime",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"xff_proto":          "false",
							"connect_timeout":    "6",
							"read_timeout":       "6",
							"write_timeout":      "6",
							"keepalive":          "true",
							"retry":              "false",
							"keepalive_requests": "100",
							"keepalive_timeout":  "16",
							"request_headers": []map[string]interface{}{
								{
									"key":   "testky1",
									"value": "testvalue2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listen": []map[string]interface{}{
						{
							"xff_header_mode":     "1",
							"protection_resource": "gslb",
							"ipv6_enabled":        "false",
							"https_ports":         []string{},
							"http_ports": []string{
								"84", "86"},
						},
					},
					"redirect": []map[string]interface{}{
						{
							"loadbalance":        "leastTime",
							"connect_timeout":    "6",
							"read_timeout":       "6",
							"write_timeout":      "6",
							"focus_http_backend": "false",
							"sni_enabled":        "false",
							"keepalive":          "false",
							"retry":              "false",
							"xff_proto":          "false",
							"backends": []string{
								"1.1.1.1", "5.5.5.5"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudWafv3DomainMap11009 = map[string]string{
	"status":    CHECKSET,
	"domain_id": CHECKSET,
}

func AlicloudWafv3DomainBasicDependence11009(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Test Wafv3 Domain. <<< Resource test cases, automatically generated.
