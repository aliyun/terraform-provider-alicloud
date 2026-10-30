package sae_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sae"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// package_type = Image
func TestAccAliCloudSAESlbInternet_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_load_balancer_internet.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSAESLBInternetApplicationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sae.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApplicationSlb")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSAESlbInternetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":          "${alicloud_sae_application.default.id}",
					"internet_slb_id": "${alicloud_slb_load_balancer.default.id}",
					"internet": []map[string]interface{}{
						{
							"protocol":    "TCP",
							"port":        "80",
							"target_port": "8080",
						},
						{
							"protocol":    "TCP",
							"port":        "89",
							"target_port": "8989",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          acctest.CHECKSET,
						"internet_slb_id": acctest.CHECKSET,
						"internet.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":          "${alicloud_sae_application.default.id}",
					"internet_slb_id": "${alicloud_slb_load_balancer.default.id}",
					"internet": []map[string]interface{}{
						{
							"protocol":    "TCP",
							"port":        "90",
							"target_port": "9090",
						},
						{
							"protocol":    "TCP",
							"port":        "99",
							"target_port": "9999",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          acctest.CHECKSET,
						"internet_slb_id": acctest.CHECKSET,
						"internet.#":      "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":          "${alicloud_sae_application.default.id}",
					"internet_slb_id": "${alicloud_slb_load_balancer.secondary.id}",
					"internet": []map[string]interface{}{
						{
							"protocol":      "HTTPS",
							"port":          "80",
							"target_port":   "8080",
							"https_cert_id": "${alicloud_slb_server_certificate.default.id}",
						},
						{
							"protocol":      "HTTPS",
							"port":          "89",
							"target_port":   "8989",
							"https_cert_id": "${alicloud_slb_server_certificate.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":          acctest.CHECKSET,
						"internet_slb_id": acctest.CHECKSET,
						"internet.#":      "2",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id"},
			},
		},
	})
}

func AlicloudSAESlbInternetBasicDependence0(name string) string {
	return fmt.Sprintf(`
data "alicloud_vpcs" "default"	{
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "internet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = data.alicloud_vswitches.default.vswitches[0].id
}
resource "alicloud_slb_load_balancer" "secondary" {
  load_balancer_name = var.name
  address_type       = "internet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = data.alicloud_vswitches.default.vswitches[0].id
}

resource "alicloud_sae_namespace" "default" {
	namespace_description = var.name
	namespace_id = "%s:%s"
	namespace_name = var.name
}
resource "alicloud_sae_application" "default" {
  app_description= var.name
  app_name=        var.name
  namespace_id=    alicloud_sae_namespace.default.namespace_id
  image_url=     "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type=    "Image"
  jdk=             "Open JDK 8"
  vswitch_id=      data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone = "Asia/Shanghai"
  replicas=        "5"
  cpu=             "500"
  memory =          "2048"
}


variable "name" {
  default = "%s"
}
`, acctest.DefaultRegionToTest, name, name) + slbServerCertificateHCL
}

var AlicloudSAESLBInternetApplicationMap0 = map[string]string{}

const slbServerCertificateHCL = `
resource "alicloud_slb_server_certificate" "default" {
  name               = var.name
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDdjCCAl4CCQCcm+erkcKN7DANBgkqhkiG9w0BAQsFADB9MQswCQYDVQQGEwJj\nbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmFsaXl1\nbjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0LmNvbTEaMBgGCSqGSIb3\nDQEJARYLMTIzQDEyMy5jb20wHhcNMTkwNDI2MDM0ODAxWhcNMjQwNDI1MDM0ODAx\nWjB9MQswCQYDVQQGEwJjbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcx\nDzANBgNVBAoMBmFsaXl1bjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0\nLmNvbTEaMBgGCSqGSIb3DQEJARYLMTIzQDEyMy5jb20wggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQDKMKF5qmN/uoMjdH3D8aPRcUOA0s8rZpYhG8zbkF1j\n8gHYoB/FDvM7G7dfVsyjbMwLOxKvAhWvHHSpEz/t7gB+QdwrAMiMJwGmtCnXrh2E\nWiXgalMe1y4a/T5R7q+m4T1zFATf+kbnHWfkSGF4W7b6UBoaH+9StQ95CnqzNf/2\np/Of7+S0XzCxFXw8GIVzZk0xFe6lHJzaq06f3mvzrD+4rpO56tTUvrgTY/n61gsF\nZP7f0CJ2JQh6eNRFOEUSfxKu/Dy/+IsQxorCJY2Q59ZAf3rXrqDN104jw9PlwnLl\nqfZz3RMODN6BWjxE8rvRtT0qMfuAfv1gjBdWZN0hUYBRAgMBAAEwDQYJKoZIhvcN\nAQELBQADggEBAABzo82TxGp5poVkd5pLWj5ACgcBv8Cs6oH9D+4Jz9BmyuBUsQXh\n2aG0hQAe1mU61C9konsl/GTW8umJQ4M4lYEztXXwMf5PlBMGwebM0ZbSGg6jKtZg\nWCgJ3eP/FMmyXGL5Jji5+e09eObhUDVle4tdi0On97zBoz85W02rgWFAqZJwiEAP\nt+c7jX7uOSBq2/38iGStlrX5yB1at/gJXXiA5CL5OtlR3Okvb0/QH37efO1Nu39m\nlFi0ODPAVyXjVypAiLguDxPn6AtDTdk9Iw9B19OD4NrzNRWgSSX5vuxo/VcRcgWk\n3gEe9Ca0ZKN20q9XgthAiFFjl1S9ZgdA6Zc=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAyjCheapjf7qDI3R9w/Gj0XFDgNLPK2aWIRvM25BdY/IB2KAf\nxQ7zOxu3X1bMo2zMCzsSrwIVrxx0qRM/7e4AfkHcKwDIjCcBprQp164dhFol4GpT\nHtcuGv0+Ue6vpuE9cxQE3/pG5x1n5EhheFu2+lAaGh/vUrUPeQp6szX/9qfzn+/k\ntF8wsRV8PBiFc2ZNMRXupRyc2qtOn95r86w/uK6TuerU1L64E2P5+tYLBWT+39Ai\ndiUIenjURThFEn8Srvw8v/iLEMaKwiWNkOfWQH96166gzddOI8PT5cJy5an2c90T\nDgzegVo8RPK70bU9KjH7gH79YIwXVmTdIVGAUQIDAQABAoIBAE1J4a/8biR5S3/W\nG+03BYQeY8tuyjqw8FqfoeOcf9agwAvqybouSNQjeCk9qOQfxq/UWQQFK/zQR9gJ\nv7pX7GBXFK5rkj3g+0SaQhRsPmRFgY0Tl8qGPt2aSKRRNVv5ZeADmwlzRn86QmiF\nMp0rkfqFfDTYWEepZszCML0ouzuxsW/9tq7rvtSjsgATNt31B3vFa3D3JBi31jUh\n5nfR9A3bATze7mQw3byEDiVl5ASRDgYyur403P1fDnMy9DBHZ8NaPOsFF6OBpJal\nBJsG5z00hll5PFN2jfmBQKlvAeU7wfwqdaSnGHOfqf2DeTTaFjIQ4gUhRn/m6pLo\n6kXttLECgYEA9sng0Qz/TcPFfM4tQ1gyvB1cKnnGIwg1FP8sfUjbbEgjaHhA224S\nk3BxtX2Kq6fhTXuwusAFc6OVMAZ76FgrQ5K4Ci7+DTsrF28z4b8td+p+lO/DxgP9\nlTgN+ddsiTOV4fUef9Z3yY0Zr0CnBUMbQYRaV2UIbCdiB0G4V/bt9TsCgYEA0bya\nOo9wGI0RJV0bYP7qwO74Ra1/i1viWbRlS7jU37Q+AZstrlKcQ5CTPzOjKFKMiUzl\n4miWacZ0/q2n+Mvd7NbXGXTLijahnyOYKaHJYyh4oBymfkgAifRstE0Ki9gdvArb\n/I+emC0GvLSyfGN8UUeDJs4NmqdEXGqjo2JOV+MCgYALFv1MR5o9Y1u/hQBRs2fs\nPiGDIx+9OUQxYloccyaxEfjNXAIGGkcpavchIbgWiJ++PJ2vdquIC8TLeK8evL+M\n9M3iX0Q5UfxYvD2HmnCvn9D6Xl/cyRcfGnq+TGjrLW9BzSMGuZt+aiHKV0xqFx7l\nbc4leTvMqGRmURS4lzcQOwKBgQCDzA/i4sYfN25h21tcHXSpnsG3D2rJyQi5NCo/\nZjunA92/JqOTGuiFcLGHEszhhtY3ZXJET1LNz18vtzKJnpqrvOnYXlOVW/U+SqDQ\n8JDb1c/PVZGuY1KrXkR9HLiW3kz5IJ3S3PFdUVYdeTN8BQxXCyg4V12nJJtJs912\ny0zN3wKBgGDS6YttCN6aI4EOABYE8fI1EYQ7vhfiYsaWGWSR1l6bQey7KR6M1ACz\nZzMASNyytVt12yXE4/Emv6/pYqigbDLfL1zQJSLJ3EHJYTh2RxjR+AaGDudYFG/T\nliQ9YXhV5Iu2x1pNwrtFnssDdaaGpfA7l3xC00BL7Z+SAJyI4QKA\n-----END RSA PRIVATE KEY-----"
}
`
