package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOtsInstanceAttachmentsDataSource_basic(t *testing.T) {
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
	resource "alicloud_vpc" "foo" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "foo" {
	  vpc_id = "${alicloud_vpc.foo.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.foo.zones.0.id}"
	}
	resource "alicloud_ots_instance_attachment" "foo" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  vpc_name = "testvpc"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, name)
}
