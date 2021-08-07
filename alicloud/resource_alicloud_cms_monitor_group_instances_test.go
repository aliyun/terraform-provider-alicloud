package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCmsMonitorGroupInstances_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_monitor_group_instances.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMonitorGroupInstancesMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMonitorGroupInstances")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMonitorGroupInstancesBasicDependence)
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
					"group_id": "${alicloud_cms_monitor_group.default.id}",
					"instances": []map[string]string{
						{
							"category":      "vpc",
							"instance_id":   "${alicloud_vpc.default.id}",
							"instance_name": "tf-testaccvpcname",
							"region_id":     os.Getenv("ALICLOUD_REGION"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":    CHECKSET,
						"instances.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instances": []map[string]string{
						{
							"category":      "vpc",
							"instance_id":   "${alicloud_vpc.default.id}",
							"instance_name": "tf-testaccvpcname",
							"region_id":     os.Getenv("ALICLOUD_REGION"),
						},
						{
							"category":      "slb",
							"instance_id":   "${alicloud_slb_load_balancer.default1.id}",
							"instance_name": "tf-testacccmsslb1",
							"region_id":     os.Getenv("ALICLOUD_REGION"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":    CHECKSET,
						"instances.#": "2",
					}),
				),
			},
		},
	})
}

var AlicloudCmsMonitorGroupInstancesMap = map[string]string{}

func AlicloudCmsMonitorGroupInstancesBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_vpc" "default" {
cidr_block = "192.168.0.0/16"
name = var.name
}
resource "alicloud_cms_monitor_group" "default" {
monitor_group_name = var.name
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  load_balancer_spec = "slb.s2.small"
  vswitch_id = data.alicloud_vswitches.default.ids.0
}
resource "alicloud_slb_load_balancer" "default1" {
  load_balancer_name = "${var.name}1"
  load_balancer_spec = "slb.s2.small"
  vswitch_id = data.alicloud_vswitches.default.ids.0
}
`, name)
}
