package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Nlb LoadbalancerCommonBandwidthPackageAttachment. >>> Resource test cases, automatically generated.
// Case 4619
func TestAccAliCloudNlbLoadbalancerCommonBandwidthPackageAttachment_basic4619(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentMap4619)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadbalancerCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancercommonbandwidthpackageattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentBasicDependence4619)
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
					"load_balancer_id":     "${alicloud_nlb_load_balancer.nlb.id}",
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.cbwp.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":     CHECKSET,
						"bandwidth_package_id": CHECKSET,
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

var AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentMap4619 = map[string]string{
	"status": CHECKSET,
}

func AlicloudNlbLoadbalancerCommonBandwidthPackageAttachmentBasicDependence4619(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "vswtich" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.0.1.0/24"
}

resource "alicloud_vswitch" "vswtich2" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "10.0.2.0/24"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich.id
    zone_id    = alicloud_vswitch.vswtich.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich2.id
    zone_id    = alicloud_vswitch.vswtich2.zone_id
  }
  load_balancer_type = "Network"
  vpc_id             = alicloud_vpc.vpc.id
  address_type       = "Internet"
  address_ip_version = "Ipv4"
}

resource "alicloud_common_bandwidth_package" "cbwp" {
  description          = "nlb-tf-test"
  bandwidth            = "1000"
  internet_charge_type = "PayByBandwidth"
}


`, name)
}

// Test Nlb LoadbalancerCommonBandwidthPackageAttachment. <<< Resource test cases, automatically generated.
