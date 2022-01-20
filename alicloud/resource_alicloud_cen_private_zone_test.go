package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenPrivateZone_basic(t *testing.T) {
	var v cbn.PrivateZoneInfo
	resourceId := "alicloud_cen_private_zone.default"
	ra := resourceAttrInit(resourceId, CenPrivateZoneMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenPrivateZone")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenPrivateZone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenPrivateZoneBasicdependence)
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
					"access_region_id": defaultRegionToTest,
					"cen_id":           "${alicloud_cen_instance_attachment.default.instance_id}",
					"host_region_id":   defaultRegionToTest,
					"host_vpc_id":      "${alicloud_cen_instance_attachment.default.child_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_region_id": defaultRegionToTest,
						"cen_id":           CHECKSET,
						"host_region_id":   defaultRegionToTest,
						"host_vpc_id":      CHECKSET,
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

var CenPrivateZoneMap = map[string]string{
	"status": CHECKSET,
}

func CenPrivateZoneBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_cen_instance" "default" {
	name = "${var.name}"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
resource "alicloud_cen_instance_attachment" "default" {
	instance_id = "${alicloud_cen_instance.default.id}"
	child_instance_id = "${data.alicloud_vpcs.default.ids.0}"
	child_instance_type = "VPC"
	child_instance_region_id = "%s"
}
`, name, defaultRegionToTest)
}
