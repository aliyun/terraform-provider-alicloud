package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudWafCertificate_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_waf_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudWAFCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafcertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWAFCertificateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_name": name,
					"instance_id":      "${alicloud_waf_domain.domain.instance_id}",
					"domain":           "${alicloud_waf_domain.domain.domain_name}",
					"certificate":      "${var.cert}",
					"private_key":      "${var.private_key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_name": name,
						"instance_id":      CHECKSET,
						"domain":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
		},
	})
}

func TestAccAliCloudWafCertificate_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	resourceId := "alicloud_waf_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudWAFCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafcertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWAFCertificateBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_waf_domain.domain.instance_id}",
					"domain":         "${alicloud_waf_domain.domain.domain_name}",
					"certificate_id": "${alicloud_ssl_certificates_service_certificate.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"domain":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
		},
	})
}

var AlicloudWAFCertificateMap0 = map[string]string{}

func AlicloudWAFCertificateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

// If cert has expired, there can generate a new one: https://myssl.com/create_test_cert.html
variable "cert" {
  default = <<EOF
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
variable "private_key" {
  default = <<EOF
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
}

data "alicloud_waf_instances" "default" {}

resource "alicloud_waf_domain" "domain" {
  domain_name       = "alicloud-provider.cn"
  instance_id       = data.alicloud_waf_instances.default.ids.0
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
`, name)
}

func AlicloudWAFCertificateBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "cert" {
  default = <<EOF
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
variable "private_key" {
  default = <<EOF
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
}

data "alicloud_waf_instances" "default" {}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
  cert = var.cert
  key = var.private_key
}

resource "alicloud_waf_domain" "domain" {
  domain_name       = "alicloud-provider.cn"
  instance_id       = data.alicloud_waf_instances.default.ids.0
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
`, name)
}

func TestUnitAlicloudWAFCertificate(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dId, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dIdCreate, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dIdCreate.MarkNewResource()
	dCreateError, _ := schema.InternalMap(p["alicloud_waf_certificate"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"certificate_name": "certificate_name",
		"instance_id":      "instance_id",
		"domain":           "domain",
		"private_key":      `-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA5KeYEdE3moKDwwB4DV+yB44BNOqJF6KOS3wSF0IhjnWRTJ13\nUoCTN0CwIJTPzgu7hPXuvoRe2Mgn/5CclHaF3x1+FIGVskyC1hm5I29nuP21MDkN\njJ0e9TLzSLfKNmSiJKbUqeLa6l46U/5rUdI3a+qOADXNIV4chZlgYXng+HykXfDK\nXDdsuGE8h9Ue2+1WLM7KlWpbFlG8JCTTxIKhuKqORvdZRPE0u20tVqELtSjrglkR\n6yc1ZXZ3MIpS9L6QynysxnCJ/CpYKxlR+SkVR5Uozj3+J4tPtYy7kQR7AsWAIxnr\n5JCAW96yFjK+LSdtiJm5fvcmXPJt8Byui67sJwIDAQABAoIBAHkIBUSZG/ebUids\nHh/mIZRCJH1gEAUEtJFcMS9CgASqUqjdSQXQrRw4D4KPgpesOiwdCayq4JGbESED\n6P3k5uO89rngBd8FBNRTdV4+YAbZtIi8CmH94R8lQwWg01aLWhBQho0OWZZpHjLL\nnv8GaPOkilmuINCiwjIth4jRPLYKM8rQUL0Y9DSLW1dZfxUMj5yPHdpV81/es21k\ns3ZUETOsygAfv9ybhxFW/1h8vldH9bCRvVq4oXHXkpF9Wp9K6z4UjR60k7hnxzTw\nUgFoEheHz58LOJNzh+dad/2HCCEVutQXMIMomoMvSmQkODhOMYsiJC8L15l6/193\ngj6Ah1ECgYEA8AE/efv7Q8DGwVRP51a4whzkn2Mws/1cURxXLflYNdrOek3Vmdda\nR4k9N6xYKTjQ4xSV+yYHpCMBdSkFcKFOOMf+DQLB5NAb3C4E+K6ea8GrfIUNnFZM\nYKwfjhDEvNFky8FqLscQGr9HCzcx+lOrwtnMK/srvS8r9IuNfXcalaUCgYEA8+Sy\n4brUb9SNqAFY68lOsovuhQpCK8UnomgbVZG52wziA368HYfdmLFSzSAMaOYA4wog\nUvC1JtOqkP6tyAFdo8qWgrGtmBpVmPSvWA2ser0zqugvRyg4uY6vd3Hwj9VEGguB\nvotYxUhzetzyapTZazzRuE7aUB9dpbnvdDhoyNsCgYBgpbABGFajfwLKoYAXwZVf\nHbF2+cOIB5PgWdBFhC5gaX9SQvusGsuRGRPc0nsiBm4fs4999l+HWk1g96boJxzP\nwsFsTdr0oFVHwgRgfDjxXakH2LCVby8MkuOWGuyOuKelYXq34ZN7oeEjBBQNIRp+\nuaP9ZgTEBzXlITlV99ttvQKBgQCEJijeslKk/XB71a8OoxbnG3bz0ykjekINd4dI\nCBvCGurjpenbvmBNedc5meHffLCCVFxLVG1zSkEjKzuSIVtRglWnHwrGXY5/wCS5\n+z63iojSU9g6IAsMT8m3WJ9V7+JkklOMoQhKbQVVTx5yrZBY0K15xg/4VeZyA5tB\nR8dO/QKBgGqkJ1AB+qi9Tl10ic8IX2blqPt3FU6MkkVmmDl8vA5R35DBuZAD4VTs\nvsc0Y79mSqP9XL3KRAfA04tbGme30gJWz4NOoNsaaF2T2fv0gNNnVma87unFS6Y4\nFv64CkXzShjd16ov4eZetsIAZYn/bVn8zp61I6V50iT6AjpX1ptX\n-----END RSA PRIVATE KEY-----`,
		"certificate":      `-----BEGIN CERTIFICATE-----\nMIID8DCCAtigAwIBAgIQTvMBGm/PRXSj352aOU7GSjANBgkqhkiG9w0BAQsFADBe\nMQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl\nc3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe\nFw0yMTA5MDcxNDM5NTBaFw0yMjA5MDcxNDM5NTBaMC0xCzAJBgNVBAYTAkNOMR4w\nHAYDVQQDExV0Zi10ZXN0YWNjLndhZnFhMy5jb20wggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQDkp5gR0TeagoPDAHgNX7IHjgE06okXoo5LfBIXQiGOdZFM\nnXdSgJM3QLAglM/OC7uE9e6+hF7YyCf/kJyUdoXfHX4UgZWyTILWGbkjb2e4/bUw\nOQ2MnR71MvNIt8o2ZKIkptSp4trqXjpT/mtR0jdr6o4ANc0hXhyFmWBheeD4fKRd\n8MpcN2y4YTyH1R7b7VYszsqValsWUbwkJNPEgqG4qo5G91lE8TS7bS1WoQu1KOuC\nWRHrJzVldncwilL0vpDKfKzGcIn8KlgrGVH5KRVHlSjOPf4ni0+1jLuRBHsCxYAj\nGevkkIBb3rIWMr4tJ22Imbl+9yZc8m3wHK6LruwnAgMBAAGjgdowgdcwDgYDVR0P\nAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME\nGDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB\nBQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov\nL2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MCAGA1UdEQQZMBeCFXRmLXRl\nc3RhY2Mud2FmcWEzLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAPjbt2H1HmEc8DzyD\npi4IF1CvaNlYgKjPojYlt/gpj2n0MfntL8Ihly3e2fdSMEeVeTnFWFd34L4uZxMa\nxE/hx6VJWfNdnYW7FGCZr0rGj/KrtAox83H1dRrZ63hynpgCMIbc5lhA7wDe0R16\nP/1l3c50ZEmidicGhK/qmzsSQIVXC0kJf6hDQCxW6LVaDrmT8mvbhRh4ZNb2pgJ5\nQIWJHnlOmZkUVsR5cMBGzK2ModADjHXHmeoHHr3Tw7mPioE4Xh5EmMTXPLG22BKN\nRBFG9gSFri+3RxqdXwi1ZJajO3Nup5mcdGaHJbUbNUf16YKIq50PJlrVxzCZV31f\n7cOGfw==\n-----END CERTIFICATE-----`,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for key, value := range map[string]interface{}{
		"certificate_name": "certificate_name",
		"instance_id":      "instance_id",
		"domain":           "domain",
	} {
		err := dCreateError.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for idKey, idValue := range map[string]interface{}{
		"instance_id":    "instance_id",
		"domain":         "domain",
		"certificate_id": "123456",
	} {
		err := dIdCreate.Set(idKey, idValue)
		assert.Nil(t, err)
		err = dId.Set(idKey, idValue)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Certificates": []interface{}{
			map[string]interface{}{
				"CertificateId":   "123456",
				"CertificateName": "certificate_name",
				"InstanceId":      "instance_id",
				"Domain":          "domain",
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_waf_certificate", "MockCertificateId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["CertificateId"] = "123456"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewWafClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudWafCertificateCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dCreateError, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("CreateCertificateByCertificateIdAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dId, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateCertificateByCertificateIdNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudWafCertificateCreate(dIdCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "domain", ":", "123456"))
	//Delete
	t.Run("DeleteNormal", func(t *testing.T) {
		err := resourceAlicloudWafCertificateDelete(d, rawClient)
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeWafCertificateNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudWafCertificateRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeWafCertificateAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudWafCertificateRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
