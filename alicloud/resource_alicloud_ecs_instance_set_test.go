package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSInstanceSet_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_instance_set.default"
	ra := resourceAttrInit(resourceId, AlicloudECSInstanceSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInstanceSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secs%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSInstanceSetBasicDependence0)
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
					"amount":                     "10",
					"image_id":                   "${data.alicloud_images.default.images[0].id}",
					"instance_type":              "${data.alicloud_instance_types.default.instance_types[0].id}",
					"instance_name":              "${var.name}",
					"security_group_ids":         []string{"${data.alicloud_security_groups.default.groups.0.id}"},
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "10",
					"zone_id":                    "${data.alicloud_zones.default.zones[0].id}",
					"instance_charge_type":       "PostPaid",
					"system_disk_category":       "cloud_efficiency",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id":                   CHECKSET,
						"instance_type":              CHECKSET,
						"instance_name":              CHECKSET,
						"security_group_ids.#":       "1",
						"internet_charge_type":       "PayByTraffic",
						"internet_max_bandwidth_out": "10",
						"zone_id":                    CHECKSET,
						"instance_charge_type":       "PostPaid",
						"system_disk_category":       "cloud_efficiency",
						"vswitch_id":                 CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "Test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "Test2",
						"Step":    "Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "3",
					}),
				),
			},
		},
	})
}

var AlicloudECSInstanceSetMap0 = map[string]string{}

func AlicloudECSInstanceSetBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_security_groups" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}
