package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCmsMetricRuleBlackList_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_black_list.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMetricRuleBlackListMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMetricRuleBlackList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsMetricRuleBlackList%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMetricRuleBlackListBasicDependence)
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
					"instances": []string{
						"{\\\"instanceId\\\":\\\"${alicloud_instance.instance.id}\\\"}"},
					"metrics": []map[string]interface{}{
						{
							"metric_name": "disk_utilization",
						},
					},
					"category":                    "ecs",
					"enable_end_time":             "1640608200000",
					"namespace":                   "acs_ecs_dashboard",
					"enable_start_time":           "1640237400000",
					"metric_rule_black_list_name": "henghai1342432432432",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instances.#":                 "1",
						"metrics.#":                   "1",
						"category":                    "ecs",
						"enable_end_time":             "1640608200000",
						"namespace":                   "acs_ecs_dashboard",
						"enable_start_time":           "1640237400000",
						"metric_rule_black_list_name": "henghai1342432432432",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			}},
	})
}

var AlicloudCmsMetricRuleBlackListMap = map[string]string{}

func AlicloudCmsMetricRuleBlackListBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
  availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_instance_types" "new" {
	eni_amount = 2
    availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
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
resource "alicloud_security_group" "default" {
  	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "instance" {
  	image_id = "${data.alicloud_images.default.images.0.id}"
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = "${alicloud_security_group.default.*.id}"
  	internet_charge_type = "PayByTraffic"
  	internet_max_bandwidth_out = "10"
  	availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  	instance_charge_type = "PostPaid"
  	system_disk_category = "cloud_efficiency"
  	vswitch_id = data.alicloud_vswitches.default.ids[0]
}`, name)
}
