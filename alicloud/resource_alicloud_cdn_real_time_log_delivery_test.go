package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCDNRealTimeLogDelivery_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_real_time_log_delivery.default"
	ra := resourceAttrInit(resourceId, AlicloudCDNRealTimeLogDeliveryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCdnRealTimeLogDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scdnrealtimelogdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDNRealTimeLogDeliveryBasicDependence0)
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
					"domain":     "${alicloud_cdn_domain_new.default.domain_name}",
					"project":    "${alicloud_log_project.default.name}",
					"logstore":   "${alicloud_log_store.default.name}",
					"sls_region": defaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":     CHECKSET,
						"project":    name,
						"logstore":   name,
						"sls_region": defaultRegionToTest,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCDNRealTimeLogDeliveryMap0 = map[string]string{}

func AlicloudCDNRealTimeLogDeliveryBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_cdn_domain_new" "default" {
  domain_name = "%s"
  cdn_type = "web"
  scope = "overseas"
  sources {
	 content = "www.aliyuntest.com"
	 type = "domain"
	 priority = 20
	 port = 80
	 weight = 10
  }
}
resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}
resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
`, name, fmt.Sprintf("%s.example.com", name))
}
