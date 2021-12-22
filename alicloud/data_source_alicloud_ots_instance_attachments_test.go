package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOtsInstanceAttachmentsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_instance_attachments.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%d", rand),
		dataSourceOtsInstanceAttachmentsConfigDependence)

	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_instance_attachment.foo.instance_name}",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_instance_attachment.foo.instance_name}",
			"name_regex":    "${alicloud_ots_instance_attachment.foo.vpc_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_instance_attachment.foo.instance_name}",
			"name_regex":    "${alicloud_ots_instance_attachment.foo.vpc_name}-fake",
		}),
	}

	var existOtsInstanceAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                     "1",
			"names.0":                     "testvpc",
			"vpc_ids.#":                   "1",
			"vpc_ids.0":                   CHECKSET,
			"attachments.#":               "1",
			"attachments.0.id":            fmt.Sprintf("tf-testAcc%d", rand),
			"attachments.0.domain":        CHECKSET,
			"attachments.0.endpoint":      CHECKSET,
			"attachments.0.region":        CHECKSET,
			"attachments.0.instance_name": fmt.Sprintf("tf-testAcc%d", rand),
			"attachments.0.vpc_name":      "testvpc",
			"attachments.0.vpc_id":        CHECKSET,
		}
	}

	var fakeOtsInstanceAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":       "0",
			"vpc_ids.#":     "0",
			"attachments.#": "0",
		}
	}

	var otsInstanceAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsInstanceAttachmentsMapFunc,
		fakeMapFunc:  fakeOtsInstanceAttachmentsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsInstanceAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceNameConf, allConf)
}

func dataSourceOtsInstanceAttachmentsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "Capacity"
	}

	data "alicloud_zones" "foo" {
	  available_resource_creation = "VSwitch"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id      = data.alicloud_zones.foo.zones.0.id
	}
	
	resource "alicloud_vswitch" "vswitch" {
	  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	  vpc_id            = data.alicloud_vpcs.default.ids.0
	  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	  zone_id           = data.alicloud_zones.foo.zones.0.id
	  vswitch_name      = var.name
	}
	
	locals {
	  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}
	resource "alicloud_ots_instance_attachment" "foo" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  vpc_name = "testvpc"
	  vswitch_id = local.vswitch_id
	}
	`, name)
}
