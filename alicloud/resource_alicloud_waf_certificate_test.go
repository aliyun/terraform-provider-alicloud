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

func TestAccAlicloudWAFCertificate_basic0(t *testing.T) {
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
					"certificate":      `-----BEGIN CERTIFICATE-----\nMIIHrjCCBZagAwIBAgIQC2/QuNyfdGuiLP17UTTLyjANBgkqhkiG9w0BAQsFADBc\nMQswCQYDVQQGEwJVUzEXMBUGA1UEChMORGlnaUNlcnQsIEluYy4xNDAyBgNVBAMT\nK1JhcGlkU1NMIEdsb2JhbCBUTFMgUlNBNDA5NiBTSEEyNTYgMjAyMiBDQTEwHhcN\nMjMwMzA4MDAwMDAwWhcNMjQwMzA3MjM1OTU5WjAfMR0wGwYDVQQDExRhbGljbG91\nZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMbN\nO7Q4hF5PYG/9fGqrI06P4bwPW1ZhRPR1Wg1f9hHIQot17X1uHqhmkKtzaqT9U7P2\nFdRkOqOynLFNujTKohuy5R7Diuv9hAnG+t33wooAQscYj+2x8OMbPw5QoN9A1Ddx\nR5c2RZHNZVjVL5aPMqssm0bARF9vKhPcZiHe6Cb/Xn1YwLLkodZ68kxwwJAGm/wO\nwVesmQVOeUPxv/Dh8R8GBVPwq0ewOs+UwJhrBNnyfI9Pey8fGAaEOz3Ucu3KMcXt\nA6LrsWGcXRHKH2zClzuXWzGWkEfrYl8FNKE8rop8bf7pFL2wNjxVgAW3tSAzzXQe\nclp3YEJ5/KfJsYp0z4ECAwEAAaOCA6cwggOjMB8GA1UdIwQYMBaAFPCchf2in32P\nyWi71dSJTR2+05D/MB0GA1UdDgQWBBQAu+tDmtZxFyGaT2ZZzNvgxMlgcjA5BgNV\nHREEMjAwghRhbGljbG91ZC1wcm92aWRlci5jboIYd3d3LmFsaWNsb3VkLXByb3Zp\nZGVyLmNuMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYB\nBQUHAwIwgZ8GA1UdHwSBlzCBlDBIoEagRIZCaHR0cDovL2NybDMuZGlnaWNlcnQu\nY29tL1JhcGlkU1NMR2xvYmFsVExTUlNBNDA5NlNIQTI1NjIwMjJDQTEuY3JsMEig\nRqBEhkJodHRwOi8vY3JsNC5kaWdpY2VydC5jb20vUmFwaWRTU0xHbG9iYWxUTFNS\nU0E0MDk2U0hBMjU2MjAyMkNBMS5jcmwwPgYDVR0gBDcwNTAzBgZngQwBAgEwKTAn\nBggrBgEFBQcCARYbaHR0cDovL3d3dy5kaWdpY2VydC5jb20vQ1BTMIGHBggrBgEF\nBQcBAQR7MHkwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBR\nBggrBgEFBQcwAoZFaHR0cDovL2NhY2VydHMuZGlnaWNlcnQuY29tL1JhcGlkU1NM\nR2xvYmFsVExTUlNBNDA5NlNIQTI1NjIwMjJDQTEuY3J0MAkGA1UdEwQCMAAwggF+\nBgorBgEEAdZ5AgQCBIIBbgSCAWoBaAB1AO7N0GTV2xrOxVy3nbTNE6Iyh0Z8vOze\nw1FIWUZxH7WbAAABhsATxUwAAAQDAEYwRAIgO0/k4INTQs3bauRtx0bz7BkmcDDE\n5xjCzPPqMbxnwYsCIHcVsaYOZO73eiosYrTolY5/XIUWVDYfLmDia/bxtc9QAHYA\nc9meiRtMlnigIH1HneayxhzQUV5xGSqMa4AQesF3crUAAAGGwBPFowAABAMARzBF\nAiBO/OkeQCRMOU8UiRToVjnnC1F4lKvjChQKnlgDi/26ggIhAJOTO3zod7m0ILNY\n8Ue/jiAwmX3x66dEU+pIXuMSqoi5AHcASLDja9qmRzQP5WoC+p0w6xxSActW3SyB\n2bu/qznYhHMAAAGGwBPFeAAABAMASDBGAiEAjsP0zTN0Gu9D10an9u+KysJAmcnN\nSG7djcEpOuf0ov0CIQD+6s3lC2Qj4ddNUQ12MRMpytJO4a440LQypfDQtEGcHTAN\nBgkqhkiG9w0BAQsFAAOCAgEAcD60o2Astoawsg8MV3ilTBVf8yYaWXQ4uFy5/xwO\nKSQb9Tb1ZWVyB+zJcHsJTJNVphC4rDzd2bGTDir+rcZlSDSTU9/xn/0g26I28dNq\n86bGJOvWVMJmo9WbtQPbEncR9K/Vc0eQteX8q8DZsTPyAXIJSwVOXGpsqDRmZPvT\n65mcg5KxTzu7MmvFkD8oPQrp9XmURxu1Z9Fm6YOr5B6/CZquY9fjNGUP52lzxU7j\n7qbQsGTGMRAQ7HcJvetanP3DXYOi/nFdUPULsmTXbEVLV0c8nrsLtjAJI0pZaKkx\npqjTr8I6jmOLrzj0O4NymXe54kb0ch65l85A9Pauc7mmTT87gpEjuCy1FXlPpsJZ\nxi3uANL6SzPKnma0GrJbBhaX7HpiroxuHU50fr/umdQVUqBldfcRAaY5OfcdB1zU\n7rCQXfeSvR98yr2pLu9xVUs1EFvMD5XzPk1WA3aZDRCEenhYykVFFrog37UVbxJA\nHvOPx9qzcy23Oz4BQu2sMNmJKV8zBxCDxTpEiGjIuDykDmvgeS4uhKsMU4tPRpXi\nsw6G6grVJ4YmiiA5uUrBbXhZFw8IYxk27av9b8kEiUOv9ytdN3Ay84E9HwWDld5B\ngvACxJTy8IyS3s4vi/HRS9k4IMgCZLUvlGFi3G8VPl/cFQSUeurMyQxuMbqYrXcL\nt9Q=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIFyzCCBLOgAwIBAgIQCgWbJfVLPYeUzGYxR3U4ozANBgkqhkiG9w0BAQsFADBh\nMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\nd3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD\nQTAeFw0yMjA1MDQwMDAwMDBaFw0zMTExMDkyMzU5NTlaMFwxCzAJBgNVBAYTAlVT\nMRcwFQYDVQQKEw5EaWdpQ2VydCwgSW5jLjE0MDIGA1UEAxMrUmFwaWRTU0wgR2xv\nYmFsIFRMUyBSU0E0MDk2IFNIQTI1NiAyMDIyIENBMTCCAiIwDQYJKoZIhvcNAQEB\nBQADggIPADCCAgoCggIBAKY5PJhwCX2UyBb1nelu9APen53D5+C40T+BOZfSFaB0\nv0WJM3BGMsuiHZX2IHtwnjUhLL25d8tgLASaUNHCBNKKUlUGRXGztuDIeXb48d64\nk7Gk7u7mMRSrj+yuLSWOKnK6OGKe9+s6oaVIjHXY+QX8p2I2S3uew0bW3BFpkeAr\nLBCU25iqeaoLEOGIa09DVojd3qc/RKqr4P11173R+7Ub05YYhuIcSv8e0d7qN1sO\n1+lfoNMVfV9WcqPABmOasNJ+ol0hAC2PTgRLy/VZo1L0HRMr6j8cbR7q0nKwdbn4\nAr+ZMgCgCcG9zCMFsuXYl/rqobiyV+8U37dDScAebZTIF/xPEvHcmGi3xxH6g+dT\nCjetOjJx8sdXUHKXGXC9ka33q7EzQIYlZISF7EkbT5dZHsO2DOMVLBdP1N1oUp0/\n1f6fc8uTDduELoKBRzTTZ6OOBVHeZyFZMMdi6tA5s/jxmb74lqH1+jQ6nTU2/Mma\nhGNxUuJpyhUHezgBA6sto5lNeyqc+3Cr5ehFQzUuwNsJaWbDdQk1v7lqRaqOlYjn\niomOl36J5txTs0wL7etCeMRfyPsmc+8HmH77IYVMUOcPJb+0gNuSmAkvf5QXbgPI\nZursn/UYnP9obhNbHc/9LYdQkB7CXyX9mPexnDNO7pggNA2jpbEarLmZGi4grMmf\nAgMBAAGjggGCMIIBfjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBTwnIX9\nop99j8lou9XUiU0dvtOQ/zAfBgNVHSMEGDAWgBQD3lA1VtFMu2bwo+IbG8OXsj3R\nVTAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC\nMHYGCCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNl\ncnQuY29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20v\nRGlnaUNlcnRHbG9iYWxSb290Q0EuY3J0MEIGA1UdHwQ7MDkwN6A1oDOGMWh0dHA6\nLy9jcmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RDQS5jcmwwPQYD\nVR0gBDYwNDALBglghkgBhv1sAgEwBwYFZ4EMAQEwCAYGZ4EMAQIBMAgGBmeBDAEC\nAjAIBgZngQwBAgMwDQYJKoZIhvcNAQELBQADggEBAAfjh/s1f5dDdfm0sNm74/dW\nMbbsxfYV1LoTpFt+3MSUWvSbiPQfUkoV57b5rutRJvnPP9mSlpFwcZ3e1nSUbi2o\nITGA7RCOj23I1F4zk0YJm42qAwJIqOVenR3XtyQ2VR82qhC6xslxtNf7f2Ndx2G7\nMem4wpFhyPDT2P6UJ2MnrD+FC//ZKH5/ERo96ghz8VqNlmL5RXo8Ks9rMr/Ad9xw\nY4hyRvAz5920myUffwdUqc0SvPlFnahsZg15uT5HkK48tHR0TLuLH8aRpzh4KJ/Y\np0sARNb+9i1R4Fg5zPNvHs2BbIve0vkwxAy+R4727qYzl3027w9jEFC6HMXRaDc=\n-----END CERTIFICATE-----`,
					"private_key":      `-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAxs07tDiEXk9gb/18aqsjTo/hvA9bVmFE9HVaDV/2EchCi3Xt\nfW4eqGaQq3NqpP1Ts/YV1GQ6o7KcsU26NMqiG7LlHsOK6/2ECcb63ffCigBCxxiP\n7bHw4xs/DlCg30DUN3FHlzZFkc1lWNUvlo8yqyybRsBEX28qE9xmId7oJv9efVjA\nsuSh1nryTHDAkAab/A7BV6yZBU55Q/G/8OHxHwYFU/CrR7A6z5TAmGsE2fJ8j097\nLx8YBoQ7PdRy7coxxe0DouuxYZxdEcofbMKXO5dbMZaQR+tiXwU0oTyuinxt/ukU\nvbA2PFWABbe1IDPNdB5yWndgQnn8p8mxinTPgQIDAQABAoIBAAeQeZbC68PMO2bS\nJgYdBZaQrPiDO74kAgQB5G+hLh4WbknGLpzY2tEc40uoHu+w3/AiC0nt7Q6bgOYn\nh0/zlzV8mhvl5DGp5M3xGfc7h9XJFWwOnYi1fkKgn8aYqiEJ43P0myfAHeeNVto0\nH7/Cf8gO8MhBZUNIWePa6Aaeklt2x3PD2iGHh9XxKJkIXykRv0MwaeScp4rzgiDU\nb0dXhHb1DzAnNEUf98DxA9X+dRb/JFaI7hDH+7zt+rV6idwR2MVlJnRAUunRUpEZ\ntEiAXNzUWwIe7Hzj+jRUj7kHgwyfgR2xalF8ZZunBsW13WPrlNCyxOhC1ZKzvwet\nktlHZU0CgYEA+8l6JtHOr/Pid2P6VUgP1edqr1QDa9QJ49hJP9MDYFsfJlvR3C1G\natGRDGuNlLE+Mht33wRUmkqqWrRq7WAS9cFt6yfIyQ8CzREr9ZY/AGri1SWL0yGD\nH9BXSRUmmwRICOwd0vOcz0gEnBQPhxPeVEaBprSBneYi/Zj+a93Hwh8CgYEAyiDL\ne5buc6IXakqtW0oPlVucvSVhOuvDBLzfexv5WGywjILpJI0Pamb9tDVjz8HiLgO+\nf9QSihidyiTh8QKs/BqcVcMW8qR/lZIyOnKgvc5snBqx8aDeAFgfWqIZjwKsq7eu\nEO7s/2ZFn2Tg+Yrh57+jRPSFgBcC8fn14/UPel8CgYEA26suyBl54qa7aQAv0iGZ\ntzw8JpT8myGWM1NqrKVTp8g4CAZJtHdHnUAS9SwMjLKKGbs/PLIRgb8smxAWzdxp\nDkHgvS2rjkZrYi/eE1guxRV7qwwjzmLJfIMO8/LhCJOWqToe1wG1SF6DrwwNdALQ\nkOQmjyOihfjXhDrrlX4bnm8CgYEAoeOGVmRac47dhiptknJM7OlFLlEkANcXiVHl\nBGsxWFslREUNilLYh/YZOR6R5LJ6/zPMfgFTH/v7VQ4ZULk952VA1Ye9d3W9IgTH\n6fwpNhyA7L6MyuR3KntqEWNHaT1RIu+ooLNcamp1VOatMfDEYqMgl3fo/OBxYXE8\nOenrDa0CgYEAtB0RHiQJMi2n0AtOGuDbV8HRxk9Mh0GypjD3LGsdduG68bWe8WeT\nxN9ZNUeHT2VXSF1ra0OZvyAdTZnzxwgZFHlKSbh5Fzta1fnoXj8x8akp1mycvJL7\nEhc7Tyxrl5vrxONKMbQdpC45rPxyKeGmk+b1LdbvoWr2iQwqL/D9Lr0=\n-----END RSA PRIVATE KEY-----`,
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

func TestAccAlicloudWAFCertificate_basic1(t *testing.T) {
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
data "alicloud_waf_instances" "default" {}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = "tf-testaccSslCertificate"
  cert = <<EOF
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
  key = <<EOF
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
