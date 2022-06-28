package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSLBAttachmentsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_attachment.default.load_balancer_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_attachment.default.load_balancer_id}"`,
			"instance_ids":     `["${alicloud_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_attachment.default.load_balancer_id}"`,
			"instance_ids":     `["${alicloud_instance.default.id}_fake"]`,
		}),
	}

	var existSLBAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_attachments.#":             "1",
			"slb_attachments.0.instance_id": CHECKSET,
			"slb_attachments.0.weight":      "42",
		}
	}

	var fakeSLBAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_attachments.#": "0",
		}
	}

	var slbAttachmentCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_attachments.default",
		existMapFunc: existSLBAttachmentsMapFunc,
		fakeMapFunc:  fakeSLBAttachmentsMapFunc,
	}

	slbAttachmentCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

func testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSlbAttachmentsDataSourceBasic-%d"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "default" {
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = local.vswitch_id
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  image_id = "${data.alicloud_images.default.images.0.id}"

  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = local.vswitch_id
}

resource "alicloud_slb_attachment" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  instance_ids = ["${alicloud_instance.default.id}"]
  weight = 42
}

data "alicloud_slb_attachments" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
