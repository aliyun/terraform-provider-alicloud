package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudVODDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudVODDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svoddomain%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVODDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckVodDomainRegistered(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":      "136.chat",
					"check_url":        "http://www.aliyun.com/index.html",
					"top_level_domain": "aliyun.com",
					"sources": []map[string]interface{}{
						{
							"source_type":     "oss",
							"source_content":  "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com",
							"source_port":     "80",
							"source_priority": "20",
						},
					},
					"scope": "domestic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":      "136.chat",
						"scope":            "domestic",
						"sources.#":        "1",
						"check_url":        "http://www.aliyun.com/index.html",
						"top_level_domain": "aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"tftestacc":    "TFTEST",
						"Tftestacc123": "Tftest123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":            "2",
						"tags.tftestacc":    "TFTEST",
						"tags.Tftestacc123": "Tftest123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_url":        "http://www.aliyun.com/index2.html",
					"top_level_domain": "aliyun.com.cn",
					"sources": []map[string]interface{}{
						{
							"source_type":     "ipaddr",
							"source_content":  "1.1.1.1",
							"source_port":     "443",
							"source_priority": "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":      "136.chat",
						"scope":            "domestic",
						"sources.#":        "1",
						"check_url":        "http://www.aliyun.com/index2.html",
						"top_level_domain": "aliyun.com.cn",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "tags", "check_url", "ssl_pri", "ssl_pub", "env", "cert_type"},
			},
		},
	})
}

func TestAccAliCloudVODDomain_ssl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudVODDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%svoddomain%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVODDomainSSLBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckVodDomainRegistered(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": name,
					"scope":       "domestic",
					"sources": []map[string]interface{}{
						{
							"source_type":    "domain",
							"source_content": "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com",
							"source_port":    "80",
						},
					},
					"ssl_protocol": "on",
					"ssl_pub":      "${var.public_key}",
					"ssl_pri":      "${var.private_key}",
					"cert_name":    "vodsslcert",
					"cert_type":    "upload",
					"env":          "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":  name,
						"ssl_protocol": "on",
						"ssl_pub":      CHECKSET,
						"ssl_pri":      CHECKSET,
						"cert_name":    "vodsslcert",
						"cert_type":    "upload",
						"env":          "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"ssl_pub":      "${var.public_key}",
					"ssl_pri":      "${var.private_key}",
					"cert_name":    "vodsslcert_update",
					"cert_type":    "upload",
					"env":          "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"ssl_pub":      CHECKSET,
						"ssl_pri":      CHECKSET,
						"cert_name":    "vodsslcert_update",
						"cert_type":    "upload",
						"env":          "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "off",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "tags", "check_url", "ssl_pri", "ssl_pub", "env", "cert_type"},
			},
		},
	})
}

func TestAccAliCloudVODDomain_certCas(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudVODDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%svodcas%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVODDomainCasBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckVodDomainRegistered(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": name,
					"scope":       "domestic",
					"sources": []map[string]interface{}{
						{
							"source_type":    "domain",
							"source_content": "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com",
							"source_port":    "80",
						},
					},
					"ssl_protocol": "on",
					"cert_type":    "cas",
					"cert_region":  "cn-hangzhou",
					"cert_id":      "${alicloud_ssl_certificates_service_certificate.default.id}",
					"cert_name":    "vodcasdomain",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":  name,
						"ssl_protocol": "on",
						"cert_type":    "cas",
						"cert_region":  "cn-hangzhou",
						"cert_name":    "vodcasdomain",
						"cert_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"cert_type":    "cas",
					"cert_region":  "ap-southeast-1",
					"cert_id":      "${alicloud_ssl_certificates_service_certificate.change.id}",
					"cert_name":    "vodcasdomain_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"cert_type":    "cas",
						"cert_region":  "ap-southeast-1",
						"cert_name":    "vodcasdomain_update",
						"cert_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "off",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "tags", "check_url", "ssl_pri", "ssl_pub", "env", "cert_type"},
			},
		},
	})
}

// testAccPreCheckVodDomainRegistered skips VOD domain acceptance tests when the
// test account has no pre-registered (ICP-filed + AddVodDomain) domain. VOD
// AddVodDomain API requires the domain to be registered beforehand, which is an
// external manual prerequisite the provider cannot self-provision. Set
// AC_TEST_VOD_DOMAIN_REGISTERED=1 to enable these cases.
func testAccPreCheckVodDomainRegistered(t *testing.T) {
	if v := os.Getenv("AC_TEST_VOD_DOMAIN_REGISTERED"); v == "" {
		t.Skip("skip: requires a pre-registered/ICP-filed VOD domain; set AC_TEST_VOD_DOMAIN_REGISTERED=1 to enable")
	}
}

var AlicloudVODDomainMap0 = map[string]string{}

func AlicloudVODDomainBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func AlicloudVODDomainSSLBasicDependence(name string) string {
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
`, name)
}

func AlicloudVODDomainCasBasicDependence(name string) string {
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

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = "${var.name}-default"
  cert             = var.public_key
  key              = var.private_key
}

resource "alicloud_ssl_certificates_service_certificate" "change" {
  certificate_name = "${var.name}-change"
  cert             = var.public_key
  key              = var.private_key
}
`, name)
}
