package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaIpSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_ip_set.default"
	ra := resourceAttrInit(resourceId, AlicloudGaIpSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaIpSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudGaIpSetBasicDependence)
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
					"accelerate_region_id": defaultRegionToTest,
					"bandwidth":            "5",
					"accelerator_id":       "${data.alicloud_ga_accelerators.default.ids.0}",
					"depends_on":           []string{"alicloud_ga_bandwidth_package_attachment.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerate_region_id": defaultRegionToTest,
						"bandwidth":            "5",
						"accelerator_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accelerator_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": `10`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
		},
	})
}

var AlicloudGaIpSetMap = map[string]string{
	"ip_address_list.#": CHECKSET,
	"ip_version":        "IPv4",
	"status":            "active",
}

func AlicloudGaIpSetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_ga_accelerators" "default"{
}
data "alicloud_ga_bandwidth_packages" "default"{
}
resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = "${data.alicloud_ga_accelerators.default.ids.0}"
  bandwidth_package_id = "${data.alicloud_ga_bandwidth_packages.default.ids.0}"
}`, name)
}
