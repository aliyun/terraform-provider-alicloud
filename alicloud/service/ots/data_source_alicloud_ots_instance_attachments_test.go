package ots_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudOtsInstanceAttachmentsDataSourceBasic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_instance_attachments.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%d", rand),
		dataSourceOtsInstanceAttachmentsConfigDependence)

	instanceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_instance_attachment.foo.instance_name}",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_instance_attachment.foo.instance_name}",
			"name_regex":    "${alicloud_ots_instance_attachment.foo.vpc_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_instance_attachment.foo.instance_name}",
			"name_regex":    "${alicloud_ots_instance_attachment.foo.vpc_name}-fake",
		}),
	}

	var existOtsInstanceAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                     "1",
			"names.0":                     "testvpc",
			"vpc_ids.#":                   "1",
			"vpc_ids.0":                   acctest.CHECKSET,
			"attachments.#":               "1",
			"attachments.0.id":            fmt.Sprintf("tf-testAcc%d", rand),
			"attachments.0.domain":        acctest.CHECKSET,
			"attachments.0.endpoint":      acctest.CHECKSET,
			"attachments.0.region":        acctest.CHECKSET,
			"attachments.0.instance_name": fmt.Sprintf("tf-testAcc%d", rand),
			"attachments.0.vpc_name":      "testvpc",
			"attachments.0.vpc_id":        acctest.CHECKSET,
		}
	}

	var fakeOtsInstanceAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":       "0",
			"vpc_ids.#":     "0",
			"attachments.#": "0",
		}
	}

	var otsInstanceAttachmentsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOtsInstanceAttachmentsMapFunc,
		FakeMapFunc:  fakeOtsInstanceAttachmentsMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsInstanceAttachmentsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceNameConf, allConf)
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
	resource "alicloud_vpc" "default" {
	  vpc_name   = "${var.name}"
	}
	data "alicloud_vswitches" "default" {
		vpc_id = alicloud_vpc.default.id
		zone_id      = data.alicloud_zones.foo.zones.0.id
	}
	
	resource "alicloud_vswitch" "vswitch" {
	  vpc_id            = alicloud_vpc.default.id
	  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
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
