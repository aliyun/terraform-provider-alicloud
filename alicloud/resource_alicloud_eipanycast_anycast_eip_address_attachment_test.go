package alicloud

import (
	"fmt"
	"os"
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
			testAccPreCheckWithSlbInstanceSetting(t)
			testAccPreCheckWithRegions(t, true, connectivity.EipanycastSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "bind_instance_region_id" must be consistent with the region of slb instance.
					"anycast_id":              "${alicloud_eipanycast_anycast_eip_address.default.id}",
					"bind_instance_id":        os.Getenv("ALICLOUD_SLB_INSTANCE_ID"),
					"bind_instance_region_id": "cn-hongkong",
					"bind_instance_type":      "SlbInstance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anycast_id":              CHECKSET,
						"bind_instance_id":        os.Getenv("ALICLOUD_SLB_INSTANCE_ID"),
						"bind_instance_region_id": "cn-hongkong",
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

resource "alicloud_eipanycast_anycast_eip_address" "default" {
  anycast_eip_address_name = "${var.name}"
  service_location = "international"
}

`, name)
}
