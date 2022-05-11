package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEipanycastAnycastEipAddressAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eipanycast_anycast_eip_address_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEipanycastAnycastEipAddressAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipanycastService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipanycastAnycastEipAddressAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEipanycastAnycastEipAddressAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence)
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
					// "bind_instance_region_id" must be consistent with the region of slb instance.
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address.default.id}",
					"bind_instance_id":        "${alicloud_slb_load_balancer.default.id}",
					"bind_instance_region_id": defaultRegionToTest,
					"bind_instance_type":      "SlbInstance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anycast_id":              CHECKSET,
						"bind_instance_id":        CHECKSET,
						"bind_instance_region_id": defaultRegionToTest,
						"bind_instance_type":      "SlbInstance",
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

var AlicloudEipanycastAnycastEipAddressAttachmentMap = map[string]string{
	"bind_time": CHECKSET,
}

func AlicloudEipanycastAnycastEipAddressAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_slb_load_balancer" "default" {
	address_type = "intranet"
	vswitch_id = data.alicloud_vswitches.default.ids[0]
	load_balancer_name = var.name
	load_balancer_spec = "slb.s1.small"
    master_zone_id = "${data.alicloud_slb_zones.default.zones.0.id}"
}

resource "alicloud_eipanycast_anycast_eip_address" "default" {
  anycast_eip_address_name = "${var.name}"
  service_location = "international"
}

`, name)
}
