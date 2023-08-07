package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Nlb LoadbalancerCommonBandwidthPackageAttachment. >>> Resource test cases, automatically generated.
// Case 3488
func TestAccAlicloudNlbLoadbalancerCommonBandwidthPackageAttachment_basic3488(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentMap3488)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadbalancerCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancercommonbandwidthpackageattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentBasicDependence3488)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"load_balancer_id":     "${alicloud_nlb_load_balancer.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": CHECKSET,
						"load_balancer_id":     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentMap3488 = map[string]string{
	"status": CHECKSET,
}

func AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentBasicDependence3488(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_vswitches" "default_1" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_nlb_zones.default.zones.0.id
}

data "alicloud_vswitches" "default_2" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_nlb_zones.default.zones.1.id
}

data "alicloud_vswitches" "default_3" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_nlb_zones.default.zones.2.id
}

locals {
	zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
	vswitch_id_1 = data.alicloud_vswitches.default_1.ids[0]
	zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
	vswitch_id_2 = data.alicloud_vswitches.default_2.ids[0]
	zone_id_3    = data.alicloud_nlb_zones.default.zones.2.id
	vswitch_id_3 = data.alicloud_vswitches.default_3.ids[0]
}

resource "alicloud_common_bandwidth_package" "default" {
	bandwidth            = 2
	internet_charge_type = "PayByBandwidth"
	name                 = "${var.name}"
	description          = "${var.name}_description"
}

resource "alicloud_nlb_load_balancer" "default" {
	load_balancer_name = var.name
	resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
	load_balancer_type = "Network"
	address_type       = "Internet"
	address_ip_version = "Ipv4"
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_mappings {
		vswitch_id = local.vswitch_id_1
		zone_id    = local.zone_id_1
	}
	zone_mappings {
		vswitch_id = local.vswitch_id_2
		zone_id    = local.zone_id_2
	}
}

`, name)
}

// Test Nlb LoadbalancerCommonBandwidthPackageAttachment. <<< Resource test cases, automatically generated.
