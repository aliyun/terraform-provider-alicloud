package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEsaCustomHostname_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_hostname.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaCustomHostnameMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomHostname")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-esa-ch%d.com", rand)
	hostname := fmt.Sprintf("custom.%s", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaCustomHostnameBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.default.id}",
					"hostname":    hostname,
					"record_id":   "${alicloud_esa_record.default.id}",
					"ssl_flag":    "on",
					"cert_type":   "free",
					"certificate": "",
					"private_key": "",
					"cas_id":      0,
					"cas_region":  "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hostname":    hostname,
						"ssl_flag":    "on",
						"cert_type":   "free",
						"status":      CHECKSET,
						"hostname_id": CHECKSET,
						"site_id":     CHECKSET,
						"record_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.default.id}",
					"hostname":    hostname,
					"record_id":   "${alicloud_esa_record.default.id}",
					"ssl_flag":    "off",
					"cert_type":   "free",
					"certificate": "",
					"private_key": "",
					"cas_id":      0,
					"cas_region":  "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_flag": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.default.id}",
					"hostname":    hostname,
					"record_id":   "${alicloud_esa_record.default.id}",
					"ssl_flag":    "on",
					"cert_type":   "free",
					"certificate": "",
					"private_key": "",
					"cas_id":      0,
					"cas_region":  "",
				}),
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
		},
	})
}

func TestAccAliCloudEsaCustomHostname_datasource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-esa-ds%d.com", rand)
	hostname := fmt.Sprintf("custom.%s", name)
	resourceId := "alicloud_esa_custom_hostname.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaCustomHostnameBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.default.id}",
					"hostname":    hostname,
					"record_id":   "${alicloud_esa_record.default.id}",
					"ssl_flag":    "on",
					"cert_type":   "free",
					"certificate": "",
					"private_key": "",
					"cas_id":      0,
					"cas_region":  "",
				}) + `

data "alicloud_esa_custom_hostnames" "default" {
  site_id = alicloud_esa_site.default.id
  ids     = [alicloud_esa_custom_hostname.default.id]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_custom_hostnames.default", "hostnames.#", "1"),
				),
			},
		},
	})
}

func AlicloudEsaCustomHostnameBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "default" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "default" {
  site_name   = var.name
  instance_id = alicloud_esa_rate_plan_instance.default.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_record" "default" {
  site_id     = alicloud_esa_site.default.id
  record_name = "www.${var.name}"
  record_type = "CNAME"
  source_type = "S3"
  data {
    value = "www.example.com"
  }
  biz_name    = "api"
  host_policy = "follow_hostname"
  ttl         = "100"
  auth_conf {
    secret_key  = "hijklmnhijklmnhijklmnhijklmn"
    version     = "v4"
    region      = "us-east-1"
    auth_type   = "private"
    access_key  = "abcdefgabcdefgabcdefgabcdefg"
  }
}
`, name)
}

var AlicloudEsaCustomHostnameMap = map[string]string{
	"hostname":    CHECKSET,
	"ssl_flag":    CHECKSET,
	"cert_type":   CHECKSET,
	"status":      CHECKSET,
	"hostname_id": CHECKSET,
	"site_id":     CHECKSET,
	"record_id":   CHECKSET,
}
