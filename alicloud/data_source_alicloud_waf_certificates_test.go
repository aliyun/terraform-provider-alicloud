package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudWAFCertificateDataSource(t *testing.T) {
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

resource "alicloud_waf_certificate" "default" {
  certificate_name = var.name
  instance_id = alicloud_waf_domain.domain.instance_id
  domain = alicloud_waf_domain.domain.domain_name
  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAxs07tDiEXk9gb/18aqsjTo/hvA9bVmFE9HVaDV/2EchCi3Xt
fW4eqGaQq3NqpP1Ts/YV1GQ6o7KcsU26NMqiG7LlHsOK6/2ECcb63ffCigBCxxiP
7bHw4xs/DlCg30DUN3FHlzZFkc1lWNUvlo8yqyybRsBEX28qE9xmId7oJv9efVjA
suSh1nryTHDAkAab/A7BV6yZBU55Q/G/8OHxHwYFU/CrR7A6z5TAmGsE2fJ8j097
Lx8YBoQ7PdRy7coxxe0DouuxYZxdEcofbMKXO5dbMZaQR+tiXwU0oTyuinxt/ukU
vbA2PFWABbe1IDPNdB5yWndgQnn8p8mxinTPgQIDAQABAoIBAAeQeZbC68PMO2bS
JgYdBZaQrPiDO74kAgQB5G+hLh4WbknGLpzY2tEc40uoHu+w3/AiC0nt7Q6bgOYn
h0/zlzV8mhvl5DGp5M3xGfc7h9XJFWwOnYi1fkKgn8aYqiEJ43P0myfAHeeNVto0
H7/Cf8gO8MhBZUNIWePa6Aaeklt2x3PD2iGHh9XxKJkIXykRv0MwaeScp4rzgiDU
b0dXhHb1DzAnNEUf98DxA9X+dRb/JFaI7hDH+7zt+rV6idwR2MVlJnRAUunRUpEZ
tEiAXNzUWwIe7Hzj+jRUj7kHgwyfgR2xalF8ZZunBsW13WPrlNCyxOhC1ZKzvwet
ktlHZU0CgYEA+8l6JtHOr/Pid2P6VUgP1edqr1QDa9QJ49hJP9MDYFsfJlvR3C1G
atGRDGuNlLE+Mht33wRUmkqqWrRq7WAS9cFt6yfIyQ8CzREr9ZY/AGri1SWL0yGD
H9BXSRUmmwRICOwd0vOcz0gEnBQPhxPeVEaBprSBneYi/Zj+a93Hwh8CgYEAyiDL
e5buc6IXakqtW0oPlVucvSVhOuvDBLzfexv5WGywjILpJI0Pamb9tDVjz8HiLgO+
f9QSihidyiTh8QKs/BqcVcMW8qR/lZIyOnKgvc5snBqx8aDeAFgfWqIZjwKsq7eu
EO7s/2ZFn2Tg+Yrh57+jRPSFgBcC8fn14/UPel8CgYEA26suyBl54qa7aQAv0iGZ
tzw8JpT8myGWM1NqrKVTp8g4CAZJtHdHnUAS9SwMjLKKGbs/PLIRgb8smxAWzdxp
DkHgvS2rjkZrYi/eE1guxRV7qwwjzmLJfIMO8/LhCJOWqToe1wG1SF6DrwwNdALQ
kOQmjyOihfjXhDrrlX4bnm8CgYEAoeOGVmRac47dhiptknJM7OlFLlEkANcXiVHl
BGsxWFslREUNilLYh/YZOR6R5LJ6/zPMfgFTH/v7VQ4ZULk952VA1Ye9d3W9IgTH
6fwpNhyA7L6MyuR3KntqEWNHaT1RIu+ooLNcamp1VOatMfDEYqMgl3fo/OBxYXE8
OenrDa0CgYEAtB0RHiQJMi2n0AtOGuDbV8HRxk9Mh0GypjD3LGsdduG68bWe8WeT
xN9ZNUeHT2VXSF1ra0OZvyAdTZnzxwgZFHlKSbh5Fzta1fnoXj8x8akp1mycvJL7
Ehc7Tyxrl5vrxONKMbQdpC45rPxyKeGmk+b1LdbvoWr2iQwqL/D9Lr0=
-----END RSA PRIVATE KEY-----
EOF

  certificate =  <<EOF
-----BEGIN CERTIFICATE-----
MIIHrjCCBZagAwIBAgIQC2/QuNyfdGuiLP17UTTLyjANBgkqhkiG9w0BAQsFADBc
MQswCQYDVQQGEwJVUzEXMBUGA1UEChMORGlnaUNlcnQsIEluYy4xNDAyBgNVBAMT
K1JhcGlkU1NMIEdsb2JhbCBUTFMgUlNBNDA5NiBTSEEyNTYgMjAyMiBDQTEwHhcN
MjMwMzA4MDAwMDAwWhcNMjQwMzA3MjM1OTU5WjAfMR0wGwYDVQQDExRhbGljbG91
ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMbN
O7Q4hF5PYG/9fGqrI06P4bwPW1ZhRPR1Wg1f9hHIQot17X1uHqhmkKtzaqT9U7P2
FdRkOqOynLFNujTKohuy5R7Diuv9hAnG+t33wooAQscYj+2x8OMbPw5QoN9A1Ddx
R5c2RZHNZVjVL5aPMqssm0bARF9vKhPcZiHe6Cb/Xn1YwLLkodZ68kxwwJAGm/wO
wVesmQVOeUPxv/Dh8R8GBVPwq0ewOs+UwJhrBNnyfI9Pey8fGAaEOz3Ucu3KMcXt
A6LrsWGcXRHKH2zClzuXWzGWkEfrYl8FNKE8rop8bf7pFL2wNjxVgAW3tSAzzXQe
clp3YEJ5/KfJsYp0z4ECAwEAAaOCA6cwggOjMB8GA1UdIwQYMBaAFPCchf2in32P
yWi71dSJTR2+05D/MB0GA1UdDgQWBBQAu+tDmtZxFyGaT2ZZzNvgxMlgcjA5BgNV
HREEMjAwghRhbGljbG91ZC1wcm92aWRlci5jboIYd3d3LmFsaWNsb3VkLXByb3Zp
ZGVyLmNuMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYB
BQUHAwIwgZ8GA1UdHwSBlzCBlDBIoEagRIZCaHR0cDovL2NybDMuZGlnaWNlcnQu
Y29tL1JhcGlkU1NMR2xvYmFsVExTUlNBNDA5NlNIQTI1NjIwMjJDQTEuY3JsMEig
RqBEhkJodHRwOi8vY3JsNC5kaWdpY2VydC5jb20vUmFwaWRTU0xHbG9iYWxUTFNS
U0E0MDk2U0hBMjU2MjAyMkNBMS5jcmwwPgYDVR0gBDcwNTAzBgZngQwBAgEwKTAn
BggrBgEFBQcCARYbaHR0cDovL3d3dy5kaWdpY2VydC5jb20vQ1BTMIGHBggrBgEF
BQcBAQR7MHkwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBR
BggrBgEFBQcwAoZFaHR0cDovL2NhY2VydHMuZGlnaWNlcnQuY29tL1JhcGlkU1NM
R2xvYmFsVExTUlNBNDA5NlNIQTI1NjIwMjJDQTEuY3J0MAkGA1UdEwQCMAAwggF+
BgorBgEEAdZ5AgQCBIIBbgSCAWoBaAB1AO7N0GTV2xrOxVy3nbTNE6Iyh0Z8vOze
w1FIWUZxH7WbAAABhsATxUwAAAQDAEYwRAIgO0/k4INTQs3bauRtx0bz7BkmcDDE
5xjCzPPqMbxnwYsCIHcVsaYOZO73eiosYrTolY5/XIUWVDYfLmDia/bxtc9QAHYA
c9meiRtMlnigIH1HneayxhzQUV5xGSqMa4AQesF3crUAAAGGwBPFowAABAMARzBF
AiBO/OkeQCRMOU8UiRToVjnnC1F4lKvjChQKnlgDi/26ggIhAJOTO3zod7m0ILNY
8Ue/jiAwmX3x66dEU+pIXuMSqoi5AHcASLDja9qmRzQP5WoC+p0w6xxSActW3SyB
2bu/qznYhHMAAAGGwBPFeAAABAMASDBGAiEAjsP0zTN0Gu9D10an9u+KysJAmcnN
SG7djcEpOuf0ov0CIQD+6s3lC2Qj4ddNUQ12MRMpytJO4a440LQypfDQtEGcHTAN
BgkqhkiG9w0BAQsFAAOCAgEAcD60o2Astoawsg8MV3ilTBVf8yYaWXQ4uFy5/xwO
KSQb9Tb1ZWVyB+zJcHsJTJNVphC4rDzd2bGTDir+rcZlSDSTU9/xn/0g26I28dNq
86bGJOvWVMJmo9WbtQPbEncR9K/Vc0eQteX8q8DZsTPyAXIJSwVOXGpsqDRmZPvT
65mcg5KxTzu7MmvFkD8oPQrp9XmURxu1Z9Fm6YOr5B6/CZquY9fjNGUP52lzxU7j
7qbQsGTGMRAQ7HcJvetanP3DXYOi/nFdUPULsmTXbEVLV0c8nrsLtjAJI0pZaKkx
pqjTr8I6jmOLrzj0O4NymXe54kb0ch65l85A9Pauc7mmTT87gpEjuCy1FXlPpsJZ
xi3uANL6SzPKnma0GrJbBhaX7HpiroxuHU50fr/umdQVUqBldfcRAaY5OfcdB1zU
7rCQXfeSvR98yr2pLu9xVUs1EFvMD5XzPk1WA3aZDRCEenhYykVFFrog37UVbxJA
HvOPx9qzcy23Oz4BQu2sMNmJKV8zBxCDxTpEiGjIuDykDmvgeS4uhKsMU4tPRpXi
sw6G6grVJ4YmiiA5uUrBbXhZFw8IYxk27av9b8kEiUOv9ytdN3Ay84E9HwWDld5B
gvACxJTy8IyS3s4vi/HRS9k4IMgCZLUvlGFi3G8VPl/cFQSUeurMyQxuMbqYrXcL
t9Q=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFyzCCBLOgAwIBAgIQCgWbJfVLPYeUzGYxR3U4ozANBgkqhkiG9w0BAQsFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0yMjA1MDQwMDAwMDBaFw0zMTExMDkyMzU5NTlaMFwxCzAJBgNVBAYTAlVT
MRcwFQYDVQQKEw5EaWdpQ2VydCwgSW5jLjE0MDIGA1UEAxMrUmFwaWRTU0wgR2xv
YmFsIFRMUyBSU0E0MDk2IFNIQTI1NiAyMDIyIENBMTCCAiIwDQYJKoZIhvcNAQEB
BQADggIPADCCAgoCggIBAKY5PJhwCX2UyBb1nelu9APen53D5+C40T+BOZfSFaB0
v0WJM3BGMsuiHZX2IHtwnjUhLL25d8tgLASaUNHCBNKKUlUGRXGztuDIeXb48d64
k7Gk7u7mMRSrj+yuLSWOKnK6OGKe9+s6oaVIjHXY+QX8p2I2S3uew0bW3BFpkeAr
LBCU25iqeaoLEOGIa09DVojd3qc/RKqr4P11173R+7Ub05YYhuIcSv8e0d7qN1sO
1+lfoNMVfV9WcqPABmOasNJ+ol0hAC2PTgRLy/VZo1L0HRMr6j8cbR7q0nKwdbn4
Ar+ZMgCgCcG9zCMFsuXYl/rqobiyV+8U37dDScAebZTIF/xPEvHcmGi3xxH6g+dT
CjetOjJx8sdXUHKXGXC9ka33q7EzQIYlZISF7EkbT5dZHsO2DOMVLBdP1N1oUp0/
1f6fc8uTDduELoKBRzTTZ6OOBVHeZyFZMMdi6tA5s/jxmb74lqH1+jQ6nTU2/Mma
hGNxUuJpyhUHezgBA6sto5lNeyqc+3Cr5ehFQzUuwNsJaWbDdQk1v7lqRaqOlYjn
iomOl36J5txTs0wL7etCeMRfyPsmc+8HmH77IYVMUOcPJb+0gNuSmAkvf5QXbgPI
Zursn/UYnP9obhNbHc/9LYdQkB7CXyX9mPexnDNO7pggNA2jpbEarLmZGi4grMmf
AgMBAAGjggGCMIIBfjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBTwnIX9
op99j8lou9XUiU0dvtOQ/zAfBgNVHSMEGDAWgBQD3lA1VtFMu2bwo+IbG8OXsj3R
VTAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
MHYGCCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNl
cnQuY29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20v
RGlnaUNlcnRHbG9iYWxSb290Q0EuY3J0MEIGA1UdHwQ7MDkwN6A1oDOGMWh0dHA6
Ly9jcmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RDQS5jcmwwPQYD
VR0gBDYwNDALBglghkgBhv1sAgEwBwYFZ4EMAQEwCAYGZ4EMAQIBMAgGBmeBDAEC
AjAIBgZngQwBAgMwDQYJKoZIhvcNAQELBQADggEBAAfjh/s1f5dDdfm0sNm74/dW
MbbsxfYV1LoTpFt+3MSUWvSbiPQfUkoV57b5rutRJvnPP9mSlpFwcZ3e1nSUbi2o
ITGA7RCOj23I1F4zk0YJm42qAwJIqOVenR3XtyQ2VR82qhC6xslxtNf7f2Ndx2G7
Mem4wpFhyPDT2P6UJ2MnrD+FC//ZKH5/ERo96ghz8VqNlmL5RXo8Ks9rMr/Ad9xw
Y4hyRvAz5920myUffwdUqc0SvPlFnahsZg15uT5HkK48tHR0TLuLH8aRpzh4KJ/Y
p0sARNb+9i1R4Fg5zPNvHs2BbIve0vkwxAy+R4727qYzl3027w9jEFC6HMXRaDc=
-----END CERTIFICATE-----
EOF
}

data "alicloud_waf_certificates" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
