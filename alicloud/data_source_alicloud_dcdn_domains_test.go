package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDcdnDomainDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_dcdn_domains.default", strconv.FormatInt(int64(rand), 10), dataSourceDcdnDomainsConfigDependence)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_dcdn_domain.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_dcdn_domain.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_dcdn_domain.default.id}"},
			"enable_details": "true",
			"status":         "online",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_dcdn_domain.default.id}"},
			"enable_details": "true",
			"status":         "offline",
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_dcdn_domain.default.id}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_dcdn_domain.default.id}-fake",
			"enable_details": "true",
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_dcdn_domain.default.id}"},
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
			"enable_details":    "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_dcdn_domain.default.id}"},
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}" + "fake",
			"enable_details":    "true",
		}),
	}

	var existDcdnDomainsMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     fmt.Sprintf("tf-testacccn-hangzhou%d.alicloud-provider.cn", rand),
			"domains.#":                   "1",
			"domains.0.gmt_modified":      CHECKSET,
			"domains.0.cname":             CHECKSET,
			"domains.0.id":                fmt.Sprintf("tf-testacccn-hangzhou%d.alicloud-provider.cn", rand),
			"domains.0.domain_name":       fmt.Sprintf("tf-testacccn-hangzhou%d.alicloud-provider.cn", rand),
			"domains.0.ssl_protocol":      "on",
			"domains.0.status":            "online",
			"domains.0.scope":             "overseas",
			"domains.0.cert_name":         CHECKSET,
			"domains.0.description":       "",
			"domains.0.ssl_pub":           CHECKSET,
			"domains.0.resource_group_id": CHECKSET,
		}
	}

	var fakeDcdnDomainsMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"domains.#": "0",
		}
	}

	var DcdnDomainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dcdn_domains.default",
		existMapFunc: existDcdnDomainsMapCheck,
		fakeMapFunc:  fakeDcdnDomainsMapCheck,
	}

	DcdnDomainsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, nameConf, resourceGroupIdConf)
}

func dataSourceDcdnDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "domain_name" {
  default = "tf-testacccn-hangzhou%s.alicloud-provider.cn"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_dcdn_domain" "default" {
  domain_name = "${var.domain_name}"
  sources {
    content = "1.1.1.1"
    port = "80"
    priority = "20"
    type = "ipaddr"
	weight   = 10
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  scope  = "overseas"
  status = "online"
  ssl_protocol = "on"
  ssl_pri = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAp7YDnKCEMAu8mIo0o4iFGhriHp7DJfcrqueCEvNILv/GooBc
tFKexN6Zp+BRp4sKxntauqrpiX+vHc5V7ON5lV6AayfgvcJnAsk9pJ+aAU3Oe/uc
+ntoaYSxCGL5SRQm+wQPLFeanZVydigMO6Khfidx3fz9EGYAPYWhr81uVqoqGXEb
EvdQnVphqus3u9u3xeu6PKyioGzXcn6+hlzGRRhB6aJMxQ4RLA+A7Ia1eTz233tq
xL5jZCgllx+Sw97zoeodT3muxBhWqPR3spRKH5nkEZ8lYSYBCUUUQlZI6Io5Rjsz
tPHNxiNDoytBSqJeOAJQYAu4Mf7nZSwigDa7RwIDAQABAoIBAAs/NpDLZvH954Dn
S85nul1czis1hGrIX6JPcjapH/8e4ghFyXHCVKlpMC7E6VTuCyPyY8w+5/hzmp/K
FZMUUjQFKWGGRBkVr2jNbBfdKCvMNvuzjPxzSDZDUsf4MzWGZ3LP++CCY3kL66gm
2WMqbeAS7xzu+V2fKYb2rjgm865WET0Hv9dpJTwVNqgP73lAFtl0jENZisFT81YR
e4gC5rsks7bVjQHHIh3p6nIcZ6Zd6SpnOv6dEsWjPSBeJhAmR4QLkBRGpBPDr1YW
uSm8AqRvKRprRHqwvSLMN9AQFjFjTGOPYJJJil7z7qefdZPFEUdXB8sEpst9llK+
6BiS0OECgYEAyTNF9ea1b6y4pl42rlSLuymqLDNp2mx8BMLF25INyAVrpnikxSYE
x5k5Z+kXE+AD4nKOO45qapIwu3Mejvapb60h+nB0Uxsx3H+9poy/Le75lIjaGU/o
UwAIFb22bNiocm18mSwavbEtVlTsR1gZoNvz3HhhKDIs2AzEjKwrS7cCgYEA1WOx
ls4wBGZTK5XyXWpd5vf5Oio9l4PfxqDDVhT144KiXSScnUU1rbL8wOHzwGzOKeaV
YaoBLde7fbwhuaJvYA527mYsMFueeg7EHKs11o7EN+j5TTDqR18xE3iSrkBfuNEv
mN1o4KOmjsuH5CY4h4PDetVoovrY4sQmajsXLPECgYAR4gXI2m2r9F0hJGSV0Bvv
Ub+3WAaDjHrlbW5qmquw6JJt5HE4uK1aFEte6f/MG3Ac83Oi5YCd4kqEjrHboR7k
Ny469T3RmSwwXgY8RGxFp+T1B8ji0RBkOC9/xzHssMEgEo0tjBcAXzwZXUj2+mSk
wIgHQ4fXK8aCmXfqzO64NwKBgQCd4yzkY1099CQ3zLPOkMQ4AGSkt9powEeT5SGD
EPE6zE6sUkmbSDlGc3f2k3jSeO82K4l+ANbsf4IXr1rYyqpTzYAMNwcdJL0mnMRY
Xgnw3iOrJrNHfRjrhDCAsqb9TV5GFml8Vt6h0BSN9WRv2CPdiQ3bVgodBTPy3aV6
1ov4UQKBgQCzGBg5hTGxSeU2giT4A3HkPLdSpWM6vfnCFsh4J8nkd0TNI7DLbBME
NAp8ZaNLnFRwjGU9vePKNsovnWD69MbEtJK23OgV7NheSYE33XEXJYU8zi8Tvuha
qhylxeghk9affB7FMAzWuGTeUgpRgwGxkWaVQyCoOzisewta9BXADQ==
-----END RSA PRIVATE KEY-----
EOF
  ssl_pub = <<EOF
-----BEGIN CERTIFICATE-----
MIID7jCCAtagAwIBAgIQBL5mm8p4TtWVOUtL8meCkzANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNDEyMjMwNTU3NDBaFw0yOTEyMjIwNTU3NDBaMCwxCzAJBgNVBAYTAkNOMR0w
GwYDVQQDExRhbGljbG91ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAKe2A5yghDALvJiKNKOIhRoa4h6ewyX3K6rnghLzSC7/xqKA
XLRSnsTemafgUaeLCsZ7Wrqq6Yl/rx3OVezjeZVegGsn4L3CZwLJPaSfmgFNznv7
nPp7aGmEsQhi+UkUJvsEDyxXmp2VcnYoDDuioX4ncd38/RBmAD2Foa/NblaqKhlx
GxL3UJ1aYarrN7vbt8XrujysoqBs13J+voZcxkUYQemiTMUOESwPgOyGtXk89t97
asS+Y2QoJZcfksPe86HqHU95rsQYVqj0d7KUSh+Z5BGfJWEmAQlFFEJWSOiKOUY7
M7TxzcYjQ6MrQUqiXjgCUGALuDH+52UsIoA2u0cCAwEAAaOB2TCB1jAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFCiBJgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEF
BQcwAYYVaHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8v
Y2EubXlzc2wuY29tL215c3NsdGVzdHJzYS5jcnQwHwYDVR0RBBgwFoIUYWxpY2xv
dWQtcHJvdmlkZXIuY24wDQYJKoZIhvcNAQELBQADggEBALKsozVNUed06FDSlH7E
VBWqUZPHwB87t7519bVd8/5jMm6S8U6v4HtsmtVy3QcXD5Abxhl6uPH2kqxfb15L
6KhU1OemKmoFNX3/7zehnuh7e4/0Ild6smxsj8jjjHKRVFFM9eXBn4PPP+MAX0QI
oLeyYZcNv8ZOtRu6li50VlwHbLOpif4flTOfA0rSBnnXwrT8WUMKliHdKu454CcI
HBwi47xKbHW9a0vyPawy5YNNQUwDG641PcYIKq89BxepBPMpMhNxQJv1yOYvkIIx
Acu4FFkxsr/BH1Gsg2CE5ecN2Dr8+D96czG7w4NGLoPnh2LXydMiRKtyjQRFFeax
HQQ=
-----END CERTIFICATE-----
EOF
}`, name)
}
