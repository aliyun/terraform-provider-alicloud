package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudWafCertificateDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafCertificateDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_waf_certificate.default.id}"]`,
			"instance_id": `"${alicloud_waf_certificate.default.instance_id}"`,
			"domain":      `"${alicloud_waf_certificate.default.domain}"`,
		}),
		fakeConfig: testAccCheckAlicloudWafCertificateDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_waf_certificate.default.id}_fake"]`,
			"instance_id": `"${alicloud_waf_certificate.default.instance_id}"`,
			"domain":      `"${alicloud_waf_certificate.default.domain}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafCertificateDataSourceName(rand, map[string]string{
			"name_regex":  `"${alicloud_waf_certificate.default.certificate_name}"`,
			"instance_id": `"${alicloud_waf_certificate.default.instance_id}"`,
			"domain":      `"${alicloud_waf_certificate.default.domain}"`,
		}),
		fakeConfig: testAccCheckAlicloudWafCertificateDataSourceName(rand, map[string]string{
			"name_regex":  `"${alicloud_waf_certificate.default.certificate_name}_fake"`,
			"instance_id": `"${alicloud_waf_certificate.default.instance_id}"`,
			"domain":      `"${alicloud_waf_certificate.default.domain}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafCertificateDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_waf_certificate.default.id}"]`,
			"instance_id": `"${alicloud_waf_certificate.default.instance_id}"`,
			"domain":      `"${alicloud_waf_certificate.default.domain}"`,
			"name_regex":  `"${alicloud_waf_certificate.default.certificate_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudWafCertificateDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_waf_certificate.default.id}_fake"]`,
			"instance_id": `"${alicloud_waf_certificate.default.instance_id}"`,
			"domain":      `"${alicloud_waf_certificate.default.domain}"`,
			"name_regex":  `"${alicloud_waf_certificate.default.certificate_name}_fake"`,
		}),
	}
	var existAlicloudWafCertificateDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"certificates.#":                  "1",
			"certificates.0.certificate_name": CHECKSET,
			"certificates.0.certificate_id":   CHECKSET,
		}
	}
	var fakeAlicloudWafCertificateSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudWafCertificateCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_waf_certificates.default",
		existMapFunc: existAlicloudWafCertificateDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudWafCertificateSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckForCleanUpInstances(t, string(connectivity.APSouthEast1), "waf", "waf", "waf", "waf")
	}
	alicloudWafCertificateCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudWafCertificateDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccWafCertificate-%d"
}

variable "is_intl" {
    default = %v
}

// Create a waf instance in ap-southeast-1
resource "alicloud_waf_instance" "default" {
  big_screen           = "0"
  exclusive_ip_package = "0"
  ext_bandwidth        = "0"
  ext_domain_package   = "0"
  package_code         = var.is_intl ? "version_pro" : "version_pro_asia"
  prefessional_service = "false"
  subscription_type    = "Subscription"
  period               = 3
  waf_log              = "false"
  log_storage          = "3"
  log_time             = "180"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

resource "alicloud_waf_domain" "domain" {
  domain_name       = "alicloud-provider.cn"
  instance_id       = alicloud_waf_instance.default.id
  is_access_product = "On"
  source_ips        = ["1.1.1.1"]
  cluster_type      = "PhysicalCluster"
  http2_port        = [443]
  http_port         = [80]
  https_port        = [443]
  http_to_user_ip   = "Off"
  https_redirect    = "Off"
  load_balancing    = "IpHash"
  log_headers {
    key   = "foo"
    value = "http"
  }
}

data "alicloud_resource_manager_resource_groups" "default" {
	name_regex="^default$"
}
// If cert has expired, there can generate a new one: https://myssl.com/create_test_cert.html
resource "alicloud_waf_certificate" "default" {
  certificate_name = var.name
  instance_id = alicloud_waf_domain.domain.instance_id
  domain = alicloud_waf_domain.domain.domain_name
  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAptPW1XETCowwrPjTpAHzAbXHt5Pw3l3coeUFZPv37rAzvqXc
z55FrjDvO1Cvrh3X/M7yEl+F/F2ziwqJkwBIA29rxM92qiyTkQPcUYbTXykPnFZ2
VylScsRFboEHGTBQmM6yl7UTM2EbyFlTRgbu/NHvi5N5uXTLDiZxC9l4zOkkIsJH
HwwlO+Rwdx9/sDapZp6CXnWNUxW/XA4Ka11OBKxv3wxxZ/HC3Dmn8tZhOluk4SKn
QwU470pP8XE9NX/+S1mx7P437v8DwCzsPhB+/kDyKDtw36M0c94poL4YqO4QMy01
4hSG4zxehLgvdMkK0FdUBL4vvKJGLBMg2JCqgQIDAQABAoIBABqYixdb/30e/GOX
B6aN6P/MyxopK61HqXTl3gZLZStYS+eI6brvsIwkoP8Dzf4kI9rZ4x2qvOGOzqDL
/ULCOLjTeorVyEU49g2YX/dfVzNEiIiUACnFHK/POoJzb35EhZTW3wHwjC8UvCkk
lHIFiPQlQ9ssl6tJQ1XuKveIpAmGn6B2QAIGfHtYpmR3vA1RAAgTTvLDLnPUhgiy
e8UWfyTk3kTRYWO44mvitiwoEd7wHsmWT58dSxYyMxu/Tp2sstC5d/U+PzRRJJqm
UL5SzDg1BRZ9cIi0wJHAG41AL3vcTEp3rTMPD514te2MrUNkmtBAQ9dZpuqIqrUM
HH6THwECgYEAyLbWKMAeITJoQDFQ5i0phZ9gLPFKb3GRNAqHkS2e+ayXRDv+DKej
1VuuMc8TG1M2zgEXxqHoa9OTWh3CP/CDP+1ZKkGXa2hGLQ0ulczam3HN3r7wkcTK
yGxLs3SjLsSmV2sWmOwOdVWexMNpV2HPsc77Vtt2syqMSLApboSWlN0CgYEA1MeD
00n+IDgDAbDhAzzy8DL4xWZznR7+njROTYoillG2OFY5we5QISJUhqIKrA2DxJqG
7W8B5IuC+8H7kxn6YSAmvKBFwkU+TftAtGrZ5Njdxq6e/L/DarGzERduf1rHU3ho
XGdp+PCY0BhGvQWZV7XRG42zuzGnsw2VXvRUT/UCgYEAsc/gO+uvBYPNfYjy/Wdn
EHyXzFjas43tsCff38qFuvSecZGgR/+/kYWN9YxYhe2bcLbhJFRNPjEKqlwQuWpO
LDUNt/SxuzJ+cOzrp0P7KpiQZFPNdaXca+Ac1FdxNT57iphRRZpiWKpIZFIloYcC
Y6hVW6ZDKwh3jPGbUtBdQnUCgYEAyMa6jVT0hlJUA/RzXirU0fJG632Q+mGUwJN6
j5AYst7HosE6HtRQp/NM7v1YNKFMBOgLCWAF2TILlqOkT1nj+GHuK27QZFASSmMa
2yM4F8jK+/8p5jTMoJCs9yO1EFcIXqrvTsIGLPBaQYGmvhXaa/kUJ3u+bJSeOUZH
/CFwrUECgYB+VBPXSMdf31h0HxA3qrlZYYiIQJGRUpcEbG9yRH0daCT7XANhsQ0H
smgbcvXNfgSkgSztCmdjJMBPOB7LzBdQN15Lr3TUAsaRwGs9Hn6VigeeVAsCLp0d
D55K6Bc+fX4vmnIk4fQ4rN+eIX2vOvoji199ED+XY++P4zAL2CMvkA==
-----END RSA PRIVATE KEY-----
EOF

  certificate =  <<EOF
-----BEGIN CERTIFICATE-----
MIID7TCCAtWgAwIBAgIPdLEpOedN9r3bJJ67LrHBMA0GCSqGSIb3DQEBCwUAMF4x
CzAJBgNVBAYTAkNOMQ4wDAYDVQQKEwVNeVNTTDErMCkGA1UECxMiTXlTU0wgVGVz
dCBSU0EgLSBGb3IgdGVzdCB1c2Ugb25seTESMBAGA1UEAxMJTXlTU0wuY29tMB4X
DTI0MTIxMDA4MzMxM1oXDTI5MTIwOTA4MzMxM1owLDELMAkGA1UEBhMCQ04xHTAb
BgNVBAMTFGFsaWNsb3VkLXByb3ZpZGVyLmNuMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAptPW1XETCowwrPjTpAHzAbXHt5Pw3l3coeUFZPv37rAzvqXc
z55FrjDvO1Cvrh3X/M7yEl+F/F2ziwqJkwBIA29rxM92qiyTkQPcUYbTXykPnFZ2
VylScsRFboEHGTBQmM6yl7UTM2EbyFlTRgbu/NHvi5N5uXTLDiZxC9l4zOkkIsJH
HwwlO+Rwdx9/sDapZp6CXnWNUxW/XA4Ka11OBKxv3wxxZ/HC3Dmn8tZhOluk4SKn
QwU470pP8XE9NX/+S1mx7P437v8DwCzsPhB+/kDyKDtw36M0c94poL4YqO4QMy01
4hSG4zxehLgvdMkK0FdUBL4vvKJGLBMg2JCqgQIDAQABo4HZMIHWMA4GA1UdDwEB
/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHwYDVR0jBBgw
FoAUKIEmBdE0Gj/Bcw+7k88VHD8Dv38wYwYIKwYBBQUHAQEEVzBVMCEGCCsGAQUF
BzABhhVodHRwOi8vb2NzcC5teXNzbC5jb20wMAYIKwYBBQUHMAKGJGh0dHA6Ly9j
YS5teXNzbC5jb20vbXlzc2x0ZXN0cnNhLmNydDAfBgNVHREEGDAWghRhbGljbG91
ZC1wcm92aWRlci5jbjANBgkqhkiG9w0BAQsFAAOCAQEAGki5WtiA+LoPzvBWXTQW
ZjMTfjjbwbh9KXCAkL1s5itlGtU4ZPzj7OrK6JxgqUyvwsOj3oHwbrD6oJmtVrmh
JUOeK6JF51b+KhVu9oOapKGWfcEwZEmjtL/vMGhZF4Bk6pTBERurZg9+HHyj5SrI
xTU4T++GJib+ixTYu2wgzu3MD+UXVb2zwaEfoPMRIaiVOdNJPy3GzaQmc3wBW+Qn
FiHTfmirK5wPDarw+J1CCc5MBki8yNh+OVma0twEliIn3H1LcKB+tkcJAjeQ+tAX
Lx+hoPqFkcdlcgN0svIIUNoN370wuH4t1fzX9AnpgcM04J5+m0LkAENzp1MXzjhf
TA==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDuzCCAqOgAwIBAgIQSEIWDPfWTDKZcWNyL2O+fjANBgkqhkiG9w0BAQsFADBf
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxLDAqBgNVBAsTI015U1NMIFRl
c3QgUm9vdCAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMTcxMTE2MDUzNTM1WhcNMjcxMTE2MDUzNTM1WjBeMQswCQYDVQQGEwJDTjEO
MAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRlc3QgUlNBIC0gRm9yIHRl
c3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAMBOtZk0uzdG4dcIIdcAdSSYDbua0Bdd6N6s4hZaCOup
q7G7lwXkCyViTYAFa3wZ0BMQ4Bl9Q4j82R5IaoqG7WRIklwYnQh4gZ14uRde6Mr8
yzvPRbAXKVoVh4NPqpE6jWMTP38mh94bKc+ITAE5QBRhCTQ0ah2Hq846ZiDAj6sY
hMJuhUWegVGd0vh0rvtzvYNx7NGyxzoj6MxkDiYfFiuBhF2R9Tmq2UW9KCZkEBVL
Q/YKQuvZZKFqR7WUU8GpCwzUm1FZbKtaCyRRvzLa5otghU2teKS5SKVI+Tpxvasp
fu4eXBvveMgyWwDpKlzLCLgvoC9YNpbmdiVxNNkjwNsCAwEAAaN0MHIwDgYDVR0P
AQH/BAQDAgGGMA8GA1UdJQQIMAYGBFUdJQAwDwYDVR0TAQH/BAUwAwEB/zAfBgNV
HSMEGDAWgBSa8Z+5JRISiexzGLmXvMX4oAp+UzAdBgNVHQ4EFgQUKIEmBdE0Gj/B
cw+7k88VHD8Dv38wDQYJKoZIhvcNAQELBQADggEBAEl01ufit9rUeL5kZ31ox2vq
648azH/r/GR1S+mXci0Mg6RrDdLzUO7VSf0JULJf98oEPr9fpIZuRTyWcxiP4yh0
wVd35OIQBTToLrMOWYWuApU4/YLKvg4A86h577kuYeSsWyf5kk0ngXsL1AFMqjOk
Tc7p8PuW68S5/88Pe+Bq3sAaG3U5rousiTIpoN/osq+GyXisgv5jd2M4YBtl/NlD
ppZs5LAOjct+Aaofhc5rNysonKjkd44K2cgBkbpOMj0dbVNKyL2/2I0zyY1FU2Mk
URUHyMW5Qd5Q9g6Y4sDOIm6It9TF7EjpwMs42R30agcRYzuUsN72ZFBYFJwnBX8=
-----END CERTIFICATE-----
EOF
}

data "alicloud_waf_certificates" "default" {	
	%s
}
`, rand, os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE") == string(IntlSite), strings.Join(pairs, " \n "))
	return config
}
