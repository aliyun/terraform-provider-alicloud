package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudAlbLoadBalancerCommonBandwidthPackageAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancerCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sAlbLoadBalancerCommonBandwidthPackageAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentBasicDependence)
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
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"load_balancer_id":     "${alicloud_alb_load_balancer.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": CHECKSET,
						"load_balancer_id":     CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentMap = map[string]string{}

func AlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_alb_zones" "default"{}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "vswitch_1" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = "${var.name}_1"
}

resource "alicloud_vswitch" "vswitch_2" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
  zone_id           = data.alicloud_zones.default.zones.1.id
  vswitch_name      = "${var.name}_2"
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id =              alicloud_vpc.default.id
  address_type =        "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name =    var.name
  load_balancer_edition = "Standard"
  load_balancer_billing_config {
    pay_type = 	"PayAsYouGo"
  }
  zone_mappings{
		vswitch_id = alicloud_vswitch.vswitch_1.id
		zone_id =  data.alicloud_alb_zones.default.zones.0.id
	}
  zone_mappings{
		vswitch_id = alicloud_vswitch.vswitch_2.id
		zone_id =   data.alicloud_alb_zones.default.zones.1.id
	}
}

resource "alicloud_common_bandwidth_package" "default" {
  		bandwidth            = 3
  		internet_charge_type = "PayByBandwidth"
}

`, name)
}
