package dcdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var dcdnDomainMap = map[string]string{
	"resource_group_id": acctest.CHECKSET,
	"ssl_protocol":      "off",
	"scope":             "overseas",
	"status":            acctest.CHECKSET,
	"cname":             acctest.CHECKSET,
}

func dcdnDomainBasicDependence(name string) string {
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

func TestAccAliCloudDcdnDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, dcdnDomainMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, dcdnDomainBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"scope":       "overseas",
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
						"ssl_pub":      acctest.CHECKSET,
						"ssl_pri":      acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, dcdnDomainMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, dcdnDomainBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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

// TestAccAliCloudDcdnDomain_noScope verifies that omitting sources.weight
// from the configuration does not produce plan drift: the server-side
// default is read back and accepted because the field is declared Computed.
//
// scope is set to "overseas" because the API defaults an unset scope to
// "domestic", which requires ICP filing (unavailable for test domains).
// The scope Computed behaviour is exercised by basic0/basic1, which set
// scope=overseas and pass without drift. The second step here is PlanOnly
// to confirm that the read-back weight produces an empty plan (no drift).
func TestAccAliCloudDcdnDomain_noScope(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"scope":  acctest.CHECKSET,
		"status": acctest.CHECKSET,
	})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, dcdnDomainBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create the domain with scope=overseas (domestic default needs
				// ICP) and WITHOUT sources.weight. The server populates weight
				// with its default; Computed prevents drift.
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
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name,
						"scope":       acctest.CHECKSET,
						"sources.#":   "1",
					}),
				),
			},
			{
				// Re-apply the identical config: the plan must be empty.
				ResourceName: resourceId,
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
						},
					},
				}),
				PlanOnly: true,
			},
		},
	})
}
