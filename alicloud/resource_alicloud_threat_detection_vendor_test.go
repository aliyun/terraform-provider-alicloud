package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionVendor_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_vendor.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionVendorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSiemService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionVendor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionVendor%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionVendorBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Shanghai})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vendor_name": "${var.name}",
					"lang":        "en",
					"role_for":    0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vendor_name": CHECKSET,
						"lang":        "en",
						"role_for":    "0",
						"vendor_id":   CHECKSET,
						"vendor_type": CHECKSET,
						"create_time": CHECKSET,
						"update_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vendor_name": "${var.name}_update",
					"lang":        "zh",
					"role_for":    0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vendor_name": CHECKSET,
						"lang":        "zh",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "role_for"},
			},
		},
	})
}

var AlicloudThreatDetectionVendorMap = map[string]string{}

func AlicloudThreatDetectionVendorBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
