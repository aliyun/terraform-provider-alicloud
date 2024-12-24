package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_domain",
		&resource.Sweeper{
			Name: "alicloud_dcdn_domain",
			F:    testSweepDcdnDomain,
		})
}

func testSweepDcdnDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeDcdnUserDomains"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var allDomains []string
	var response map[string]interface{}
	for {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_domains", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains.PageData", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if strings.HasPrefix(item["DomainName"].(string), "tf-testacc") {
				allDomains = append(allDomains, item["DomainName"].(string))
			} else {
				log.Printf("Skip %#v", item["DomainName"].(string))
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	action = "DeleteDcdnDomain"
	for _, domain := range allDomains {
		query := make(map[string]interface{})
		request = make(map[string]interface{})
		query["DomainName"] = domain

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, domain, action, AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}

func TestAccAliCloudDcdnDomain_basic0(t *testing.T) {
	var v dcdn.DomainDetail
	resourceId := "alicloud_dcdn_domain.default"
	ra := resourceAttrInit(resourceId, DcdnDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DcdnDomainBasicdependence)
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
					"domain_name": name,
					"scope":       "overseas",
					"status":      "online",
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name,
						"sources.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"ssl_pub":      "${var.public_key}",
					"ssl_pri":      "${var.private_key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"ssl_pub":      CHECKSET,
						"ssl_pri":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
						},
						{
							"content": "2.2.2.2",
							"type":    "ipaddr",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
						},
						{
							"content":  "3.3.3.3",
							"port":     "8080",
							"priority": "30",
							"type":     "ipaddr",
							"weight":   "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "cert_type", "check_url", "ssl_pri"},
			},
		},
	})
}

func TestAccAliCloudDcdnDomain_basic1(t *testing.T) {
	var v dcdn.DomainDetail
	resourceId := "alicloud_dcdn_domain.default"
	ra := resourceAttrInit(resourceId, DcdnDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DcdnDomainBasicdependence)
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
					"domain_name":       name,
					"scope":             "overseas",
					"status":            "online",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Domain",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":       name,
						"sources.#":         "1",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Domain",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "Domain_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "Domain_Update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "cert_type", "check_url", "ssl_pri"},
			},
		},
	})
}

var DcdnDomainMap = map[string]string{
	"resource_group_id": CHECKSET,
	"ssl_protocol":      "off",
	"scope":             "overseas",
	"status":            CHECKSET,
	"cname":             CHECKSET,
}

func DcdnDomainBasicdependence(name string) string {
	return fmt.Sprintf(`
	
variable "name" {
	default = "%s"
}

variable "private_key" {
  default = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAu1c9Uv+nLlTzRbSwII2F/sMKv7hJZt8Dw2L4Xx7OTrKIeAxj
d58+iRIgcuoidjHmAbRHZBFET0A7ZjjSmtsUPziHDFr9b3MaaZEI2hl6To0DJ+IO
i+DTUBbeZni61U54m8NLax24zSFdvqvi7xUSlGtJ9Fn5BYJDiavD3ykvJlFUX2vo
EAAScJX7OUoSzuNIduWZ11zkDKM6hTZtnEn1AorOh/zNzIRqwEzX1nTqLcJmH+fl
TJA4OnOGg4g6484eLaxn0ucTYQNLYgHkaJHazaS77c6UfWKwRoFhZegQEZXEyGr8
IK2CpNbdc36P2oah9Oe8T09Y236QHMUH/yHNHQIDAQABAoIBAQCx3tLKyxDgXKfd
twDC55whlu3Nuht3IKdiC8XmCkm3Tqtjz99g5EFrw1orwUGXFyla1OAzknFZDZNY
KvtLLFa8797JTFr0RkT9lkbhTO9jRV+JrohBJuV7VTsz78z0Wd0JhxNEUKP1n4hy
UKDWfxt076j357UYFeYqAHuolmG97jdUBW5Imv279DRCQK7Qc68FbeLO/gu5IJsP
eKG68BO+q14VNV+S6ekXCHn4wnxwHc0oduh9CmP0eofcx76nwmMMiPFMBg9iLxh4
o0MoXZ4z7ZsY2yfFFB58unI5+lbVIo+gwGz03BHrWfhYmJpzQ0KzvMdUeFd7bqx/
q+4y2JwBAoGBAMICYULTEqqz1L8UI/RTQ0NzTvaOO1b0j1gtV7ilpG1JuGXl7U2S
3242JC39rYz2dy1PggpoDBovVwUgDjwIEvCAb6oaDVNxTUJxTmN9/jMbd0FhMTZu
PpbAmhMCbr2psmpJ87DLQjHoZwsP14D6TOut1hZNrG0UskwoDsMmZbONAoGBAPcz
Yvmzo/t9WxIv27sDnNywyoqMNHcUE80UticXnktuC+OHc80Tn5rqMYvCR9YhAJXu
2+Dd2xIAdCLD+ucpqmn6d+2fm9NjHIyP0mTCrr+miydzDFlAubWHYXwesvuVI2Id
xDB0Yu1wW6nlMhJQ3URNhafHULhZsqBRS6Nx0dPRAoGBAIj2GSeNzu3Hindihodj
iGbDrokMnAOlHtUHHZhzB4NHue/lxAMxnp41hpEZNz3+eN/58znZfkG2Dd7GZIYo
xQYYBby2K5YutHYle0ttlNkLmMMFFDLy3Siby6mD3B31AMlcb7btp0uIX8ZFZsPc
8BSpYivYpdNT+xMcbF+EaeO5AoGBAISplC1TheZ6cLyC6JYlqzIYwqm18pYRNUsz
GUpDd5Udes3hrHjbViVKF8rcObclwO214VR9W4r+qVTa/jS+fJEhdOkWZgb8wp6A
tLWUcTmzBCzopjDj9oYAIIX+56jycam/NcGXRFwOl3LG6KdBtG1qeRcAdUZqBN3a
oxAVDjlxAoGBAKg9wGcU1OgnCyOzUJwksMRuSxZT8Cc2Lqo5QN1jJjsx1bOhE2kU
fLiVkG9Qo44dx/cs9EbYIKlUxfkzjrUcIUMKSvi8fCJ751Q2Mf6NpurD1tNFqdjf
D9z9Rp1EGnjVjphgyysISwgunr0g78220JP/ZJOmPGqacQsvqzthveiX
-----END RSA PRIVATE KEY-----
EOF
}

variable "public_key" {
  default = <<EOF
-----BEGIN CERTIFICATE-----
MIID7zCCAtegAwIBAgIRAJzNPvPgpE3Bg7DjYcTQ17gwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjQxMjIzMDU0OTM0WhcNMjkxMjIyMDU0OTM0WjAsMQswCQYDVQQGEwJDTjEd
MBsGA1UEAxMUYWxpY2xvdWQtcHJvdmlkZXIuY24wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQC7Vz1S/6cuVPNFtLAgjYX+wwq/uElm3wPDYvhfHs5Osoh4
DGN3nz6JEiBy6iJ2MeYBtEdkEURPQDtmONKa2xQ/OIcMWv1vcxppkQjaGXpOjQMn
4g6L4NNQFt5meLrVTnibw0trHbjNIV2+q+LvFRKUa0n0WfkFgkOJq8PfKS8mUVRf
a+gQABJwlfs5ShLO40h25ZnXXOQMozqFNm2cSfUCis6H/M3MhGrATNfWdOotwmYf
5+VMkDg6c4aDiDrjzh4trGfS5xNhA0tiAeRokdrNpLvtzpR9YrBGgWFl6BARlcTI
avwgrYKk1t1zfo/ahqH057xPT1jbfpAcxQf/Ic0dAgMBAAGjgdkwgdYwDgYDVR0P
AQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME
GDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB
BQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov
L2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MB8GA1UdEQQYMBaCFGFsaWNs
b3VkLXByb3ZpZGVyLmNuMA0GCSqGSIb3DQEBCwUAA4IBAQABSp2RLQD+NEudQ1Z2
yhCD9ADbdrWQHPBZgtUV0EjN4gMYucz7dzWo1xjg5BhKd8naku21U2ZUa8TgnIgt
IK+GL8gLex4iXq9CiZqZsFhYnuopR0ISULtC+Oz+YfrKfzMHDK9UU3AZT8bKT4mm
T9nAWV5Fa4Ik1HlA0kykNVrNCef+zLT4W7x/YMSPIMUDHRMeGXOEPnqIOBnR0ha+
KDhZPviYhN2M4u0tVVb/2NBQLYgVLspj28dQShBlrXC51SurAwmnw5gcVlJG3r1H
b494lL9Ycx+Q3rlziqYLMYq3+8x+bNQhI1iDjWeYtVoG2qyX4Q8l5IOWQ/mKtYz1
nf8k
-----END CERTIFICATE-----
EOF
}

data "alicloud_resource_manager_resource_groups" "default" {
}
`, name)
}

// Case SetDcdnDomainSSLCertificate替换 7246
var AlicloudDcdnDomainMap7246 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudDcdnDomainBasicDependence7246(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = "casdcdnterr"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID2zCCAsOgAwIBAgIRAN2UzXkVHEkUn0xJezHFFoMwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjQxMjIzMDU1MjU1WhcNMjkxMjIyMDU1MjU1WjAiMQswCQYDVQQGEwJDTjET
MBEGA1UEAxMKcGZ5dGxtLnh5ejCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAMA7x0zoPGObPQYsrAoVJzXOeete55XG2YJWVdxfcivWbdIeG144BRTcrah4
pheJTjTVS4MiS5Q3iV5Zq9n/cPUFyrMWF75NUylOHGss9MNoemFeji0WPnRWij4c
s4WfvHb7W++YUyNFIqQxt0yOscrbtmu0gUsZZyF7lXWe2FO+tbrNv17jxHSduiey
86T806IfDcTDabOEzEGnbO3P2BjlhRxexKJf7Q4MNJ5+K+N5wGNZzml1hz3/vSwh
T6etbdpcFqy5UoA/c9pS4RZfL0udx4XHoORtyP/KzlvQ0yw7wBN+mIDfCj6Is6sw
X4qi+ZVZ1DvoZBCl45H78NzZVrUCAwEAAaOBzzCBzDAOBgNVHQ8BAf8EBAMCBaAw
HQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiBJgXR
NBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYVaHR0
cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlzc2wu
Y29tL215c3NsdGVzdHJzYS5jcnQwFQYDVR0RBA4wDIIKcGZ5dGxtLnh5ejANBgkq
hkiG9w0BAQsFAAOCAQEAjIJ2L6Y9F4RqObYd1DwiK7UfEn2tAyP3LntYLdTD9RZv
OymX1ZFHDvPdM1ulCFHI5wkvi1PCTxgGdr4Pv3SdKRw2cwiqElPTSHcoEasT5Hv1
RcDaiKYo4zop6YfI2qeFc3tDAAZTgubZl05KpRKDfbw+s4XZ8fS4/k0RRW+OORUj
BFlZFjx9jP8TWTbsD08i5dQ0Jk/mbbmkYWH1UMONd3JG7xHkgfTG+og5+WANRiXX
6ohwYMoGzw4RMXstC4AXG1JnICJffNKrZRz6zSFgXdu3+Ue3jYqSbQ/6M19Esb27
kTyVMey/vBCc1ocVu7OVRLPoNLbaoIlAX+iqGn8HcQ==
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwDvHTOg8Y5s9BiysChUnNc55617nlcbZglZV3F9yK9Zt0h4b
XjgFFNytqHimF4lONNVLgyJLlDeJXlmr2f9w9QXKsxYXvk1TKU4cayz0w2h6YV6O
LRY+dFaKPhyzhZ+8dvtb75hTI0UipDG3TI6xytu2a7SBSxlnIXuVdZ7YU761us2/
XuPEdJ26J7LzpPzToh8NxMNps4TMQads7c/YGOWFHF7Eol/tDgw0nn4r43nAY1nO
aXWHPf+9LCFPp61t2lwWrLlSgD9z2lLhFl8vS53Hhceg5G3I/8rOW9DTLDvAE36Y
gN8KPoizqzBfiqL5lVnUO+hkEKXjkfvw3NlWtQIDAQABAoIBABpgGVhUBPUlt5nB
R1mazWZ0jgXdX6kNP4rCjcVO0ztwkGDkAJ1M0mWqYalb5G4WSMS2/0Vezz/m3tIz
O4ENq1HzGXy460kREvf3365U3MBy9VemwZsuEiOkPBOJnJgY8qLgmhylqcKNGdOt
fpjie0J6Iu1kNtk3Aw91BWy9/rB+nNcJsUx+5KgSRtQ5hnYEkUMP5yoGF1Tr8pSV
o/meDBCQTr3Y6GZpmHs7IPsFiErexsd8sfUS52evDrJtClPg0M6Gpjv65gU5kq3v
+jbHEcRe98YIn7noXfLOUCDGYQVSfyxou1czYCOdx7enIE0tNqk46j8mBVc6OhMA
XkVCRYECgYEA+1OZfO3gcinnKtFMpVCmDi64DU4l0D0vwXyCt6TdBcxPgMh51tOr
XWOQTU4uf941kKDIVMdzW1agG242HUAz+73AS0H/+QjzoYf//PSsWje/xCccfVS3
sjXBKhMoVwIrWempeakGIfOi7yjCv1kkQstPJ8Qokf3/h3cgWyHdXncCgYEAw87g
NQ4pGLLUMyNfK0SJlPLp5jdghY+/a1MHxFaOcjmOv5jqN57Rvub3ADEmlMcda7W3
A5qU8qmx5gbsM0oBDl9bbmqxF95BHqJrQcaPDw7V7HOcZf6gdMNPtSn0DK9uG31s
IZQDD7kxD4PT4/FWhPRqsFZ9L/zMIWwZCxfk4zMCgYAvFLEbIyC5ojno3n6CNYJ2
A7B85ZfV07B/iYifSGYTMPvvvx577PkcLIuav7ucPo9AQa5lm1tzz918ZgADKMTU
Mu6z6nA+QbwKFYUR6O/kkq782urOW7Fx0/oUnLQg4IoodMpHvS8l6xMpxDP/Tn6p
eJaid2+2MaPNx7Yq/EQQ7wKBgQCdvDeNRd0BUn3yvBncRxf170FQ/Wc58LSpBngJ
SBj0Fz3RRqPXLo+Uk4aClxWXYFdo/zdxJcO7P8xZm1YHcyQqqdKDvlru+VHIFdsF
X6i63p6iHfftihNEPFonfKZm2aN/baf/3LYionLNJss4op+p9yNC7klmsOTYP7Zk
41i1VQKBgFTVLZpG5vE3x13bxUez3ydg9NLClrxwRHiHPBUlQx/McQsBpnWrDBmU
rJyNbR+eJR+mjbbT75w137lc2lOV/K86SnhEa1E7LS/mzI7jlxw3dDS30rRYw3e0
CiHrbj05hznS76XoNzFKtTvMgrynZKIAtmrORPss7H/RlhhhAgDm
-----END RSA PRIVATE KEY-----
EOF
}

variable "cert" {
	default = <<EOF
-----BEGIN CERTIFICATE-----
MIID2jCCAsKgAwIBAgIQe1Ygh7CARyeg18X8898L7zANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNDEyMjMwNTUzMjJaFw0yOTEyMjIwNTUzMjJaMCIxCzAJBgNVBAYTAkNOMRMw
EQYDVQQDEwpwZnl0bG0ueHl6MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAmKQ1/ZoZkk2VhVTkPz350ieLvXJTFPuY/yTRzJ0CQGgS6nfm6bW5zZtavKG4
CQvqKyfFZruygBz4eZfaz2qtrq21bmB6XUE4ZYroUN6kiD4LCPxN9FLU3ARqsSo5
BqEwJ5jAJmU7u0sfSoxTAcEtZ7qNwxlC8Znyo9OL92vu/520fzf3d11Ui+LSkVIx
5TUtTg8D+vmgH9FSJw4KwDh66tOdJRH5hPkQQa066b0y+DtEtY9OShuFcKZPJhWU
YfCj52wT7xyhvwo1InbZMbsWys9S5BEx9cdOe7+oPR/DUqIuMWpN1GBDE163InEE
nQ0oeiGlOkgnFLfAPfvzVZ2s8QIDAQABo4HPMIHMMA4GA1UdDwEB/wQEAwIFoDAd
BgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHwYDVR0jBBgwFoAUKIEmBdE0
Gj/Bcw+7k88VHD8Dv38wYwYIKwYBBQUHAQEEVzBVMCEGCCsGAQUFBzABhhVodHRw
Oi8vb2NzcC5teXNzbC5jb20wMAYIKwYBBQUHMAKGJGh0dHA6Ly9jYS5teXNzbC5j
b20vbXlzc2x0ZXN0cnNhLmNydDAVBgNVHREEDjAMggpwZnl0bG0ueHl6MA0GCSqG
SIb3DQEBCwUAA4IBAQCQCzlOVoBY/RSoX9+LwhKdmPfg52muosYq5B1HOPwQbUKX
7CFWj9xNZG43yTLGzmx7PTR/qv4WjDkqC5w3hZTG6Hrd7TPCTw4WMOIF0LKyhGzd
IxUunkK09BfpTmzn6PTASX6uDjNir0/QAouxX0Lw3OiY+0JmavgShLcOwTvM5Nwm
efLIabQCstRhyaL2ZaR2T9feXLwnJTEEPNDm1J1rZAXBOlIPscJB8WGj5H5YjABc
BPFkRxW7eVxVqV7saP57bc+TfEqXizE9Ky5dYATA//diQ63aq6k5ef4UHgo87dSG
F0lkcx10kC5Xa6mNKG2Q9BW9wSBo2e1xiv73nyNN
-----END CERTIFICATE-----
EOF
}

variable "private_key" {
	default = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAmKQ1/ZoZkk2VhVTkPz350ieLvXJTFPuY/yTRzJ0CQGgS6nfm
6bW5zZtavKG4CQvqKyfFZruygBz4eZfaz2qtrq21bmB6XUE4ZYroUN6kiD4LCPxN
9FLU3ARqsSo5BqEwJ5jAJmU7u0sfSoxTAcEtZ7qNwxlC8Znyo9OL92vu/520fzf3
d11Ui+LSkVIx5TUtTg8D+vmgH9FSJw4KwDh66tOdJRH5hPkQQa066b0y+DtEtY9O
ShuFcKZPJhWUYfCj52wT7xyhvwo1InbZMbsWys9S5BEx9cdOe7+oPR/DUqIuMWpN
1GBDE163InEEnQ0oeiGlOkgnFLfAPfvzVZ2s8QIDAQABAoIBAEjQdMzwWOh4yC3d
bDBbATRmFvwdcFKfHsH/r3E7KNrOis98uROd0++n/2Xig2cVXvSNOVajjSgeKc3f
ScsOKaIdTWJE9bpMpXmTBPWm77fqWNtFeG5noRD/rmGrMZ7e/5iz/l2Shyb2VAv8
2pAItf84d+2svEmCVcQe4zL5Mv6OYPCruA9EYqWATXNc1vG9iuPzAtMb/zrSVKbc
9KpoSIorBcq0d1Kb8WovRhGQgQxHP/3Ccm/VqMZpgfNNxQSg9Yn7iX3pHcmZbOdy
fAeH31NmqIRTGL/F54Hb+ILxQAtCe7uY5DCRjF3R5zC4AX8dwRNBm4bs884EuWOR
pYNgGY0CgYEAwsChJOkTz1QcVYafH33H1Z5CCFB487BycgSTxZy3NH+CRr4fhP/U
BFSi6z8AiipqahQJwKCEvgDLOBzDZjVh7//UFYpjruQPeoKM75eYbeF6lNUhsFwG
btSjLqyc3QquKeRZQq4Ouuxa1mZkSxr1e/xbb03ICgl+pmnJNcd+UhsCgYEAyKVB
pZ3ulB0lfi7IxOcEjl26LPwk9WOCSY2UVpXqmvqnigZ3i0zWd/kvv1PzFX4yGXyL
7kC3hS60foJPTbR/GAsfl8t+FVTp/MXuIC/n03DWBDO0nQAt3qXDrtLl8NDL8R9o
+XNLHM4cAFqUSxFPFZ3LVI2Q7148FAVSP95GjeMCgYBftGX8S4XupvjdlrBvu1IO
yhzNFS67IoS7P0CXJfJqHBcbSKcYpte74RPG40kSnNF6m6pHPRq+fIlhY9EqUyVz
2ZaRl1ZxRaXNoIY935OKu/mPVkWd8zs+D8S5VR4pCeyYrZynxf17IldpcRvsRK1K
ZrNQOTsKo6vXf7jfcs/C7QKBgDp9umpuZNt2t7RWLR8BfZmHBzwP8TI75QJOLJ0l
LPQq9+ZLxlOsfaUR1nJ/JZDxbedyIFS/NwCzQdjTYgzz/kzjCT22C7ZqP5/5j1aA
wKMp9Kna7N8L61NJnYb8Yh3WsG1FS9PUYWQvTYho32wWyqgxjNHERykQnpDzkCug
P48jAoGARm6kA+MugcTgU1dU7O0/W2NGKKgeTKlzj72itPliBaTIg8GgIFeMSfdO
SGPpkuSYQzClpS5mOkkNHL7YTbKkPpX08/O1V7nri1Tabr9VLIxNRSqTpxCMksu2
pSFqMSjrbbQkxmcUyomKRMWXPovJbxbdL5k6Ht8zOD9U015GXTI=
-----END RSA PRIVATE KEY-----
EOF
}
resource "alicloud_ssl_certificates_service_certificate" "change" {
  		certificate_name = "casdcdnterr_update"
  		cert             = var.cert
  		key              = var.private_key
}

`, name)
}

// Case SetDcdnDomainSSLCertificate替换 7246  raw
func TestAccAliCloudDcdnDomain_basic7246_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnDomainMap7246)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "test02.pfytlm.xyz"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnDomainBasicDependence7246)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "off",
					"domain_name":  name,
					"sources": []map[string]interface{}{
						{
							"type":     "ipaddr",
							"content":  "1.1.1.1",
							"priority": "20",
							"port":     "80",
							"weight":   "20",
						},
					},
					"scope":     "domestic",
					"check_url": "http://test02.pfytlm.xyz/test.html",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "off",
						"domain_name":  name,
						"sources.#":    "1",
						"scope":        "domestic",
						"check_url":    "http://test02.pfytlm.xyz/test.html",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"cert_region":  "cn-hangzhou",
					"cert_type":    "cas",
					"cert_id":      "${alicloud_ssl_certificates_service_certificate.default.id}",
					"cert_name":    "casdcdnterr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"cert_region":  "cn-hangzhou",
						"cert_type":    "cas",
						"cert_id":      CHECKSET,
						"cert_name":    "casdcdnterr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_id":   "${alicloud_ssl_certificates_service_certificate.change.id}",
					"cert_name": "casdcdnterr_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_id":   CHECKSET,
						"cert_name": "casdcdnterr_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url", "env", "ssl_pri", "top_level_domain"},
			},
		},
	})
}

func TestAccAliCloudDcdnDomain_basic7246_change(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnDomainMap7246)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "test02.pfytlm.xyz"
	rand := acctest.RandIntRange(1000000, 9999999)
	certName := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnDomainBasicDependence7246)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":  name,
					"scope":        "domestic",
					"ssl_protocol": "on",
					"sources": []map[string]interface{}{
						{
							"type":     "ipaddr",
							"content":  "1.1.1.1",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
					},
					"cert_type":     "upload",
					"cert_name":     certName,
					"ssl_pub":       "${var.cert}",
					"ssl_pri":       "${var.private_key}",
					"scene":         "apiscene",
					"function_type": "routine",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol":  "on",
						"cert_type":     "upload",
						"cert_name":     CHECKSET,
						"ssl_pub":       CHECKSET,
						"ssl_pri":       CHECKSET,
						"scene":         "apiscene",
						"function_type": "routine",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "off",
					"domain_name":  name,
					"sources": []map[string]interface{}{
						{
							"type":     "ipaddr",
							"content":  "1.1.1.1",
							"priority": "20",
							"port":     "80",
							"weight":   "20",
						},
					},
					"scope":       "domestic",
					"check_url":   "http://test02.pfytlm.xyz/test.html",
					"cert_region": "",
					"cert_id":     "",
					"cert_name":   "",
					"cert_type":   "cas",
					"ssl_pub":     "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "off",
						"domain_name":  name,
						"sources.#":    "1",
						"scope":        "domestic",
						"check_url":    "http://test02.pfytlm.xyz/test.html",
						"cert_region":  "",
						"cert_id":      "",
						"cert_name":    "",
						"cert_type":    "cas",
						"ssl_pub":      "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"cert_region":  "cn-hangzhou",
					"cert_type":    "cas",
					"cert_id":      "${alicloud_ssl_certificates_service_certificate.default.id}",
					"cert_name":    "casdcdnterr",
					"env":          "no-staging",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"cert_region":  "cn-hangzhou",
						"cert_type":    "cas",
						"cert_id":      CHECKSET,
						"cert_name":    "casdcdnterr",
						"ssl_pub":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_id":   "${alicloud_ssl_certificates_service_certificate.change.id}",
					"cert_name": certName + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_id":   CHECKSET,
						"cert_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url", "env", "ssl_pri", "top_level_domain"},
			},
		},
	})
}

// Test Dcdn Domain. <<< Resource test cases, automatically generated.
